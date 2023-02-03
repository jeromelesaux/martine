package menu

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
)

var SpriteSize float32 = 80.

type SpriteExportFormat string

var (
	SpriteFlatExport  SpriteExportFormat = "Flat"
	SpriteFilesExport SpriteExportFormat = "Files"
	SpriteImpCatcher  SpriteExportFormat = "Impcatcher"
	SpriteCompiled    SpriteExportFormat = "Compiled"
	SpriteHard        SpriteExportFormat = "Sprite Hard"
)

type SpriteMenu struct {
	IsHardSprite    bool
	originalBoard   *canvas.Image
	originalPalette *canvas.Image
	palette         color.Palette
	paletteImage    *canvas.Image

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
	ExportFormat           SpriteExportFormat
	ExportDsk              bool
	ExportText             bool
	ExportWithAmsdosHeader bool
	ExportZigzag           bool
	ExportJson             bool
	ExportCompression      compression.CompressionMethod
	ExportFolderPath       string
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
		path, err := directory.DefaultDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		// d.Resize(dialogSize)
		d.Show()
	})
}
