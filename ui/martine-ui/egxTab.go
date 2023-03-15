package ui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"

	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/screen"
	"github.com/jeromelesaux/martine/convert/screen/overscan"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/effect"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

func (m *MartineUI) newEgxTab(di *menu.DoubleImageMenu) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Image 1", m.newEgxImageTransfertTab(di.LeftImage)),
		container.NewTabItem("Image 2", m.newEgxImageTransfertTab(di.RightImage)),
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
		im = di.LeftImage
		palette = di.LeftImage.Palette()
		secondIm = di.RightImage
	} else {
		im = di.RightImage
		palette = di.RightImage.Palette()
		secondIm = di.LeftImage
	}
	cfg := m.NewConfig(im, false)
	if cfg == nil {
		return
	}
	out, downgraded, _, _, err := gfx.ApplyOneImage(secondIm.CpcImage().Image, cfg, secondIm.Mode, palette, uint8(secondIm.Mode))
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	secondIm.Data = out
	secondIm.SetCpcImage(downgraded)
	secondIm.SetPalette(palette)

	im.SetPaletteImage(png.PalToImage(secondIm.Palette()))
	out, downgraded, _, _, err = gfx.ApplyOneImage(im.CpcImage().Image, cfg, im.Mode, palette, uint8(im.Mode))
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	im.Data = out
	im.SetCpcImage(downgraded)
	im.SetPalette(palette)
	im.SetPaletteImage(png.PalToImage(im.Palette()))

	pi := custom_widget.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	res, _, egxType, err := effect.EgxRaw(di.LeftImage.Data, di.RightImage.Data, palette, di.LeftImage.Mode, di.RightImage.Mode, cfg)
	pi.Hide()
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	di.ResultImage.Data = res
	di.ResultImage.Palette = palette
	di.ResultImage.EgxType = egxType
	var img image.Image
	if cfg.Overscan {
		img, err = overscan.OverscanRawToImg(di.ResultImage.Data, 0, di.ResultImage.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	} else {
		img, err = screen.ScrRawToImg(di.ResultImage.Data, 0, di.ResultImage.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	}
	di.ResultImage.CpcResultImage.Image = img
	di.ResultImage.CpcResultImage.Refresh()
	di.ResultImage.PaletteImage.Image = png.PalToImage(di.ResultImage.Palette)
	di.ResultImage.PaletteImage.Refresh()
}

func (m *MartineUI) newEgxTabItem(di *menu.DoubleImageMenu) fyne.CanvasObject {
	di.ResultImage = menu.NewMergedImageMenu()
	di.ResultImage.CpcLeftImage = di.LeftImage.CpcImage()
	di.ResultImage.CpcRightImage = di.RightImage.CpcImage()
	di.ResultImage.LeftPalette = di.LeftImage.Palette()
	di.ResultImage.RightPalette = di.RightImage.Palette()
	di.ResultImage.LeftPaletteImage = di.LeftImage.PaletteImage()
	di.ResultImage.RightPaletteImage = di.RightImage.PaletteImage()

	return container.New(
		layout.NewGridLayoutWithRows(3),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			di.LeftImage.CpcImage(),
			di.RightImage.CpcImage(),
		),

		container.New(
			layout.NewVBoxLayout(),

			container.New(
				layout.NewGridLayoutWithColumns(2),
				container.New(
					layout.NewGridLayoutWithColumns(1),
					di.LeftImage.PaletteImage(),
				),
				container.New(
					layout.NewGridLayoutWithColumns(1),
					di.RightImage.PaletteImage(),
				),

				widget.NewButtonWithIcon("Merge image", theme.MediaPlayIcon(), func() {
					m.MergeImages(m.egx)
				}),
				widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
					di.ResultImage.CpcLeftImage.Image = di.LeftImage.CpcImage().Image
					di.ResultImage.CpcRightImage.Image = di.RightImage.CpcImage().Image
					di.ResultImage.CpcLeftImage.Refresh()
					di.ResultImage.CpcResultImage.Refresh()
					di.ResultImage.LeftPalette = di.LeftImage.Palette()
					di.ResultImage.RightPalette = di.RightImage.Palette()
					di.ResultImage.LeftPaletteImage.Image = di.LeftImage.PaletteImage().Image
					di.ResultImage.LeftPaletteImage.Refresh()
					di.ResultImage.RightPaletteImage.Image = di.RightImage.PaletteImage().Image
					di.ResultImage.LeftPaletteImage.Refresh()
					di.ResultImage.PaletteImage.Image = png.PalToImage(di.ResultImage.Palette)
					di.ResultImage.PaletteImage.Refresh()
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
					log.GetLogger().Info("%s\n", di.CmdLine())
					size := m.window.Content().Size()
					size = fyne.Size{Width: size.Width / 2, Height: size.Height / 2}
					d.Resize(size)
					d.Show()
				}),
			),
		),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			di.ResultImage.CpcResultImage,
			container.New(
				layout.NewGridLayoutWithRows(3),
				di.ResultImage.PaletteImage,
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

	openFileWidget := widget.NewButton("Image", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			directory.SetDefaultDirectoryURI(reader.URI())
			me.SetOriginalImagePath(reader.URI())
			img, err := openImage(me.OriginalImagePath())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			me.SetOriginalImage(img)
		}, m.window)
		dir, err := directory.DefaultDirectoryURI()
		if err != nil {
			d.SetLocation(dir)
		}
		d.SetFilter(imagesFilesFilter)
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportDialog(m.imageExport, m.window)
	})

	applyButton := widget.NewButtonWithIcon("Apply", theme.VisibilityIcon(), func() {
		log.GetLogger().Infoln("apply.")
		m.ApplyOneImage(me)
	})

	openFileWidget.Icon = theme.FileImageIcon()

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
			log.GetLogger().Error("Error %s cannot be cast in int\n", s)
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
				me.OriginalImage()),
			container.NewScroll(
				me.CpcImage()),
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
					me.PaletteImage(),
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, me.Palette(), m.window, func() {
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
								cfg := config.NewMartineConfig(filepath.Base(paletteExportPath), paletteExportPath)
								cfg.NoAmsdosHeader = false
								if err := impPalette.SaveKit(paletteExportPath+".kit", me.Palette(), false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := ocpartstudio.SavePal(paletteExportPath+".pal", me.Palette(), uint8(me.Mode), false); err != nil {
									dialog.ShowError(err, m.window)
								}
							}, m.window)
							dir, err := directory.DefaultDirectoryURI()
							if err != nil {
								d.SetLocation(dir)
							}
							d.Show()
						}),
						widget.NewButton("Gray", func() {
							if me.IsCpcPlus {
								me.SetPalette(ci.MonochromePalette(me.Palette()))
								me.SetPaletteImage(png.PalToImage(me.Palette()))
								forcePalette.SetChecked(true)
								forcePalette.Refresh()
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
						log.GetLogger().Info("%s\n", me.CmdLine())
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
