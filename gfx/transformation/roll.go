package transformation

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/log"
)

// roll is the generic implementation shared by all Roll* functions.
// copyMain copies the shifted pixels into im for each iteration.
// copyWrap copies the wrapped-around pixels when keep mode is active (may be nil for lost mode).
func roll(
	keep, lost, iterations int,
	label string,
	downgraded *image.NRGBA,
	copyMain func(im *image.NRGBA, nbPixels int),
	copyWrap func(im *image.NRGBA, nbPixels int),
) []*image.NRGBA {
	images := make([]*image.NRGBA, 0)
	if keep == -1 && lost == -1 {
		return images
	}
	log.GetLogger().Info("%s: Iterations (%d)\n", label, iterations)
	for i := 0; i < iterations; i++ {
		nbPixels := 0
		if keep != -1 {
			nbPixels = keep * (1 + i)
		} else {
			nbPixels = lost * (1 + i)
		}
		im := image.NewNRGBA(image.Rectangle{
			image.Point{0, 0},
			image.Point{downgraded.Bounds().Max.X, downgraded.Bounds().Max.Y},
		})
		copyMain(im, nbPixels)
		if keep != -1 && copyWrap != nil {
			copyWrap(im, nbPixels)
		}
		images = append(images, im)
	}
	return images
}

func RollLeft(rla, sla, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {

	return roll(rla, sla, iterations, "RLA/SLA", downgraded,
		func(im *image.NRGBA, nbPixels int) {
			y2 := 0
			for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
				x2 := 0
				for x := downgraded.Bounds().Min.X + nbPixels; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
		},
		func(im *image.NRGBA, nbPixels int) {
			y2 := 0
			for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
				x2 := downgraded.Bounds().Max.X - nbPixels
				for x := downgraded.Bounds().Min.X; x < nbPixels; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
		},
	)
}

func RollRight(rra, sra, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {

	return roll(rra, sra, iterations, "RRA/SRA", downgraded,
		func(im *image.NRGBA, nbPixels int) {
			y2 := 0
			for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
				x2 := nbPixels
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X-nbPixels; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
		},
		func(im *image.NRGBA, nbPixels int) {
			y2 := 0
			for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
				x2 := 0
				for x := downgraded.Bounds().Max.X - nbPixels; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
		},
	)
}

func RollUp(keephigh, losthigh, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {

	return roll(keephigh, losthigh, iterations, "keephigh/losthigh", downgraded,
		func(im *image.NRGBA, nbPixels int) {
			y2 := 0
			for y := downgraded.Bounds().Min.Y + nbPixels; y < downgraded.Bounds().Max.Y; y++ {
				x2 := 0
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
		},
		func(im *image.NRGBA, nbPixels int) {
			y2 := downgraded.Bounds().Max.Y - nbPixels
			for y := downgraded.Bounds().Min.Y; y < nbPixels; y++ {
				x2 := 0
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
		},
	)
}

func RollLow(keeplow, lostlow, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {

	return roll(keeplow, lostlow, iterations, "keeplow/lostlow", downgraded,
		func(im *image.NRGBA, nbPixels int) {
			y2 := downgraded.Bounds().Max.Y - 1
			for y := downgraded.Bounds().Max.Y - nbPixels; y >= downgraded.Bounds().Min.Y; y-- {
				x2 := 0
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2--
			}
		},
		func(im *image.NRGBA, nbPixels int) {
			y2 := nbPixels - 1
			for y := downgraded.Bounds().Max.Y - 1; y >= downgraded.Bounds().Max.Y-nbPixels; y-- {
				x2 := 0
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2--
			}
		},
	)
}
