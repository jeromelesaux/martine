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
)

var SpriteSize float32 = 80.

type SpriteExportFormat string

var (
	SpriteFlatExport  SpriteExportFormat = "Flat"
	SpriteFilesExport SpriteExportFormat = "Files"
	SpriteImpCatcher  SpriteExportFormat = "Impcatcher"
	SpriteCompiled    SpriteExportFormat = "Compiled"
)

type SpriteMenu struct {
	IsHardSprite    bool
	OriginalBoard   canvas.Image
	OriginalPalette canvas.Image
	Palette         color.Palette
	PaletteImage    canvas.Image

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
	ExportCompression      int
	ExportFolderPath       string
}

func (s *SpriteMenu) SetPalette(p color.Palette) {
	s.Palette = p
}

func (s *SpriteMenu) SetPaletteImage(c canvas.Image) {
	s.PaletteImage = c
}
func NewSpriteMenu() *SpriteMenu {
	return &SpriteMenu{
		OriginalBoard:     canvas.Image{},
		OriginalImages:    custom_widget.NewEmptyImageTable(fyne.NewSize(SpriteSize, SpriteSize)),
		SpritesCollection: make([][]*image.NRGBA, 0),
		SpritesData:       make([][][]byte, 0),
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
			//filename := reader.URI()

		}, win)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		//d.Resize(dialogSize)
		d.Show()
	})
}
