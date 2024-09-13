package ui

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/convert/export"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/diskimage"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/m4"
	"github.com/jeromelesaux/martine/log"

	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/snapshot"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

var egxFilename = "aa.scr"

// nolint:funlen
func (m *MartineUI) exportEgxDialog(cfg *config.MartineConfig, w fyne.Window) {
	m2host := widget.NewEntry()
	m2host.SetPlaceHolder("Set your M2 IP here.")

	cont := container.NewVBox(
		container.NewHBox(
			widget.NewCheck("import all file in Dsk", func(b bool) {
				if b {
					cfg.ContainerCfg.AddExport(config.DskContainer)
				} else {
					cfg.ContainerCfg.RemoveExport(config.DskContainer)
				}

			}),
			widget.NewCheck("export text file", func(b bool) {
				if b {
					cfg.ScrCfg.AddExport(config.AssemblyExport)
				} else {
					cfg.ScrCfg.RemoveExport(config.AssemblyExport)
				}
			}),
			widget.NewCheck("export Json file", func(b bool) {
				if b {
					cfg.ScrCfg.AddExport(config.JsonExport)
				} else {
					cfg.ScrCfg.RemoveExport(config.JsonExport)
				}
			}),
			widget.NewCheck("add amsdos header", func(b bool) {
				cfg.ScrCfg.NoAmsdosHeader = b == false

			}),
			widget.NewCheck("apply zigzag", func(b bool) {
				cfg.ZigZag = b
			}),
			widget.NewCheck("export to M2", func(b bool) {
				cfg.M4cfg.Enabled = true
				cfg.M4cfg.Host = m2host.Text
			}),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					cfg.ScrCfg.Compression = compression.NONE
				case "rle":
					cfg.ScrCfg.Compression = compression.RLE
				case "rle 16bits":
					cfg.ScrCfg.Compression = compression.RLE16
				case "Lz4 Classic":
					cfg.ScrCfg.Compression = compression.LZ4
				case "Lz4 Raw":
					cfg.ScrCfg.Compression = compression.RawLZ4
				case "zx0 crunch":
					cfg.ScrCfg.Compression = compression.ZX0
				}
			}),
		m2host,
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
				cfg.ScrCfg.OutputPath = lu.Path()
				m.egx.ResultImage.Path = lu.Path()
				log.GetLogger().Infoln(cfg.ScrCfg.OutputPath)
				m.ExportEgxImage(m.egx)

				// apply and export
			}, m.window)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(savingDialogSize)
			m.CheckAmsdosHeaderExport(cfg.ContainerCfg.HasExport(config.DskContainer),
				cfg.ScrCfg.NoAmsdosHeader == false, fo, m.window)
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

// nolint: funlen
func (m *MartineUI) ExportEgxImage(me *menu.DoubleImageMenu) {
	pi := wgt.NewProgressInfinite("Saving...., please wait.", m.window)
	pi.Show()
	cfg := me.LeftImage.Cfg
	if cfg == nil {
		pi.Hide()
		return
	}
	cfg.ScrCfg.OutputPath = me.ResultImage.Path

	if me.ResultImage.EgxType == 1 {
		cfg.ScrCfg.Type = config.Egx1Format
	} else {
		cfg.ScrCfg.Type = config.Egx2Format
	}
	cfg.EgxMode1 = uint8(me.LeftImage.Cfg.ScrCfg.Mode)
	cfg.EgxMode2 = uint8(me.RightImage.Cfg.ScrCfg.Mode)

	// palette export
	defer func() {
		os.Remove("temporary_palette.kit")
	}()
	if err := impPalette.SaveKit("temporary_palette.kit", me.ResultImage.Palette, false); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
	}
	cfg.PalCfg.Path = "temporary_palette.kit"

	if !cfg.ScrCfg.Type.IsFullScreen() {
		if err := ocpartstudio.EgxLoader(me.ResultImage.Path+string(filepath.Separator)+egxFilename, me.ResultImage.Palette, uint8(me.LeftImage.Cfg.ScrCfg.Mode), uint8(me.RightImage.Cfg.ScrCfg.Mode), cfg); err != nil {
			pi.Hide()
			dialog.ShowError(err, m.window)
			return
		}
	} else {
		if err := export.Export(me.ResultImage.Path+string(filepath.Separator)+egxFilename, me.ResultImage.Data, me.ResultImage.Palette, uint8(me.ResultImage.EgxType), cfg); err != nil {
			pi.Hide()
			dialog.ShowError(err, m.window)
			return
		}
	}
	if cfg.HasContainerExport(config.DskContainer) {
		if err := diskimage.ImportInDsk(me.ResultImage.Path, cfg); err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
	}
	if cfg.HasContainerExport(config.SnaContainer) {
		if cfg.ScrCfg.Type.IsFullScreen() {
			var gfxFile string
			for _, v := range cfg.DskFiles {
				if filepath.Ext(v) == ".SCR" {
					gfxFile = v
					break
				}
			}
			cfg.ContainerCfg.Path = filepath.Join(me.ResultImage.Path, "test.sna")
			if err := snapshot.ImportInSna(gfxFile, cfg.ContainerCfg.Path, 0); err != nil {
				dialog.NewError(err, m.window).Show()
				return
			}
		}
	}
	if cfg.M4cfg.Enabled {
		if err := m4.ImportInM4(cfg); err != nil {
			dialog.NewError(err, m.window).Show()
			log.GetLogger().Error("Cannot send to M4 error :%v\n", err)
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in folder \n"+cfg.ScrCfg.OutputPath, m.window)
}
