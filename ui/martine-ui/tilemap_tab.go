package ui

import (
	"fmt"
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
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

func (m *MartineUI) TilemapApply(me *menu.TilemapMenu) {
	context := m.NewContext(&me.ImageMenu)
	if context == nil {
		return
	}

	context.CustomDimension = true

	pi := dialog.NewProgressInfinite("Computing", "Please wait.", m.window)
	pi.Show()
	analyze, tiles, palette, err := gfx.TilemapMemory(uint8(me.Mode), context.Size, me.OriginalImage.Image, context)
	pi.Hide()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	me.Result = analyze
	me.Palette = palette
	tilesCanvas := make([][]canvas.Image, len(tiles))
	for i, v := range tiles {
		tilesCanvas[i] = make([]canvas.Image, len(v))
		for i2, v2 := range v {
			tilesCanvas[i][i2] = *canvas.NewImageFromImage(v2)
		}
	}
	me.TileImages.Update(&tilesCanvas, len(tiles)-1, len(tiles[0])-1)
	me.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(me.Palette))
	refreshUI.OnTapped()
}

func (m *MartineUI) newTilemapTab(tm *menu.TilemapMenu) fyne.CanvasObject {
	tm.IsSprite = true
	importOpen := NewImportButton(m, &tm.ImageMenu)

	paletteOpen := NewOpenPaletteButton(&tm.ImageMenu, m.window)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		tm.UsePalette = b
	})

	forceUIRefresh := widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
		s := m.window.Content().Size()
		s.Height += 10.
		s.Width += 10.
		m.window.Resize(s)
		m.window.Canvas().Refresh(&tm.OriginalImage)
		m.window.Canvas().Refresh(&tm.PaletteImage)
		m.tilemap.TileImages.Refresh()
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

			tm.OriginalImagePath = reader.URI()
			img, err := openImage(tm.OriginalImagePath.Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			canvasImg := canvas.NewImageFromImage(img)
			tm.OriginalImage = *canvas.NewImageFromImage(canvasImg.Image)
			tm.OriginalImage.FillMode = canvas.ImageFillContain
			m.window.Canvas().Refresh(&tm.OriginalImage)
			m.window.Resize(m.window.Content().Size())
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".gif", ".png", ".jpeg"}))
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportDialog(m.window)
	})

	applyButton := widget.NewButtonWithIcon("Compute", theme.VisibilityIcon(), func() {
		fmt.Println("compute.")
		m.TilemapApply(tm)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	tm.OriginalImage = canvas.Image{}
	tm.PaletteImage = canvas.Image{}

	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		tm.IsCpcPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", s)
		}
		tm.Mode = mode
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")
	tm.Width = widget.NewEntry()
	tm.Width.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	tm.Height = widget.NewEntry()
	tm.Height.Validator = validation.NewRegexp("\\d+", "Must contain a number")

	tm.TileImages = custom_widget.NewEmptyImageTable(fyne.NewSize(20, 20))

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.NewScroll(
				&tm.OriginalImage),
			container.NewScroll(
				tm.TileImages),
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
						tm.Width,
					),
					container.New(
						layout.NewHBoxLayout(),
						heightLabel,
						tm.Height,
					),
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(6),

				container.New(
					layout.NewGridLayoutWithColumns(2),
					&tm.PaletteImage,
					container.New(
						layout.NewHBoxLayout(),
						forcePalette,
						widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
							w2.SwapColor(m.SetPalette, tm.Palette, m.window)
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
								if err := file.SaveKit(paletteExportPath+".kit", tm.Palette, false); err != nil {
									dialog.ShowError(err, m.window)
								}
								if err := file.SavePal(paletteExportPath+".pal", tm.Palette, uint8(tm.Mode), false); err != nil {
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
						e.SetText(tm.CmdLine())

						d := dialog.NewCustom("Command line generated",
							"Ok",
							e,
							m.window)
						fmt.Printf("%s\n", tm.CmdLine())
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
