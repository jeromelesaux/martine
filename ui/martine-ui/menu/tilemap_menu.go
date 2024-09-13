package menu

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
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
}

func (tm *TilemapMenu) ResetExport() {
	tm.Cfg.Reset()
}

func NewTilemapMenu() *TilemapMenu {
	return &TilemapMenu{
		ImageMenu:  NewImageMenu(),
		Result:     &transformation.AnalyzeBoard{},
		TileImages: w.NewEmptyImageTable(fyne.NewSize(TileSize, TileSize)),
	}
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
