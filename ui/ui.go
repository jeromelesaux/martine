package main

import (
	"os"

	"fyne.io/fyne/v2/app"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

func main() {
	os.Setenv("FYNE_SCALE", "1.0")
	/* main application */
	app := app.NewWithID("Martine @IMPact")
	martineUI := ui.NewMartineUI()
	martineUI.Load(app)
	app.Run()
}
