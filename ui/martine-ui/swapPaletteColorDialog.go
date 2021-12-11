package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func swapColor(p color.Palette, w fyne.Window) {
	var colorIndex int
	pt := NewPaletteTable(p, nil, nil, nil)
	var cont *fyne.Container
	var colorToChange color.Color

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
			npt := NewPaletteTable(p, nil, nil, nil)
			pt = npt
			cont.Refresh()
		}))
	cont.Resize(fyne.NewSize(200, 200))
	d := dialog.NewCustom("Swap color", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}
