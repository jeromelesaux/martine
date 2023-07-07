package image

import (
	"image"
	"image/color"

	"github.com/mccutchen/palettor"
)

func Kmeans(nbColors, iteration int, img image.Image) (*image.NRGBA, error) {
	newPalette, err := palettor.Extract(nbColors, iteration, img)
	if err != nil {
		return &image.NRGBA{}, err
	}
	var p color.Palette

	for _, c := range newPalette.Colors() {
		p = append(p, c)
	}

	newImg := image.NewNRGBA(img.Bounds())
	for x := 0; x <= img.Bounds().Max.X; x++ {
		for y := 0; y <= img.Bounds().Max.Y; y++ {
			c := img.At(x, y)
			nc := p.Convert(c)
			newImg.Set(x, y, nc)
		}
	}
	return newImg, nil
}
