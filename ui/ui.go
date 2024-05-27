package main

import (
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

func main() {

	/* main application */
	scale := os.Getenv("FYNE_SCALE")
	if scale == "" {
		os.Setenv("FYNE_SCALE", "0.75")
	}
	mapp := app.NewWithID("Martine @IMPact")
	mapp.Settings().SetTheme(theme.DarkTheme()) // nolint:
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
