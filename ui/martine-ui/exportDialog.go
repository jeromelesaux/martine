package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (m *MartineUI) exportDialog(w fyne.Window) {

	cont := container.NewVBox(
		container.NewHBox(
			widget.NewCheck("import all file in Dsk", func(b bool) {
				m.exportDsk = b
			}),
			widget.NewCheck("export text file", func(b bool) {
				m.exportText = b
			}),
			widget.NewCheck("export Json file", func(b bool) {
				m.exportJson = b
			}),
			widget.NewCheck("add amsdos header", func(b bool) {
				m.exportWithAmsdosHeader = b
			}),
			widget.NewCheck("apply zigzag", func(b bool) {
				m.exportZigzag = b
			}),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					m.exportCompression = 0
				case "rle":
					m.exportCompression = 1
				case "rle 16bits":
					m.exportCompression = 2
				case "Lz4 Classic":
					m.exportCompression = 3
				case "Lz4 Raw":
					m.exportCompression = 4
				case "zx0 crunch":
					m.exportCompression = 5
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
				m.exportFolderPath = lu.Path()
				fmt.Println(m.exportFolderPath)
				m.ExportOneImage(m.main)
				// apply and export
			}, m.window)
			fo.Show()
		}),
	)
	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}
