package ui

import "fyne.io/fyne/v2/widget"

func NewWinFormatRadio(me *ImageMenu) *widget.RadioGroup {
	winFormat := widget.NewRadioGroup([]string{"Normal", "Fullscreen", "Sprite", "Sprite Hard"}, func(s string) {
		switch s {
		case "Normal":
			me.isFullScreen = false
			me.isSprite = false
			me.isHardSprite = false
		case "Fullscreen":
			me.isFullScreen = true
			me.isSprite = false
			me.isHardSprite = false
		case "Sprite":
			me.isFullScreen = false
			me.isSprite = true
			me.isHardSprite = false
		case "Sprite Hard":
			me.isFullScreen = false
			me.isSprite = false
			me.isHardSprite = true
		}
	})
	winFormat.SetSelected("Normal")
	return winFormat
}
