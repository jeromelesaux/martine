package gfx

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

func AnalyzeTilemap(mode uint8, filename, picturePath string, in image.Image, exportType *export.ExportType, criteria common.AnalyseTilemapOption) error {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := convert.Resize(in, mapSize, exportType.ResizingAlgo)
	var palette color.Palette
	var err error
	palette, m, err = convert.DowngradingPalette(m, mapSize, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", exportType.InputPath)
	}
	refPalette := constants.CpcOldPalette
	if exportType.CpcPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = convert.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = convert.DowngradingWithPalette(m, palette)
	file.PalToPng(exportType.OutputPath+"/palette.png", palette)
	file.Png(exportType.OutputPath+"/map.png", m)
	size := constants.Size{Width: 8, Height: 8}
	var sizeIteration int
	switch mode {
	case 0:
		sizeIteration += 8
	case 1:
		sizeIteration += 4
	case 2:
		sizeIteration += 2
	}
	boards := make([]*transformation.AnalyzeBoard, 0)
	for size.Width < 30 || size.Height < 30 {
		fmt.Printf("Analyse the image for size [width:%d,height:%d]", size.Width, size.Height)
		board := transformation.AnalyzeTilesBoard(m, size)
		fmt.Printf(" found [%d] tiles\n", len(board.BoardTiles))
		boards = append(boards, board)
		size.Width += sizeIteration
		size.Height += sizeIteration
	}

	// analyze the results and apply criteria
	var lowerSizeIndex, numberTilesIndex int
	lowerSizeValue := math.MaxInt64
	numberTilesValue := math.MaxInt64
	for i, v := range boards {
		if len(v.BoardTiles) < numberTilesValue {
			numberTilesIndex = i
			numberTilesValue = len(v.BoardTiles)
		}
		tilesSize := sizeOctet(v.TileSize, mode) * len(v.BoardTiles)
		if tilesSize < lowerSizeValue {
			tilesSize = lowerSizeValue
			lowerSizeIndex = 1
		}
	}
	var choosenBoard *transformation.AnalyzeBoard
	switch criteria {
	case common.NumberTilemapOption:
		choosenBoard = boards[numberTilesIndex]
		fmt.Printf("choose the [%d]board with number tiles [%d] and size [width:%d, height:%d]\n", numberTilesIndex, len(choosenBoard.BoardTiles), choosenBoard.TileSize.Width, choosenBoard.TileSize.Height)
	case common.SizeTilemapOption:
		choosenBoard = boards[lowerSizeIndex]
		fmt.Printf("choose the [%d]board with number tiles [%d] and size [width:%d, height:%d]\n", numberTilesIndex, len(choosenBoard.BoardTiles), choosenBoard.TileSize.Width, choosenBoard.TileSize.Height)
	default:
		return errors.ErrorCriteriaNotFound
	}
	if err := choosenBoard.SaveSchema(filepath.Join(exportType.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := choosenBoard.SaveTilemap(filepath.Join(exportType.OutputPath, "tilesmap.map")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap csv file error :%v\n", err)
		return err
	}

	// applyOneImage
	// sort tiles
	// check < 256 tiles
	// finally export
	// 20 tiles large 25 tiles height
	tiles := choosenBoard.Sort()
	data := make([]byte, 0)

	finalFile := strings.ReplaceAll(filename, "?", "")
	if err := file.Kit(finalFile, palette, mode, false, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	os.Mkdir(exportType.OutputPath+string(filepath.Separator)+"tiles", os.ModePerm)
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, err := ApplyOneImage(tile,
				exportType,
				int(mode),
				palette,
				mode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while transforming sprite error : %v\n", err)
			}
			data = append(data, d...)
			scenePath := filepath.Join(exportType.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
			f, err := os.Create(scenePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create scene tile-%.2d error %v\n", i, err)
				return err
			}

			if err := png.Encode(f, tile); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot encode in png scene tile-%.2d error %v\n", i, err)
				return err
			}
			f.Close()
			nbFrames++
		}
	}
	// save the file sprites
	finalFile = strings.ReplaceAll(filename, "?", "")
	if err := file.Imp(data, uint(nbFrames), uint(choosenBoard.TileSize.Width), uint(choosenBoard.TileSize.Height), uint(mode), finalFile, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", picturePath, err)
	}

	// save the tilemap
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	switch choosenBoard.TileSize.Width {
	case 8:
		nbTilePixelLarge = 20
	case 4:
		nbTilePixelLarge = 40
	}
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(exportType.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
	index := 0
	for y := 0; y < m.Bounds().Max.Y; y += (nbTilePixelHigh * choosenBoard.TileSize.Height) {
		for x := 0; x < m.Bounds().Max.X; x += (nbTilePixelLarge * choosenBoard.TileSize.Width) {
			m1 := image.NewNRGBA(image.Rect(0, 0, nbTilePixelLarge*choosenBoard.TileSize.Width, nbTilePixelHigh*choosenBoard.TileSize.Height))
			// copy of the map
			for i := 0; i < nbTilePixelLarge*choosenBoard.TileSize.Width; i++ {
				for j := 0; j < nbTilePixelHigh*choosenBoard.TileSize.Height; j++ {
					var c color.Color = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					if x+i < m.Bounds().Max.X && y+j < m.Bounds().Max.Y {
						c = m.At(x+i, y+j)
					}
					m1.Set(i, j, c)
				}
			}
			// store the map in the slice
			scenes = append(scenes, m1)
			scenePath := filepath.Join(exportType.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
			f, err := os.Create(scenePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create scene scence-%.2d error %v\n", index, err)
				return err
			}

			if err := png.Encode(f, m1); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot encode in png scene scene-%.2d error %v\n", index, err)
				return err
			}
			f.Close()
			index++
		}
	}

	// now thread all maps images
	tileMaps := make([]byte, 0)
	for _, v := range scenes {
		for y := 0; y < v.Bounds().Max.Y; y += choosenBoard.TileSize.Height {
			for x := 0; x < v.Bounds().Max.X; x += choosenBoard.TileSize.Width {
				sprt, err := transformation.ExtractTile(v, choosenBoard.TileSize, x, y)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
					break
				}
				index := choosenBoard.TileIndex(sprt, tiles)
				tileMaps = append(tileMaps, byte(index))
			}
		}
	}

	if err := file.TileMap(tileMaps, finalFile, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-TileMap the image %s error %v", picturePath, err)
	}
	return err
}

func sizeOctet(size constants.Size, mode uint8) int {
	switch mode {
	case 0:
		return size.Height * (size.Width / 2)
	case 1:
		return size.Height * (size.Width / 4)
	case 2:
		return size.Height * (size.Width / 8)
	}
	return 0
}

func Tilemap(mode uint8, filename, picturePath string, size constants.Size, in image.Image, exportType *export.ExportType) error {
	/*
		8x8 : 40x25
		16x8 : 20x25
		16x16 : 20x24
	*/

	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	maxTiles := 255
	nbPixelWidth := 0
	switch mode {
	case 0:
		nbPixelWidth = exportType.Size.Width / 2
	case 1:
		nbPixelWidth = exportType.Size.Width / 4
	case 2:
		nbPixelWidth = exportType.Size.Width / 8
	default:
		fmt.Fprintf(os.Stderr, "Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 8 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	if exportType.Size.Height != 16 && exportType.Size.Height != 8 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	switch exportType.Size.Width {
	case 8:
		nbTilePixelLarge = 20
		if exportType.Size.Height == 16 {
			maxTiles = 240
		}
	case 4:
		nbTilePixelLarge = 40
	}

	if !exportType.CustomDimension {
		fmt.Fprintf(os.Stderr, "You must set height and width to define the tile dimensions (options -h and -w) error:%v\n", errors.ErrorCustomDimensionMustBeSet)
		return errors.ErrorCustomDimensionMustBeSet
	}
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := convert.Resize(in, mapSize, exportType.ResizingAlgo)
	var palette color.Palette
	var err error
	palette, m, err = convert.DowngradingPalette(m, mapSize, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", exportType.InputPath)
	}
	refPalette := constants.CpcOldPalette
	if exportType.CpcPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = convert.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = convert.DowngradingWithPalette(m, palette)
	file.PalToPng(exportType.OutputPath+"/palette.png", palette)
	file.Png(exportType.OutputPath+"/map.png", m)

	analyze := transformation.AnalyzeTilesBoard(m, exportType.Size)
	if err := analyze.SaveSchema(filepath.Join(exportType.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := analyze.SaveTilemap(filepath.Join(exportType.OutputPath, "tilesmap.map")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap csv file error :%v\n", err)
		return err
	}

	// applyOneImage
	// sort tiles
	// check < 256 tiles
	// finally export
	// 20 tiles large 25 tiles height
	tiles := analyze.Sort()
	data := make([]byte, 0)

	finalFile := strings.ReplaceAll(filename, "?", "")
	if err := file.Kit(finalFile, palette, mode, false, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	os.Mkdir(exportType.OutputPath+string(filepath.Separator)+"tiles", os.ModePerm)
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, err := ApplyOneImage(tile,
				exportType,
				int(mode),
				palette,
				mode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while transforming sprite error : %v\n", err)
			}
			data = append(data, d...)
			scenePath := filepath.Join(exportType.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
			f, err := os.Create(scenePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create scene tile-%.2d error %v\n", i, err)
				return err
			}

			if err := png.Encode(f, tile); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot encode in png scene tile-%.2d error %v\n", i, err)
				return err
			}
			f.Close()
			if i >= maxTiles {
				fmt.Fprintf(os.Stderr, "Maximum of %d tiles accepted, skipping...\n", maxTiles)
				break
			}
			nbFrames++
		}
	}
	// save the file sprites
	finalFile = strings.ReplaceAll(filename, "?", "")
	if err := file.Imp(data, uint(nbFrames), uint(analyze.TileSize.Width), uint(analyze.TileSize.Height), uint(mode), finalFile, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", picturePath, err)
	}

	// save the tilemap
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(exportType.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
	index := 0
	for y := 0; y < m.Bounds().Max.Y; y += (nbTilePixelHigh * analyze.TileSize.Height) {
		for x := 0; x < m.Bounds().Max.X; x += (nbTilePixelLarge * analyze.TileSize.Width) {
			m1 := image.NewNRGBA(image.Rect(0, 0, nbTilePixelLarge*analyze.TileSize.Width, nbTilePixelHigh*analyze.TileSize.Height))
			// copy of the map
			for i := 0; i < nbTilePixelLarge*analyze.TileSize.Width; i++ {
				for j := 0; j < nbTilePixelHigh*analyze.TileSize.Height; j++ {
					var c color.Color = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					if x+i < m.Bounds().Max.X && y+j < m.Bounds().Max.Y {
						c = m.At(x+i, y+j)
					}
					m1.Set(i, j, c)
				}
			}
			// store the map in the slice
			scenes = append(scenes, m1)
			scenePath := filepath.Join(exportType.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
			f, err := os.Create(scenePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create scene scence-%.2d error %v\n", index, err)
				return err
			}

			if err := png.Encode(f, m1); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot encode in png scene scene-%.2d error %v\n", index, err)
				return err
			}
			f.Close()
			index++
		}
	}

	// now thread all maps images
	tileMaps := make([]byte, 0)
	for _, v := range scenes {
		for y := 0; y < v.Bounds().Max.Y; y += analyze.TileSize.Height {
			for x := 0; x < v.Bounds().Max.X; x += analyze.TileSize.Width {
				sprt, err := transformation.ExtractTile(v, analyze.TileSize, x, y)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
					break
				}
				index := analyze.TileIndex(sprt, tiles)
				tileMaps = append(tileMaps, byte(index))
			}
		}
	}

	if err := file.TileMap(tileMaps, finalFile, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-TileMap the image %s error %v", picturePath, err)
	}
	return err
}
