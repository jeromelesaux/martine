package ui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/filter"
)

type MartineUI struct {
	window            fyne.Window
	originalImage     canvas.Image
	cpcImage          canvas.Image
	originalImagePath fyne.URI
	isCpcPlus         bool
	isFullScreen      bool
	isSprite          bool
	isHardSprite      bool
	mode              int
	width             widget.Entry
	height            widget.Entry
	palette           color.Palette
	data              []byte
	downgraded        *image.NRGBA
	ditheringMatrix   [][]float32
	ditheringType     constants.DitheringType
	applyDithering    bool
	resizeAlgo        imaging.ResampleFilter
}

func NewMartineUI() *MartineUI {
	return &MartineUI{}
}

func (m *MartineUI) Load(app fyne.App) {
	m.window = app.NewWindow("Martine @IMPact v" + common.AppVersion)
	m.window.SetContent(m.NewTabs())
	m.window.Resize(fyne.NewSize(1200, 800))
	m.window.SetTitle("Martine @IMPact v" + common.AppVersion)
	m.window.Show()
}

func (m *MartineUI) NewTabs() *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Image treatment", m.newImageTransfertTab()),
		container.NewTabItem("Animation treatment", widget.NewLabel("Animation")),
	)
}

func (m *MartineUI) ApplyOneImage() {
	m.cpcImage = canvas.Image{}
	m.window.Canvas().Refresh(&m.cpcImage)
	m.window.Resize(m.window.Content().Size())

	context := export.NewMartineContext(m.originalImagePath.Path(), "")
	context.CpcPlus = m.isCpcPlus
	context.Overscan = m.isFullScreen
	var size constants.Size
	switch m.mode {
	case 0:
		size = constants.Mode0
		if m.isFullScreen {
			size = constants.OverscanMode0
		}
	case 1:
		size = constants.Mode1
		if m.isFullScreen {
			size = constants.OverscanMode1
		}
	case 2:
		size = constants.Mode2
		if m.isFullScreen {
			size = constants.OverscanMode2
		}
	}
	context.Size = size
	if m.isSprite {
		width, err := strconv.Atoi(m.width.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
		height, err := strconv.Atoi(m.height.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
		context.Size.Height = height
		context.Size.Width = width
	}
	if m.isHardSprite {
		context.Size.Height = 16
		context.Size.Width = 16
	}

	if m.applyDithering {
		context.DitheringMatrix = m.ditheringMatrix
		context.DitheringType = m.ditheringType
	}
	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	out, downgraded, palette, _, err := gfx.ApplyOneImage(m.originalImage.Image, context, m.mode, color.Palette{}, uint8(m.mode))
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	m.data = out
	m.downgraded = downgraded
	m.palette = palette

	m.cpcImage = *canvas.NewImageFromImage(m.downgraded)
	m.cpcImage.FillMode = canvas.ImageFillContain
	m.window.Canvas().Refresh(&m.cpcImage)
	m.window.Resize(m.window.Content().Size())
}

func (m *MartineUI) newImageTransfertTab() fyne.CanvasObject {
	openFileWidget := widget.NewButton("Open image", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}

			m.originalImagePath = reader.URI()
			img, err := openImage(m.originalImagePath.Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			canvasImg := canvas.NewImageFromImage(img)
			m.originalImage = *canvas.NewImageFromImage(canvasImg.Image)
			m.originalImage.FillMode = canvas.ImageFillContain
			m.window.Canvas().Refresh(&m.originalImage)
			m.window.Resize(m.window.Content().Size())
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".gif", ".png", ".jpeg"}))
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {

	})

	applyButton := widget.NewButtonWithIcon("Apply", theme.VisibilityIcon(), func() {
		fmt.Println("apply.")
		m.ApplyOneImage()
	})

	openFileWidget.Icon = theme.FileImageIcon()

	m.cpcImage = canvas.Image{}
	m.originalImage = canvas.Image{}

	winFormat := widget.NewRadioGroup([]string{"Normal", "Fullscreen", "Sprite", "Sprite Hard"}, func(s string) {
		switch s {
		case "Normal":
			m.isFullScreen = false
			m.isSprite = false
			m.isHardSprite = false
		case "Fullscreen":
			m.isFullScreen = true
			m.isSprite = false
			m.isHardSprite = false
		case "Sprite":
			m.isFullScreen = false
			m.isSprite = true
			m.isHardSprite = false
		case "Sprite Hard":
			m.isFullScreen = false
			m.isSprite = false
			m.isHardSprite = true
		}
	})
	winFormat.SetSelected("Normal")

	resize := widget.NewSelect([]string{"NearestNeighbor",
		"CatmullRom",
		"Lanczos",
		"Linear",
		"Box",
		"Hermite",
		"BSpline",
		"Hamming",
		"Hann",
		"Gaussian",
		"Blackman",
		"Bartlett",
		"Welch",
		"Cosine",
		"MitchellNetravali",
	}, func(s string) {
		switch s {
		case "NearestNeighbor":
			m.resizeAlgo = imaging.NearestNeighbor
		case "CatmullRom":
			m.resizeAlgo = imaging.CatmullRom
		case "Lanczos":
			m.resizeAlgo = imaging.Lanczos
		case "Linear":
			m.resizeAlgo = imaging.Linear
		case "Box":
			m.resizeAlgo = imaging.Box
		case "Hermite":
			m.resizeAlgo = imaging.Hermite
		case "BSpline":
			m.resizeAlgo = imaging.BSpline
		case "Hamming":
			m.resizeAlgo = imaging.Hamming
		case "Hann":
			m.resizeAlgo = imaging.Hann
		case "Gaussian":
			m.resizeAlgo = imaging.Gaussian
		case "Blackman":
			m.resizeAlgo = imaging.Blackman
		case "Bartlett":
			m.resizeAlgo = imaging.Bartlett
		case "Welch":
			m.resizeAlgo = imaging.Welch
		case "Cosine":
			m.resizeAlgo = imaging.Cosine
		case "MitchellNetravali":
			m.resizeAlgo = imaging.MitchellNetravali
		}
	})

	resize.SetSelected("NearestNeighbor")
	resizeLabel := widget.NewLabel("Resize algorithm")

	dithering := widget.NewSelect([]string{"FloydSteinberg",
		"JarvisJudiceNinke",
		"Stucki",
		"Atkinson",
		"Sierra",
		"SierraLite",
		"Sierra3",
		"Bayer2",
		"Bayer3",
		"Bayer4",
		"Bayer8",
	}, func(s string) {
		switch s {
		case "FloydSteinberg":
			m.ditheringMatrix = filter.FloydSteinberg
			m.ditheringType = constants.ErrorDiffusionDither
		case "JarvisJudiceNinke":
			m.ditheringMatrix = filter.JarvisJudiceNinke
			m.ditheringType = constants.ErrorDiffusionDither
		case "Stucki":
			m.ditheringMatrix = filter.Stucki
			m.ditheringType = constants.ErrorDiffusionDither
		case "Atkinson":
			m.ditheringMatrix = filter.Atkinson
			m.ditheringType = constants.ErrorDiffusionDither
		case "Sierra":
			m.ditheringMatrix = filter.Sierra
			m.ditheringType = constants.ErrorDiffusionDither
		case "SierraLite":
			m.ditheringMatrix = filter.SierraLite
			m.ditheringType = constants.ErrorDiffusionDither
		case "Sierra3":
			m.ditheringMatrix = filter.Sierra3
			m.ditheringType = constants.ErrorDiffusionDither
		case "Bayer2":
			m.ditheringMatrix = filter.Bayer2
			m.ditheringType = constants.OrderedDither
		case "Bayer3":
			m.ditheringMatrix = filter.Bayer3
			m.ditheringType = constants.OrderedDither
		case "Bayer4":
			m.ditheringMatrix = filter.Bayer4
			m.ditheringType = constants.OrderedDither
		case "Bayer8":
			m.ditheringMatrix = filter.Bayer8
			m.ditheringType = constants.OrderedDither
		}
	})
	dithering.SetSelected("FloydSteinberg")

	enableDithering := widget.NewCheck("Enable dithering", func(b bool) {
		m.applyDithering = b
	})
	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		m.isCpcPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		m.mode = mode
	})
	modes.SetSelected("0")
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	m.width = *widget.NewEntry()

	heightLabel := widget.NewLabel("Height")
	m.height = *widget.NewEntry()

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.NewScroll(
				&m.originalImage),
			container.NewScroll(
				&m.cpcImage),
		),
		container.New(
			layout.NewGridLayoutWithRows(3),
			container.New(
				layout.NewHBoxLayout(),
				openFileWidget,
				applyButton,
				exportButton,
			),
			container.New(
				layout.NewHBoxLayout(),
				isPlus,
				winFormat,

				container.New(
					layout.NewGridLayoutWithRows(3),
					container.New(
						layout.NewGridLayoutWithRows(2),
						modeLabel,
						modes,
					),
					container.New(
						layout.NewGridLayoutWithColumns(2),
						widthLabel,
						&m.width,
					),
					container.New(
						layout.NewGridLayoutWithColumns(2),
						heightLabel,
						&m.height,
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithColumns(2),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					resizeLabel,
					resize,
				),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					enableDithering,
					dithering,
				),
			),
		),
	)
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	return i, err
}