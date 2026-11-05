package menu

import (
	"fmt"
	"image"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	w "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/sprite"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
)

var TileSize float32 = 20.

type TilemapMenu struct {
	*ImageMenu
	Result       *transformation.AnalyzeBoard
	TileImages   *w.ImageTable
	ExportZigzag bool
	Historic     *sprite.TilesHistorical
	Col          int
	Row          int
	tileImage    image.Image
	currentTile  transformation.Tile
}

func (me *TilemapMenu) CurrentTile() transformation.Tile {
	return me.currentTile
}

func (me *TilemapMenu) ResetExport() {
	me.Cfg.Reset()
}

func NewTilemapMenu() *TilemapMenu {
	t := &TilemapMenu{
		ImageMenu: NewImageMenu(),
		Result:    &transformation.AnalyzeBoard{},
	}
	t.TileImages = w.NewEmptyImageTable(
		fyne.NewSize(TileSize, TileSize),
		t.TileSelected,
	)
	return t
}

func (me *TilemapMenu) CmdLine() string {
	exec, err := os.Executable()
	if err != nil {
		log.GetLogger().Error("error while getting executable path :%v\n", err)
		return exec
	}
	if me.OriginalImagePath() != "" {
		exec += constants.InOptionLabel + me.OriginalImagePath()
	}
	if me.Cfg.ScrCfg.IsPlus {
		exec += constants.PlusOptionLabel
	}
	if me.Cfg.ScrCfg.Type.IsFullScreen() {
		exec += constants.FullscreenOptionLabel
	}
	if me.Cfg.ScrCfg.Type.IsSprite() {
		width, widthText, err := me.GetWidth()
		if err != nil {
			log.GetLogger().Error("cannot convert width value :%s error :%v\n", widthText, err)
		} else {
			exec += constants.WidthOptionLabel + strconv.Itoa(width)
		}
		height, heightText, err := me.GetHeight()
		if err != nil {
			log.GetLogger().Error("cannot convert height value :%s error :%v\n", heightText, err)
		} else {
			exec += constants.HeightOptionLabel + strconv.Itoa(height)
		}
	}
	if me.Cfg.ScrCfg.Type.IsSpriteHard() {
		exec += constants.SpritehardOptionLabel
	}
	if me.Cfg.ScrCfg.Process.ApplyDithering {
		if me.Cfg.ScrCfg.Process.Dithering.WithQuantification {
			exec += constants.QuantizationOptionLabel
		} else {
			exec += constants.MultiplierOptionLabel + fmt.Sprintf("%.2f", me.Cfg.ScrCfg.Process.Dithering.Multiplier)
		}
		exec += constants.DitheringOptionLabel + strconv.Itoa(me.Cfg.ScrCfg.Process.Dithering.Algo)
		// stockage du numéro d'algo
	}
	exec += constants.ModeOptionLabel + strconv.Itoa(int(me.Cfg.ScrCfg.Mode))
	if me.Cfg.ScrCfg.Process.Reducer != 0 {
		exec += constants.ReducerOptionLabel + strconv.Itoa(me.Cfg.ScrCfg.Process.Reducer)
	}
	// resize algo
	if me.ResizeAlgoNumber != 0 {
		exec += " -algo " + strconv.Itoa(me.ResizeAlgoNumber)
	}
	if me.Cfg.ScrCfg.Process.Brightness != 0 {
		exec += " -brightness " + fmt.Sprintf("%.2f", me.Cfg.ScrCfg.Process.Brightness)
	}
	if me.Cfg.ScrCfg.Process.Saturation != 0 {
		exec += " -saturation " + fmt.Sprintf("%.2f", me.Cfg.ScrCfg.Process.Saturation)
	}

	exec += " -tilemap"
	me.CmdLineGenerate = exec
	return exec
}

func (me *TilemapMenu) TileSelected(row, col int) {
	if row < 0 || col < 0 {
		return
	}
	me.Row = row
	me.Col = col
	if me.Result == nil {
		return
	}
	if row < len(me.Result.Tiles) && col < len(me.Result.Tiles[0]) {
		tile := me.Result.Tiles[row][col]
		me.currentTile = *transformation.TileFromImage(tile.(*image.NRGBA))
		me.tileImage = tile
	}
}

func (me *TilemapMenu) TileImage() image.Image {
	return me.tileImage
}

func (me *TilemapMenu) SetNewTilesImages(tiles [][]image.Image) {
	tilesCanvas := w.NewImageTableCache(len(tiles), len(tiles[0]), fyne.NewSize(50, 50), nil)
	for i, v := range tiles {
		for i2, v2 := range v {
			tilesCanvas.Set(i, i2, canvas.NewImageFromImage(v2))
		}
	}
	me.TileImages.Update(tilesCanvas, len(tiles)-1, len(tiles[0])-1)
	me.TileImages.Refresh()
}
