package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) newEgxTab(di *menu.DoubleImageMenu) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Image 1", m.newImageTransfertTab(&di.LeftImage)),
		container.NewTabItem("Image 2", m.newImageTransfertTab(&di.RightImage)),
		container.NewTabItem("Egx", m.newEgxTabItem(di)),
	)
}

func (m *MartineUI) newEgxTabItem(di *menu.DoubleImageMenu) fyne.CanvasObject {
	di.ResultImage.CpcLeftImage = di.LeftImage.CpcImage
	di.ResultImage.CpcRightImage = di.RightImage.CpcImage
	di.ResultImage.LeftPalette = di.LeftImage.Palette
	di.ResultImage.RightPalette = di.RightImage.Palette
	return container.New(
		layout.NewGridLayoutWithRows(1),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			&di.ResultImage.CpcLeftImage,
			&di.RightImage.CpcImage,
			&di.ResultImage.CpcResultImage,
			container.New(
				layout.NewVBoxLayout(),
				widget.NewButtonWithIcon("Merge image", theme.MediaPlayIcon(), func() {

				}),
				widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
					di.ResultImage.CpcLeftImage = di.LeftImage.CpcImage
					di.ResultImage.CpcRightImage = di.RightImage.CpcImage
					di.ResultImage.LeftPalette = di.LeftImage.Palette
					di.ResultImage.RightPalette = di.RightImage.Palette
					s := m.window.Content().Size()
					s.Height += 10.
					s.Width += 10.
					m.window.Resize(s)
					m.window.Canvas().Refresh(&di.ResultImage.CpcLeftImage)
					m.window.Canvas().Refresh(&di.ResultImage.CpcRightImage)
					m.window.Canvas().Refresh(&di.ResultImage.CpcResultImage)

					m.window.Resize(m.window.Content().Size())
					m.window.Content().Refresh()
				}),
				widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {

				}),
			),
		),
	)
}
