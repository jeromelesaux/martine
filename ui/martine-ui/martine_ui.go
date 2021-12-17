package ui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx"
)

var (
	refreshUI        *widget.Button
	modeSelection    *widget.Select
	paletteSelection *widget.Select
	dialogSize       = fyne.NewSize(800, 800)
)

type MartineUI struct {
	window fyne.Window
	main   *ImageMenu

	exportDsk              bool
	exportText             bool
	exportWithAmsdosHeader bool
	exportZigzag           bool
	exportJson             bool
	exportCompression      int
	exportFolderPath       string
}

func NewMartineUI() *MartineUI {
	return &MartineUI{main: &ImageMenu{}}
}

func (m *MartineUI) SetPalette(p color.Palette) {

	m.main.palette = p
	m.main.paletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
	refreshUI.OnTapped()
}

func (m *MartineUI) Load(app fyne.App) {
	m.window = app.NewWindow("Martine @IMPact v" + common.AppVersion)
	m.window.SetContent(m.NewTabs())
	m.window.Resize(fyne.NewSize(1400, 1000))
	m.window.SetTitle("Martine @IMPact v" + common.AppVersion)
	m.window.Show()
}

func (m *MartineUI) NewTabs() *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Image", m.newImageTransfertTab(m.main)),
		//container.NewTabItem("Animation", widget.NewLabel("Animation")),
	)
}

func (m *MartineUI) NewContext(me *ImageMenu) *export.MartineContext {
	if m.main.originalImagePath == nil {
		return nil
	}
	context := export.NewMartineContext(m.main.originalImagePath.Path(), "")
	context.CpcPlus = m.main.isCpcPlus
	context.Overscan = m.main.isFullScreen
	context.DitheringMultiplier = m.main.ditheringMultiplier
	context.Brightness = m.main.brightness
	context.Saturation = m.main.saturation
	if m.main.brightness > 0 && m.main.saturation == 0 {
		context.Saturation = me.brightness
	}
	if me.brightness == 0 && me.saturation > 0 {
		context.Brightness = me.saturation
	}
	context.Reducer = me.reducer
	var size constants.Size
	switch me.mode {
	case 0:
		size = constants.Mode0
		if me.isFullScreen {
			size = constants.OverscanMode0
		}
	case 1:
		size = constants.Mode1
		if me.isFullScreen {
			size = constants.OverscanMode1
		}
	case 2:
		size = constants.Mode2
		if me.isFullScreen {
			size = constants.OverscanMode2
		}
	}
	context.Size = size
	if me.isSprite {
		width, err := strconv.Atoi(me.width.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return nil
		}
		height, err := strconv.Atoi(me.height.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return nil
		}
		context.Size.Height = height
		context.Size.Width = width
	}
	if me.isHardSprite {
		context.Size.Height = 16
		context.Size.Width = 16
	}

	if me.applyDithering {
		context.DitheringAlgo = 0
		context.DitheringMatrix = me.ditheringMatrix
		context.DitheringType = me.ditheringType
	} else {
		context.DitheringAlgo = -1
	}
	context.DitheringWithQuantification = me.withQuantification
	context.OutputPath = m.exportFolderPath
	context.InputPath = me.originalImagePath.Path()
	context.Json = m.exportJson
	context.Ascii = m.exportText
	context.NoAmsdosHeader = !m.exportWithAmsdosHeader
	context.ZigZag = m.exportZigzag
	context.Compression = m.exportCompression
	context.Dsk = m.exportDsk
	return context
}

func (m *MartineUI) ExportOneImage(me *ImageMenu) {
	pi := dialog.NewProgressInfinite("Saving....", "Please wait.", m.window)
	pi.Show()
	context := m.NewContext(me)
	// palette export
	defer func() {
		os.Remove("temporary_palette.kit")
	}()
	if err := file.SaveKit("temporary_palette.kit", me.palette, false); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
	}
	context.KitPath = "temporary_palette.kit"
	if err := gfx.ApplyOneImageAndExport(me.originalImage.Image, context, filepath.Base(m.exportFolderPath), m.exportFolderPath, me.mode, uint8(me.mode)); err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	if context.Dsk {
		if err := file.ImportInDsk(me.originalImagePath.Path(), context); err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
	}
	if context.Sna {
		if context.Overscan {
			var gfxFile string
			for _, v := range context.DskFiles {
				if filepath.Ext(v) == ".SCR" {
					gfxFile = v
					break
				}
			}
			context.SnaPath = filepath.Join(m.exportFolderPath, "test.sna")
			if err := file.ImportInSna(gfxFile, context.SnaPath, uint8(me.mode)); err != nil {
				dialog.NewError(err, m.window).Show()
				return
			}
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in foler \n"+m.exportFolderPath, m.window)

}

func (m *MartineUI) ApplyOneImage(me *ImageMenu) {
	me.cpcImage = canvas.Image{}
	context := m.NewContext(me)
	if context == nil {
		return
	}

	var inPalette color.Palette
	if me.usePalette {
		inPalette = me.palette
		maxPalette := len(inPalette)
		switch me.mode {
		case 1:
			if maxPalette > 4 {
				maxPalette = 4
			}
			inPalette = inPalette[0:maxPalette]
		case 2:
			if maxPalette > 2 {
				maxPalette = 2
			}
			inPalette = inPalette[0:maxPalette]
		}

	}
	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	out, downgraded, palette, _, err := gfx.ApplyOneImage(me.originalImage.Image, context, me.mode, inPalette, uint8(me.mode))
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	me.data = out
	me.downgraded = downgraded
	if !me.usePalette {
		me.palette = palette
	}
	if me.isSprite || me.isHardSprite {
		newSize := constants.Size{Width: context.Size.Width * 50, Height: context.Size.Height * 50}
		me.downgraded = convert.Resize(me.downgraded, newSize, me.resizeAlgo)
	}
	me.cpcImage = *canvas.NewImageFromImage(me.downgraded)
	me.cpcImage.FillMode = canvas.ImageFillStretch
	me.paletteImage = *canvas.NewImageFromImage(file.PalToImage(me.palette))
	refreshUI.OnTapped()
}

func (m *MartineUI) newImageTransfertTab(me *ImageMenu) fyne.CanvasObject {
	importOpen := NewImportButton(m, me)

	paletteOpen := NewOpenPaletteButton(me, m.window)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		me.usePalette = b
	})

	forceUIRefresh := widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
		s := m.window.Content().Size()
		s.Height += 10.
		s.Width += 10.
		m.window.Resize(s)
		m.window.Canvas().Refresh(&me.originalImage)
		m.window.Canvas().Refresh(&me.paletteImage)
		m.window.Canvas().Refresh(&me.cpcImage)
		m.window.Resize(m.window.Content().Size())
		m.window.Content().Refresh()
	})
	refreshUI = forceUIRefresh
	openFileWidget := widget.NewButton("Image", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}

			me.originalImagePath = reader.URI()
			img, err := openImage(me.originalImagePath.Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			canvasImg := canvas.NewImageFromImage(img)
			me.originalImage = *canvas.NewImageFromImage(canvasImg.Image)
			me.originalImage.FillMode = canvas.ImageFillContain
			m.window.Canvas().Refresh(&me.originalImage)
			m.window.Resize(m.window.Content().Size())
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".gif", ".png", ".jpeg"}))
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportDialog(m.window)
	})

	applyButton := widget.NewButtonWithIcon("Apply", theme.VisibilityIcon(), func() {
		fmt.Println("apply.")
		m.ApplyOneImage(me)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	me.cpcImage = canvas.Image{}
	me.originalImage = canvas.Image{}
	me.paletteImage = canvas.Image{}

	winFormat := NewWinFormatRadio(me)

	colorReducerLabel := widget.NewLabel("Color reducer")
	colorReducer := widget.NewSelect([]string{"none", "Lower", "Medium", "Strong"}, func(s string) {
		switch s {
		case "none":
			me.reducer = 0
		case "Lower":
			me.reducer = 1
		case "Medium":
			me.reducer = 2
		case "Strong":
			me.reducer = 3
		}
	})
	colorReducer.SetSelected("none")

	resize := NewResizeAlgorithmSelect(me)
	resizeLabel := widget.NewLabel("Resize algorithm")

	ditheringMultiplier := widget.NewSlider(0., 2.5)
	ditheringMultiplier.Step = 0.1
	ditheringMultiplier.SetValue(1.18)
	ditheringMultiplier.OnChanged = func(f float64) {
		me.ditheringMultiplier = f
	}
	dithering := NewDitheringSelect(me)

	ditheringWithQuantification := widget.NewCheck("With quantification", func(b bool) {
		me.withQuantification = b
	})

	enableDithering := widget.NewCheck("Enable dithering", func(b bool) {
		me.applyDithering = b
	})
	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		me.isCpcPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		me.mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	me.width = widget.NewEntry()
	me.width.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	me.height = widget.NewEntry()
	me.height.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	brightness := widget.NewSlider(0.0, 1.0)
	brightness.SetValue(1.)
	brightness.Step = .01
	brightness.OnChanged = func(f float64) {
		me.brightness = f
	}
	saturationLabel := widget.NewLabel("Saturation")
	saturation := widget.NewSlider(0.0, 1.0)
	saturation.SetValue(1.)
	saturation.Step = .01
	saturation.OnChanged = func(f float64) {
		me.saturation = f
	}
	brightnessLabel := widget.NewLabel("Brightness")
	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.NewScroll(
				&me.originalImage),
			container.NewScroll(
				&me.cpcImage),
		),
		container.New(
			layout.NewVBoxLayout(),
			container.New(
				layout.NewHBoxLayout(),
				openFileWidget,
				paletteOpen,
				applyButton,
				exportButton,
				importOpen,
				forceUIRefresh,
			),
			container.New(
				layout.NewHBoxLayout(),
				isPlus,
				winFormat,

				container.New(
					layout.NewVBoxLayout(),
					container.New(
						layout.NewVBoxLayout(),
						modeLabel,
						modes,
					),
					container.New(
						layout.NewHBoxLayout(),
						widthLabel,
						me.width,
					),
					container.New(
						layout.NewHBoxLayout(),
						heightLabel,
						me.height,
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(5),
				container.New(
					layout.NewGridLayoutWithRows(2),
					container.New(
						layout.NewGridLayoutWithColumns(2),
						resizeLabel,
						resize,
					),
					container.New(
						layout.NewGridLayoutWithColumns(4),
						enableDithering,
						dithering,
						ditheringMultiplier,
						ditheringWithQuantification,
					),
				),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					&me.paletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							swapColor(m.SetPalette, me.palette, m.window)
						}),
						widget.NewButtonWithIcon("export", theme.DocumentSaveIcon(), func() {
							d := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
								if err != nil {
									dialog.ShowError(err, m.window)
									return
								}
								if uc == nil {
									return
								}

								paletteExportPath := uc.URI().Path()
								uc.Close()
								os.Remove(uc.URI().Path())
								context := export.NewMartineContext(filepath.Base(paletteExportPath), paletteExportPath)
								context.NoAmsdosHeader = false
								if err := file.SaveKit(paletteExportPath+".kit", me.palette, false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := file.SavePal(paletteExportPath+".pal", me.palette, uint8(me.mode), false); err != nil {
									dialog.ShowError(err, m.window)
								}
							}, m.window)
							d.Show()
						}),
					),
				),
				container.New(
					layout.NewVBoxLayout(),
					colorReducerLabel,
					colorReducer,
				),
				container.New(
					layout.NewVBoxLayout(),
					brightnessLabel,
					brightness,
				),
				container.New(
					layout.NewVBoxLayout(),
					saturationLabel,
					saturation,
					widget.NewButton("show cmd", func() {
						e := widget.NewMultiLineEntry()
						e.SetText(me.CmdLine())

						d := dialog.NewCustom("Command line generated",
							"Ok",
							e,
							m.window)
						fmt.Printf("%s\n", me.CmdLine())
						size := m.window.Content().Size()
						size = fyne.Size{Width: size.Width / 2, Height: size.Height / 2}
						d.Resize(size)
						d.Show()
					}),
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
