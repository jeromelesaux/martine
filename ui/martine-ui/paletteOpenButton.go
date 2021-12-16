package ui

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/export/file"
)

func NewOpenPaletteButton(m *ImageMenu, win fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("Palette", theme.ColorChromaticIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}
			palettePath := reader.URI().Path()
			switch strings.ToLower(filepath.Ext(palettePath)) {
			case ".pal":
				p, _, err := file.OpenPal(palettePath)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				m.palette = p
				m.paletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
			case ".kit":
				p, _, err := file.OpenKit(palettePath)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				m.palette = p
				m.paletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
			}
			refreshUI.OnTapped()
		}, win)

		d.SetFilter(storage.NewExtensionFileFilter([]string{".pal", ".kit"}))
		d.Resize(dialogSize)
		d.Show()
	})
}
