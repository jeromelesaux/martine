package ui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) exportSpriteBoard(s *menu.SpriteMenu, w fyne.Window) {

	formatSelect := widget.NewSelect([]string{string(menu.SpriteImpCatcher), string(menu.SpriteFlatExport), string(menu.SpriteFilesExport)}, func(v string) {
		switch menu.SpriteExportFormat(v) {
		case menu.SpriteFlatExport:
			s.ExportFormat = menu.SpriteFlatExport
		case menu.SpriteFilesExport:
			s.ExportFormat = menu.SpriteFilesExport
		case menu.SpriteImpCatcher:
			s.ExportFormat = menu.SpriteImpCatcher
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

}
