package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/log"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
)

var appPrefix = fmt.Sprintf("Martine (%v)", common.AppVersion)

func main() {
	_, err := log.InitLoggerWithFile(appPrefix)
	if err != nil {
		panic(err)
	}
	/* main application */
	mapp := app.NewWithID("Martine @IMPact")
	martineUI := ui.NewMartineUI()
	martineUI.Load(mapp)

	mapp.Run()
}
