package ui

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (m *MartineUI) newGreedings() *fyne.Container {
	return container.New(
		layout.NewVBoxLayout(),
		widget.NewRichTextFromMarkdown("# Some greetings.\n"+
			"\n## Thanks a lot to all the Impact members.\n"+
			"\n **Ast, CMP, Demoniak, Kris and Drill**\n"+
			"\n Specials thanks for support to :\n"+
			"\n - **AST**\n"+
			"\n - **Tronic**\n"+
			"\n - **Siko**\n"+
			"\n - **Roudoudou**\n"+
			"\n - **Hwikaa**\n"+
			"\n - and thanks a lot to all users^^\n"+
			"\nfor more informations about this tool, go to [Martine's Github page](https://github.com/jeromelesaux/martine)\n"+
			"\nfor more informations about my tool go to [github](https://github.com/jeromelesaux)\n"+
			"\nto follow me on my old website [sidhome](https://koaks.amstrad.free.fr/amstrad/)\n"+
			"\nto chat with us got to our [Impact forum](https://amstradplus.forumforever.com/index.php)  or\n"+
			"\n[discord chat](https://discord.com/channels/453480213032992768/454619697485447169)\n"),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Change  application size "),
			widget.NewButton("+", func() {
				current := os.Getenv("FYNE_SCALE")
				c, err := strconv.ParseFloat(current, 32)
				if err != nil {
					log.Default().Printf("error while getting FYNE_SCALE error [%v]", err)
				}
				os.Setenv("FYNE_SCALE", fmt.Sprintf("%f", c+0.1))
			}),
			widget.NewButton("-", func() {
				current := os.Getenv("FYNE_SCALE")
				c, err := strconv.ParseFloat(current, 32)
				if err != nil {
					log.Default().Printf("error while getting FYNE_SCALE error [%v]", err)
				}
				os.Setenv("FYNE_SCALE", fmt.Sprintf("%f", c-0.1))
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
