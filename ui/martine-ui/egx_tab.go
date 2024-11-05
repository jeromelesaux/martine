package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/fyne-io/widget/editor"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/log"

	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/screen"
	"github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/effect"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	pal "github.com/jeromelesaux/martine/ui/martine-ui/palette"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

func (m *MartineUI) newEgxTab(di *menu.DoubleImageMenu) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Image 1", m.newEgxImageTransfertTab(di.LeftImage)),
		container.NewTabItem("Image 2", m.newEgxImageTransfertTab(di.RightImage)),
		container.NewTabItem("Egx", m.newEgxTabItem(di)),
	)
}

// nolint:funlen
func (m *MartineUI) MergeImages(di *menu.DoubleImageMenu) {
	if di.RightImage.Cfg.ScrCfg.Mode == di.LeftImage.Cfg.ScrCfg.Mode {
		dialog.ShowError(fmt.Errorf("mode between the images must differ"), m.window)
		return
	}
	if di.RightImage.Cfg.ScrCfg.IsPlus != di.LeftImage.Cfg.ScrCfg.IsPlus {
		dialog.ShowError(fmt.Errorf("plus mode between the images differs, set the same mode"), m.window)
		return
	}
	if di.RightImage.Cfg.ScrCfg.Type != di.LeftImage.Cfg.ScrCfg.Type {
		dialog.ShowError(fmt.Errorf("the size of both image differs, set the same size for both"), m.window)
		return
	}
	var im *menu.ImageMenu
	var im2 *menu.ImageMenu
	var palette2 color.Palette
	if di.LeftImage.Cfg.ScrCfg.Mode < di.RightImage.Cfg.ScrCfg.Mode {
		im = di.LeftImage
		palette2 = append(palette2, di.LeftImage.Palette()...)
		im2 = di.RightImage
	} else {
		im = di.RightImage
		palette2 = append(palette2, di.RightImage.Palette()...)
		im2 = di.LeftImage
	}

	di.ResultImage.Cfg.ScrCfg.Mode = im.Cfg.ScrCfg.Mode
	di.ResultImage.Cfg.ScrCfg.Size = im.Cfg.ScrCfg.Size
	di.ResultImage.Cfg.ScrCfg.Type = im.Cfg.ScrCfg.Type
	di.ResultImage.Cfg.ScrCfg.Export = append(di.ResultImage.Cfg.ScrCfg.Export, im.Cfg.ScrCfg.Export...)

	out, downgraded, pal, _, err := gfx.ApplyOneImage(im.CpcImage().Image, im.Cfg, im.Palette(), im.Cfg.ScrCfg.Mode)
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	im.Cfg.PalCfg.Palette = pal
	im.Data = out
	im.SetCpcImage(downgraded)
	im.SetPalette(im.Palette())
	im.SetPaletteImage(png.PalToImage(im.Palette()))

	out, downgraded, pal, _, err = gfx.ApplyOneImage(im2.CpcImage().Image, im2.Cfg, palette2, im2.Cfg.ScrCfg.Mode)
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	im2.Cfg.PalCfg.Palette = pal
	im2.Data = out
	im2.SetCpcImage(downgraded)
	im2.SetPalette(im2.Palette())
	im2.SetPaletteImage(png.PalToImage(im.Palette()))

	pi := wgt.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	di.ResultImage.Cfg.PalCfg.Palette = im.Palette()
	di.ResultImage.Cfg.ScrCfg.ResetExport()
	di.ResultImage.Cfg.ScrCfg.Type = im.Cfg.ScrCfg.Type
	di.ResultImage.Cfg.ScrCfg.Export = append(di.ResultImage.Cfg.ScrCfg.Export, im.Cfg.ScrCfg.Export...)
	di.ResultImage.Cfg.ScrCfg.Size = im.Cfg.ScrCfg.Size

	res, pal, egxType, err := effect.EgxRaw(im.CpcImage().Image, im2.CpcImage().Image, di.ResultImage.Cfg.PalCfg.Palette, int(im.Cfg.ScrCfg.Mode), int(im2.Cfg.ScrCfg.Mode), di.ResultImage.Cfg)
	pi.Hide()
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	di.ResultImage.EgxType = egxType
	di.ResultImage.Data = res
	di.ResultImage.Cfg.PalCfg.Palette = pal
	var img image.Image
	if di.ResultImage.Cfg.ScrCfg.Type == config.FullscreenFormat {
		img, err = overscan.OverscanRawToImg(di.ResultImage.Data, im.Cfg.ScrCfg.Mode, di.ResultImage.Cfg.PalCfg.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	} else {
		img, err = screen.ScrRawToImg(di.ResultImage.Data, 0, di.ResultImage.Cfg.PalCfg.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	}
	if di.ResultImage.Cfg.ScrCfg.Type.IsFullScreen() {
		if im.Cfg.ScrCfg.Mode == 2 || im2.Cfg.ScrCfg.Mode == 2 {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx2FullscreenFormat
		} else {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx1FullscreenFormat
		}
	} else {
		if im.Cfg.ScrCfg.Mode == 2 || im2.Cfg.ScrCfg.Mode == 2 {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx2Format
		} else {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx1Format
		}
	}
	di.ResultImage.CpcResultImage.Image = img
	di.ResultImage.CpcResultImage.Refresh()
	di.ResultImage.PaletteImage.Image = png.PalToImage(im.Palette())
	di.ResultImage.PaletteImage.Refresh()
}

// nolint: funlen
func (m *MartineUI) newEgxTabItem(di *menu.DoubleImageMenu) *fyne.Container {
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
					di.ResultImage.PaletteImage.Image = png.PalToImage(di.ResultImage.Cfg.PalCfg.Palette)
					di.ResultImage.PaletteImage.Refresh()
				}),
				widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
					// export the egx image
					di.ResultImage.Cfg.ScrCfg.IsPlus = di.LeftImage.Cfg.ScrCfg.IsPlus
					m.exportEgxDialog(di.ResultImage.Cfg, m.window)
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

// nolint: funlen
func (m *MartineUI) newEgxImageTransfertTab(me *menu.ImageMenu) *fyne.Container {
	me.SetWindow(m.window)
	importOpen := newImportButton(m, me)

	paletteOpen := pal.NewOpenPaletteButton(me, m.window, nil)

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
			directory.SetImportDirectoryURI(reader.URI())
			me.SetOriginalImagePath(reader.URI())
			img, err := openImage(me.OriginalImagePath())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			me.SetOriginalImage(img)
		}, m.window)
		dir, err := directory.ImportDirectoryURI()
		if err != nil {
			d.SetLocation(dir)
		}
		d.SetFilter(imagesFilesFilter)
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		me.ExportDialog(me.Cfg, me.GetConfig)
	})

	applyButton := widget.NewButtonWithIcon("Apply", theme.VisibilityIcon(), func() {
		log.GetLogger().Infoln("apply.")
		m.ApplyOneImage(me)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	winFormat := me.NewFormatRadio()

	colorReducerLabel := widget.NewLabel("Color reducer")
	colorReducer := widget.NewSelect([]string{"none", "Lower", "Medium", "Strong"}, func(s string) {
		switch s {
		case "none":
			me.Cfg.ScrCfg.Process.Reducer = 0
		case "Lower":
			me.Cfg.ScrCfg.Process.Reducer = 1
		case "Medium":
			me.Cfg.ScrCfg.Process.Reducer = 2
		case "Strong":
			me.Cfg.ScrCfg.Process.Reducer = 3
		}
	})
	colorReducer.SetSelected("none")

	resize := w2.NewResizeAlgorithmSelect(me)
	resizeLabel := widget.NewLabel("Resize algorithm")

	ditheringMultiplier := widget.NewSlider(0., 2.5)
	ditheringMultiplier.Step = 0.1
	ditheringMultiplier.SetValue(1.18)
	ditheringMultiplier.OnChanged = func(f float64) {
		me.Cfg.ScrCfg.Process.Dithering.Multiplier = f
	}
	dithering := w2.NewDitheringSelect(me)

	ditheringWithQuantification := widget.NewCheck("With quantification", func(b bool) {
		me.Cfg.ScrCfg.Process.Dithering.WithQuantification = b
	})

	enableDithering := widget.NewCheck("Enable dithering", func(b bool) {
		me.Cfg.ScrCfg.Process.ApplyDithering = b
	})
	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		me.Cfg.ScrCfg.IsPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", s)
		}
		me.Cfg.ScrCfg.Mode = uint8(mode)
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	brightness := widget.NewSlider(0.0, 1.0)
	brightness.SetValue(1.)
	brightness.Step = .01
	brightness.OnChanged = func(f float64) {
		me.Cfg.ScrCfg.Process.Brightness = f
	}
	saturationLabel := widget.NewLabel("Saturation")
	saturation := widget.NewSlider(0.0, 1.0)
	saturation.SetValue(1.)
	saturation.Step = .01
	saturation.OnChanged = func(f float64) {
		me.Cfg.ScrCfg.Process.Saturation = f
	}
	brightnessLabel := widget.NewLabel("Brightness")

	warningLabel := widget.NewLabel("Setting thoses parameters will affect your palette, you can't force palette.")
	warningLabel.TextStyle = fyne.TextStyle{Bold: true}

	oneLine := widget.NewCheck("Every other line", func(b bool) {
		me.Cfg.ScrCfg.Process.OneLine = b
	})
	oneRow := widget.NewCheck("Every other row", func(b bool) {
		me.Cfg.ScrCfg.Process.OneRow = b
	})

	editButton := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		p := constants.CpcOldPalette
		if me.Cfg.ScrCfg.IsPlus {
			p = constants.CpcPlusPalette
		}
		if me.CpcImage().Image == nil || me.PaletteImage().Image == nil {
			return
		}
		edit := editor.NewEditor(me.CpcImage().Image,
			editor.MagnifyX2,
			me.Palette(),
			p, me.SetImagePalette,
			m.window)

		d := dialog.NewCustom("Editor", "Ok", edit.NewEditor(), m.window)
		size := m.window.Content().Size()
		size = fyne.Size{Width: size.Width, Height: size.Height}
		d.Resize(size)
		d.Show()
		// after the me.CpcImage().Image must be used to export
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
				applyButton,
				exportButton,
				importOpen,
				editButton,
			),
			container.New(
				layout.NewVBoxLayout(),
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
						paletteOpen,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, me.Palette(), m.window, func() {
								forcePalette.SetChecked(true)
							})
						}),
						m.newImageMenuExportButton(me),
						widget.NewButton("Gray", func() {
							if me.Cfg.ScrCfg.IsPlus {
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
