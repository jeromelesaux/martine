package transformation

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func Rotate(in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	mode uint8,
	rollIteration int,
	resizeAlgo imaging.ResampleFilter) ([]*image.NRGBA, error) {
	images := make([]*image.NRGBA, 0)
	if rollIteration == -1 {
		return images, errors.ErrorMissingNumberOfImageToGenerate
	}

	var indice int
	angle := 360. / float64(rollIteration)
	var maxSize constants.Size
	for i := 0.; i < 360.; i += angle {
		rin := imaging.Rotate(in, float64(i), color.Transparent)
		if maxSize.Width < rin.Bounds().Max.X {
			maxSize.Width = rin.Bounds().Max.X
		}
		if maxSize.Height < rin.Bounds().Max.Y {
			maxSize.Height = rin.Bounds().Max.Y
		}
	}
	background := image.NewRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: maxSize.Width, Y: maxSize.Height}})
	draw.Draw(background, background.Bounds(), &image.Uniform{p[0]}, image.Point{0, 0}, draw.Src)

	for i := 0.; i < 360.; i += angle {
		rin := imaging.Rotate(in, float64(i), p[0])

		if rin.Bounds().Max.X < maxSize.Width || rin.Bounds().Max.Y < maxSize.Height {
			rin = imaging.PasteCenter(
				background,
				rin,
			)
		}
		_, rin = convert.DowngradingWithPalette(rin, p)
		images = append(images, rin)
		/*	newFilename := cont.OsFullPath(filePath, fmt.Sprintf("%.2d", indice)+".png")
			if err := file.Png(newFilename, rin); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create image (%s) error :%v\n", newFilename, err)
			}
			if err := common.ToSpriteAndExport(rin, p, maxSize, mode, newFilename, false, cont); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create sprite image (%s) error %v\n", newFilename, err)
			}*/
		indice++
	}
	return images, nil
}
