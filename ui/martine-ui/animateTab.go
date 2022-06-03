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
					context := m.NewContext(&a.ImageMenu, false)
					if context == nil {
						return
					}
					context.Compression = m.animateExport.ExportCompression
					if a.ExportVersion == 0 {
						a.ExportVersion = animate.DeltaExportV1
					}
					address, err := strconv.ParseUint(a.InitialAddress.Text, 16, 64)
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
						a.IsSprite,
						context,
						uint16(address),
						uint8(a.Mode),
						a.ExportVersion,
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
	m.animate.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(m.animate.Palette))
	refreshUI.OnTapped()
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
	context := m.NewContext(&a.ImageMenu, false)
	if context == nil {
		return
	}
	context.Compression = m.animateExport.ExportCompression
	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	address, err := strconv.ParseUint(a.InitialAddress.Text, 16, 64)
	if err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
		return
	}
	// controle de de la taille de la largeur en fonction du mode
	width := context.Size.Width
	mode := a.Mode
	// get all images from widget imagetable
	if !CheckWidthSize(width, mode) {
		pi.Hide()
		dialog.ShowError(fmt.Errorf("the width in not a multiple of color per pixel, increase the width"), m.window)
		return
	}
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

func (m *MartineUI) ImageIndexToRemove(row, col int) {
	m.animate.ImageToRemoveIndex = col
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
				gifCfg, err := gif.DecodeConfig(fr)
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
					return
				}
				fmt.Println(gifCfg.Height)
				fr.Seek(0, io.SeekStart)
				gifImages, err := gif.DecodeAll(fr)
				if err != nil {
					pi.Hide()
					dialog.ShowError(err, m.window)
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
		d.SetFilter(imagesFilesFilter)
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

	removeButton := widget.NewButtonWithIcon("Remove", theme.DeleteIcon(), func() {
		fmt.Printf("image index to remove %d\n", a.ImageToRemoveIndex)
		images := a.AnimateImages.Images()
		if len(images[0]) <= a.ImageToRemoveIndex {
			return
		}
		images[0] = append(images[0][:a.ImageToRemoveIndex], images[0][a.ImageToRemoveIndex+1:]...)
		canvasImages := make([][]canvas.Image, 0)
		for i := 0; i < len(images); i++ {
			canvasImagesRow := make([]canvas.Image, 0)
			for y := 0; y < len(images[i]); y++ {
				canvasImagesRow = append(canvasImagesRow, *canvas.NewImageFromImage(images[i][y]))
			}
			canvasImages = append(canvasImages, canvasImagesRow)
		}
		a.AnimateImages.Update(&canvasImages, len(canvasImages), len(canvasImages[0]))
		a.AnimateImages.Refresh()
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
	a.AnimateImages.IndexCallbackFunc = m.ImageIndexToRemove

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

	oneLine := widget.NewCheck("Every other line", func(b bool) {
		a.ImageMenu.OneLine = b
	})
	oneRow := widget.NewCheck("Every other row", func(b bool) {
		a.ImageMenu.OneRow = b
	})

	return container.New(
		layout.NewGridLayout(1),
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
				removeButton,
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
				layout.NewGridLayoutWithRows(3),
				container.New(
					layout.NewVBoxLayout(),
					oneLine,
					oneRow,
				),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					&a.PaletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, a.Palette, m.window, m.refreshAnimatePalette)
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
