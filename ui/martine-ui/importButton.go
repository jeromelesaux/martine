package ui

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/convert/screen"
	ovs "github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/png"
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
			SetDefaultDirectoryURI(reader.URI())
			me.SetOriginalImagePath(reader.URI())
			if me.IsFullScreen {

				// open palette widget to get palette
				p, mode, err := overscan.OverscanPalette(me.OriginalImagePath())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(fmt.Errorf("no palette found in selected file, try to normal option and open the associated palette"), m.window)
					return
				}
				img, err := ovs.OverscanToImg(me.OriginalImagePath(), mode, p)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(errors.New("palette is empty"), m.window)
					return
				}
				me.SetPalette(p)
				me.Mode = int(mode)
				modeSelection.SetSelectedIndex(me.Mode)
				me.SetPaletteImage(png.PalToImage(p))
				me.SetOriginalImage(img)
			} else if me.IsSprite {
				// loading sprite file
				if len(me.Palette()) == 0 {
					dialog.ShowError(errors.New("palette is empty, please import palette first"), m.window)
					return
				}
				img, size, err := sprite.SpriteToImg(me.OriginalImagePath(), uint8(me.Mode), me.Palette())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				me.Width().SetText(strconv.Itoa(size.Width))
				me.Height().SetText(strconv.Itoa(size.Height))
				me.SetOriginalImage(img)
			} else {
				// loading classical screen
				if len(me.Palette()) == 0 {
					dialog.ShowError(errors.New("palette is empty,  please import palette first, or select fullscreen option to open a fullscreen option"), m.window)
					return
				}
				img, err := screen.ScrToImg(me.OriginalImagePath(), uint8(me.Mode), me.Palette())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				me.SetOriginalImage(img)
			}
		}, m.window)
		path, err := DefaultDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		d.Resize(dialogSize)
		d.Show()
	})
}
