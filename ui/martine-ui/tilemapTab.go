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
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/config"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	w2 "github.com/jeromelesaux/martine/ui/martine-ui/widget"
)

func (m *MartineUI) IsClassicalTilemap(width, height int) bool {
	if width == 4 || width == 8 {
		if height == 8 || height == 16 {
			return false
		}
	}
	return true
}

func (m *MartineUI) TilemapApply(me *menu.TilemapMenu) {
	cfg := m.NewConfig(me.ImageMenu, true)
	if cfg == nil {
		return
	}

	cfg.CustomDimension = true

	pi := custom_widget.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	var analyze *transformation.AnalyzeBoard
	var palette color.Palette
	var tiles [][]image.Image
	var err error
	if m.IsClassicalTilemap(cfg.Size.Width, cfg.Size.Height) {
		filename := filepath.Base(me.OriginalImagePath())
		analyze, tiles, palette = gfx.TilemapClassical(uint8(me.Mode), me.IsCpcPlus, filename, me.OriginalImagePath(), me.OriginalImage().Image, cfg.Size, cfg)
		pi.Hide()
	} else {
		analyze, tiles, palette, err = gfx.TilemapRaw(uint8(me.Mode), me.IsCpcPlus, cfg.Size, me.OriginalImage().Image, cfg)
		pi.Hide()
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
	}

	me.Result = analyze
	me.SetPalette(palette)
	tilesCanvas := custom_widget.NewImageTableCache(len(tiles), len(tiles[0]), fyne.NewSize(50, 50))
	for i, v := range tiles {
		for i2, v2 := range v {
			tilesCanvas.Set(i, i2, canvas.NewImageFromImage(v2))
		}
	}
	me.TileImages.Update(tilesCanvas, len(tiles)-1, len(tiles[0])-1)
	me.SetPaletteImage(png.PalToImage(me.Palette()))
}

func (m *MartineUI) newTilemapTab(tm *menu.TilemapMenu) fyne.CanvasObject {
	tm.IsSprite = true
	importOpen := NewImportButton(m, tm.ImageMenu)

	paletteOpen := NewOpenPaletteButton(tm.ImageMenu, m.window)

	forcePalette := widget.NewCheck("use palette", func(b bool) {
		tm.UsePalette = b
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
			tm.SetOriginalImagePath(reader.URI())
			img, err := openImage(tm.OriginalImagePath())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			tm.SetOriginalImage(img)
			// m.window.Canvas().Refresh(&tm.OriginalImage)
			// m.window.Resize(m.window.Content().Size())
		}, m.window)
		path, err := directory.DefaultDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(imagesFilesFilter)
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.tilemap.ResetExport()
		m.exportTilemapDialog(m.window)
	})

	applyButton := widget.NewButtonWithIcon("Compute", theme.VisibilityIcon(), func() {
		fmt.Println("compute.")
		m.TilemapApply(tm)
	})

	openFileWidget.Icon = theme.FileImageIcon()

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

	tm.Width().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	tm.Height().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	tm.TileImages = custom_widget.NewEmptyImageTable(fyne.NewSize(menu.TileSize, menu.TileSize))

	return container.New(
		layout.NewGridLayoutWithRows(2),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			container.New(
				layout.NewGridLayoutWithColumns(1),
				container.NewScroll(
					tm.OriginalImage()),
			),
			container.New(
				layout.NewVBoxLayout(),
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
								tm.Width(),
							),
							container.New(
								layout.NewHBoxLayout(),
								heightLabel,
								tm.Height(),
							),
						),
					),
					container.New(
						layout.NewGridLayoutWithRows(6),

						container.New(
							layout.NewGridLayoutWithColumns(2),
							tm.PaletteImage(),
							container.New(
								layout.NewHBoxLayout(),
								forcePalette,
								widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
									w2.SwapColor(m.SetPalette, tm.Palette(), m.window, nil)
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
										if err := impPalette.SaveKit(paletteExportPath+".kit", tm.Palette(), false); err != nil {
											dialog.ShowError(err, m.window)
										}
										if err := ocpartstudio.SavePal(paletteExportPath+".pal", tm.Palette(), uint8(tm.Mode), false); err != nil {
											dialog.ShowError(err, m.window)
										}
									}, m.window)
									dir, err := directory.DefaultDirectoryURI()
									if err != nil {
										d.SetLocation(dir)
									}
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
			),
		),
		container.NewScroll(
			tm.TileImages),
	)
}
