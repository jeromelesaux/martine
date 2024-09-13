package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
)

var colorSelectorIndex int
var colorSelectorToChange color.Color

func ColorSelector(setColor func(color.Color), p color.Palette, w fyne.Window, performActionAfter func()) {
	pt := wgt.NewPaletteTable(p, colorSelectorChanged, indexSelectorColor, nil)
	var cont *fyne.Container

	cont = container.NewVBox(
		pt,
		widget.NewButton("select", func() {
			if p == nil {
				return
			}
			p[colorSelectorIndex] = colorSelectorToChange
			npt := wgt.NewPaletteTable(p, colorSelectorChanged, indexSelectorColor, nil)
			pt = npt
			if setColor != nil {
				setColor(colorSelectorToChange)
			}
			if performActionAfter != nil {
				performActionAfter()
			}
			cont.Refresh()
		}))
	cont.Resize(fyne.NewSize(200, 200))
	d := dialog.NewCustom("Select a  color", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func indexSelectorColor(index int) {
	colorSelectorIndex = index
}

func colorSelectorChanged(c color.Color) {
	colorSelectorToChange = c
}
