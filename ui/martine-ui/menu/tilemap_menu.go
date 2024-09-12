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
	Result                 *transformation.AnalyzeBoard
	TileImages             *w.ImageTable
	ExportDsk              bool
	ExportText             bool
	ExportWithAmsdosHeader bool
	ExportZigzag           bool
	ExportJson             bool
	ExportCompression      int
	ExportFolderPath       string
	ExportImpdraw          bool
	ExportFlat             bool
	Historic               *sprite.TilesHistorical
}

func (tm *TilemapMenu) ResetExport() {
	tm.ExportDsk = false
	tm.ExportText = false
	tm.ExportWithAmsdosHeader = false
	tm.ExportZigzag = false
	tm.ExportJson = false
	tm.ExportCompression = -1
	tm.ExportImpdraw = false
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
	if i.IsCpcPlus {
		exec += " -plus"
	}
	if i.Format.IsFullScreen() {
		exec += " -fullscreen"
	}
	if i.Format.IsSprite() {
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
	if i.Format.IsSpriteHard() {
		exec += " -spritehard"
	}
	if i.ApplyDithering {
		if i.WithQuantification {
			exec += " -quantization"
		} else {
			exec += " -multiplier " + fmt.Sprintf("%.2f", i.DitheringMultiplier)
		}
		exec += " -dithering " + strconv.Itoa(i.DitheringAlgoNumber)
		// stockage du numéro d'algo
	}
	exec += " -mode " + strconv.Itoa(i.Mode)
	if i.Reducer != 0 {
		exec += " -reducer " + strconv.Itoa(i.Reducer)
	}
	// resize algo
	if i.ResizeAlgoNumber != 0 {
		exec += " -algo " + strconv.Itoa(i.ResizeAlgoNumber)
	}
	if i.Brightness != 0 {
		exec += " -brightness " + fmt.Sprintf("%.2f", i.Brightness)
	}
	if i.Saturation != 0 {
		exec += " -saturation " + fmt.Sprintf("%.2f", i.Saturation)
	}

	exec += " -tilemap"
	i.CmdLineGenerate = exec
	return exec
}
