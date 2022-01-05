package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func colorPicked(c color.Color, w fyne.Window) {
	rectangle := canvas.NewRectangle(c)
	size := 2 * theme.IconInlineSize()
	rectangle.SetMinSize(fyne.NewSize(size, size))
	dialog.ShowCustom("Color", "Ok", rectangle, w)
}

func NewColorPicker(win fyne.Window) *widget.Button {
	return widget.NewButton("Change Color", func() {
		picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
			colorPicked(c, win)
		}, win)
		picker.Advanced = true
		picker.Show()
	})
}
