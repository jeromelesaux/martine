package main

import (
	"os"

	"fyne.io/fyne/v2/app"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

func main() {

	/* main application */
	os.Setenv("FYNE_SCALE", "0.7")
	mapp := app.NewWithID("Martine @IMPact")
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
