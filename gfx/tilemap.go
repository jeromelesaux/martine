package gfx

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/log"

	ci "github.com/jeromelesaux/martine/convert/image"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/impdraw/tile"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/export/sprite"
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

// nolint: funlen, gocognit
func AnalyzeTilemap(mode uint8, isCpcPlus bool, filename, picturePath string, in image.Image, cfg *config.MartineConfig, criteria common.AnalyseTilemapOption) error {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := ci.Resize(in, mapSize, cfg.ScrCfg.Process.ResizingAlgo)
	var palette color.Palette
	var err error
	palette = ci.ExtractPalette(m, isCpcPlus, cfg.ScrCfg.Size.ColorsAvailable)
	refPalette := constants.CpcOldPalette
	if cfg.ScrCfg.IsPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = ci.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = ci.DowngradingWithPalette(m, palette)
	err = png.PalToPng(cfg.ScrCfg.OutputPath+"/palette.png", palette)
	if err != nil {
		return err
	}
	err = png.Png(cfg.ScrCfg.OutputPath+"/map.png", m)
	if err != nil {
		return err
	}
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
			log.GetLogger().Info("Analyse the image for size [width:%d,height:%d]", size.Width, size.Height)
			board := transformation.AnalyzeTilesBoard(m, size, nil)
			tilesSize := sizeOctet(board.TileSize, mode) * len(board.BoardTiles)
			log.GetLogger().Info(" found [%d] tiles full length [#%X]\n", len(board.BoardTiles), tilesSize)
			boards = append(boards, board)
			size.Width += sizeIteration
			size.Height += sizeIteration
		}
	} else {
		for _, s := range impTilesSizes {
			log.GetLogger().Info("Analyse the image for size [width:%d,height:%d]", s.Width, s.Height)
			board := transformation.AnalyzeTilesBoard(m, s, nil)
			tilesSize := sizeOctet(board.TileSize, mode) * len(board.BoardTiles)
			log.GetLogger().Info(" found [%d] tiles full length [#%X]\n", len(board.BoardTiles), tilesSize)
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
		log.GetLogger().Info("choose the [%d]board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", numberTilesIndex, len(choosenBoard.BoardTiles), choosenBoard.TileSize.Width, choosenBoard.TileSize.Height, tilesSize)
	case common.SizeTilemapOption:
		choosenBoard = boards[lowerSizeIndex]
		tilesSize := sizeOctet(choosenBoard.TileSize, mode) * len(choosenBoard.BoardTiles)
		log.GetLogger().Info("choose the [%d]board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", lowerSizeIndex, len(choosenBoard.BoardTiles), choosenBoard.TileSize.Width, choosenBoard.TileSize.Height, tilesSize)
	default:
		return errors.ErrorCriteriaNotFound
	}
	cfg.ScrCfg.Size.Width = choosenBoard.BoardTiles[0].Tile.Size.Width
	cfg.ScrCfg.Size.Height = choosenBoard.BoardTiles[0].Tile.Size.Height
	cfg.CustomDimension = true
	if err := choosenBoard.SaveSchema(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap_schema.png")); err != nil {
		log.GetLogger().Error("Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := choosenBoard.SaveAssemblyTiles(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap.asm")); err != nil {
		log.GetLogger().Error("Cannot save tilemap csv file error :%v\n", err)
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
	if err := impPalette.Kit(finalFile, palette, mode, false, cfg); err != nil {
		log.GetLogger().Error("Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "tiles"), os.ModePerm)
	if err != nil {
		return err
	}
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, _, err := ApplyOneImage(tile,
				cfg,
				palette,
				mode)
			if err != nil {
				log.GetLogger().Error("Error while transforming sprite error : %v\n", err)
			}
			if nbFrames < 255 {
				data = append(data, d...)
				scenePath := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
				if err := png.Png(scenePath, tile); err != nil {
					log.GetLogger().Error("Cannot create scene tile-%.2d error %v\n", i, err)
					return err
				}
			} else {
				log.GetLogger().Infoln("skipping number of frames exceed 255.")
				break
			}
			nbFrames++
		}
	}
	// save the file sprites
	finalFile = strings.ReplaceAll(filename, "?", "")
	if err := tile.Imp(data, uint(nbFrames), uint(choosenBoard.TileSize.Width), uint(choosenBoard.TileSize.Height), uint(mode), finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", picturePath, err)
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
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "scenes"), os.ModePerm)
	if err != nil {
		return err
	}
	scenes, err = fillImage(m, scenes, nbTilePixelHigh, nbTilePixelLarge, choosenBoard.TileSize.Width, choosenBoard.TileSize.Height, cfg.ScrCfg.OutputPath)
	if err != nil {
		return err
	}

	// now thread all maps images
	tileMaps := make([]byte, 0)
	for _, v := range scenes {
		for y := 0; y < v.Bounds().Max.Y; y += choosenBoard.TileSize.Height {
			for x := 0; x < v.Bounds().Max.X; x += choosenBoard.TileSize.Width {
				sprt, err := transformation.ExtractTile(v, choosenBoard.TileSize, x, y)
				if err != nil {
					log.GetLogger().Error("Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
					break
				}
				index := choosenBoard.TileIndex(sprt, tiles)
				tileMaps = append(tileMaps, byte(index))
			}
		}
	}

	if err := tile.TileMap(tileMaps, finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-TileMap the image %s error %v", picturePath, err)
	}
	return err
}

func ExportTilemapClassical(m image.Image, filename string, palette color.Palette, board *transformation.AnalyzeBoard, size constants.Size, cfg *config.MartineConfig) error {
	finalFile := strings.ReplaceAll(filename, "?", "")
	if err := board.SaveSchema(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap_schema.png")); err != nil {
		log.GetLogger().Error("Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := board.SaveAssemblyTiles(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap.asm")); err != nil {
		log.GetLogger().Error("Cannot save tilemap csv file error :%v\n", err)
		return err
	}

	if err := impPalette.Kit(finalFile, palette, cfg.ScrCfg.Mode, false, cfg); err != nil {
		log.GetLogger().Error("Error while saving file %s error :%v", finalFile, err)
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
	err := os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "scenes"), os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fillImage(m, scenes, nbTilePixelHigh, nbTilePixelLarge, board.TileSize.Width, board.TileSize.Height, cfg.ScrCfg.OutputPath)
	if err != nil {
		return err
	}

	if cfg.ScrCfg.Type == config.ScreenFormat(config.SpriteFlatExport) {
		if err := board.SaveFlatFile(cfg.ScrCfg.OutputPath, palette, cfg.ScrCfg.Mode, cfg); err != nil {
			return err
		}
	}

	if cfg.ScrCfg.Type == config.ScreenFormat(config.SpriteExport) {
		if err := board.SaveOcpWindowFile(cfg.ScrCfg.OutputPath, palette, cfg.ScrCfg.Mode, cfg); err != nil {
			return err
		}
	}

	// now save the tiles index table
	tileMaps := board.TilesIndexTable()

	if err := tile.TileMap(tileMaps, finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-TileMap the image %s error %v", cfg.ScrCfg.OutputPath, err)
		return err
	}
	return board.SaveHistoric(cfg.ScrCfg.OutputPath)
}

func TilemapClassical(mode uint8, isCpcPlus bool, filename, picturePath string, in image.Image, size constants.Size, cfg *config.MartineConfig, tilesHistoric *sprite.TilesHistorical) (*transformation.AnalyzeBoard, [][]image.Image, color.Palette) {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := ci.Resize(in, mapSize, cfg.ScrCfg.Process.ResizingAlgo)
	var palette color.Palette
	palette = ci.ExtractPalette(m, isCpcPlus, cfg.ScrCfg.Size.ColorsAvailable)
	refPalette := constants.CpcOldPalette
	if cfg.ScrCfg.IsPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = ci.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = ci.DowngradingWithPalette(m, palette)
	tilemap := transformation.AnalyzeTilesBoard(m, size, tilesHistoric)
	var tilesImagesTilemap [][]image.Image
	for y := 0; y < m.Bounds().Max.Y; y += tilemap.TileSize.Height {
		tilesmap := make([]image.Image, 0)
		for x := 0; x < m.Bounds().Max.X; x += tilemap.TileSize.Width {
			sprt, err := transformation.ExtractTile(m, tilemap.TileSize, x, y)
			if err != nil {
				log.GetLogger().Error("Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
				break
			}
			tilesmap = append(tilesmap, sprt.Image())
		}
		tilesImagesTilemap = append(tilesImagesTilemap, tilesmap)
	}
	tilemap.Tiles = tilesImagesTilemap
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

// nolint: funlen, gocognit
func Tilemap(mode uint8, filename, picturePath string, size constants.Size, in image.Image, cfg *config.MartineConfig, tilesHistoric *sprite.TilesHistorical) error {
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
		nbPixelWidth = cfg.ScrCfg.Size.Width / 2
	case 1:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 4
	case 2:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 8
	default:
		log.GetLogger().Error("Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		log.GetLogger().Error("%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	if cfg.ScrCfg.Size.Height != 16 && cfg.ScrCfg.Size.Height != 8 {
		log.GetLogger().Error("%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	switch cfg.ScrCfg.Size.Width {
	case 4:
		nbTilePixelLarge = 20
		if cfg.ScrCfg.Size.Height == 16 {
			maxTiles = 240
		}
	case 2:
		nbTilePixelLarge = 40
	}

	if !cfg.CustomDimension {
		log.GetLogger().Error("You must set height and width to define the tile dimensions (options -h and -w) error:%v\n", errors.ErrorCustomDimensionMustBeSet)
		return errors.ErrorCustomDimensionMustBeSet
	}
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	m := ci.Resize(in, mapSize, cfg.ScrCfg.Process.ResizingAlgo)
	var palette color.Palette
	var err error
	palette, m, err = ci.DowngradingPalette(m, mapSize, true)
	if err != nil {
		log.GetLogger().Error("Cannot downgrade colors palette for this image %s\n", cfg.ScrCfg.InputPath)
	}
	refPalette := constants.CpcOldPalette
	if cfg.ScrCfg.IsPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = ci.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = ci.DowngradingWithPalette(m, palette)
	err = png.PalToPng(cfg.ScrCfg.OutputPath+"/palette.png", palette)
	if err != nil {
		return err
	}
	err = png.Png(cfg.ScrCfg.OutputPath+"/map.png", m)
	if err != nil {
		return err
	}
	analyze := transformation.AnalyzeTilesBoard(m, cfg.ScrCfg.Size, tilesHistoric)
	tilesSize := sizeOctet(analyze.TileSize, mode) * len(analyze.BoardTiles)
	log.GetLogger().Info("board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", len(analyze.BoardTiles), analyze.TileSize.Width, analyze.TileSize.Height, tilesSize)
	if err := analyze.SaveSchema(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap_schema.png")); err != nil {
		log.GetLogger().Error("Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err := analyze.SaveAssemblyTiles(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap.asm")); err != nil {
		log.GetLogger().Error("Cannot save tilemap csv file error :%v\n", err)
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
	if err := impPalette.Kit(finalFile, palette, mode, false, cfg); err != nil {
		log.GetLogger().Error("Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "tiles"), os.ModePerm)
	if err != nil {
		return err
	}
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, _, err := ApplyOneImage(tile,
				cfg,
				palette,
				mode)
			if err != nil {
				log.GetLogger().Error("Error while transforming sprite error : %v\n", err)
			}
			data = append(data, d...)
			scenePath := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
			if err := png.Png(scenePath, tile); err != nil {
				log.GetLogger().Error("Cannot encode in png scene tile-%.2d error %v\n", i, err)
				return err
			}
			if i >= maxTiles {
				log.GetLogger().Error("Maximum of %d tiles accepted, skipping...\n", maxTiles)
				break
			}
			nbFrames++
		}
	}
	// save the file sprites
	finalFile = strings.ReplaceAll(filename, "?", "")
	if err := tile.Imp(data, uint(nbFrames), uint(analyze.TileSize.Width), uint(analyze.TileSize.Height), uint(mode), finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", picturePath, err)
	}

	// save the tilemap
	scenes := make([]*image.NRGBA, 0)
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "scenes"), os.ModePerm)
	if err != nil {
		return err
	}
	scenes, err = fillImage(m, scenes, nbTilePixelHigh, nbTilePixelLarge, analyze.TileSize.Width, analyze.TileSize.Height, cfg.ScrCfg.OutputPath)
	if err != nil {
		return err
	}

	// now thread all maps images
	tileMaps := make([]byte, 0)
	for _, v := range scenes {
		for y := 0; y < v.Bounds().Max.Y; y += analyze.TileSize.Height {
			for x := 0; x < v.Bounds().Max.X; x += analyze.TileSize.Width {
				sprt, err := transformation.ExtractTile(v, analyze.TileSize, x, y)
				if err != nil {
					log.GetLogger().Error("Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
					break
				}
				index := analyze.TileIndex(sprt, tiles)
				tileMaps = append(tileMaps, byte(index))
			}
		}
	}

	if err := tile.TileMap(tileMaps, finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-TileMap the image %s error %v", picturePath, err)
	}
	return err
}

func fillImage(m image.Image, scenes []*image.NRGBA, nbTilePixelHigh, nbTilePixelLarge, tileWidth, tileHeight int, output string) ([]*image.NRGBA, error) {
	index := 0
	for y := 0; y < m.Bounds().Max.Y; y += (nbTilePixelHigh * tileHeight) {
		for x := 0; x < m.Bounds().Max.X; x += (nbTilePixelLarge * tileWidth) {
			m1 := image.NewNRGBA(image.Rect(0, 0, nbTilePixelLarge*tileWidth, nbTilePixelHigh*tileHeight))
			// copy of the map
			for i := 0; i < nbTilePixelLarge*tileWidth; i++ {
				for j := 0; j < nbTilePixelHigh*tileHeight; j++ {
					var c color.Color = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					if x+i < m.Bounds().Max.X && y+j < m.Bounds().Max.Y {
						c = m.At(x+i, y+j)
					}
					m1.Set(i, j, c)
				}
			}
			// store the map in the slice
			scenes = append(scenes, m1)
			scenePath := filepath.Join(output, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
			if err := png.Png(scenePath, m1); err != nil {
				log.GetLogger().Error("Cannot encode in png scene scene-%.2d error %v\n", index, err)
				return scenes, err
			}
			index++
		}
	}
	return scenes, nil
}

// nolint: funlen
func TilemapRaw(mode uint8, isCpcPlus bool, size constants.Size, in image.Image, cfg *config.MartineConfig, tilesHistoric *sprite.TilesHistorical) (*transformation.AnalyzeBoard, [][]image.Image, color.Palette, error) {
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
		nbPixelWidth = cfg.ScrCfg.Size.Width / 2
	case 1:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 4
	case 2:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 8
	default:
		log.GetLogger().Error("Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		log.GetLogger().Error("%v\n", errors.ErrorWidthSizeNotAccepted)
		return nil, tilesImagesTilemap, palette, errors.ErrorWidthSizeNotAccepted
	}
	if cfg.ScrCfg.Size.Height != 16 && cfg.ScrCfg.Size.Height != 8 {
		log.GetLogger().Error("%v\n", errors.ErrorWidthSizeNotAccepted)
		return nil, tilesImagesTilemap, palette, errors.ErrorWidthSizeNotAccepted
	}
	if !cfg.CustomDimension {
		log.GetLogger().Error("You must set height and width to define the tile dimensions (options -h and -w) error:%v\n", errors.ErrorCustomDimensionMustBeSet)
		return nil, tilesImagesTilemap, palette, errors.ErrorCustomDimensionMustBeSet
	}
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: cfg.ScrCfg.Size.ColorsAvailable}
	m := ci.Resize(in, mapSize, cfg.ScrCfg.Process.ResizingAlgo)
	if isCpcPlus {
		palette, _, _ = ci.DowngradingPalette(m, mapSize, isCpcPlus)
	} else {
		palette = ci.ExtractPalette(m, isCpcPlus, cfg.ScrCfg.Size.ColorsAvailable)
	}
	refPalette := constants.CpcOldPalette
	if cfg.ScrCfg.IsPlus {
		refPalette = constants.CpcPlusPalette
	}
	palette = ci.ToCPCPalette(palette, refPalette)
	palette = constants.SortColorsByDistance(palette)
	_, m = ci.DowngradingWithPalette(m, palette)

	analyze := transformation.AnalyzeTilesBoard(m, cfg.ScrCfg.Size, tilesHistoric)

	// now display all maps images
	for y := 0; y < len(analyze.TileMap); y++ {
		tilesmap := make([]image.Image, 0)
		for x := 0; x < len(analyze.TileMap[y]); x++ {
			i := analyze.TileMap[y][x]
			if i < len(analyze.BoardTiles) {
				tl := analyze.BoardTiles[i]
				tilesmap = append(tilesmap, tl.Tile.Image())
			} else {
				bl := image.NewNRGBA(image.Rect(0, 0, analyze.TileSize.Width, analyze.TileSize.Height))
				draw.Draw(bl, bl.Bounds(), &image.Uniform{color.Alpha16{0xFFF}}, image.Pt(0, 0), draw.Src)
				tilesmap = append(tilesmap, bl)
			}
		}
		tilesImagesTilemap = append(tilesImagesTilemap, tilesmap)
	}
	analyze.Tiles = tilesImagesTilemap
	// applyOneImage
	// sort tiles
	// check < 256 tiles
	// finally export
	// 20 tiles large 25 tiles height
	return analyze, tilesImagesTilemap, palette, nil
}

// nolint: funlen
func ExportTilemap(analyze *transformation.AnalyzeBoard, filename string, palette color.Palette, mode uint8, in image.Image, flatExport bool, cfg *config.MartineConfig) (err error) {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: 16}
	tilesSize := sizeOctet(analyze.TileSize, mode) * len(analyze.BoardTiles)
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25

	log.GetLogger().Info("board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", len(analyze.BoardTiles), analyze.TileSize.Width, analyze.TileSize.Height, tilesSize)
	if err = analyze.SaveSchema(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap_schema.png")); err != nil {
		log.GetLogger().Error("Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err = analyze.SaveAssemblyTiles(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap.asm")); err != nil {
		log.GetLogger().Error("Cannot save tilemap csv file error :%v\n", err)
		return err
	}

	finalFile := strings.ReplaceAll(filename, "?", "")
	if err = impPalette.Kit(finalFile, palette, mode, false, cfg); err != nil {
		log.GetLogger().Error("Error while saving file %s error :%v", finalFile, err)
		return err
	}

	if flatExport {
		if err = analyze.SaveFlatFile(cfg.ScrCfg.OutputPath, palette, mode, cfg); err != nil {
			log.GetLogger().Error("Error while saving sprites in folder %s error :%v", cfg.ScrCfg.OutputPath, err)
		}
	} else {
		if err = analyze.SaveTiles(cfg.ScrCfg.OutputPath, palette, mode, cfg); err != nil {
			log.GetLogger().Error("Error while saving sprites in folder %s error :%v", cfg.ScrCfg.OutputPath, err)
		}
	}
	// scenes := make([]*image.NRGBA, 0)
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "scenes"), os.ModePerm)
	if err != nil {
		return err
	}
	index := 0
	m := ci.Resize(in, mapSize, cfg.ScrCfg.Process.ResizingAlgo)
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
			//	scenes = append(scenes, m1)
			scenePath := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("%sscenes%sscene-%.2d.png", string(filepath.Separator), string(filepath.Separator), index))
			if err := png.Png(scenePath, m1); err != nil {
				log.GetLogger().Error("Cannot encode in png scene scene-%.2d error %v\n", index, err)
				return err
			}
			index++
		}
	}
	// save historic
	return analyze.SaveHistoric(cfg.ScrCfg.OutputPath)
}

// nolint: funlen, gocognit
func ExportImpdrawTilemap(analyze *transformation.AnalyzeBoard, filename string, palette color.Palette, mode uint8, size constants.Size, in image.Image, cfg *config.MartineConfig) (err error) {
	mapSize := constants.Size{Width: in.Bounds().Max.X, Height: in.Bounds().Bounds().Max.Y, ColorsAvailable: cfg.ScrCfg.Size.ColorsAvailable}
	nbTilePixelLarge := 20
	nbTilePixelHigh := 25
	maxTiles := 255
	nbPixelWidth := 0
	switch mode {
	case 0:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 2
	case 1:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 4
	case 2:
		nbPixelWidth = cfg.ScrCfg.Size.Width / 8
	default:
		log.GetLogger().Error("Mode %d  not available\n", mode)
	}

	if nbPixelWidth != 4 && nbPixelWidth != 2 {
		log.GetLogger().Error("%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	if cfg.ScrCfg.Size.Height != 16 && cfg.ScrCfg.Size.Height != 8 {
		log.GetLogger().Error("%v\n", errors.ErrorWidthSizeNotAccepted)
		return errors.ErrorWidthSizeNotAccepted
	}
	switch cfg.ScrCfg.Size.Width {
	case 4:
		nbTilePixelLarge = 20
		if cfg.ScrCfg.Size.Height == 16 {
			maxTiles = 240
		}
	case 2:
		nbTilePixelLarge = 40
	}
	tilesSize := sizeOctet(analyze.TileSize, mode) * len(analyze.BoardTiles)
	log.GetLogger().Info("board with number of tiles [%d] and size [width:%d, height:%d] size:#%X\n", len(analyze.BoardTiles), analyze.TileSize.Width, analyze.TileSize.Height, tilesSize)
	if err = analyze.SaveSchema(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap_schema.png")); err != nil {
		log.GetLogger().Error("Cannot save tilemap schema error :%v\n", err)
		return err
	}
	if err = analyze.SaveAssemblyTiles(filepath.Join(cfg.ScrCfg.OutputPath, "tilesmap.asm")); err != nil {
		log.GetLogger().Error("Cannot save tilemap csv file error :%v\n", err)
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
	if err = impPalette.Kit(finalFile, palette, mode, false, cfg); err != nil {
		log.GetLogger().Error("Error while saving file %s error :%v", finalFile, err)
		return err
	}
	nbFrames := 0
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "tiles"), os.ModePerm)
	if err != nil {
		return err
	}
	for i, v := range tiles {
		if v.Occurence > 0 {
			tile := v.Tile.Image()
			d, _, _, _, err := ApplyOneImage(tile,
				cfg,
				palette,
				mode)
			if err != nil {
				log.GetLogger().Error("Error while transforming sprite error : %v\n", err)
			}
			data = append(data, d...)
			scenePath := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("%stiles%stile-%.2d.png", string(filepath.Separator), string(filepath.Separator), i))
			if err := png.Png(scenePath, tile); err != nil {
				log.GetLogger().Error("Cannot encode in png scene tile-%.2d error %v\n", i, err)
				return err
			}
			if i >= maxTiles {
				log.GetLogger().Error("Maximum of %d tiles accepted, skipping...\n", maxTiles)
				break
			}
			nbFrames++
		}
	}
	// save the file sprites
	finalFile = strings.ReplaceAll(filename, "?", "")
	if err = tile.Imp(data, uint(nbFrames), uint(analyze.TileSize.Width), uint(analyze.TileSize.Height), uint(mode), finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", cfg.ScrCfg.OutputPath, err)
	}
	m := ci.Resize(in, mapSize, cfg.ScrCfg.Process.ResizingAlgo)
	// save the tilemap
	scenes := make([]*image.NRGBA, 0)
	err = os.Mkdir(filepath.Join(cfg.ScrCfg.OutputPath, "scenes"), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = fillImage(m, scenes, nbTilePixelHigh, nbTilePixelLarge, analyze.TileSize.Width, analyze.TileSize.Height, cfg.ScrCfg.OutputPath)
	if err != nil {
		return err
	}

	// now save all the sprites indexes
	tileMaps := analyze.TilesIndexTable()

	if err = tile.TileMap(tileMaps, finalFile, cfg); err != nil {
		log.GetLogger().Error("Cannot export to Imp-TileMap the image %s error %v", cfg.ScrCfg.OutputPath, err)
	}
	return err
}
