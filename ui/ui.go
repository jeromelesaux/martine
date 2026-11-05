package main

import (
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

var frPatch = []byte(`{"(Select one)": "(Sélectionner)"}`)

func main() {
	scale := os.Getenv("FYNE_SCALE")
	if scale == "" {
		_ = os.Setenv("FYNE_SCALE", "0.7")
	}
	mapp := app.NewWithID("Martine @IMPact")
	_ = lang.AddTranslationsForLocale(frPatch, "fr")
	// nolint: staticcheck
	mapp.Settings().SetTheme(theme.DarkTheme())
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
