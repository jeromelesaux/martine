package ui

import (
	"encoding/json"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/fyne-io/widget/editor"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	exspr "github.com/jeromelesaux/martine/export/sprite"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
	pal "github.com/jeromelesaux/martine/ui/martine-ui/palette"
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
	cfg := me.Cfg
	if cfg == nil {
		return
	}
	var err error
	cfg.ScrCfg.Size.Height, _, err = me.GetHeight()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	cfg.ScrCfg.Size.Width, _, err = me.GetWidth()
	if err != nil {
		dialog.NewError(err, m.window).Show()
		return
	}
	cfg.CustomDimension = true

	pi := wgt.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	var analyze *transformation.AnalyzeBoard
	var palette color.Palette
	var tiles [][]image.Image
	if m.IsClassicalTilemap(cfg.ScrCfg.Size.Width, cfg.ScrCfg.Size.Height) {
		filename := filepath.Base(me.OriginalImagePath())
		analyze, tiles, palette = gfx.TilemapClassical(me.Cfg.ScrCfg.Mode, me.Cfg.ScrCfg.IsPlus, filename, me.OriginalImagePath(), me.OriginalImage().Image, cfg.ScrCfg.Size, cfg, me.Historic)
		pi.Hide()
	} else {
		analyze, tiles, palette, err = gfx.TilemapRaw(me.Cfg.ScrCfg.Mode, me.Cfg.ScrCfg.IsPlus, cfg.ScrCfg.Size, me.OriginalImage().Image, cfg, me.Historic)
		pi.Hide()
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
	}

	me.Result = analyze
	me.SetPalette(palette)
	me.SetNewTilesImages(tiles)
	me.SetPaletteImage(png.PalToImage(me.Palette()))

}

func (m *MartineUI) newImageMenuExportButton(tm *menu.ImageMenu) *widget.Button {
	return widget.NewButtonWithIcon("export", theme.DocumentSaveIcon(), func() {
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
			cfg.ScrCfg.NoAmsdosHeader = false
			if err := impPalette.SaveKit(paletteExportPath+".kit", tm.Palette(), false); err != nil {
				dialog.ShowError(err, m.window)
			}
			if err := ocpartstudio.SavePal(paletteExportPath+".pal", tm.Palette(), tm.Cfg.ScrCfg.Mode, false); err != nil {
				dialog.ShowError(err, m.window)
			}
		}, m.window)
		dir, err := directory.ExportDirectoryURI()
		if err != nil {
			d.SetLocation(dir)
		}
		d.Show()
	})
}

func (m *MartineUI) SetTileImage(newTile image.Image, p color.Palette) {
	if m.tilemap.Result == nil {
		return
	}
	// set the tile to replace in result
	prevTile := m.tilemap.CurrentTile()
	m.tilemap.Result.ReplaceTileAt(m.tilemap.Row, m.tilemap.Col, newTile, prevTile)
	// get all the positions of the new Tile

	// display in tiles display widget
	m.tilemap.Result.Tiles = m.tilemap.Result.TilesImage()
	m.tilemap.SetNewTilesImages(m.tilemap.Result.Tiles)

	m.window.Content().Refresh()
}

// nolint: funlen, gocognit
func (m *MartineUI) newTilemapTab(tm *menu.TilemapMenu) *fyne.Container {
	tm.Cfg.ScrCfg.Size.ColorsAvailable = constants.Mode0.ColorsAvailable // default value
	tm.ImageMenu.SetWindow(m.window)
	tm.Cfg.ScrCfg.Type = config.SpriteFormat
	importOpen := newImportButton(m, tm.ImageMenu)

	paletteOpen := pal.NewOpenPaletteButton(tm.ImageMenu, m.window, nil)

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
			directory.SetImportDirectoryURI(reader.URI())
			tm.SetOriginalImagePath(reader.URI())
			img, err := openImage(tm.OriginalImagePath())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			tm.SetOriginalImage(img)
			tm.Historic = nil
			// m.window.Canvas().Refresh(&tm.OriginalImage)
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
		m.tilemap.ResetExport()
		m.exportTilemapDialog(m.window)
	})

	applyButton := widget.NewButtonWithIcon("Compute", theme.VisibilityIcon(), func() {
		log.GetLogger().Infoln("compute.")
		m.TilemapApply(tm)
	})

	historicOpen := widget.NewButtonWithIcon("Open Historic", theme.FolderOpenIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			directory.SetImportDirectoryURI(reader.URI())

			f, err := os.Open(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			var h exspr.TilesHistorical
			if err := json.NewDecoder(f).Decode(&h); err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			tm.Historic = &h

		}, m.window)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".th", ".TH"}))
		d.Resize(dialogSize)
		d.Show()
	})
	openFileWidget.Icon = theme.FileImageIcon()

	isPlus := widget.NewCheck("CPC Plus", func(b bool) {
		tm.Cfg.ScrCfg.IsPlus = b
	})

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", s)
		}
		tm.Cfg.ScrCfg.Mode = uint8(mode)
		switch mode {
		case 0:
			tm.Cfg.ScrCfg.Size.ColorsAvailable = constants.Mode0.ColorsAvailable
		case 1:
			tm.Cfg.ScrCfg.Size.ColorsAvailable = constants.Mode1.ColorsAvailable
		case 2:
			tm.Cfg.ScrCfg.Size.ColorsAvailable = constants.Mode2.ColorsAvailable
		}
	})
	modes.SetSelected("0")
	modeSelection = modes
	modeLabel := widget.NewLabel("Mode:")

	widthLabel := widget.NewLabel("Width")

	tm.Width().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	heightLabel := widget.NewLabel("Height")
	tm.Height().Validator = validation.NewRegexp("\\d+", "Must contain a number")

	tm.TileImages = wgt.NewEmptyImageTable(fyne.NewSize(menu.TileSize, menu.TileSize), tm.TileSelected)

	editButton := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		p := constants.CpcOldPalette
		if tm.Cfg.ScrCfg.IsPlus {
			p = constants.CpcPlusPalette
		}
		if tm.TileImage() == nil || tm.PaletteImage().Image == nil {
			return
		}
		edit := editor.NewEditor(tm.TileImage(),
			editor.MagnifyX8,
			tm.Palette(),
			p, m.SetTileImage,
			m.window)

		d := dialog.NewCustom("Editor Selected Tile", "Ok", edit.NewEditor(), m.window)
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
				tm.OriginalImage()),

			container.NewScroll(
				tm.TileImages),
		),
		container.New(

			layout.NewVBoxLayout(),
			container.New(
				layout.NewVBoxLayout(),
				container.New(
					layout.NewHBoxLayout(),
					openFileWidget,
					exportButton,
					importOpen,
				),
				container.New(
					layout.NewHBoxLayout(),
					historicOpen,
				),
				container.New(
					layout.NewHBoxLayout(),

					container.New(
						layout.NewVBoxLayout(),
						container.New(
							layout.NewVBoxLayout(),
							isPlus,
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
						layout.NewGridLayoutWithRows(2),
						tm.PaletteImage(),
						container.New(
							layout.NewHBoxLayout(),
							forcePalette,
							paletteOpen,
							widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
								w2.SwapColor(tm.SetPalette, tm.Palette(), m.window, nil)
							}),
							m.newImageMenuExportButton(tm.ImageMenu),
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
							log.GetLogger().Info("%s\n", tm.CmdLine())
							size := m.window.Content().Size()
							size = fyne.Size{Width: size.Width / 2, Height: size.Height / 2}
							d.Resize(size)
							d.Show()
						}),

						applyButton,
						editButton,
					),
				),
			),
		),
	)

}
