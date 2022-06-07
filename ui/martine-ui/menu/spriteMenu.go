package menu

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/jeromelesaux/fyne-io/custom_widget"
)

var SpriteSize float32 = 80.

type SpriteMenu struct {
	IsHardSprite    bool
	OriginalBoard   canvas.Image
	OriginalPalette canvas.Image

	Palette               color.Palette
	SpritesData           [][][]byte
	CompileSprite         bool
	IsCpcPlus             bool
	OriginalImages        *custom_widget.ImageTable
	SpritesCollection     [][]*image.NRGBA
	SpriteNumberPerRow    int
	SpriteNumberPerColumn int
	Mode                  int
	SpriteWidth           int
	SpriteHeight          int
}

func NewSpriteMenu() *SpriteMenu {
	return &SpriteMenu{
		OriginalBoard:     canvas.Image{},
		OriginalImages:    custom_widget.NewEmptyImageTable(fyne.NewSize(SpriteSize, SpriteSize)),
		SpritesCollection: make([][]*image.NRGBA, 0),
		SpritesData:       make([][][]byte, 0),
	}
}
