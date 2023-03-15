package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/jeromelesaux/martine/log"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

type MartineTheme struct{}

func main() {
	_, err := log.InitLoggerWithFile()
	if err != nil {
		panic(err)
	}
	/* main application */
	mapp := app.NewWithID("Martine @IMPact")
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
