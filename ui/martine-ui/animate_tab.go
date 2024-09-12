package ui

import (
	"fmt"
	"image/gif"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	pal "github.com/jeromelesaux/martine/ui/martine-ui/palette"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

// nolint: funlen
func (m *MartineUI) exportAnimationDialog(a *menu.AnimateMenu, w fyne.Window) {
	cont := container.NewVBox(
		container.NewHBox(
			widget.NewButtonWithIcon("Export into folder", theme.DocumentSaveIcon(), func() {
				fo := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					if lu == nil {
						// cancel button
						return
					}
					directory.SetExportDirectoryURI(lu)
					cfg := a.Cfg
					if cfg == nil {
						return
					}
					if a.ExportVersion == 0 {
						a.ExportVersion = animate.DeltaExportV1
					}
					address, err := strconv.ParseUint(a.InitialAddress.Text, 16, 64)
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					a.Cfg.ScreenCfg.OutputPath = lu.Path()
					log.GetLogger().Infoln(a.Cfg.ScreenCfg.OutputPath)
					pi := wgt.NewProgressInfinite("Exporting, please wait.", m.window)
					pi.Show()
					code, err := animate.ExportDeltaAnimate(
						a.RawImages[0],
						a.DeltaCollection,
						a.Palette(),
						a.Cfg.ScreenCfg.Type.IsSprite(),
						cfg,
						uint16(address),
						uint8(a.Mode),
						"",
						a.ExportVersion,
					)
					pi.Hide()
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					err = amsdos.SaveOSFile(a.Cfg.ScreenCfg.OutputPath+string(filepath.Separator)+"code.asm", []byte(code))
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					dialog.ShowInformation("Save", "Your files are save in folder \n"+a.Cfg.ScreenCfg.OutputPath, m.window)
				}, m.window)
				d, err := directory.ExportDirectoryURI()
				if err == nil {
					fo.SetLocation(d)
				}
				fo.Resize(savingDialogSize)
				fo.Show()
			}),
			widget.NewLabel("Export version (V1 not optimized, V2 optimized)"),
			widget.NewSelect([]string{"Version 1", "Version 2"}, func(v string) {
				switch v {
				case "Version 1":
					a.ExportVersion = animate.DeltaExportV1
				case "Version 2":
					a.ExportVersion = animate.DeltaExportV2
				default:
					a.ExportVersion = animate.DeltaExportV1
				}
			}),
		),
	)

	d := dialog.NewCustom("Export  animation", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) refreshAnimatePalette() {
	m.animate.SetPaletteImage(png.PalToImage(m.animate.Palette()))
}

func CheckWidthSize(width, mode int) bool {
	var colorPerPixel int

	switch mode {
	case 0:
		colorPerPixel = 2
	case 1:
		colorPerPixel = 4
	case 2:
		colorPerPixel = 8
	}
	remain := width % colorPerPixel
	return remain == 0
}

func (m *MartineUI) AnimateApply(a *menu.AnimateMenu) {
	cfg := a.Cfg
	if cfg == nil {
		return
	}

	pi := wgt.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	address, err := strconv.ParseUint(a.InitialAddress.Text, 16, 64)
	if err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
		return
	}
	// controle de de la taille de la largeur en fonction du mode
	width := cfg.ScreenCfg.Size.Width
	mode := a.Mode
	// get all images from widget imagetable
	if !CheckWidthSize(width, mode) {
		pi.Hide()
		dialog.ShowError(fmt.Errorf("the width in not a multiple of color per pixel, increase the width"), m.window)
		return
	}
	imgs := a.AnimateImages.Images()
	deltaCollection, rawImages, palette, err := animate.DeltaPackingMemory(imgs, cfg, uint16(address), uint8(a.Mode))
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	a.DeltaCollection = deltaCollection
	a.SetPalette(palette)
	a.RawImages = rawImages
	a.SetPaletteImage(png.PalToImage(a.Palette()))
}

// nolint:funlen, gocognit
func (m *MartineUI) newAnimateTab(a *menu.AnimateMenu) *fyne.Container {
	a.ImageMenu.SetWindow(m.window)
	importOpen := newImportButton(m, a.ImageMenu)

	paletteOpen := pal.NewOpenPaletteButton(a.ImageMenu, m.window, nil)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		a.UsePalette = b
	})

	openFileWidget := widget.NewButton("Add image", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			directory.SetImportDirectoryURI(reader.URI())
			pi := wgt.NewProgressInfinite("Opening file, Please wait.", m.window)
			pi.Show()
			path := reader.URI()
			directory.SetImportDirectoryURI(reader.URI())
			if strings.ToUpper(filepath.Ext(path.Path())) != ".GIF" {
				img, err := openImage(path.Path())
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				a.AnimateImages.Append(canvas.NewImageFromImage(img))
				pi.Hide()
			} else {
				fr, err := os.Open(path.Path())
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				defer fr.Close()
				gifCfg, err := gif.DecodeConfig(fr)
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				log.GetLogger().Infoln(gifCfg.Height)
				_, err = fr.Seek(0, io.SeekStart)
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				gifImages, err := gif.DecodeAll(fr)
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				imgs := image.GifToImages(*gifImages)
				for _, img := range imgs {
					a.AnimateImages.Append(canvas.NewImageFromImage(img))
				}
				canvas.Refresh(m.window.Content())
				pi.Hide()
			}
			m.window.Resize(m.window.Content().Size())
		}, m.window)

		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(imagesFilesFilter)
		d.Resize(dialogSize)
		d.Show()
	})

	resetButton := widget.NewButtonWithIcon("Reset", theme.CancelIcon(), func() {
		a.AnimateImages.Reset()
		canvas.Refresh(a.AnimateImages)
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportAnimationDialog(a, m.window)
	})

	applyButton := widget.NewButtonWithIcon("Compute", theme.VisibilityIcon(), func() {
		log.GetLogger().Infoln("compute.")
		m.AnimateApply(a)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		a.Cfg.ScreenCfg.IsPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", s)
		}
		a.Mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	a.Width().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	a.Height().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	initalAddressLabel := widget.NewLabel("initial address")
	a.InitialAddress = widget.NewEntry()
	a.InitialAddress.SetText("c000")

	isSprite := widget.NewCheck("Is sprite", func(b bool) {
		a.Cfg.ScreenCfg.Type = config.SpriteFormat
	})
	a.Cfg.ScreenCfg.Compression = compression.NONE
	compressData := widget.NewCheck("Compress data", func(b bool) {
		if b {
			a.Cfg.ScreenCfg.Compression = compression.LZ4
		} else {
			a.Cfg.ScreenCfg.Compression = compression.NONE
		}
	})

	oneLine := widget.NewCheck("Every other line", func(b bool) {
		a.ImageMenu.OneLine = b
	})
	oneRow := widget.NewCheck("Every other row", func(b bool) {
		a.ImageMenu.OneRow = b
	})

	return container.New(
		layout.NewGridLayout(1),
		container.New(
			layout.NewGridLayoutWithColumns(1),
			container.NewScroll(a.AnimateImages.Container),
		),

		container.New(
			layout.NewVBoxLayout(),
			container.New(
				layout.NewHBoxLayout(),
				openFileWidget,
				resetButton,
				paletteOpen,
				applyButton,
				exportButton,
				importOpen,
			),
			container.New(
				layout.NewGridLayoutWithColumns(2),
				container.New(
					layout.NewVBoxLayout(),
					isPlus,
					container.New(
						layout.NewVBoxLayout(),
						initalAddressLabel,
						a.InitialAddress,
					),
				),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					container.New(
						layout.NewVBoxLayout(),
						isSprite,
						compressData,
					),
					container.New(
						layout.NewVBoxLayout(),
						container.New(
							layout.NewHBoxLayout(),
							modeLabel,
							modes,
						),
						container.New(
							layout.NewHBoxLayout(),
							widthLabel,
							a.Width(),
						),
						container.New(
							layout.NewHBoxLayout(),
							heightLabel,
							a.Height(),
						),
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(3),
				container.New(
					layout.NewVBoxLayout(),
					oneLine,
					oneRow,
				),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					a.PaletteImage(),
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, a.Palette(), m.window, m.refreshAnimatePalette)
						}),
						m.newImageMenuExportButton(a.ImageMenu),
					),
				),

				container.New(
					layout.NewVBoxLayout(),
					widget.NewButton("show cmd", func() {
						e := widget.NewMultiLineEntry()
						e.SetText(a.CmdLine())

						d := dialog.NewCustom("Command line generated",
							"Ok",
							e,
							m.window)
						log.GetLogger().Info("%s\n", a.CmdLine())
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
