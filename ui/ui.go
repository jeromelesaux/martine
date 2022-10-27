package main

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

type MartineTheme struct {
}

var _ fyne.Theme = (*MartineTheme)(nil)

func (m MartineTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m MartineTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m MartineTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) * 2
}
func (m MartineTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func main() {
	//os.Setenv("FYNE_SCALE", "3.0")
	os.Setenv("FYNE_THEME", "light")
	/* main application */
	mapp := app.NewWithID("Martine @IMPact")
	mapp.Settings().SetTheme(&MartineTheme{})
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
