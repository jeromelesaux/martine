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
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/diskimage"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

// nolint: funlen, gocognit
func (m *MartineUI) exportTilemapDialog(w fyne.Window) {
	cont := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("export type:"),
			widget.NewSelect([]string{"sprite", "impdraw", "flat"}, func(s string) {
				switch s {
				case "sprite":
					m.tilemap.Cfg.ScrCfg.Type = config.SpriteFormat
				case "flat":
					m.tilemap.Cfg.ScrCfg.Type = config.WindowFormat
				case "impdraw":
					m.tilemap.Cfg.ScrCfg.Type = config.ImpdrawTile
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
						m.tilemap.Cfg.ScrCfg.Type = config.ImpdrawTile
					} else {
						err = fmt.Errorf("can not apply impdraw for tiles size. (8x8, 4x8, 4x16 or 8x16)")
						dialog.NewError(err, m.window).Show()
						return
					}
				}
			}),
			widget.NewCheck("import all file in Dsk", func(b bool) {
				if b {
					m.tilemap.Cfg.ContainerCfg.AddExport(config.DskContainer)
				} else {
					m.tilemap.Cfg.ContainerCfg.RemoveExport(config.DskContainer)
				}
			}),
			widget.NewCheck("export text file", func(b bool) {
				if b {
					m.tilemap.Cfg.ScrCfg.AddExport(config.AssemblyExport)
				} else {
					m.tilemap.Cfg.ScrCfg.RemoveExport(config.AssemblyExport)
				}
			}),
			widget.NewCheck("export Json file", func(b bool) {
				if b {
					m.tilemap.Cfg.ScrCfg.AddExport(config.JsonExport)
				} else {
					m.tilemap.Cfg.ScrCfg.RemoveExport(config.JsonExport)
				}
			}),
			widget.NewCheck("add amsdos header", func(b bool) {
				m.tilemap.Cfg.ScrCfg.NoAmsdosHeader = !b
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
					m.tilemap.Cfg.ScrCfg.Compression = compression.NONE
				case "rle":
					m.tilemap.Cfg.ScrCfg.Compression = compression.RLE
				case "rle 16bits":
					m.tilemap.Cfg.ScrCfg.Compression = compression.RLE16
				case "Lz4 Classic":
					m.tilemap.Cfg.ScrCfg.Compression = compression.LZ4
				case "Lz4 Raw":
					m.tilemap.Cfg.ScrCfg.Compression = compression.RawLZ4
				case "zx0 crunch":
					m.tilemap.Cfg.ScrCfg.Compression = compression.ZX0
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
				m.tilemap.Cfg.ScrCfg.OutputPath = lu.Path()
				log.GetLogger().Infoln(m.tilemap.Cfg.ScrCfg.OutputPath)
				m.ExportTilemap(m.tilemap)
				// apply and export
			}, m.window)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(savingDialogSize)
			m.CheckAmsdosHeaderExport(m.tilemap.Cfg.ContainerCfg.HasExport(config.DskContainer), !m.tilemap.Cfg.ScrCfg.NoAmsdosHeader, fo, m.window)
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) ExportTilemap(t *menu.TilemapMenu) {
	pi := wgt.NewProgressInfinite("Saving...., Please wait.", m.window)
	pi.Show()
	if m.IsClassicalTilemap(t.Cfg.ScrCfg.Size.Width, t.Cfg.ScrCfg.Size.Height) && !t.Cfg.ScrCfg.Type.IsSprite() {
		filename := filepath.Base(t.OriginalImagePath())
		if err := gfx.ExportTilemapClassical(t.OriginalImage().Image, filename, t.Result, t.Cfg.ScrCfg.Size, t.Cfg); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}
		if t.Cfg.HasContainerExport(config.DskContainer) {
			if err := diskimage.ImportInDsk(t.OriginalImagePath(), t.Cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				return
			}
		}
		pi.Hide()
	} else {
		if t.Cfg.ContainerCfg.HasExport(config.DskContainer) {
			if err := gfx.ExportImpdrawTilemap(t.Result, "tilemap", t.Palette(), t.Cfg.ScrCfg.Mode, t.Cfg.ScrCfg.Size, t.OriginalImage().Image, t.Cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
			}
			if t.Cfg.HasContainerExport(config.DskContainer) {
				if err := diskimage.ImportInDsk(t.OriginalImagePath(), t.Cfg); err != nil {
					pi.Hide()
					dialog.NewError(err, m.window).Show()
					return
				}
			}
			pi.Hide()
		} else {

			if err := gfx.ExportTilemap(t.Result, "tilemap", t.Palette(), t.Cfg.ScrCfg.Mode, t.OriginalImage().Image, t.Cfg.ScrCfg.Type == config.SpriteFormat, m.tilemap.Cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
			}
			pi.Hide()
		}
	}
}
