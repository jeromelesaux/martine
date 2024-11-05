package ui

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/fyne-io/widget/editor"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/log"

	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	pal "github.com/jeromelesaux/martine/ui/martine-ui/palette"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

type dialogIface interface {
	Show()
}

func (m *MartineUI) CheckAmsdosHeaderExport(inDsk, addAmsdosHeader bool, d dialogIface, win fyne.Window) {
	if inDsk && !addAmsdosHeader {
		dialog.NewConfirm("Warning",
			"You are about to export files in DSK without Amsdos header, continue ? ",
			func(b bool) {
				if b {
					d.Show()
				} else {
					return
				}
			},
			win).Show()

	} else {
		d.Show()
	}
}

func (m *MartineUI) monochromeColor(c color.Color) {
	m.main.SetPalette(image.ColorMonochromePalette(c, m.main.Palette()))
	m.main.SetPaletteImage(png.PalToImage(m.main.Palette()))
}

func (m *MartineUI) ApplyOneImage(me *menu.ImageMenu) {
	me.Edited = false
	cfg := me.NewConfig(true)
	if cfg == nil {
		return
	}

	var inPalette color.Palette
	if me.UsePalette {
		inPalette = me.Palette()
		maxPalette := len(inPalette)
		switch me.Cfg.ScrCfg.Mode {
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
	pi := wgt.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	out, downgraded, palette, _, err := gfx.ApplyOneImage(me.OriginalImage().Image, cfg, inPalette, me.Cfg.ScrCfg.Mode)
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	me.Data = out
	me.Downgraded = downgraded
	if !me.UsePalette {
		me.SetPalette(palette)
	}
	if me.Cfg.ScrCfg.Type.IsSprite() || me.Cfg.ScrCfg.Type.IsSpriteHard() {
		newSize := constants.Size{Width: cfg.ScrCfg.Size.Width * 50, Height: cfg.ScrCfg.Size.Height * 50}
		me.Downgraded = image.Resize(me.Downgraded, newSize, me.Cfg.ScrCfg.Process.ResizingAlgo)
	}
	me.SetCpcImage(me.Downgraded)
	me.SetPaletteImage(png.PalToImage(me.Palette()))
}

// nolint: funlen
func (m *MartineUI) newImageTransfertTab(me *menu.ImageMenu) *fyne.Container {
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
			m.SetImage(img)
			// m.window.Canvas().Refresh(&me.OriginalImage)
			// m.window.Resize(m.window.Content().Size())
		}, m.window)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
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

	kmeansLabel := widget.NewLabel("Reduce palette with Kmeans")
	useKmeans := widget.NewCheck("Use Kmeans", func(b bool) {
		me.Cfg.ScrCfg.Process.Kmeans.Used = b
	})
	kmeansIteration := widget.NewSlider(0.01, .9)
	kmeansIteration.SetValue(.01)
	kmeansIteration.Step = .01
	kmeansIteration.OnChanged = func(f float64) {
		me.Cfg.ScrCfg.Process.Kmeans.Threshold = f
	}

	ditheringMultiplier := widget.NewSlider(0., 5.)
	ditheringMultiplier.Step = 0.1
	ditheringMultiplier.SetValue(.1)
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

	oneLine := widget.NewCheck("Every other line", func(b bool) {
		me.Cfg.ScrCfg.Process.OneLine = b
	})
	oneRow := widget.NewCheck("Every other row", func(b bool) {
		me.Cfg.ScrCfg.Process.OneRow = b
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

	widthLabel := widget.NewLabel("Width")
	me.Width().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	me.Height().Validator = validation.NewRegexp("\\d+", "Must contain a number")

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
					container.New(
						layout.NewHBoxLayout(),
						widthLabel,
						me.Width(),
					),
					container.New(
						layout.NewHBoxLayout(),
						heightLabel,
						me.Height(),
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(7),
				container.New(
					layout.NewGridLayoutWithRows(3),
					container.New(
						layout.NewGridLayoutWithColumns(3),
						kmeansLabel,
						useKmeans,
						kmeansIteration,
					),
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
					oneLine,
					oneRow,
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
								me.SetPalette(image.MonochromePalette(me.Palette()))
								me.SetPaletteImage(png.PalToImage(me.Palette()))
								forcePalette.SetChecked(true)
								forcePalette.Refresh()
							}
						}),

						widget.NewButton("Monochome", func() {
							if me.Cfg.ScrCfg.IsPlus {
								w2.ColorSelector(m.monochromeColor, me.Palette(), m.window, func() {
									forcePalette.SetChecked(true)
								})
							}
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
