package main

import (
	"fyne.io/fyne/v2/app"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)



func main() {

	/* main application */
	mapp := app.NewWithID("Martine @IMPact")
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
