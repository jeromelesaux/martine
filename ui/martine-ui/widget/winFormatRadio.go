package widget

import (
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func NewWinFormatRadio(me *menu.ImageMenu) *widget.RadioGroup {
	winFormat := widget.NewRadioGroup([]string{"Normal", "Fullscreen", "Sprite", "Sprite Hard"}, func(s string) {
		switch s {
		case "Normal":
			me.IsFullScreen = false
			me.IsSprite = false
			me.IsHardSprite = false
		case "Fullscreen":
			me.IsFullScreen = true
			me.IsSprite = false
			me.IsHardSprite = false
		case "Sprite":
			me.IsFullScreen = false
			me.IsSprite = true
			me.IsHardSprite = false
		case "Sprite Hard":
			me.IsFullScreen = false
			me.IsSprite = false
			me.IsHardSprite = true
		}
	})
	winFormat.SetSelected("Normal")
	return winFormat
}
