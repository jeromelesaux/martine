package menu

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/canvas"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"
)

type DoubleImageMenu struct {
	LeftImage   *ImageMenu
	RightImage  *ImageMenu
	ResultImage *MergedImageMenu
	Cfg         *config.MartineConfig
}

func NewDoubleImageMenu() *DoubleImageMenu {
	return &DoubleImageMenu{
		LeftImage:   NewImageMenu(),
		RightImage:  NewImageMenu(),
		ResultImage: NewMergedImageMenu(),
		Cfg:         config.NewMartineConfig("", ""),
	}
}

type MergedImageMenu struct {
	CpcLeftImage      *canvas.Image
	CpcRightImage     *canvas.Image
	CpcResultImage    *canvas.Image
	LeftPalette       color.Palette
	RightPalette      color.Palette
	LeftPaletteImage  *canvas.Image
	RightPaletteImage *canvas.Image
	Data              []byte
	Palette           color.Palette
	PaletteImage      *canvas.Image
	CmdLineGenerate   string
	Path              string
	EgxType           int
}

func NewMergedImageMenu() *MergedImageMenu {
	return &MergedImageMenu{
		CpcLeftImage:      &canvas.Image{},
		CpcRightImage:     &canvas.Image{},
		CpcResultImage:    &canvas.Image{},
		LeftPaletteImage:  &canvas.Image{},
		RightPaletteImage: &canvas.Image{},
		PaletteImage:      &canvas.Image{},
	}
}

func (m *MergedImageMenu) CmdLine() string {
	return ""
}

func (d *DoubleImageMenu) CmdLine() string {
	palFilename := config.AmsdosFilename(d.LeftImage.OriginalImagePath(), ".PAL")
	scrFilename := config.AmsdosFilename(d.LeftImage.OriginalImagePath(), ".SCR")

	cmd := "\n" + d.LeftImage.CmdLine() + " -out mode0"
	cmd += "\n" + d.RightImage.CmdLine() + " -pal " + palFilename + " -out mode1"
	exec, err := os.Executable()
	if err != nil {
		log.GetLogger().Error("error while getting executable path :%v\n", err)
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
	if d.LeftImage.Cfg.ScrCfg.Type.IsFullScreen() {
		cmd += " -fullscreen"
	}
	return cmd
}
