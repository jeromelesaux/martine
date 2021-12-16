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
)

func NewImportButton(m *MartineUI, me *ImageMenu) *widget.Button {
	return widget.NewButtonWithIcon("Import", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			me.originalImagePath = reader.URI()
			if me.isFullScreen {

				// open palette widget to get palette
				p, mode, err := file.OverscanPalette(me.originalImagePath.Path())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(fmt.Errorf("no palette found"), m.window)
					return
				}
				img, err := cgfx.OverscanToImg(me.originalImagePath.Path(), mode, p)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(errors.New("palette is empty"), m.window)
					return
				}
				me.palette = p
				me.mode = int(mode)
				modeSelection.SetSelectedIndex(me.mode)
				me.paletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
				me.originalImage = *canvas.NewImageFromImage(img)
				me.originalImage.FillMode = canvas.ImageFillContain
			} else if me.isSprite {
				// loading sprite file
				//	paletteDialog.OnTapped()
				if len(me.palette) == 0 {
					dialog.ShowError(errors.New("palette is empty, please import palette first"), m.window)
					return
				}
				img, size, err := cgfx.SpriteToImg(me.originalImagePath.Path(), uint8(me.mode), me.palette)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				me.width.SetText(strconv.Itoa(size.Width))
				me.height.SetText(strconv.Itoa(size.Height))
				me.originalImage = *canvas.NewImageFromImage(img)
				me.originalImage.FillMode = canvas.ImageFillContain
			} else {
				//loading classical screen
				//	paletteDialog.OnTapped()
				if len(me.palette) == 0 {
					dialog.ShowError(errors.New("palette is empty,  please import palette first"), m.window)
					return
				}
				img, err := cgfx.ScrToImg(me.originalImagePath.Path(), uint8(me.mode), me.palette)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				me.originalImage = *canvas.NewImageFromImage(img)
				me.originalImage.FillMode = canvas.ImageFillContain
			}
			refreshUI.OnTapped()
		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		d.Resize(dialogSize)
		d.Show()
	})
}
