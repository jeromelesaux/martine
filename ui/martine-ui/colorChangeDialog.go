package ui

import (
	"image/color"

	"fyne.io/fyne/v2/widget"
)

type ColorChangeDialog struct {
	p       color.Palette
	buttons []*widget.Button
}

func NewColorChangeDialog(palette color.Palette, maxColorsDisplayed int) *ColorChangeDialog {
	c := &ColorChangeDialog{p: palette}
	i := 0
	for i < maxColorsDisplayed && i < len(palette) {
	}
	return c
}
