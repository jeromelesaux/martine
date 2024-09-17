package ui

import (
	"errors"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/diskimage"
	"github.com/jeromelesaux/martine/export/impdraw/palette"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/impdraw/tile"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/spritehard"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

// nolint: funlen
func (m *MartineUI) exportSpriteBoard(s *menu.SpriteMenu, w fyne.Window) {
	formatSelect := widget.NewSelect(
		[]string{
			string(config.SpriteImpCatcherExport),
			string(config.SpriteFlatExport),
			string(config.OcpWindowExport),
			string(config.SpriteHardExport),
			string(config.SpriteCompiledExport),
		}, func(v string) {
			s.Cfg.ScrCfg.ResetExport()
			s.Cfg.ScrCfg.AddExport(config.ScreenExport(v))

		})
	cont := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("export type:"),
		formatSelect,
		widget.NewCheck("import all file in Dsk", func(b bool) {
			if b {
				s.Cfg.ContainerCfg.AddExport(config.DskContainer)
			} else {
				s.Cfg.ContainerCfg.RemoveExport(config.DskContainer)
			}
		}),
		widget.NewCheck("export text file", func(b bool) {
			if b {
				s.Cfg.ScrCfg.AddExport(config.AssemblyExport)
			} else {
				s.Cfg.ScrCfg.RemoveExport(config.AssemblyExport)
			}
		}),
		widget.NewCheck("export Json file", func(b bool) {
			if b {
				s.Cfg.ScrCfg.AddExport(config.JsonExport)
			} else {
				s.Cfg.ScrCfg.RemoveExport(config.JsonExport)
			}
		}),
		widget.NewCheck("add amsdos header", func(b bool) {
			s.Cfg.ScrCfg.NoAmsdosHeader = !b
		}),
		widget.NewCheck("apply zigzag", func(b bool) {
			s.Cfg.ZigZag = b
		}),
		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(v string) {
				switch v {
				case "none":
					s.Cfg.ScrCfg.Compression = compression.NONE
				case "rle":
					s.Cfg.ScrCfg.Compression = compression.RLE
				case "rle 16bits":
					s.Cfg.ScrCfg.Compression = compression.RLE16
				case "Lz4 Classic":
					s.Cfg.ScrCfg.Compression = compression.LZ4
				case "Lz4 Raw":
					s.Cfg.ScrCfg.Compression = compression.RawLZ4
				case "zx0 crunch":
					s.Cfg.ScrCfg.Compression = compression.ZX0
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
				s.Cfg.ScrCfg.OutputPath = lu.Path()
				m.ExportSpriteBoard(s)
				// apply and export
			}, m.window)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(savingDialogSize)
			m.CheckAmsdosHeaderExport(s.Cfg.ContainerCfg.HasExport(config.DskContainer), !s.Cfg.ScrCfg.NoAmsdosHeader, fo, m.window)
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

// nolint: funlen, gocognit
func (m *MartineUI) ExportSpriteBoard(s *menu.SpriteMenu) {
	pi := wgt.NewProgressInfinite("Saving...., Please wait.", m.window)
	pi.Show()
	if err := impPalette.SaveKit(
		filepath.Join(s.Cfg.ScrCfg.OutputPath, "SPRITES.KIT"),
		s.Palette(),
		!s.Cfg.ScrCfg.NoAmsdosHeader); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
		return
	}
	cfg := s.Cfg

	if s.Cfg.ScrCfg.IsExport(config.SpriteCompiledExport) {
		spr := make([][]byte, 0)
		for _, v := range s.SpritesData {
			spr = append(spr, v...)
		}
		diffs := animate.AnalyzeSpriteBoard(spr)
		var code string
		for idx, diff := range diffs {
			var routine string
			if s.Cfg.ScrCfg.Type.IsSpriteHard() {
				routine = animate.ExportCompiledSpriteHard(diff)
			} else {
				pi.Hide()
				dialog.NewError(errors.New("no yet implemented, try another option"), m.window).Show()
				// routine = animate.ExportCompiledSprite(diff)
				return
			}
			code += fmt.Sprintf("spr_%.2d:\n", idx)
			code += routine
		}

		kitPalette := palette.KitPalette{}
		for i := 0; i < len(s.Palette()); i++ {
			kitPalette.Colors[i] = constants.NewCpcPlusColor(s.Palette()[i])
		}
		code += "palette:\n"
		code += kitPalette.ToString()

		if err := amsdos.SaveStringOSFile(
			filepath.Join(s.Cfg.ScrCfg.OutputPath, "compiled_sprites.asm"),
			code); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}
	} else {
		if s.Cfg.ScrCfg.IsExport(config.OcpWindowExport) {
			for idxX, v := range s.SpritesData {
				for idxY, v0 := range v {
					filename := filepath.Join(s.Cfg.ScrCfg.OutputPath, fmt.Sprintf("L%.2dC%.2d.WIN", idxX, idxY))
					if err := window.Win(filename, v0, s.Cfg.ScrCfg.Mode, s.Cfg.ScrCfg.Size.Width, s.Cfg.ScrCfg.Size.Height, s.Cfg.ContainerCfg.HasExport(config.DskContainer), cfg); err != nil {
						log.GetLogger().Error("error while exporting sprites error %s\n", err.Error())
					}
					if s.Cfg.ContainerCfg.HasExport(config.DskContainer) {
						s.Cfg.DskFiles = append(s.Cfg.DskFiles, filename)
					}
				}
			}

		} else {
			if s.Cfg.ScrCfg.IsExport(config.SpriteFlatExport) {
				buf := make([]byte, 0)
				if s.Cfg.ZigZag {
					for x := 0; x < len(s.SpritesCollection); x++ {
						for y := 0; y < len(s.SpritesCollection[x]); y++ {
							z := transformation.Zigzag(s.SpritesCollection[x][y])
							sp, _, _, err := sprite.ToSprite(z,
								s.Palette(),
								constants.Size{
									Width:  s.Cfg.ScrCfg.Size.Width,
									Height: s.Cfg.ScrCfg.Size.Height,
								},
								s.Cfg.ScrCfg.Mode,
								s.Cfg,
							)
							if err != nil {
								pi.Hide()
								dialog.NewError(err, m.window).Show()
								return
							}
							buf = append(buf, sp...)
						}
					}
				} else {
					for _, v := range s.SpritesData {
						for _, v0 := range v {
							buf = append(buf, v0...)
						}
					}
				}
				filename := filepath.Join(s.Cfg.ScrCfg.OutputPath, "SPRITES.BIN")
				buf, _ = compression.Compress(buf, s.Cfg.ScrCfg.Compression)
				var err error
				// TODO add amsdos header
				if s.Cfg.ScrCfg.NoAmsdosHeader {
					err = amsdos.SaveAmsdosFile(filename, ".WIN", buf, 2, 0, 0x4000, 0x4000)
					if err != nil {
						pi.Hide()
						dialog.NewError(err, m.window).Show()
						log.GetLogger().Error("Error while saving flat sprites file error %s\n", err.Error())
						return
					}
				} else {
					err = amsdos.SaveOSFile(filename, buf)
					if err != nil {
						pi.Hide()
						dialog.NewError(err, m.window).Show()
						log.GetLogger().Error("Error while saving flat sprites file error %s\n", err.Error())
						return
					}
				}
				if s.Cfg.ContainerCfg.HasExport(config.DskContainer) {
					s.Cfg.DskFiles = append(s.Cfg.DskFiles, filename)
				}
			} else {
				if s.Cfg.ScrCfg.IsExport(config.SpriteImpCatcherExport) {
					buf := make([]byte, 0)
					for _, v := range s.SpritesData {
						for _, v0 := range v {
							buf = append(buf, v0...)
						}
					}
					filename := filepath.Join(s.Cfg.ScrCfg.OutputPath, "sprites.imp")
					if err := tile.Imp(buf, uint(s.SpriteRows*s.SpriteColumns), uint(s.Cfg.ScrCfg.Size.Width), uint(s.Cfg.ScrCfg.Size.Height), uint(s.Cfg.ScrCfg.Mode), filename, cfg); err != nil {
						pi.Hide()
						dialog.NewError(err, m.window).Show()
						log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
						return
					}
					if s.Cfg.ContainerCfg.HasExport(config.DskContainer) {
						s.Cfg.DskFiles = append(s.Cfg.DskFiles, filename)
					}
				} else {
					if s.Cfg.ScrCfg.IsExport(config.SpriteHardExport) {
						data := spritehard.SprImpdraw{}
						for _, v := range s.SpritesData {
							sh := spritehard.SpriteHard{}
							for _, v0 := range v {
								copy(sh.Data[:], v0[:256])
								data.Data = append(data.Data, sh)
							}
						}
						filename := filepath.Join(s.Cfg.ScrCfg.OutputPath, "sprites.spr")

						if err := spritehard.Spr(filename, data, cfg); err != nil {
							pi.Hide()
							dialog.NewError(err, m.window).Show()
							log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
							return
						}
						if s.Cfg.ContainerCfg.HasExport(config.DskContainer) {
							s.Cfg.DskFiles = append(s.Cfg.DskFiles, filename)
						}
					}
				}
			}
		}
	}

	if s.Cfg.ContainerCfg.HasExport(config.DskContainer) {
		if err := diskimage.ImportInDsk("sprites_dsk", s.Cfg); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			log.GetLogger().Error("error while saving files in dsk sprites_dsk.dsk error %v", err)
			return
		}
	}

	if s.Cfg.ScrCfg.IsExport(config.AssemblyExport) {
		data := make([][]byte, 0)
		if s.Cfg.ZigZag {
			for x := 0; x < len(s.SpritesCollection); x++ {
				for y := 0; y < len(s.SpritesCollection[x]); y++ {
					z := transformation.Zigzag(s.SpritesCollection[x][y])
					sp, _, _, err := sprite.ToSprite(z,
						s.Palette(),
						constants.Size{
							Width:  s.Cfg.ScrCfg.Size.Width,
							Height: s.Cfg.ScrCfg.Size.Height,
						},
						s.Cfg.ScrCfg.Mode,
						config.NewMartineConfig("", ""),
					)
					if err != nil {
						pi.Hide()
						dialog.NewError(err, m.window).Show()
						return
					}
					data = append(data, sp)
				}
			}
		} else {
			for _, v := range s.SpritesData {
				data = append(data, v...)
			}
		}
		kitPalette := palette.KitPalette{}
		for i := 0; i < len(s.Palette()); i++ {
			kitPalette.Colors[i] = constants.NewCpcPlusColor(s.Palette()[i])
		}
		header := fmt.Sprintf("' from file %s\n", s.FilePath)
		code := header + ascii.SpritesHardText(data, s.Cfg.ScrCfg.Compression)
		code += "Palette\n"
		code += kitPalette.ToCode()

		filename := filepath.Join(s.Cfg.ScrCfg.OutputPath, "SPRITES.ASM")
		err := amsdos.SaveStringOSFile(filename, code)
		if err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}
	}
	pi.Hide()
	dialog.ShowInformation("Saved", "Your export ended in the folder : "+s.Cfg.ScrCfg.OutputPath, m.window)
}
