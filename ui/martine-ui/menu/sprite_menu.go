package menu

import (
	"image"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	w "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
)

var SpriteSize float32 = 80.

type SpriteMenu struct {
	Cfg             *config.MartineConfig
	originalBoard   *canvas.Image
	originalPalette *canvas.Image
	palette         color.Palette
	paletteImage    *canvas.Image
	FilePath        string

	SpritesData       [][][]byte
	CompileSprite     bool
	OriginalImages    *w.ImageTable
	SpritesCollection [][]*image.NRGBA
	SpriteColumns     int
	SpriteRows        int
	Mode              int
	UsePalette        bool

	CmdLineGenerate string
}

func (s *SpriteMenu) SetPalette(p color.Palette) {
	s.UsePalette = true
	s.palette = p
}

func (s *SpriteMenu) Palette() color.Palette {
	return s.palette
}

func (s *SpriteMenu) SetPaletteImage(img image.Image) {
	s.paletteImage.Image = img
	s.paletteImage.Refresh()
}

func (s *SpriteMenu) PaletteImage() *canvas.Image {
	return s.paletteImage
}

func (s *SpriteMenu) SetOriginalPalette(img image.Image) {
	s.originalPalette.Image = img
	s.originalPalette.Refresh()
}

func (s *SpriteMenu) SetOriginalBoard(img image.Image) {
	s.originalBoard.Image = img
	s.originalBoard.Refresh()
}

func (s *SpriteMenu) OriginalBoard() *canvas.Image {
	return s.originalBoard
}

func NewSpriteMenu() *SpriteMenu {
	return &SpriteMenu{
		originalBoard:     &canvas.Image{},
		OriginalImages:    w.NewEmptyImageTable(fyne.NewSize(SpriteSize, SpriteSize)),
		SpritesCollection: make([][]*image.NRGBA, 0),
		SpritesData:       make([][][]byte, 0),
		originalPalette:   &canvas.Image{},
		paletteImage:      &canvas.Image{},
		Cfg:               config.NewMartineConfig("", ""),
	}
}

func (s *SpriteMenu) ImportSprite(win fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("Import", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}
		}, win)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		// d.Resize(dialogSize)
		d.Show()
	})
}

func (s *SpriteMenu) CmdLine() string {
	exec, err := os.Executable()
	if err != nil {
		log.GetLogger().Error("error while getting executable path :%v\n", err)
		return exec
	}

	exec += " -in " + s.FilePath
	exec += " -split"
	exec += " -mode " + strconv.Itoa(s.Mode)
	exec += " -spritesrow " + strconv.Itoa(s.SpriteRows)
	exec += " -spritescolumn " + strconv.Itoa(s.SpriteColumns)
	if s.Cfg.ScreenCfg.Type.IsSpriteHard() {
		exec += " -height 16 -width 16"
	} else {
		exec += " -height " + strconv.Itoa(s.Cfg.ScreenCfg.Size.Height)
		exec += " -width " + strconv.Itoa(s.Cfg.ScreenCfg.Size.Width)
	}

	if s.Cfg.ScreenCfg.IsPlus {
		exec += " -plus"
	}

	if s.Cfg.ScreenCfg.Compression != compression.NONE {
		exec += " -z " + strconv.Itoa((int(s.Cfg.ScreenCfg.Compression)))
	}

	if s.Cfg.ContainerCfg.HasExport(config.DskContainer) {
		exec += " -dsk"
	}

	if s.Cfg.ScreenCfg.NoAmsdosHeader {
		exec += " -noheader"
	}

	if s.Cfg.ScreenCfg.IsExport(config.SpriteCompiledExport) {
		exec += " -compiled"
	}

	if s.Cfg.ScreenCfg.IsExport(config.OcpWindowExport) {
		exec += " -ocpwin"
	}

	if s.Cfg.ScreenCfg.IsExport(config.SpriteImpCatcherExport) {
		exec += " -imp"
	}

	if s.Cfg.ScreenCfg.IsExport(config.SpriteFlatExport) {
		exec += " -flat"
	}

	if s.Cfg.ScreenCfg.IsExport(config.SpriteHardExport) {
		exec += " -spritehard"
	}

	s.CmdLineGenerate = exec
	return exec
}
