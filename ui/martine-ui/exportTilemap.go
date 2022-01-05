package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) exportTilemapDialog(w fyne.Window) {

	cont := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("export type:"),
			widget.NewSelect([]string{"sprite", "impdraw"}, func(s string) {
				switch s {
				case "sprite":
					m.tilemap.ExportImpdraw = false
				case "impdraw":
					m.tilemap.ExportImpdraw = true
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
				m.tilemap.ExportFolderPath = lu.Path()
				fmt.Println(m.exportFolderPath)
				m.ExportTilemap(m.tilemap)
				// apply and export
			}, m.window)
			fo.Show()
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) ExportTilemap(t *menu.TilemapMenu) {

}
