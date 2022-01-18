package ui

import (
	"fmt"
	"image/color"
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
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/export/net"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

func (m *MartineUI) ExportOneImage(me *menu.ImageMenu) {
	pi := dialog.NewProgressInfinite("Saving....", "Please wait.", m.window)
	pi.Show()
	context := m.NewContext(me, true)
	if context == nil {
		return
	}
	// palette export
	defer func() {
		os.Remove("temporary_palette.kit")
	}()
	if err := file.SaveKit("temporary_palette.kit", me.Palette, false); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
	}
	context.KitPath = "temporary_palette.kit"
	if err := gfx.ApplyOneImageAndExport(me.OriginalImage.Image, context, filepath.Base(m.imageExport.ExportFolderPath), m.imageExport.ExportFolderPath, me.Mode, uint8(me.Mode)); err != nil {
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
			context.SnaPath = filepath.Join(m.imageExport.ExportFolderPath, "test.sna")
			if err := file.ImportInSna(gfxFile, context.SnaPath, uint8(me.Mode)); err != nil {
				dialog.NewError(err, m.window).Show()
				return
			}
		}
	}
	if m.imageExport.ExportToM2 {
		if err := net.ImportInM4(context); err != nil {
			dialog.NewError(err, m.window).Show()
			fmt.Fprintf(os.Stderr, "Cannot send to M4 error :%v\n", err)
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in folder \n"+m.imageExport.ExportFolderPath, m.window)

}

func (m *MartineUI) ApplyOneImage(me *menu.ImageMenu) {
	me.CpcImage = canvas.Image{}
	context := m.NewContext(me, true)
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
		m.exportDialog(m.imageExport, m.window)
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
