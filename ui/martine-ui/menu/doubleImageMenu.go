package menu

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/canvas"
	"github.com/jeromelesaux/martine/export"
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
	Path              string
	EgxType           int
}

func (m *MergedImageMenu) CmdLine() string {
	return ""
}

func (d *DoubleImageMenu) CmdLine() string {
	palFilename := export.AmsdosFilename(d.LeftImage.OriginalImagePath.Path(), ".PAL")
	scrFilename := export.AmsdosFilename(d.LeftImage.OriginalImagePath.Path(), ".SCR")

	cmd := "\n" + d.LeftImage.CmdLine() + " -out mode0"
	cmd += "\n" + d.RightImage.CmdLine() + " -pal " + palFilename + " -out mode1"
	exec, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while getting executable path :%v\n", err)
		return cmd
	}

	cmd += "\n " + exec + "  -in mode0" + string(filepath.Separator) + scrFilename + " -mode " + fmt.Sprintf("%d", d.LeftImage.Mode)
	cmd += " -in2 mode1" + string(filepath.Separator) + scrFilename + " -mode2 " + fmt.Sprintf("%d", d.RightImage.Mode)
	cmd += " -out egx"
	if d.ResultImage.EgxType == 1 {
		cmd += " -egx1"
	} else {
		cmd += " -egx2"
	}
	if d.LeftImage.IsFullScreen {
		cmd += " -fullscreen"
	}
	return cmd
}
