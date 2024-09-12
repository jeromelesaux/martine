package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// nolint: ireturn, funlen
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
		widget.NewRichTextFromMarkdown("# Shorcuts: #\n"+
			"\n## In editor : ## \n"+
			"\nArrow UP : move up the cursor by one pixel\n"+
			"\nArrow DOWN : move down the cursor by one pixel\n"+
			"\nArrow LEFT : move left the cursor by one pixel\n"+
			"\nArrow RIGHT : move right the cursor by one pixel\n"+
			"\nA or Q key : move up the cursor by 10 pixel\n"+
			"\nW key : move down the cursor by 10 pixel\n"+
			"\nO key : move left the cursor by 10 pixel\n"+
			"\nP key : move right the cursor by 10 pixel\n"+
			"\nM key : change the magnify\n"+
			"\nESCAPE key : undo the last color change from palette\n"+
			"\nSPACE key : set the current pixel with the selected palette color \n"),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Change  application size "),
			widget.NewButton("Increase", func() {
				m.scale += .03
				theme := myTheme{foregroundColor: m.variant, scale: m.scale}
				fyne.CurrentApp().Settings().SetTheme(theme)
				canvas.Refresh(m.window.Canvas().Content())

			}),
			widget.NewButton("Decrease", func() {

				m.scale -= .03
				theme := myTheme{foregroundColor: m.variant, scale: m.scale}
				fyne.CurrentApp().Settings().SetTheme(theme)
				canvas.Refresh(m.window.Canvas().Content())

			}),
			widget.NewButton("Change foreground color", func() {

				if m.variant == theme.VariantDark {
					m.variant = theme.VariantLight
				} else {
					m.variant = theme.VariantDark
				}
				theme := myTheme{foregroundColor: m.variant, scale: fyne.CurrentApp().Settings().Scale()}
				fyne.CurrentApp().Settings().SetTheme(theme)
				canvas.Refresh(m.window.Canvas().Content())

			}),
		),
	)
}

type myTheme struct {
	foregroundColor fyne.ThemeVariant
	scale           float32
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if m.foregroundColor == theme.VariantLight {
			return color.White
		}
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}

	return theme.DefaultTheme().Color(name, m.foregroundColor)
}

// nolint: ireturn
func (myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m myTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameCaptionText:
		return 11 * m.scale
	case theme.SizeNameInlineIcon:
		return 20 * m.scale
	case theme.SizeNamePadding:
		return 4 * m.scale
	case theme.SizeNameScrollBar:
		return 16 * m.scale
	case theme.SizeNameScrollBarSmall:
		return 3 * m.scale
	case theme.SizeNameSeparatorThickness:
		return 1 * m.scale
	case theme.SizeNameText:
		return 14 * m.scale
	case theme.SizeNameInputBorder:
		return 2 * m.scale
	default:
		return theme.DefaultTheme().Size(s) * m.scale
	}
}

// nolint: ireturn
func (myTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return theme.DefaultTheme().Font(s)
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return theme.DefaultTheme().Font(s)
}
