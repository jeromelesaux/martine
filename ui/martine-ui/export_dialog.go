package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

// nolint: funlen
func (m *MartineUI) exportDialog(ie *menu.ImageExport, w fyne.Window) {
	m2host := widget.NewEntry()
	m2host.SetPlaceHolder("Set your M2 IP here.")
	ie.Reset()

	cont := container.NewVBox(
		container.NewHBox(
			widget.NewCheck("import all file in Dsk", func(b bool) {
				ie.ExportDsk = b
			}),
			widget.NewCheck("export as Go1 and Go2 files", func(b bool) {
				ie.ExportAsGoFiles = b
			}),
			widget.NewCheck("export text file", func(b bool) {
				ie.ExportText = b
			}),
			widget.NewCheck("export Json file", func(b bool) {
				ie.ExportJson = b
			}),
			widget.NewCheck("add amsdos header", func(b bool) {
				ie.ExportWithAmsdosHeader = b
			}),
			widget.NewCheck("apply zigzag", func(b bool) {
				ie.ExportZigzag = b
			}),
			widget.NewCheck("export to M2", func(b bool) {
				ie.ExportToM2 = b
				ie.M2IP = m2host.Text
			}),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					ie.ExportCompression = 0
				case "rle":
					ie.ExportCompression = 1
				case "rle 16bits":
					ie.ExportCompression = 2
				case "Lz4 Classic":
					ie.ExportCompression = 3
				case "Lz4 Raw":
					ie.ExportCompression = 4
				case "zx0 crunch":
					ie.ExportCompression = 5
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
				ie.ExportFolderPath = lu.Path()
				log.GetLogger().Infoln(ie.ExportFolderPath)
				// m.ExportOneImage(m.main)
				m.main.ExportImage(m.imageExport, m.window, m.NewConfig)
				// apply and export
			}, m.window)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(savingDialogSize)

			m.CheckAmsdosHeaderExport(m.imageExport.ExportDsk, m.imageExport.ExportWithAmsdosHeader, fo, m.window)
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}
