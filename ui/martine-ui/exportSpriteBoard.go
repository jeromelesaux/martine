package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/diskimage"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/impdraw/tile"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/spritehard"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) exportSpriteBoard(s *menu.SpriteMenu, w fyne.Window) {
	formatSelect := widget.NewSelect(
		[]string{
			string(menu.SpriteImpCatcher),
			string(menu.SpriteFlatExport),
			string(menu.SpriteFilesExport),
			string(menu.SpriteCompiled),
			string(menu.SpriteHard),
		}, func(v string) {
			switch menu.SpriteExportFormat(v) {
			case menu.SpriteFlatExport:
				s.ExportFormat = menu.SpriteFlatExport
			case menu.SpriteFilesExport:
				s.ExportFormat = menu.SpriteFilesExport
			case menu.SpriteImpCatcher:
				s.ExportFormat = menu.SpriteImpCatcher
			case menu.SpriteCompiled:
				s.ExportFormat = menu.SpriteCompiled
			case menu.SpriteHard:
				s.ExportFormat = menu.SpriteHard
			default:
				fmt.Fprintf(os.Stderr, "error while getting sprite export format %s\n", v)
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
				s.ExportFolderPath = lu.Path()
				m.ExportSpriteBoard(s)
				// apply and export
			}, m.window)
			fo.Show()
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) ExportSpriteBoard(s *menu.SpriteMenu) {
	pi := custom_widget.NewProgressInfinite("Saving...., Please wait.", m.window)
	pi.Show()
	if err := impPalette.SaveKit(s.ExportFolderPath+string(filepath.Separator)+"SPRITES.KIT", s.Palette(), s.ExportWithAmsdosHeader); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
		return
	}
	switch s.ExportFormat {
	case menu.SpriteCompiled:
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

		if err := amsdos.SaveStringOSFile(s.ExportFolderPath+string(filepath.Separator)+"compiled_sprites.asm", code); err != nil {
			pi.Hide()
			dialog.NewError(err, m.window).Show()
			return
		}

		pi.Hide()
	case menu.SpriteFilesExport:
		for idxX, v := range s.SpritesData {
			for idxY, v0 := range v {
				filename := s.ExportFolderPath + string(filepath.Separator) + fmt.Sprintf("L%.2dC%.2d.WIN", idxX, idxY)
				cfg := config.NewMartineConfig("", s.ExportFolderPath)
				cfg.Compression = s.ExportCompression
				cfg.NoAmsdosHeader = !s.ExportWithAmsdosHeader
				if err := window.Win(filename, v0, uint8(s.Mode), s.SpriteWidth, s.SpriteHeight, s.ExportDsk, cfg); err != nil {
					fmt.Fprintf(os.Stderr, "error while exporting sprites error %s\n", err.Error())
				}
			}
		}
		pi.Hide()
	case menu.SpriteFlatExport:
		buf := make([]byte, 0)
		for _, v := range s.SpritesData {
			for _, v0 := range v {
				buf = append(buf, v0...)
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
				fmt.Fprintf(os.Stderr, "Error while saving flat sprites file error %s\n", err.Error())
				return
			}
		} else {
			err = amsdos.SaveOSFile(filename, buf)
			if err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				fmt.Fprintf(os.Stderr, "Error while saving flat sprites file error %s\n", err.Error())
				return
			}
		}
		if s.ExportDsk {
			cfg := config.NewMartineConfig("", s.ExportFolderPath)
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", filename, err)
				return
			}
		}
		pi.Hide()
	case menu.SpriteImpCatcher:
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
			fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", filename, err)
			return
		}
		if s.ExportDsk {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", filename, err)
				return
			}
		}
		pi.Hide()
	case menu.SpriteHard:
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
			fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", filename, err)
			return
		}
		if s.ExportDsk {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				pi.Hide()
				dialog.NewError(err, m.window).Show()
				fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", filename, err)
				return
			}
		}
		pi.Hide()
	}
	if s.ExportText {
		data := make([][]byte, 0)
		for _, v := range s.SpritesData {
			data = append(data, v...)
		}
		code := ascii.SpritesHardText(data, s.ExportCompression)
		filename := s.ExportFolderPath + string(filepath.Separator) + "SPRITES.ASM"
		amsdos.SaveStringOSFile(filename, code)

	}
	dialog.ShowInformation("Saved", "Your export ended in the folder : "+s.ExportFolderPath, m.window)
}
