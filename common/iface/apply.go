package iface

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
)

type ImageAction func(cfg ImageIface) error

type ImageIface interface {
	ImagePath() string
	Img() *image.NRGBA
	SetImg(i *image.NRGBA)
	Palette() color.Palette
	SetPalette(p color.Palette)
	Mode() uint8
	Brightness() float64
	Saturation() float64
	Size() constants.Size
	CpcPlus() bool
	Reducer() int
	ResizingAlgo() imaging.ResampleFilter
	KmeansThreshold() float64
}
