package ui

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
	originalImage       canvas.Image
	cpcImage            canvas.Image
	originalImagePath   fyne.URI
	isCpcPlus           bool
	isFullScreen        bool
	isSprite            bool
	isHardSprite        bool
	mode                int
	width               *widget.Entry
	height              *widget.Entry
	palette             color.Palette
	data                []byte
	downgraded          *image.NRGBA
	ditheringMatrix     [][]float32
	ditheringType       constants.DitheringType
	ditheringAlgoNumber int
	applyDithering      bool
	resizeAlgo          imaging.ResampleFilter
	resizeAlgoNumber    int
	paletteImage        canvas.Image
	usePalette          bool
	ditheringMultiplier float64
	withQuantification  bool
	brightness          float64
	saturation          float64
	reducer             int
	cmdLineGenerate     string
}

func (i *ImageMenu) CmdLine() string {
	exec, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while getting executable path :%v\n", err)
		return exec
	}
	if i.originalImagePath != nil {
		exec += " -in " + i.originalImagePath.Path()
	}
	if i.isCpcPlus {
		exec += " -plus"
	}
	if i.isFullScreen {
		exec += " -fullscreen"
	}
	if i.isSprite {
		width, err := strconv.Atoi(i.width.Text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot convert width value :%s error :%v\n", i.width.Text, err)
		} else {
			exec += " -width " + strconv.Itoa(width)
		}
		height, err := strconv.Atoi(i.height.Text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot convert height value :%s error :%v\n", i.height.Text, err)
		} else {
			exec += " -height " + strconv.Itoa(height)
		}
	}
	if i.isHardSprite {
		exec += " -spritehard"
	}
	if i.applyDithering {
		if i.withQuantification {
			exec += " -quantization"
		} else {
			exec += " -multiplier " + fmt.Sprintf("%.2f", i.ditheringMultiplier)
		}
		exec += " -dithering " + strconv.Itoa(i.ditheringAlgoNumber)
		// stockage du numéro d'algo
	}
	exec += " -mode " + strconv.Itoa(i.mode)
	if i.reducer != 0 {
		exec += " -reducer " + strconv.Itoa(i.reducer)
	}
	// resize algo
	if i.resizeAlgoNumber != 0 {
		exec += " -algo " + strconv.Itoa(i.resizeAlgoNumber)
	}
	if i.brightness != 0 {
		exec += " -brightness " + fmt.Sprintf("%.2f", i.brightness)
	}
	if i.saturation != 0 {
		exec += " -saturation " + fmt.Sprintf("%.2f", i.saturation)
	}
	i.cmdLineGenerate = exec
	return exec
}
