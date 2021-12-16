package ui

import (
	"image"
	"image/color"

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
	applyDithering      bool
	resizeAlgo          imaging.ResampleFilter
	paletteImage        canvas.Image
	usePalette          bool
	ditheringMultiplier float64
	withQuantification  bool
	brightness          float64
	saturation          float64
	reducer             int
}
