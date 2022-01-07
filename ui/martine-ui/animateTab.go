package ui

import (
	"fmt"
	"image"
	"image/draw"
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

func (m *MartineUI) AnimateApply(a *menu.AnimateMenu) {

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
	openFileWidget := widget.NewButton("Add mage", func() {
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
				a.AnimateImages.AppendImage(*canvas.NewImageFromImage(img), 0)
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
				for _, img := range imgs {
					a.AnimateImages.AppendImage(*canvas.NewImageFromImage(img), 0)
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
		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(menu.AnimateSize), int(menu.AnimateSize)}})
		bg := theme.BackgroundColor()
		draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0}, draw.Src)
		canvasImg := canvas.NewImageFromImage(img)
		images := make([][]canvas.Image, 1)
		images[0] = make([]canvas.Image, 1)
		images[0][0] = *canvasImg
		a.AnimateImages.Update(&images, 1, 1)
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportTilemapDialog(m.window)
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
	a.InitialAddress.Validator = validation.NewRegexp("#\\s+", "Must be a hexadecimal address")

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
				layout.NewVBoxLayout(),
				container.New(
					layout.NewHBoxLayout(),
					isPlus,
					initalAddressLabel,
					a.InitialAddress,
				),
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
						a.Width,
					),
					container.New(
						layout.NewHBoxLayout(),
						heightLabel,
						a.Height,
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
