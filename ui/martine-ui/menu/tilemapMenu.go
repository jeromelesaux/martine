package menu

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

var TileSize float32 = 20.

type TilemapMenu struct {
	ImageMenu
	Result                 *transformation.AnalyzeBoard
	TileImages             *custom_widget.ImageTable
	ExportDsk              bool
	ExportText             bool
	ExportWithAmsdosHeader bool
	ExportZigzag           bool
	ExportJson             bool
	ExportCompression      int
	ExportFolderPath       string
	ExportImpdraw          bool
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
		Result:     &transformation.AnalyzeBoard{},
		TileImages: custom_widget.NewEmptyImageTable(fyne.NewSize(TileSize, TileSize)),
	}
}

func (i *TilemapMenu) CmdLine() string {
	exec, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while getting executable path :%v\n", err)
		return exec
	}
	if i.OriginalImagePath != nil {
		exec += " -in " + i.OriginalImagePath.Path()
	}
	if i.IsCpcPlus {
		exec += " -plus"
	}
	if i.IsFullScreen {
		exec += " -fullscreen"
	}
	if i.IsSprite {
		width, err := strconv.Atoi(i.Width.Text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot convert width value :%s error :%v\n", i.Width.Text, err)
		} else {
			exec += " -width " + strconv.Itoa(width)
		}
		height, err := strconv.Atoi(i.Height.Text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot convert height value :%s error :%v\n", i.Height.Text, err)
		} else {
			exec += " -height " + strconv.Itoa(height)
		}
	}
	if i.IsHardSprite {
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
