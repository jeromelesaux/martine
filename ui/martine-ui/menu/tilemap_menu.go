package menu

import (
	"fmt"
	"image"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	w "github.com/jeromelesaux/fyne-io/widget"
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

func (tm *TilemapMenu) CurrentTile() transformation.Tile {
	return tm.currentTile
}

func (tm *TilemapMenu) ResetExport() {
	tm.Cfg.Reset()
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

func (i *TilemapMenu) CmdLine() string {
	exec, err := os.Executable()
	if err != nil {
		log.GetLogger().Error("error while getting executable path :%v\n", err)
		return exec
	}
	if i.OriginalImagePath() != "" {
		exec += " -in " + i.OriginalImagePath()
	}
	if i.Cfg.ScrCfg.IsPlus {
		exec += " -plus"
	}
	if i.Cfg.ScrCfg.Type.IsFullScreen() {
		exec += " -fullscreen"
	}
	if i.Cfg.ScrCfg.Type.IsSprite() {
		width, widthText, err := i.GetWidth()
		if err != nil {
			log.GetLogger().Error("cannot convert width value :%s error :%v\n", widthText, err)
		} else {
			exec += " -width " + strconv.Itoa(width)
		}
		height, heightText, err := i.GetHeight()
		if err != nil {
			log.GetLogger().Error("cannot convert height value :%s error :%v\n", heightText, err)
		} else {
			exec += " -height " + strconv.Itoa(height)
		}
	}
	if i.Cfg.ScrCfg.Type.IsSpriteHard() {
		exec += " -spritehard"
	}
	if i.Cfg.ScrCfg.Process.ApplyDithering {
		if i.Cfg.ScrCfg.Process.Dithering.WithQuantification {
			exec += " -quantization"
		} else {
			exec += " -multiplier " + fmt.Sprintf("%.2f", i.Cfg.ScrCfg.Process.Dithering.Multiplier)
		}
		exec += " -dithering " + strconv.Itoa(i.Cfg.ScrCfg.Process.Dithering.Algo)
		// stockage du numéro d'algo
	}
	exec += " -mode " + strconv.Itoa(int(i.Cfg.ScrCfg.Mode))
	if i.Cfg.ScrCfg.Process.Reducer != 0 {
		exec += " -reducer " + strconv.Itoa(i.Cfg.ScrCfg.Process.Reducer)
	}
	// resize algo
	if i.ResizeAlgoNumber != 0 {
		exec += " -algo " + strconv.Itoa(i.ResizeAlgoNumber)
	}
	if i.Cfg.ScrCfg.Process.Brightness != 0 {
		exec += " -brightness " + fmt.Sprintf("%.2f", i.Cfg.ScrCfg.Process.Brightness)
	}
	if i.Cfg.ScrCfg.Process.Saturation != 0 {
		exec += " -saturation " + fmt.Sprintf("%.2f", i.Cfg.ScrCfg.Process.Saturation)
	}

	exec += " -tilemap"
	i.CmdLineGenerate = exec
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
