package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ButtonColored struct {
	fyne.Container
	Color  color.Color
	Button *widget.Button
	Tapped func(color.Color)
}

func NewButtonColored(c color.Color, tapped func()) *fyne.Container {

	btn := widget.NewButton("", tapped)
	btnColored := canvas.NewRectangle(c)
	contain := container.New(
		layout.NewMaxLayout(),
		btnColored,
		btn,
	)
	return contain
}
