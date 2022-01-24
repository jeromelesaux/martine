package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (m *MartineUI) newGreedings() fyne.CanvasObject {
	return container.New(
		layout.NewHBoxLayout(),
		widget.NewLabel(`Some greedings.
		Thanks a lot to all the Impact members.
		Specials thanks for support to : 
		*** Tronic ***
		*** Siko ***`),
		layout.NewSpacer(),
	)
}
