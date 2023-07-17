package image

import (
	"image"
	"image/color"

	"github.com/Baldomo/paletter"
)

func Kmeans(nbColors int, threshold float64, img image.Image) (*image.NRGBA, error) {

	if threshold != 0. {
		paletter.DeltaThreshold = threshold
	}
	obs := paletter.ImageToObservation(img)
	cs, err := paletter.CalculatePalette(obs, nbColors)
	if err != nil {
		return nil, err
	}
	colors := paletter.ColorsFromClusters(cs)

	var p color.Palette

	for _, c := range colors {
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
