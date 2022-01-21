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

var (
	impTilesSizes = []constants.Size{
		{Width: 4, Height: 8},
		{Width: 8, Height: 8},
		{Width: 8, Height: 16},
	}
	impTileFlag = true
)

func AnalyzeTilemap(mode uint8, isCpcPlus bool, filename, picturePath string, in image.Image, cont *export.MartineContext, criteria common.AnalyseTilemapOption) error {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := convert.Resize(in, mapSize, cont.ResizingAlgo)
	var palette color.Palette
	var err error
	palette = convert.ExtractPalette(m, isCpcPlus, cont.Size.ColorsAvailable)
	refPalette := constants.CpcOldPalette
	if cont.CpcPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = convert.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = convert.DowngradingWithPalette(m, palette)
	file.PalToPng(cont.OutputPath+"/palette.png", palette)
	file.Png(cont.OutputPath+"/map.png", m)
	size := constants.Size{Width: 2, Height: 8}
	var sizeIteration int
	switch mode {
	case 0:
		sizeIteration += 2
	case 1:
		sizeIteration += 4
	case 2:
		sizeIteration += 8
	}
	boards := make([]*transformation.AnalyzeBoard, 0)
	if !impTileFlag {
		for size.Width <= 32 || size.Height <= 32 {
			fmt.Printf("Analyse the image for size [width:%d,height:%d]", size.Width, size.Height)
			board := transformation.AnalyzeTilesBoard(m, size)
			tilesSize := sizeOctet(board.TileSize, mode) * len(board.BoardTiles)
			fmt.Printf(" found [%d] tiles full length [#%X]\n", len(board.BoardTiles), tilesSize)
			boards = append(boards, board)
			size.Width += sizeIteration
			size.Height += sizeIteration
		}
	} else {
		for _, s := range impTilesSizes {
			fmt.Printf("Analyse the image for size [width:%d,height:%d]", s.Width, s.Height)
			board := transformation.AnalyzeTilesBoard(m, s)
			tilesSize := sizeOctet(board.TileSize, mode) * len(board.BoardTiles)
			fmt.Printf(" found [%d] tiles full length [#%X]\n", len(board.BoardTiles), tilesSize)
			boards = append(boards, board)
		}
	}
	// analyze the results and apply criteria
	var lowerSizeIndex, numberTilesIndex int
	lowerSizeValue := math.MaxInt32
	numberTilesValue := math.MaxInt32
	for i, v := range boards {
		if len(v.BoardTiles) < numberTilesValue {
			numberTilesIndex = i
			numberTilesValue = len(v.BoardTiles)
		}
		tilesSize := sizeOctet(v.TileSize, mode) * len(v.BoardTiles)
		if tilesSize < lowerSizeValue {
			lowerSizeValue = tilesSize
			lowerSizeIndex = i
		}
	}
	var choosenBoard *transformation.AnalyzeBoard
	switch criteria {
	case common.NumberTilemapOption:
		choosenBoard = boards[numberTilesIndex]
		tilesSize := sizeOctet(choosenBoard.TileSize, mode) * len(choosenBoard.BoardTiles)
		fmt.Printf("choose the [%d]board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", numberTilesIndex, len(choosenBoard.BoardTiles), choosenBoard.TileSize.Width, choosenBoard.TileSize.Height, tilesSize)
	case common.SizeTilemapOption:
		choosenBoard = boards[lowerSizeIndex]
		tilesSize := sizeOctet(choosenBoard.TileSize, mode) * len(choosenBoard.BoardTiles)
		fmt.Printf("choose the [%d]board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", lowerSizeIndex, len(choosenBoard.BoardTiles), choosenBoard.TileSize.Width, choosenBoard.TileSize.Height, tilesSize)
	default:
		return errors.ErrorCriteriaNotFound
	}
	cont.Size.Width = choosenBoard.BoardTiles[0].Tile.Size.Width
	cont.Size.Height = choosenBoard.BoardTiles[0].Tile.Size.Height
	cont.CustomDimension = true
	if err := choosenBoard.SaveSchema(filepath.Join(cont.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := choosenBoard.SaveTilemap(filepath.Join(cont.OutputPath, "tilesmap.map")); err != nil {
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
	if err := file.Kit(finalFile, palette, mode, false, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"tiles", os.ModePerm)
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, _, err := ApplyOneImage(tile,
				cont,
				int(mode),
				palette,
				mode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while transforming sprite error : %v\n", err)
			}
			if nbFrames < 255 {
				data = append(data, d...)
				scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
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
			} else {
				fmt.Println("skipping number of frames exceed 255.")
				break
			}
			nbFrames++
		}
	}
	// save the file sprites
	finalFile = strings.ReplaceAll(filename, "?", "")
	if err := file.Imp(data, uint(nbFrames), uint(choosenBoard.TileSize.Width), uint(choosenBoard.TileSize.Height), uint(mode), finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", picturePath, err)
	}

	// save the tilemap
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	switch choosenBoard.TileSize.Width {
	case 4:
		nbTilePixelLarge = 20
	case 2:
		nbTilePixelLarge = 40
	}
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
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
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
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

	if err := file.TileMap(tileMaps, finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-TileMap the image %s error %v", picturePath, err)
	}
	return err
}

func ExportTilemapClassical(m image.Image, filename string, board *transformation.AnalyzeBoard, size constants.Size, cont *export.MartineContext) error {
	finalFile := strings.ReplaceAll(filename, "?", "")
	if err := board.SaveSchema(filepath.Join(cont.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := board.SaveTilemap(filepath.Join(cont.OutputPath, "tilesmap.map")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap csv file error :%v\n", err)
		return err
	}
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	switch board.TileSize.Width {
	case 4:
		nbTilePixelLarge = 20
	case 2:
		nbTilePixelLarge = 40
	}
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
	index := 0
	for y := 0; y < m.Bounds().Max.Y; y += (nbTilePixelHigh * board.TileSize.Height) {
		for x := 0; x < m.Bounds().Max.X; x += (nbTilePixelLarge * board.TileSize.Width) {
			m1 := image.NewNRGBA(image.Rect(0, 0, nbTilePixelLarge*board.TileSize.Width, nbTilePixelHigh*board.TileSize.Height))
			// copy of the map
			for i := 0; i < nbTilePixelLarge*board.TileSize.Width; i++ {
				for j := 0; j < nbTilePixelHigh*board.TileSize.Height; j++ {
					var c color.Color = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					if x+i < m.Bounds().Max.X && y+j < m.Bounds().Max.Y {
						c = m.At(x+i, y+j)
					}
					m1.Set(i, j, c)
				}
			}
			// store the map in the slice
			scenes = append(scenes, m1)
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
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
	tiles := board.Sort()
	tileMaps := make([]byte, 0)
	for _, v := range scenes {
		for y := 0; y < v.Bounds().Max.Y; y += board.TileSize.Height {
			for x := 0; x < v.Bounds().Max.X; x += board.TileSize.Width {
				sprt, err := transformation.ExtractTile(v, board.TileSize, x, y)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
					break
				}
				index := board.TileIndex(sprt, tiles)
				tileMaps = append(tileMaps, byte(index))
			}
		}
	}

	if err := file.TileMap(tileMaps, finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-TileMap the image %s error %v", cont.OutputPath, err)
		return err
	}
	return nil
}

func TilemapClassical(mode uint8, isCpcPlus bool, filename, picturePath string, in image.Image, size constants.Size, cont *export.MartineContext) (*transformation.AnalyzeBoard, [][]image.Image, color.Palette) {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := convert.Resize(in, mapSize, cont.ResizingAlgo)
	var palette color.Palette
	palette = convert.ExtractPalette(m, isCpcPlus, cont.Size.ColorsAvailable)
	refPalette := constants.CpcOldPalette
	if cont.CpcPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = convert.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = convert.DowngradingWithPalette(m, palette)
	tilemap := transformation.AnalyzeTilesBoard(m, size)
	var tilesImagesTilemap [][]image.Image
	for y := 0; y < m.Bounds().Max.Y; y += tilemap.TileSize.Height {
		tilesmap := make([]image.Image, 0)
		for x := 0; x < m.Bounds().Max.X; x += tilemap.TileSize.Width {
			sprt, err := transformation.ExtractTile(m, tilemap.TileSize, x, y)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
				break
			}
			tilesmap = append(tilesmap, sprt.Image())
		}
		tilesImagesTilemap = append(tilesImagesTilemap, tilesmap)
	}

	return tilemap, tilesImagesTilemap, palette
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

func Tilemap(mode uint8, filename, picturePath string, size constants.Size, in image.Image, cont *export.MartineContext) error {
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
		nbPixelWidth = cont.Size.Width / 2
	case 1:
		nbPixelWidth = cont.Size.Width / 4
	case 2:
		nbPixelWidth = cont.Size.Width / 8
	default:
		fmt.Fprintf(os.Stderr, "Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	if cont.Size.Height != 16 && cont.Size.Height != 8 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	switch cont.Size.Width {
	case 4:
		nbTilePixelLarge = 20
		if cont.Size.Height == 16 {
			maxTiles = 240
		}
	case 2:
		nbTilePixelLarge = 40
	}

	if !cont.CustomDimension {
		fmt.Fprintf(os.Stderr, "You must set height and width to define the tile dimensions (options -h and -w) error:%v\n", errors.ErrorCustomDimensionMustBeSet)
		return errors.ErrorCustomDimensionMustBeSet
	}
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := convert.Resize(in, mapSize, cont.ResizingAlgo)
	var palette color.Palette
	var err error
	palette, m, err = convert.DowngradingPalette(m, mapSize, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", cont.InputPath)
	}
	refPalette := constants.CpcOldPalette
	if cont.CpcPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = convert.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = convert.DowngradingWithPalette(m, palette)
	file.PalToPng(cont.OutputPath+"/palette.png", palette)
	file.Png(cont.OutputPath+"/map.png", m)

	analyze := transformation.AnalyzeTilesBoard(m, cont.Size)
	tilesSize := sizeOctet(analyze.TileSize, mode) * len(analyze.BoardTiles)
	fmt.Printf("board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", len(analyze.BoardTiles), analyze.TileSize.Width, analyze.TileSize.Height, tilesSize)
	if err := analyze.SaveSchema(filepath.Join(cont.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := analyze.SaveTilemap(filepath.Join(cont.OutputPath, "tilesmap.map")); err != nil {
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
	if err := file.Kit(finalFile, palette, mode, false, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"tiles", os.ModePerm)
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, _, err := ApplyOneImage(tile,
				cont,
				int(mode),
				palette,
				mode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while transforming sprite error : %v\n", err)
			}
			data = append(data, d...)
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
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
	if err := file.Imp(data, uint(nbFrames), uint(analyze.TileSize.Width), uint(analyze.TileSize.Height), uint(mode), finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", picturePath, err)
	}

	// save the tilemap
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
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
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
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

	if err := file.TileMap(tileMaps, finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-TileMap the image %s error %v", picturePath, err)
	}
	return err
}

func TilemapRaw(mode uint8, isCpcPlus bool, size constants.Size, in image.Image, cont *export.MartineContext) (*transformation.AnalyzeBoard, [][]image.Image, color.Palette, error) {
	/*
		8x8 : 40x25
		16x8 : 20x25
		16x16 : 20x24
	*/
	var palette color.Palette
	var tilesImagesTilemap [][]image.Image
	nbPixelWidth := 0

	switch mode {
	case 0:
		nbPixelWidth = cont.Size.Width / 2
	case 1:
		nbPixelWidth = cont.Size.Width / 4
	case 2:
		nbPixelWidth = cont.Size.Width / 8
	default:
		fmt.Fprintf(os.Stderr, "Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return nil, tilesImagesTilemap, palette, errors.ErrorWidthSizeNotAccepted
	}
	if cont.Size.Height != 16 && cont.Size.Height != 8 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return nil, tilesImagesTilemap, palette, errors.ErrorWidthSizeNotAccepted
	}
	if !cont.CustomDimension {
		fmt.Fprintf(os.Stderr, "You must set height and width to define the tile dimensions (options -h and -w) error:%v\n", errors.ErrorCustomDimensionMustBeSet)
		return nil, tilesImagesTilemap, palette, errors.ErrorCustomDimensionMustBeSet
	}
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := convert.Resize(in, mapSize, cont.ResizingAlgo)

	palette = convert.ExtractPalette(m, isCpcPlus, cont.Size.ColorsAvailable)

	refPalette := constants.CpcOldPalette
	if cont.CpcPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = convert.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = convert.DowngradingWithPalette(m, palette)

	analyze := transformation.AnalyzeTilesBoard(m, cont.Size)

	// now thread all maps images
	for y := 0; y < m.Bounds().Max.Y; y += analyze.TileSize.Height {
		tilesmap := make([]image.Image, 0)
		for x := 0; x < m.Bounds().Max.X; x += analyze.TileSize.Width {
			sprt, err := transformation.ExtractTile(m, analyze.TileSize, x, y)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
				break
			}
			tilesmap = append(tilesmap, sprt.Image())
		}
		tilesImagesTilemap = append(tilesImagesTilemap, tilesmap)
	}

	// applyOneImage
	// sort tiles
	// check < 256 tiles
	// finally export
	// 20 tiles large 25 tiles height
	return nil, tilesImagesTilemap, palette, nil
}

func ExportTilemap(analyze *transformation.AnalyzeBoard, filename string, palette color.Palette, mode uint8, in image.Image, cont *export.MartineContext) (err error) {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	tilesSize := sizeOctet(analyze.TileSize, mode) * len(analyze.BoardTiles)
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	nbPixelWidth := 0
	switch mode {
	case 0:
		nbPixelWidth = cont.Size.Width / 2
	case 1:
		nbPixelWidth = cont.Size.Width / 4
	case 2:
		nbPixelWidth = cont.Size.Width / 8
	default:
		fmt.Fprintf(os.Stderr, "Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	if cont.Size.Height != 16 && cont.Size.Height != 8 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	switch cont.Size.Width {
	case 4:
		nbTilePixelLarge = 20
	case 2:
		nbTilePixelLarge = 40
	}
	fmt.Printf("board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", len(analyze.BoardTiles), analyze.TileSize.Width, analyze.TileSize.Height, tilesSize)
	if err = analyze.SaveSchema(filepath.Join(cont.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err = analyze.SaveTilemap(filepath.Join(cont.OutputPath, "tilesmap.map")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap csv file error :%v\n", err)
		return err
	}

	finalFile := strings.ReplaceAll(filename, "?", "")
	if err = file.Kit(finalFile, palette, mode, false, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", finalFile, err)
		return err
	}

	if err = analyze.SaveSprites(cont.OutputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving sprites in folder %s error :%v", cont.OutputPath, err)
	}
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
	index := 0
	m := convert.Resize(in, mapSize, cont.ResizingAlgo)
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
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
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
	return
}

func ExportImpdrawTilemap(analyze *transformation.AnalyzeBoard, filename string, palette color.Palette, mode uint8, size constants.Size, in image.Image, cont *export.MartineContext) (err error) {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	maxTiles := 255
	nbPixelWidth := 0
	switch mode {
	case 0:
		nbPixelWidth = cont.Size.Width / 2
	case 1:
		nbPixelWidth = cont.Size.Width / 4
	case 2:
		nbPixelWidth = cont.Size.Width / 8
	default:
		fmt.Fprintf(os.Stderr, "Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	if cont.Size.Height != 16 && cont.Size.Height != 8 {
		fmt.Fprintf(os.Stderr, "%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	switch cont.Size.Width {
	case 4:
		nbTilePixelLarge = 20
		if cont.Size.Height == 16 {
			maxTiles = 240
		}
	case 2:
		nbTilePixelLarge = 40
	}
	tilesSize := sizeOctet(analyze.TileSize, mode) * len(analyze.BoardTiles)
	fmt.Printf("board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", len(analyze.BoardTiles), analyze.TileSize.Width, analyze.TileSize.Height, tilesSize)
	if err = analyze.SaveSchema(filepath.Join(cont.OutputPath, "tilesmap_schema.png")); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err = analyze.SaveTilemap(filepath.Join(cont.OutputPath, "tilesmap.map")); err != nil {
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
	if err = file.Kit(finalFile, palette, mode, false, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"tiles", os.ModePerm)
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, _, err := ApplyOneImage(tile,
				cont,
				int(mode),
				palette,
				mode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while transforming sprite error : %v\n", err)
			}
			data = append(data, d...)
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
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
	if err = file.Imp(data, uint(nbFrames), uint(analyze.TileSize.Width), uint(analyze.TileSize.Height), uint(mode), finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", cont.OutputPath, err)
	}
	m := convert.Resize(in, mapSize, cont.ResizingAlgo)
	// save the tilemap
	scenes := make([]*image.NRGBA, 0)
	os.Mkdir(cont.OutputPath+string(filepath.Separator)+"scenes", os.ModePerm)
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
			scenePath := filepath.Join(cont.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
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

	if err = file.TileMap(tileMaps, finalFile, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot export to Imp-TileMap the image %s error %v", cont.OutputPath, err)
	}
	return err
}
