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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
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
	var secondIm *menu.ImageMenu
	var palette color.Palette
	if di.LeftImage.Mode < di.RightImage.Mode {
		im = &di.LeftImage
		palette = di.LeftImage.Palette
		secondIm = &di.RightImage
	} else {
		im = &di.RightImage
		palette = di.RightImage.Palette
		secondIm = &di.LeftImage
	}
	context := m.NewContext(im, false)
	if context == nil {
		return
	}
	out, downgraded, _, _, err := gfx.ApplyOneImage(secondIm.CpcImage.Image, context, secondIm.Mode, palette, uint8(secondIm.Mode))
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	secondIm.Data = out
	secondIm.CpcImage = *canvas.NewImageFromImage(downgraded)
	secondIm.Palette = palette
	secondIm.PaletteImage = *canvas.NewImageFromImage(png.PalToImage(secondIm.Palette))
	out, downgraded, _, _, err = gfx.ApplyOneImage(im.CpcImage.Image, context, im.Mode, palette, uint8(im.Mode))
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	im.Data = out
	im.CpcImage = *canvas.NewImageFromImage(downgraded)
	im.Palette = palette
	im.PaletteImage = *canvas.NewImageFromImage(png.PalToImage(im.Palette))

	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	res, _, egxType, err := effect.EgxRaw(di.LeftImage.Data, di.RightImage.Data, palette, di.LeftImage.Mode, di.RightImage.Mode, context)
	pi.Hide()
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	di.ResultImage.Data = res
	di.ResultImage.Palette = palette
	di.ResultImage.EgxType = egxType
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
	di.ResultImage.PaletteImage = *canvas.NewImageFromImage(png.PalToImage(di.ResultImage.Palette))
	refreshUI.OnTapped()
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
					di.ResultImage.PaletteImage = *canvas.NewImageFromImage(png.PalToImage(di.ResultImage.Palette))
					s := m.window.Content().Size()
					s.Height += 10.
					s.Width += 10.
					m.window.Resize(s)
					m.window.Canvas().Refresh(&di.ResultImage.CpcLeftImage)
					m.window.Canvas().Refresh(&di.ResultImage.CpcRightImage)
					m.window.Canvas().Refresh(&di.ResultImage.CpcResultImage)
					m.window.Canvas().Refresh(&di.ResultImage.LeftPaletteImage)
					m.window.Canvas().Refresh(&di.ResultImage.RightPaletteImage)
					m.window.Canvas().Refresh(&di.ResultImage.PaletteImage)
					m.window.Resize(m.window.Content().Size())
					m.window.Content().Refresh()
				}),
				widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
					// export the egx image
					m.exportEgxDialog(m.egxExport, m.window)
				}),
				widget.NewButton("show cmd", func() {
					e := widget.NewMultiLineEntry()
					e.SetText(di.CmdLine())
					d := dialog.NewCustom("Command line generated",
						"Ok",
						e,
						m.window)
					fmt.Printf("%s\n", di.CmdLine())
					size := m.window.Content().Size()
					size = fyne.Size{Width: size.Width / 2, Height: size.Height / 2}
					d.Resize(size)
					d.Show()
				}),
			),
		),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			&di.ResultImage.CpcResultImage,
			container.New(
				layout.NewGridLayoutWithRows(3),
				&di.ResultImage.PaletteImage,
			),
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
		d.SetFilter(imagesFilesFilter)
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

	oneLine := widget.NewCheck("Every other line", func(b bool) {
		me.OneLine = b
	})
	oneRow := widget.NewCheck("Every other row", func(b bool) {
		me.OneRow = b
	})

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
				layout.NewGridLayoutWithRows(7),
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
					layout.NewGridLayoutWithRows(2),
					&me.PaletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, me.Palette, m.window, func() {
								forcePalette.SetChecked(true)
							})
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
								if err := ocpartstudio.SaveKit(paletteExportPath+".kit", me.Palette, false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := ocpartstudio.SavePal(paletteExportPath+".pal", me.Palette, uint8(me.Mode), false); err != nil {
									dialog.ShowError(err, m.window)
								}
							}, m.window)
							d.Show()
						}),
						widget.NewButton("Gray", func() {
							if me.IsCpcPlus {
								me.Palette = convert.MonochromePalette(me.Palette)
								me.PaletteImage = *canvas.NewImageFromImage(png.PalToImage(me.Palette))
								forcePalette.SetChecked(true)
								refreshUI.OnTapped()
							}
						}),
					),
				),
				container.New(
					layout.NewVBoxLayout(),
					oneLine,
					oneRow,
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
