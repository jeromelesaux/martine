package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
)

var colorIndex int
var colorToChange color.Color

var selectedColor = canvas.NewRectangle(color.White)
var selectedColorContainer *fyne.Container
var swapColorContainer *fyne.Container

func SwapColor(setPalette func(color.Palette), p color.Palette, w fyne.Window, performActionAfter func()) {
	selectedColor.SetMinSize(fyne.NewSize(30, 30))
	selectedColorContainer = container.New(layout.NewMaxLayout(), selectedColor)
	pt := wgt.NewPaletteTable(p, colorChanged, indexColor, nil)

	swapColorContainer = container.NewGridWithColumns(
		2,
		container.NewVBox(
			widget.NewLabel("your palette"),
			pt,
		),
		container.NewVBox(
			container.NewVBox(
				widget.NewLabel("selected color"),
				selectedColorContainer,
			),
			container.NewVBox(
				widget.NewButton("Color", func() {
					picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
						colorToChange = c
					}, w)
					picker.Advanced = true
					picker.Show()
				}),
				widget.NewButton("swap", func() {
					if p == nil {
						return
					}
					p[colorIndex] = colorToChange
					npt := wgt.NewPaletteTable(p, colorChanged, indexColor, nil)
					pt = npt
					if setPalette != nil {
						setPalette(pt.Palette)
					}
					if performActionAfter != nil {
						performActionAfter()
					}

					swapColorContainer.Refresh()
				}))))
	swapColorContainer.Resize(fyne.NewSize(200, 200))
	d := dialog.NewCustom("Swap color", "Ok", swapColorContainer, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func indexColor(index int) {
	colorIndex = index
}

func colorChanged(c color.Color) {
	colorToChange = c
	selectedColor = canvas.NewRectangle(colorToChange)
	selectedColorContainer.Add(selectedColor)
	selectedColorContainer.Refresh()
}
