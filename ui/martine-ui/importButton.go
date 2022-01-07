package ui

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/export/file"
	cgfx "github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func NewImportButton(m *MartineUI, me *menu.ImageMenu) *widget.Button {
	return widget.NewButtonWithIcon("Import", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			me.OriginalImagePath = reader.URI()
			if me.IsFullScreen {

				// open palette widget to get palette
				p, mode, err := file.OverscanPalette(me.OriginalImagePath.Path())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(fmt.Errorf("no palette found"), m.window)
					return
				}
				img, err := cgfx.OverscanToImg(me.OriginalImagePath.Path(), mode, p)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(errors.New("palette is empty"), m.window)
					return
				}
				me.Palette = p
				me.Mode = int(mode)
				modeSelection.SetSelectedIndex(me.Mode)
				me.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
				me.OriginalImage = *canvas.NewImageFromImage(img)
				me.OriginalImage.FillMode = canvas.ImageFillContain
			} else if me.IsSprite {
				// loading sprite file
				//	paletteDialog.OnTapped()
				if len(me.Palette) == 0 {
					dialog.ShowError(errors.New("palette is empty, please import palette first"), m.window)
					return
				}
				img, size, err := cgfx.SpriteToImg(me.OriginalImagePath.Path(), uint8(me.Mode), me.Palette)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				me.Width.SetText(strconv.Itoa(size.Width))
				me.Height.SetText(strconv.Itoa(size.Height))
				me.OriginalImage = *canvas.NewImageFromImage(img)
				me.OriginalImage.FillMode = canvas.ImageFillContain
			} else {
				//loading classical screen
				//	paletteDialog.OnTapped()
				if len(me.Palette) == 0 {
					dialog.ShowError(errors.New("palette is empty,  please import palette first"), m.window)
					return
				}
				img, err := cgfx.ScrToImg(me.OriginalImagePath.Path(), uint8(me.Mode), me.Palette)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				me.OriginalImage = *canvas.NewImageFromImage(img)
				me.OriginalImage.FillMode = canvas.ImageFillContain
			}
			refreshUI.OnTapped()
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		d.Resize(dialogSize)
		d.Show()
	})
}