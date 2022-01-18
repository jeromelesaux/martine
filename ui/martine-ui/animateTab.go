package ui

import (
	"fmt"
	"image/gif"
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
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

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

					context := m.NewContext(&m.animate.ImageMenu, false)
					if context == nil {
						return
					}
					context.Compression = m.animateExport.ExportCompression
					hexa := fmt.Sprintf("%x", a.InitialAddress.Text)
					address, err := strconv.Atoi(hexa)
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					m.animateExport.ExportFolderPath = lu.Path()
					fmt.Println(m.animateExport.ExportFolderPath)
					pi := dialog.NewProgressInfinite("Exporting", "Please wait.", m.window)
					pi.Show()
					code, err := animate.ExportDeltaAnimate(
						a.RawImages[0],
						a.DeltaCollection,
						a.Palette,
						context,
						uint16(address),
						uint8(a.Mode),
					)
					pi.Hide()
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					f, err := os.Create(m.animateExport.ExportFolderPath + string(filepath.Separator) + "code.asm")
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					defer f.Close()
					_, err = f.Write([]byte(code))
					if err != nil {
						dialog.ShowError(err, m.window)
						return
					}
					dialog.ShowInformation("Save", "Your files are save in folder \n"+m.animateExport.ExportFolderPath, m.window)
				}, m.window)
				fo.Show()
			}),
		),
	)

	d := dialog.NewCustom("Export  animation", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) AnimateApply(a *menu.AnimateMenu) {
	context := m.NewContext(&a.ImageMenu, false)
	if context == nil {
		return
	}
	context.Compression = m.animateExport.ExportCompression
	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	hexa := fmt.Sprintf("%x", a.InitialAddress.Text)
	address, err := strconv.Atoi(hexa)
	if err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
		return
	}
	// get all images from widget imagetable
	imgs := a.AnimateImages.Images()[0]
	deltaCollection, rawImages, palette, err := animate.DeltaPackingMemory(imgs, context, uint16(address), uint8(a.Mode))
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	a.DeltaCollection = deltaCollection
	a.Palette = palette
	a.RawImages = rawImages
	a.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(a.Palette))
	refreshUI.OnTapped()
}

func (m *MartineUI) newAnimateTab(a *menu.AnimateMenu) fyne.CanvasObject {
	importOpen := NewImportButton(m, &a.ImageMenu)

	paletteOpen := NewOpenPaletteButton(&a.ImageMenu, m.window)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		a.UsePalette = b
	})

	forceUIRefresh := widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
		s := m.window.Content().Size()
		s.Height += 10.
		s.Width += 10.
		m.window.Resize(s)
		m.window.Canvas().Refresh(&a.OriginalImage)
		m.window.Canvas().Refresh(&a.PaletteImage)
		m.tilemap.TileImages.Refresh()
		m.window.Resize(m.window.Content().Size())
		m.window.Content().Refresh()
	})
	refreshUI = forceUIRefresh
	openFileWidget := widget.NewButton("Add image", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			pi := dialog.NewProgressInfinite("Opening file", "Please wait.", m.window)
			pi.Show()
			path := reader.URI()
			if strings.ToUpper(filepath.Ext(path.Path())) != ".GIF" {
				img, err := openImage(path.Path())
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				if a.IsEmpty {
					a.AnimateImages.SubstitueImage(0, 0, *canvas.NewImageFromImage(img))
				} else {
					a.AnimateImages.AppendImage(*canvas.NewImageFromImage(img), 0)
				}
				a.IsEmpty = false
				pi.Hide()
			} else {
				fr, err := os.Open(path.Path())
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				defer fr.Close()
				gifImages, err := gif.DecodeAll(fr)
				if err != nil {
					dialog.ShowError(err, m.window)
					pi.Hide()
					return
				}
				imgs := animate.ConvertToImage(*gifImages)
				for index, img := range imgs {
					if index == 0 {
						a.AnimateImages.SubstitueImage(0, 0, *canvas.NewImageFromImage(img))
					} else {
						a.AnimateImages.AppendImage(*canvas.NewImageFromImage(img), 0)
					}
					a.IsEmpty = false
				}
				pi.Hide()
			}
			m.window.Resize(m.window.Content().Size())
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".gif", ".png", ".jpeg"}))
		d.Resize(dialogSize)
		d.Show()
	})

	resetButton := widget.NewButtonWithIcon("Reset", theme.CancelIcon(), func() {
		a.AnimateImages.Reset()
		a.IsEmpty = true
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportAnimationDialog(a, m.window)
	})

	applyButton := widget.NewButtonWithIcon("Compute", theme.VisibilityIcon(), func() {
		fmt.Println("compute.")
		m.AnimateApply(a)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	a.OriginalImage = canvas.Image{}
	a.PaletteImage = canvas.Image{}

	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		a.IsCpcPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		a.Mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	a.Width = widget.NewEntry()
	a.Width.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	a.Height = widget.NewEntry()
	a.Height.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	a.AnimateImages = custom_widget.NewEmptyImageTable(fyne.NewSize(menu.AnimateSize, menu.AnimateSize))

	initalAddressLabel := widget.NewLabel("initial address")
	a.InitialAddress = widget.NewEntry()
	a.InitialAddress.SetText("c000")

	isSprite := widget.NewCheck("Is sprite", func(b bool) {
		a.IsSprite = b
	})
	m.animateExport.ExportCompression = -1
	compressData := widget.NewCheck("Compress data", func(b bool) {
		if b {
			m.animateExport.ExportCompression = 0
		} else {
			m.animateExport.ExportCompression = -1
		}
	})

	return container.New(
		layout.NewGridLayoutWithRows(2),
		container.New(
			layout.NewGridLayoutWithRows(1),
			container.NewScroll(
				a.AnimateImages),
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
				forceUIRefresh,
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
							a.Width,
						),
						container.New(
							layout.NewHBoxLayout(),
							heightLabel,
							a.Height,
						),
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(6),

				container.New(
					layout.NewGridLayoutWithColumns(2),
					&a.PaletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, a.Palette, m.window)
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
								if err := file.SaveKit(paletteExportPath+".kit", a.Palette, false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := file.SavePal(paletteExportPath+".pal", a.Palette, uint8(a.Mode), false); err != nil {
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
						e.SetText(a.CmdLine())

						d := dialog.NewCustom("Command line generated",
							"Ok",
							e,
							m.window)
						fmt.Printf("%s\n", a.CmdLine())
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
