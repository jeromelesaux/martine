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
	"github.com/jeromelesaux/martine/export/net"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

var (
	refreshUI        *widget.Button
	modeSelection    *widget.Select
	paletteSelection *widget.Select
	dialogSize       = fyne.NewSize(800, 800)
)

type MartineUI struct {
	window  fyne.Window
	main    *menu.ImageMenu
	tilemap *menu.TilemapMenu

	exportDsk              bool
	exportText             bool
	exportWithAmsdosHeader bool
	exportZigzag           bool
	exportJson             bool
	exportCompression      int
	exportFolderPath       string
	m2IP                   string
	exportToM2             bool
}

func NewMartineUI() *MartineUI {
	return &MartineUI{
		main:    &menu.ImageMenu{},
		tilemap: &menu.TilemapMenu{},
	}
}

func (m *MartineUI) SetPalette(p color.Palette) {

	m.main.Palette = p
	m.main.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
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
		container.NewTabItem("Tile", m.newTilemapTab(m.tilemap)),
	)
}

func (m *MartineUI) ComputeTilemap(tm *menu.TilemapMenu) {

}

func (m *MartineUI) newTilemapTab(tm *menu.TilemapMenu) fyne.CanvasObject {
	importOpen := NewImportButton(m, &tm.ImageMenu)

	paletteOpen := NewOpenPaletteButton(&tm.ImageMenu, m.window)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		tm.UsePalette = b
	})

	forceUIRefresh := widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
		s := m.window.Content().Size()
		s.Height += 10.
		s.Width += 10.
		m.window.Resize(s)
		m.window.Canvas().Refresh(&tm.OriginalImage)
		m.window.Canvas().Refresh(&tm.PaletteImage)
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

			tm.OriginalImagePath = reader.URI()
			img, err := openImage(tm.OriginalImagePath.Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			canvasImg := canvas.NewImageFromImage(img)
			tm.OriginalImage = *canvas.NewImageFromImage(canvasImg.Image)
			tm.OriginalImage.FillMode = canvas.ImageFillContain
			m.window.Canvas().Refresh(&tm.OriginalImage)
			m.window.Resize(m.window.Content().Size())
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".gif", ".png", ".jpeg"}))
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportDialog(m.window)
	})

	applyButton := widget.NewButtonWithIcon("Compute", theme.VisibilityIcon(), func() {
		fmt.Println("compute.")
		m.ComputeTilemap(tm)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	tm.OriginalImage = canvas.Image{}
	tm.PaletteImage = canvas.Image{}

	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		tm.IsCpcPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		tm.Mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	tm.Width = widget.NewEntry()
	tm.Width.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	tm.Height = widget.NewEntry()
	tm.Height.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.NewScroll(
				&tm.OriginalImage),
			container.NewScroll(
				&tm.TileImages),
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
						tm.Width,
					),
					container.New(
						layout.NewHBoxLayout(),
						heightLabel,
						tm.Height,
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(6),

				container.New(
					layout.NewGridLayoutWithColumns(2),
					&tm.PaletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, tm.Palette, m.window)
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
								if err := file.SaveKit(paletteExportPath+".kit", tm.Palette, false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := file.SavePal(paletteExportPath+".pal", tm.Palette, uint8(tm.Mode), false); err != nil {
									dialog.ShowError(err, m.window)
								}
							}, m.window)
							d.Show()
						}),
					),
				),

				container.New(
					layout.NewVBoxLayout(),
					widget.NewButton("show cmd", func() {
						e := widget.NewMultiLineEntry()
						e.SetText(tm.CmdLine())

						d := dialog.NewCustom("Command line generated",
							"Ok",
							e,
							m.window)
						fmt.Printf("%s\n", tm.CmdLine())
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

func (m *MartineUI) NewContext(me *menu.ImageMenu) *export.MartineContext {
	if m.main.OriginalImagePath == nil {
		return nil
	}
	context := export.NewMartineContext(m.main.OriginalImagePath.Path(), "")
	context.CpcPlus = m.main.IsCpcPlus
	context.Overscan = m.main.IsFullScreen
	context.DitheringMultiplier = m.main.DitheringMultiplier
	context.Brightness = m.main.Brightness
	context.Saturation = m.main.Saturation
	if m.main.Brightness > 0 && m.main.Saturation == 0 {
		context.Saturation = me.Brightness
	}
	if me.Brightness == 0 && me.Saturation > 0 {
		context.Brightness = me.Saturation
	}
	context.Reducer = me.Reducer
	var size constants.Size
	switch me.Mode {
	case 0:
		size = constants.Mode0
		if me.IsFullScreen {
			size = constants.OverscanMode0
		}
	case 1:
		size = constants.Mode1
		if me.IsFullScreen {
			size = constants.OverscanMode1
		}
	case 2:
		size = constants.Mode2
		if me.IsFullScreen {
			size = constants.OverscanMode2
		}
	}
	context.Size = size
	if me.IsSprite {
		width, err := strconv.Atoi(me.Width.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return nil
		}
		height, err := strconv.Atoi(me.Height.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return nil
		}
		context.Size.Height = height
		context.Size.Width = width
	}
	if me.IsHardSprite {
		context.Size.Height = 16
		context.Size.Width = 16
	}

	if me.ApplyDithering {
		context.DitheringAlgo = 0
		context.DitheringMatrix = me.DitheringMatrix
		context.DitheringType = me.DitheringType
	} else {
		context.DitheringAlgo = -1
	}
	context.DitheringWithQuantification = me.WithQuantification
	context.OutputPath = m.exportFolderPath
	context.InputPath = me.OriginalImagePath.Path()
	context.Json = m.exportJson
	context.Ascii = m.exportText
	context.NoAmsdosHeader = !m.exportWithAmsdosHeader
	context.ZigZag = m.exportZigzag
	context.Compression = m.exportCompression
	context.Dsk = m.exportDsk
	return context
}

func (m *MartineUI) ExportOneImage(me *menu.ImageMenu) {
	pi := dialog.NewProgressInfinite("Saving....", "Please wait.", m.window)
	pi.Show()
	context := m.NewContext(me)
	// palette export
	defer func() {
		os.Remove("temporary_palette.kit")
	}()
	if err := file.SaveKit("temporary_palette.kit", me.Palette, false); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
	}
	context.KitPath = "temporary_palette.kit"
	if err := gfx.ApplyOneImageAndExport(me.OriginalImage.Image, context, filepath.Base(m.exportFolderPath), m.exportFolderPath, me.Mode, uint8(me.Mode)); err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	if context.Dsk {
		if err := file.ImportInDsk(me.OriginalImagePath.Path(), context); err != nil {
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
			if err := file.ImportInSna(gfxFile, context.SnaPath, uint8(me.Mode)); err != nil {
				dialog.NewError(err, m.window).Show()
				return
			}
		}
	}
	if m.exportToM2 {
		if err := net.ImportInM4(context); err != nil {
			dialog.NewError(err, m.window).Show()
			fmt.Fprintf(os.Stderr, "Cannot send to M4 error :%v\n", err)
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in foler \n"+m.exportFolderPath, m.window)

}

func (m *MartineUI) ApplyOneImage(me *menu.ImageMenu) {
	me.CpcImage = canvas.Image{}
	context := m.NewContext(me)
	if context == nil {
		return
	}

	var inPalette color.Palette
	if me.UsePalette {
		inPalette = me.Palette
		maxPalette := len(inPalette)
		switch me.Mode {
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
	out, downgraded, palette, _, err := gfx.ApplyOneImage(me.OriginalImage.Image, context, me.Mode, inPalette, uint8(me.Mode))
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	me.Data = out
	me.Downgraded = downgraded
	if !me.UsePalette {
		me.Palette = palette
	}
	if me.IsSprite || me.IsHardSprite {
		newSize := constants.Size{Width: context.Size.Width * 50, Height: context.Size.Height * 50}
		me.Downgraded = convert.Resize(me.Downgraded, newSize, me.ResizeAlgo)
	}
	me.CpcImage = *canvas.NewImageFromImage(me.Downgraded)
	me.CpcImage.FillMode = canvas.ImageFillStretch
	me.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(me.Palette))
	refreshUI.OnTapped()
}

func (m *MartineUI) newImageTransfertTab(me *menu.ImageMenu) fyne.CanvasObject {
	importOpen := NewImportButton(m, me)

	paletteOpen := NewOpenPaletteButton(me, m.window)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		me.UsePalette = b
	})

	forceUIRefresh := widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
		s := m.window.Content().Size()
		s.Height += 10.
		s.Width += 10.
		m.window.Resize(s)
		m.window.Canvas().Refresh(&me.OriginalImage)
		m.window.Canvas().Refresh(&me.PaletteImage)
		m.window.Canvas().Refresh(&me.CpcImage)
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

			me.OriginalImagePath = reader.URI()
			img, err := openImage(me.OriginalImagePath.Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			canvasImg := canvas.NewImageFromImage(img)
			me.OriginalImage = *canvas.NewImageFromImage(canvasImg.Image)
			me.OriginalImage.FillMode = canvas.ImageFillContain
			m.window.Canvas().Refresh(&me.OriginalImage)
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

	me.CpcImage = canvas.Image{}
	me.OriginalImage = canvas.Image{}
	me.PaletteImage = canvas.Image{}

	winFormat := w2.NewWinFormatRadio(me)

	colorReducerLabel := widget.NewLabel("Color reducer")
	colorReducer := widget.NewSelect([]string{"none", "Lower", "Medium", "Strong"}, func(s string) {
		switch s {
		case "none":
			me.Reducer = 0
		case "Lower":
			me.Reducer = 1
		case "Medium":
			me.Reducer = 2
		case "Strong":
			me.Reducer = 3
		}
	})
	colorReducer.SetSelected("none")

	resize := w2.NewResizeAlgorithmSelect(me)
	resizeLabel := widget.NewLabel("Resize algorithm")

	ditheringMultiplier := widget.NewSlider(0., 2.5)
	ditheringMultiplier.Step = 0.1
	ditheringMultiplier.SetValue(1.18)
	ditheringMultiplier.OnChanged = func(f float64) {
		me.DitheringMultiplier = f
	}
	dithering := w2.NewDitheringSelect(me)

	ditheringWithQuantification := widget.NewCheck("With quantification", func(b bool) {
		me.WithQuantification = b
	})

	enableDithering := widget.NewCheck("Enable dithering", func(b bool) {
		me.ApplyDithering = b
	})
	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		me.IsCpcPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		me.Mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	me.Width = widget.NewEntry()
	me.Width.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	me.Height = widget.NewEntry()
	me.Height.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	brightness := widget.NewSlider(0.0, 1.0)
	brightness.SetValue(1.)
	brightness.Step = .01
	brightness.OnChanged = func(f float64) {
		me.Brightness = f
	}
	saturationLabel := widget.NewLabel("Saturation")
	saturation := widget.NewSlider(0.0, 1.0)
	saturation.SetValue(1.)
	saturation.Step = .01
	saturation.OnChanged = func(f float64) {
		me.Saturation = f
	}
	brightnessLabel := widget.NewLabel("Brightness")

	warningLabel := widget.NewLabel("Setting thoses parameters will affect your palette, you can't force palette.")
	warningLabel.TextStyle = fyne.TextStyle{Bold: true}

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.NewScroll(
				&me.OriginalImage),
			container.NewScroll(
				&me.CpcImage),
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
						me.Width,
					),
					container.New(
						layout.NewHBoxLayout(),
						heightLabel,
						me.Height,
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(6),
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
					&me.PaletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, me.Palette, m.window)
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
								if err := file.SaveKit(paletteExportPath+".kit", me.Palette, false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := file.SavePal(paletteExportPath+".pal", me.Palette, uint8(me.Mode), false); err != nil {
									dialog.ShowError(err, m.window)
								}
							}, m.window)
							d.Show()
						}),
					),
				),
				container.New(
					layout.NewVBoxLayout(),
					warningLabel,
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
