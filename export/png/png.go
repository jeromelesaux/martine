package png

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/jeromelesaux/martine/log"
)

func Png(filePath string, im *image.NRGBA) error {
	fwd, err := os.Create(filePath)
	if err != nil {
		log.GetLogger().Error("Cannot create new image (%s) error %v\n", filePath, err)
		return err
	}

	if err := png.Encode(fwd, im); err != nil {
		fwd.Close()
		log.GetLogger().Error("Cannot create new image (%s) as png error %v\n", filePath, err)
		return err
	}
	fwd.Close()
	log.GetLogger().Info("Create output file (%s)\n", filePath)
	return nil
}

func PalToPng(filePath string, palette color.Palette) error {
	colorWidth := 20

	im := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: (16*5 + 5 + (colorWidth * 16)), Y: 30},
	})

	for i := 0; i < len(palette); i++ {
		if i >= 16 {
			break
		}
		contour := image.Rectangle{
			Min: image.Point{X: 5 + (i*colorWidth + i*5), Y: 5},
			Max: image.Point{X: colorWidth + 5 + (i*colorWidth + i*5), Y: colorWidth + 5},
		}
		draw.Draw(im, contour, &image.Uniform{palette[i]}, image.Point{0, 0}, draw.Src)
	}

	return Png(filePath, im)
}

func PalToImage(palette color.Palette) *image.NRGBA {
	colorWidth := 20

	im := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: (16*5 + 5 + (colorWidth * 16)), Y: 30},
	})

	for i := 0; i < len(palette); i++ {
		if i >= 16 {
			break
		}
		contour := image.Rectangle{
			Min: image.Point{X: 5 + (i*colorWidth + i*5), Y: 5},
			Max: image.Point{X: colorWidth + 5 + (i*colorWidth + i*5), Y: colorWidth + 5},
		}
		draw.Draw(im, contour, &image.Uniform{palette[i]}, image.Point{0, 0}, draw.Src)
	}
	return im
}
