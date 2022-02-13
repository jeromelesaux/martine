package menu

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

var AnimateSize float32 = 150.

type AnimateMenu struct {
	ImageMenu
	Originalmages      []canvas.Image
	AnimateImages      *custom_widget.ImageTable
	DeltaCollection    []*transformation.DeltaCollection
	InitialAddress     *widget.Entry
	RawImages          [][]byte
	IsEmpty            bool
	OneLine            bool
	OneRow             bool
	ImageToRemoveIndex int
}

func NewAnimateMenu() *AnimateMenu {
	return &AnimateMenu{
		Originalmages:   make([]canvas.Image, 0),
		AnimateImages:   custom_widget.NewEmptyImageTable(fyne.NewSize(AnimateSize, AnimateSize)),
		DeltaCollection: make([]*transformation.DeltaCollection, 1),
	}
}

func (i *AnimateMenu) CmdLine() string {
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

	if i.IsSprite {
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
