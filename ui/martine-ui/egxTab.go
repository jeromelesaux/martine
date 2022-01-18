package ui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/gfx/effect"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

func (m *MartineUI) newEgxTab(di *menu.DoubleImageMenu) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Image 1", m.newEgxImageTransfertTab(&di.LeftImage)),
		container.NewTabItem("Image 2", m.newEgxImageTransfertTab(&di.RightImage)),
		container.NewTabItem("Egx", m.newEgxTabItem(di)),
	)
}

func (m *MartineUI) MergeImages(di *menu.DoubleImageMenu) {
	if di.RightImage.Mode == di.LeftImage.Mode {
		dialog.ShowError(fmt.Errorf("mode between the images must differ"), m.window)
		return
	}
	if di.RightImage.IsCpcPlus != di.LeftImage.IsCpcPlus {
		dialog.ShowError(fmt.Errorf("plus mode between the images must differ"), m.window)
		return
	}
	if di.RightImage.IsHardSprite != di.LeftImage.IsHardSprite {
		dialog.ShowError(fmt.Errorf("sprite hard mode between the images must differ"), m.window)
		return
	}
	if di.RightImage.IsFullScreen != di.LeftImage.IsFullScreen {
		dialog.ShowError(fmt.Errorf("fullscreen mode between the images must  differ"), m.window)
		return
	}

	var im *menu.ImageMenu
	var palette color.Palette
	if di.LeftImage.Mode == 0 {
		im = &di.LeftImage
		palette = di.LeftImage.Palette
	} else {
		im = &di.RightImage
		palette = di.RightImage.Palette
	}
	context := m.NewContext(im, false)
	if context == nil {
		return
	}
	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	res, palette, err := effect.EgxRaw(di.LeftImage.Data, di.RightImage.Data, palette, di.LeftImage.Mode, di.RightImage.Mode, context)
	pi.Hide()
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	di.ResultImage.Data = res
	di.ResultImage.Palette = palette
	var img image.Image
	if context.Overscan {
		img, err = common.OverscanRawToImg(di.ResultImage.Data, 0, di.ResultImage.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	} else {
		img, err = common.ScrRawToImg(di.ResultImage.Data, 0, di.ResultImage.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	}
	di.ResultImage.CpcResultImage = *canvas.NewImageFromImage(img)
	di.ResultImage.CpcResultImage.Refresh()
	m.window.Canvas().Refresh(&di.ResultImage.CpcResultImage)
	m.window.Resize(m.window.Content().Size())
	m.window.Content().Refresh()
}

func (m *MartineUI) newEgxTabItem(di *menu.DoubleImageMenu) fyne.CanvasObject {
	di.ResultImage.CpcLeftImage = di.LeftImage.CpcImage
	di.ResultImage.CpcRightImage = di.RightImage.CpcImage
	di.ResultImage.LeftPalette = di.LeftImage.Palette
	di.ResultImage.RightPalette = di.RightImage.Palette
	di.ResultImage.LeftPaletteImage = di.LeftImage.PaletteImage
	di.ResultImage.RightPaletteImage = di.RightImage.PaletteImage
	return container.New(
		layout.NewGridLayoutWithRows(3),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			&di.LeftImage.CpcImage,
			&di.RightImage.CpcImage,
		),

		container.New(
			layout.NewVBoxLayout(),

			container.New(
				layout.NewGridLayoutWithColumns(2),
				container.New(
					layout.NewGridLayoutWithColumns(1),
					&di.LeftImage.PaletteImage,
				),
				container.New(
					layout.NewGridLayoutWithColumns(1),
					&di.RightImage.PaletteImage,
				),

				widget.NewButtonWithIcon("Merge image", theme.MediaPlayIcon(), func() {
					m.MergeImages(m.egx)
				}),
				widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
					di.ResultImage.CpcLeftImage = di.LeftImage.CpcImage
					di.ResultImage.CpcRightImage = di.RightImage.CpcImage
					di.ResultImage.LeftPalette = di.LeftImage.Palette
					di.ResultImage.RightPalette = di.RightImage.Palette
					di.ResultImage.LeftPaletteImage = di.LeftImage.PaletteImage
					di.ResultImage.RightPaletteImage = di.RightImage.PaletteImage
					s := m.window.Content().Size()
					s.Height += 10.
					s.Width += 10.
					m.window.Resize(s)
					m.window.Canvas().Refresh(&di.ResultImage.CpcLeftImage)
					m.window.Canvas().Refresh(&di.ResultImage.CpcRightImage)
					m.window.Canvas().Refresh(&di.ResultImage.CpcResultImage)
					m.window.Canvas().Refresh(&di.ResultImage.LeftPaletteImage)
					m.window.Canvas().Refresh(&di.ResultImage.RightPaletteImage)
					m.window.Resize(m.window.Content().Size())
					m.window.Content().Refresh()
				}),
				widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {

				}),
			),
		),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			&di.ResultImage.CpcResultImage,
		),
	)
}

func (m *MartineUI) newEgxImageTransfertTab(me *menu.ImageMenu) fyne.CanvasObject {
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

	winFormat := w2.NewScreenFormatRadio(me)

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

	modes := widget.NewSelect([]string{"0", "1"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		me.Mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

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
