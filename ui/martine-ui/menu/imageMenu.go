package menu

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
)

type ImageMenu struct {
	OriginalImage       canvas.Image
	CpcImage            canvas.Image
	OriginalImagePath   fyne.URI
	IsCpcPlus           bool
	IsFullScreen        bool
	IsSprite            bool
	IsHardSprite        bool
	Mode                int
	Width               *widget.Entry
	Height              *widget.Entry
	Palette             color.Palette
	Data                []byte
	Downgraded          *image.NRGBA
	DitheringMatrix     [][]float32
	DitheringType       constants.DitheringType
	DitheringAlgoNumber int
	ApplyDithering      bool
	ResizeAlgo          imaging.ResampleFilter
	ResizeAlgoNumber    int
	PaletteImage        canvas.Image
	UsePalette          bool
	DitheringMultiplier float64
	WithQuantification  bool
	Brightness          float64
	Saturation          float64
	Reducer             int
	OneLine             bool
	OneRow              bool
	CmdLineGenerate     string
}

func (i *ImageMenu) SetPalette(p color.Palette) {
	i.Palette = p
}

func (i *ImageMenu) SetPaletteImage(c canvas.Image) {
	i.PaletteImage = c
}

func (i *ImageMenu) CmdLine() string {
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
	if i.OneLine {
		exec += " -oneline"
	}
	if i.OneRow {
		exec += " -onerow"
	}
	i.CmdLineGenerate = exec
	return exec
}
