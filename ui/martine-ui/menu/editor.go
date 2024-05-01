package menu

import (
	"image"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/jeromelesaux/fyne-io/widget/editor"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/palette"
)

type Editor struct {
	im  *ImageMenu
	ex  *ImageExport
	ca  color.Palette
	e   *editor.Editor
	w   fyne.Window
	sel *widget.Select
}

func NewEditor() *Editor {
	e := &Editor{
		im: NewImageMenu(),
		ca: constants.CpcOldPalette,
		ex: &ImageExport{},
	}
	e.im.SetOriginalImage(image.NewNRGBA(image.Rect(0, 0, constants.Mode2.Width, constants.Mode2.Height)))
	e.im.UsePalette = true
	return e
}

func (e *Editor) imageNPalette(i image.Image, p color.Palette) {
	e.im.SetOriginalImage(i)
	e.im.SetPalette(p)
	e.im.ExportDialog(e.ex)
}
func (e *Editor) refreshEditor() {
	// set new image and palette in editor here
	e.e.NewImageAndPalette(e.im.OriginalImage().Image, e.im.Palette())
}
func (e *Editor) New(w fyne.Window) *fyne.Container {
	e.w = w

	e.e = editor.NewEditor(
		e.im.OriginalImage().Image,
		editor.MagnifyX2,
		e.im.Palette(),
		e.ca,
		e.imageNPalette,
		e.w,
	)
	modes := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {
		mode, err := strconv.Atoi(s)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", s)
		}
		e.im.Mode = mode
	})
	modes.SetSelected("0")
	e.sel = modes
	e.im.SetWindow(e.w)

	return container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			e.im.NewImportButton(e.sel, e.refreshEditor),
			palette.NewOpenPaletteButton(e.im, e.w, e.refreshEditor),
			widget.NewCheck("CPC Plus", func(b bool) {
				e.im.IsCpcPlus = b
				if !b {
					e.e.NewAvailablePalette(constants.CpcOldPalette)
				} else {
					e.e.NewAvailablePalette(constants.CpcPlusPalette)
				}
			}),
			widget.NewLabel("Mode:"),
			modes,
			widget.NewLabel("Format:"),
			e.im.NewFormatRadio(),
		),
		e.e.NewEmbededEditor("Export"),
	)
}
