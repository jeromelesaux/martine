package transformation

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/proc"
)

func Rotate(in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	mode uint8,
	rollIteration int,
	resizeAlgo imaging.ResampleFilter) ([]*image.NRGBA, error) {
	if rollIteration == -1 {
		return nil, errors.ErrorMissingNumberOfImageToGenerate
	}

	n := rollIteration
	angle := 360. / float64(n)
	angles := make([]float64, n)
	for i := range angles {
		angles[i] = float64(i) * angle
	}

	var mu sync.Mutex
	var maxSize constants.Size
	proc.Parallel(0, n, func(yc <-chan int) {
		for i := range yc {
			rin := imaging.Rotate(in, angles[i], color.Transparent)
			mu.Lock()
			if maxSize.Width < rin.Bounds().Max.X {
				maxSize.Width = rin.Bounds().Max.X
			}
			if maxSize.Height < rin.Bounds().Max.Y {
				maxSize.Height = rin.Bounds().Max.Y
			}
			mu.Unlock()
		}
	})

	background := image.NewRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: maxSize.Width, Y: maxSize.Height}})
	draw.Draw(background, background.Bounds(), &image.Uniform{p[0]}, image.Point{0, 0}, draw.Src)

	images := make([]*image.NRGBA, n)
	proc.Parallel(0, n, func(yc <-chan int) {
		for i := range yc {
			rin := imaging.Rotate(in, angles[i], p[0])
			if rin.Bounds().Max.X < maxSize.Width || rin.Bounds().Max.Y < maxSize.Height {
				rin = imaging.PasteCenter(background, rin)
			}
			_, rin = ci.DowngradingWithPalette(rin, p)
			images[i] = rin
		}
	})
	return images, nil
}
