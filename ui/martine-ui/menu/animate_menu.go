package menu

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	w "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
)

const (
	AnimateWidthSize  float32 = 350.
	AnimateHeightSize float32 = 250.
)

type AnimateMenu struct {
	*ImageMenu
	Originalmages      []*canvas.Image
	AnimateImages      *w.ImageSelectionTable
	DeltaCollection    []*transformation.DeltaCollection
	InitialAddress     *widget.Entry
	RawImages          [][]byte
	OneLine            bool
	OneRow             bool
	ImageToRemoveIndex int
	ExportVersion      animate.DeltaExportFormat
}

func NewAnimateMenu() *AnimateMenu {
	return &AnimateMenu{
		ImageMenu:       NewImageMenu(),
		Originalmages:   make([]*canvas.Image, 0),
		AnimateImages:   w.NewImageSelectionTable(fyne.NewSize(AnimateWidthSize, AnimateHeightSize)),
		DeltaCollection: make([]*transformation.DeltaCollection, 1),
	}
}

// nolint: funlen
func (i *AnimateMenu) CmdLine() string {
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

	if i.Cfg.ScrCfg.Type.IsSprite() {
		exec += " -address 0x" + i.InitialAddress.Text
	}

	if i.OneLine {
		exec += " -oneline"
	}
	if i.OneRow {
		exec += " -onerow"
	}

	exec += " -animate"

	i.CmdLineGenerate = exec
	return exec
}
