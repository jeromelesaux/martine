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
	"github.com/jeromelesaux/martine/export"
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
			string(export.SpriteImpCatcher),
			string(export.SpriteFlatExport),
			string(export.OcpWinExport),
			string(export.SpriteCompiled),
			string(export.SpriteHard),
		}, func(v string) {
			switch export.ExportFormat(v) {
			case export.SpriteFlatExport:
				s.ExportFormat = export.SpriteFlatExport
			case export.OcpWinExport:
				s.ExportFormat = export.OcpWinExport
			case export.SpriteImpCatcher:
				s.ExportFormat = export.SpriteImpCatcher
			case export.SpriteCompiled:
				s.ExportFormat = export.SpriteCompiled
			case export.SpriteHard:
				s.ExportFormat = export.SpriteHard
			default:
				log.GetLogger().Error("error while getting sprite export format %s\n", v)
			}
		})
	cont := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("export type:"),
		formatSelect,
		widget.NewCheck("import all file in Dsk", func(b bool) {
			s.ExportDsk = b
		}),
		widget.NewCheck("export text file", func(b bool) {
			s.ExportText = b
		}),
		widget.NewCheck("export Json file", func(b bool) {
			s.ExportJson = b
		}),
		widget.NewCheck("add amsdos header", func(b bool) {
			s.ExportWithAmsdosHeader = b
		}),
		widget.NewCheck("apply zigzag", func(b bool) {
			s.ExportZigzag = b
		}),
		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(v string) {
				switch v {
				case "none":
					s.ExportCompression = 0
				case "rle":
					s.ExportCompression = 1
				case "rle 16bits":
					s.ExportCompression = 2
				case "Lz4 Classic":
					s.ExportCompression = 3
				case "Lz4 Raw":
					s.ExportCompression = 4
				case "zx0 crunch":
					s.ExportCompression = 5
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
				s.ExportFolderPath = lu.Path()
				m.ExportSpriteBoard(s)
				// apply and export
			}, m.window)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(savingDialogSize)
			m.CheckAmsdosHeaderExport(s.ExportDsk, s.ExportWithAmsdosHeader, fo, m.window)
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
	if err := impPalette.SaveKit(s.ExportFolderPath+string(filepath.Separator)+"SPRITES.KIT", s.Palette(), s.ExportWithAmsdosHeader); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
		return
	}
	switch s.ExportFormat {
	case export.SpriteCompiled:
		spr := make([][]byte, 0)
		for _, v := range s.SpritesData {
			spr = append(spr, v...)
		}
		diffs := animate.AnalyzeSpriteBoard(spr)
		var code string
		for idx, diff := range diffs {
			var routine string
			if s.IsHardSprite {
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

		if err := amsdos.SaveStringOSFile(s.ExportFolderPath+string(filepath.Separator)+"compiled_sprites.asm", code); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}

		pi.Hide()
	case export.OcpWinExport:
		for idxX, v := range s.SpritesData {
			for idxY, v0 := range v {
				filename := s.ExportFolderPath + string(filepath.Separator) + fmt.Sprintf("L%.2dC%.2d.WIN", idxX, idxY)
				cfg := config.NewMartineConfig("", s.ExportFolderPath)
				cfg.Compression = s.ExportCompression
				cfg.NoAmsdosHeader = !s.ExportWithAmsdosHeader
				if err := window.Win(filename, v0, uint8(s.Mode), s.SpriteWidth, s.SpriteHeight, s.ExportDsk, cfg); err != nil {
					log.GetLogger().Error("error while exporting sprites error %s\n", err.Error())
				}
			}
		}
		pi.Hide()
	case export.SpriteFlatExport:
		buf := make([]byte, 0)
		if s.ExportZigzag {
			for x := 0; x < len(s.SpritesCollection); x++ {
				for y := 0; y < len(s.SpritesCollection[x]); y++ {
					z := transformation.Zigzag(s.SpritesCollection[x][y])
					sp, _, _, err := sprite.ToSprite(z,
						s.Palette(),
						constants.Size{
							Width:  s.SpriteWidth,
							Height: s.SpriteHeight,
						},
						uint8(s.Mode),
						config.NewMartineConfig("", ""),
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
		filename := s.ExportFolderPath + string(filepath.Separator) + "SPRITES.BIN"
		buf, _ = compression.Compress(buf, s.ExportCompression)
		var err error
		// TODO add amsdos header
		if s.ExportWithAmsdosHeader {
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
		if s.ExportDsk {
			cfg := config.NewMartineConfig("", s.ExportFolderPath)
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
				return
			}
		}
		pi.Hide()
	case export.SpriteImpCatcher:
		buf := make([]byte, 0)
		for _, v := range s.SpritesData {
			for _, v0 := range v {
				buf = append(buf, v0...)
			}
		}
		filename := s.ExportFolderPath + string(filepath.Separator) + "sprites.imp"
		cfg := config.NewMartineConfig("", s.ExportFolderPath)
		cfg.Compression = s.ExportCompression
		cfg.NoAmsdosHeader = !s.ExportWithAmsdosHeader
		if err := tile.Imp(buf, uint(s.SpriteRows*s.SpriteColumns), uint(s.SpriteWidth), uint(s.SpriteHeight), uint(s.Mode), filename, cfg); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
			return
		}
		if s.ExportDsk {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
				return
			}
		}
		pi.Hide()
	case export.SpriteHard:
		data := spritehard.SprImpdraw{}
		for _, v := range s.SpritesData {
			sh := spritehard.SpriteHard{}
			for _, v0 := range v {
				copy(sh.Data[:], v0[:256])
				data.Data = append(data.Data, sh)
			}
		}
		filename := s.ExportFolderPath + string(filepath.Separator) + "sprites.spr"
		cfg := config.NewMartineConfig("", s.ExportFolderPath)
		cfg.Compression = s.ExportCompression
		cfg.NoAmsdosHeader = !s.ExportWithAmsdosHeader
		if err := spritehard.Spr(filename, data, cfg); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
			return
		}
		if s.ExportDsk {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				log.GetLogger().Error("Cannot export to Imp-Catcher the image %s error %v", filename, err)
				return
			}
		}
		pi.Hide()
	}
	if s.ExportText {
		data := make([][]byte, 0)
		if s.ExportZigzag {
			for x := 0; x < len(s.SpritesCollection); x++ {
				for y := 0; y < len(s.SpritesCollection[x]); y++ {
					z := transformation.Zigzag(s.SpritesCollection[x][y])
					sp, _, _, err := sprite.ToSprite(z,
						s.Palette(),
						constants.Size{
							Width:  s.SpriteWidth,
							Height: s.SpriteHeight,
						},
						uint8(s.Mode),
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
		code := header + ascii.SpritesHardText(data, s.ExportCompression)
		code += "Palette\n"
		code += kitPalette.ToCode()

		filename := s.ExportFolderPath + string(filepath.Separator) + "SPRITES.ASM"
		err := amsdos.SaveStringOSFile(filename, code)
		if err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}
		pi.Hide()
	}
	dialog.ShowInformation("Saved", "Your export ended in the folder : "+s.ExportFolderPath, m.window)
}
