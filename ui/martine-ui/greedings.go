package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (m *MartineUI) newGreedings() fyne.CanvasObject {
	return container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel(`Some greedings.
		Thanks a lot to all the Impact members.
		Ast, CMP, Demoniak, Kris and Drill
		Specials thanks for support to :
		***        AST        ***
		***      Tronic        ***
		***        Siko          ***
		*** Roudoudou ***
		and thanks a lot to all users^^
		for more informations about this tool, go to https://github.com/jeromelesaux/martine
		for more informations about my tool go to https://github.com/jeromelesaux
		to follow me on my old website https://http://koaks.amstrad.free.fr/amstrad/
		to chat with us got to https://amstradplus.forumforever.com/index.php  or
		https://discord.com/channels/453480213032992768/454619697485447169 on discord
		`),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Change color scheme"),
			widget.NewSelect([]string{"Black", "White"}, func(s string) {
				a := fyne.CurrentApp()
				switch s {
				case "Black":
					a.Settings().SetTheme(theme.DarkTheme())
				case "White":
					a.Settings().SetTheme(theme.LightTheme())
				}
			}),
		),
	/*	container.New(
		layout.NewHBoxLayout(),
		widget.NewLabel("Font"),
		widget.NewButtonWithIcon("Increase", theme.ContentAddIcon(), func() {
			a := fyne.CurrentApp()
			v := a.Settings().Scale()
			v += .2
			os.Setenv("FYNE_SCALE", fmt.Sprintf("%.2f", v))
		}),
		widget.NewButtonWithIcon("Decrease", theme.ContentRemoveIcon(), func() {
			a := fyne.CurrentApp()
			v := a.Settings().Scale()
			v -= .2
			os.Setenv("FYNE_SCALE", fmt.Sprintf("%.2f", v))
		}),
	),*/
	)
}
