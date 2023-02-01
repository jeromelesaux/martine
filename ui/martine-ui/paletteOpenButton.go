package ui

import (
	"image"
	"image/color"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
)

type PaletteInterface interface {
	SetPalette(color.Palette)
	SetPaletteImage(image.Image)
}

func NewOpenPaletteButton(m PaletteInterface, win fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("Palette", theme.ColorChromaticIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}
			SetDefaultDirectoryURI(reader.URI())
			palettePath := reader.URI().Path()
			switch strings.ToLower(filepath.Ext(palettePath)) {
			case ".pal":
				p, _, err := ocpartstudio.OpenPal(palettePath)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				m.SetPalette(p)
				m.SetPaletteImage(png.PalToImage(p))
			case ".kit":
				p, _, err := impPalette.OpenKit(palettePath)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				m.SetPalette(p)
				m.SetPaletteImage(png.PalToImage(p))
			}
		}, win)
		path, err := DefaultDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".pal", ".kit"}))
		d.Resize(dialogSize)
		d.Show()
	})
}
