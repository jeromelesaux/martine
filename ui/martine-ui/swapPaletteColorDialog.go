package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	ui "github.com/jeromelesaux/fyne-io/custom_widget"
)

var colorIndex int
var colorToChange color.Color

func swapColor(setPalette func(color.Palette), p color.Palette, w fyne.Window) {

	pt := ui.NewPaletteTable(p, colorChanged, indexColor, nil)
	var cont *fyne.Container

	cont = container.NewVBox(
		pt,
		widget.NewButton("Color", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorToChange = c
			}, w)
			picker.Advanced = true
			picker.Show()
		}),
		widget.NewButton("swap", func() {
			p[colorIndex] = colorToChange
			npt := ui.NewPaletteTable(p, colorChanged, indexColor, nil)
			pt = npt
			if setPalette != nil {
				setPalette(pt.Palette)
			}
			cont.Refresh()
		}))
	cont.Resize(fyne.NewSize(200, 200))
	d := dialog.NewCustom("Swap color", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func indexColor(index int) {
	colorIndex = index
}

func colorChanged(c color.Color) {
	colorToChange = c
}
