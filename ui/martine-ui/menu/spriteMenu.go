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
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
)

var SpriteSize float32 = 80.

type SpriteMenu struct {
	IsHardSprite    bool
	originalBoard   *canvas.Image
	originalPalette *canvas.Image
	palette         color.Palette
	paletteImage    *canvas.Image
	FilePath        string

	SpritesData            [][][]byte
	CompileSprite          bool
	IsCpcPlus              bool
	OriginalImages         *custom_widget.ImageTable
	SpritesCollection      [][]*image.NRGBA
	SpriteColumns          int
	SpriteRows             int
	Mode                   int
	SpriteWidth            int
	SpriteHeight           int
	ExportFormat           export.ExportFormat
	ExportDsk              bool
	ExportText             bool
	ExportWithAmsdosHeader bool
	ExportZigzag           bool
	ExportJson             bool
	ExportCompression      compression.CompressionMethod
	ExportFolderPath       string

	CmdLineGenerate string
}

func (s *SpriteMenu) SetPalette(p color.Palette) {
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
		OriginalImages:    custom_widget.NewEmptyImageTable(fyne.NewSize(SpriteSize, SpriteSize)),
		SpritesCollection: make([][]*image.NRGBA, 0),
		SpritesData:       make([][][]byte, 0),
		originalPalette:   &canvas.Image{},
		paletteImage:      &canvas.Image{},
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
	if s.IsHardSprite {
		exec += " -height 16 -width 16"
	} else {
		exec += " -height " + strconv.Itoa(s.SpriteHeight)
		exec += " -width " + strconv.Itoa(s.SpriteWidth)
	}

	if s.IsCpcPlus {
		exec += " -plus"
	}

	if s.ExportCompression != 0 {
		exec += " -z " + strconv.Itoa((int(s.ExportCompression)))
	}

	if s.ExportDsk {
		exec += " -dsk"
	}

	if s.ExportWithAmsdosHeader {
		exec += " -noheader"
	}

	switch s.ExportFormat {
	case export.SpriteCompiled:
		exec += " -compiled"
	case export.OcpWinExport:
		exec += " -ocpwin"
	case export.SpriteImpCatcher:
		exec += " -imp"
	case export.SpriteFlatExport:
		exec += " -flat"
	case export.SpriteHard:
		exec += " -spritehard"
	}

	s.CmdLineGenerate = exec
	return exec
}
