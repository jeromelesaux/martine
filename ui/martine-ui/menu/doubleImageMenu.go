package menu

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

type DoubleImageMenu struct {
	LeftImage   ImageMenu
	RightImage  ImageMenu
	ResultImage MergedImageMenu
}

type MergedImageMenu struct {
	CpcLeftImage      canvas.Image
	CpcRightImage     canvas.Image
	CpcResultImage    canvas.Image
	LeftPalette       color.Palette
	RightPalette      color.Palette
	LeftPaletteImage  canvas.Image
	RightPaletteImage canvas.Image
	Data              []byte
	Palette           color.Palette
	PaletteImage      canvas.Image
	CmdLineGenerate   string
}

func (m *MergedImageMenu) CmdLine() string {
	return ""
}

func (d *DoubleImageMenu) CmdLine() string {
	cmd := d.LeftImage.CmdLine()
	cmd += "\n" + d.RightImage.CmdLine()
	cmd += "\n" + d.ResultImage.CmdLine()
	return cmd
}
