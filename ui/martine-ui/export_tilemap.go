package ui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/export/diskimage"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

// nolint: funlen
func (m *MartineUI) exportTilemapDialog(w fyne.Window) {
	cont := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("export type:"),
			widget.NewSelect([]string{"sprite", "impdraw", "flat"}, func(s string) {
				switch s {
				case "sprite":
					m.tilemap.ExportImpdraw = false
					m.tilemap.ExportFlat = false
				case "flat":
					m.tilemap.ExportImpdraw = false
					m.tilemap.ExportFlat = true
				case "impdraw":
					m.tilemap.ExportFlat = false
					var width, height int
					var err error
					width, _, err = m.tilemap.GetWidth()
					if err != nil {
						dialog.NewError(err, m.window).Show()
						return
					}
					height, _, err = m.tilemap.GetHeight()
					if err != nil {
						dialog.NewError(err, m.window).Show()
						return
					}

					if !m.IsClassicalTilemap(width, height) {
						m.tilemap.ExportImpdraw = true
					} else {
						err = fmt.Errorf("can not apply impdraw for tiles size. (8x8, 4x8, 4x16 or 8x16)")
						dialog.NewError(err, m.window).Show()
						return
					}
				}
			}),
			widget.NewCheck("import all file in Dsk", func(b bool) {
				m.tilemap.ExportDsk = b
			}),
			widget.NewCheck("export text file", func(b bool) {
				m.tilemap.ExportText = b
			}),
			widget.NewCheck("export Json file", func(b bool) {
				m.tilemap.ExportJson = b
			}),
			widget.NewCheck("add amsdos header", func(b bool) {
				m.tilemap.ExportWithAmsdosHeader = b
			}),
			widget.NewCheck("apply zigzag", func(b bool) {
				m.tilemap.ExportZigzag = b
			}),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					m.tilemap.ExportCompression = 0
				case "rle":
					m.tilemap.ExportCompression = 1
				case "rle 16bits":
					m.tilemap.ExportCompression = 2
				case "Lz4 Classic":
					m.tilemap.ExportCompression = 3
				case "Lz4 Raw":
					m.tilemap.ExportCompression = 4
				case "zx0 crunch":
					m.tilemap.ExportCompression = 5
				}
			}),
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
				m.tilemap.ExportFolderPath = lu.Path()
				log.GetLogger().Infoln(m.tilemapExport.ExportFolderPath)
				m.ExportTilemap(m.tilemap)
				// apply and export
			}, m.window)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(savingDialogSize)
			m.CheckAmsdosHeaderExport(m.tilemap.ExportDsk, m.tilemap.ExportWithAmsdosHeader, fo, m.window)
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) ExportTilemap(t *menu.TilemapMenu) {
	pi := wgt.NewProgressInfinite("Saving...., Please wait.", m.window)
	pi.Show()
	cfg := m.NewConfig(t.ImageMenu, true)
	cfg.OutputPath = t.ExportFolderPath
	if m.IsClassicalTilemap(cfg.Size.Width, cfg.Size.Height) && !m.tilemap.IsSprite {
		filename := filepath.Base(t.OriginalImagePath())
		if err := gfx.ExportTilemapClassical(t.OriginalImage().Image, filename, t.Result, cfg.Size, cfg); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}
		if cfg.Dsk {
			if err := diskimage.ImportInDsk(t.OriginalImagePath(), cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				return
			}
		}
		pi.Hide()
	} else {
		if t.ExportImpdraw {
			if err := gfx.ExportImpdrawTilemap(t.Result, "tilemap", t.Palette(), uint8(t.Mode), cfg.Size, t.OriginalImage().Image, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
			}
			if cfg.Dsk {
				if err := diskimage.ImportInDsk(t.OriginalImagePath(), cfg); err != nil {
					pi.Hide()
					dialog.NewError(err, m.window).Show()
					return
				}
			}
			pi.Hide()
		} else {

			if err := gfx.ExportTilemap(t.Result, "tilemap", t.Palette(), uint8(t.Mode), t.OriginalImage().Image, t.ExportFlat, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
			}
			if cfg.Dsk {
				if err := diskimage.ImportInDsk(t.OriginalImagePath(), cfg); err != nil {
					pi.Hide()
					dialog.NewError(err, m.window).Show()
					return
				}
			}
			pi.Hide()
		}
	}
}
