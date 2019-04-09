package gfx

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"image"
	"image/color"
	"os"
)

var (
	ErrorMissingNumberOfImageToGenerate = errors.New("Iteration is not set, cannot define the number of images to generate.")
	ErrorSizeMismatch                   = errors.New("Error width and height mismatch cannot perform rotation.")
)

func Rotate(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filePath string, resizeAlgo imaging.ResampleFilter, exportType *ExportType) error {
	if exportType.RollIteration == -1 {
		return ErrorMissingNumberOfImageToGenerate
	}

	var indice int
	angle := 360. / float64(exportType.RollIteration)

	for i := 0.; i < 360.; i += angle {
		rin := imaging.Rotate(in, float64(i), color.Transparent)


		/*if rin.Bounds().Max.X%2 != 0 {
			newSize := constants.Size{
				ColorsAvailable: size.ColorsAvailable,
				ColumnsNumber:   size.ColumnsNumber,
				LinesNumber:     size.LinesNumber,
				Width:           rin.Bounds().Max.X + 1,
				Height:          rin.Bounds().Max.Y + 1,
			}
			rin = convert.Resize(rin, newSize, resizeAlgo)
		} else {

			if rin.Bounds().Max.Y%2 != 0 {
				newSize := constants.Size{
					ColorsAvailable: size.ColorsAvailable,
					ColumnsNumber:   size.ColumnsNumber,
					LinesNumber:     size.LinesNumber,
					Width:           rin.Bounds().Max.X + 1,
					Height:          rin.Bounds().Max.Y + 1,
				}
				rin = convert.Resize(rin, newSize, resizeAlgo)
			}
		}*/
		newSize := constants.Size{
			ColorsAvailable: size.ColorsAvailable,
			ColumnsNumber:   size.ColumnsNumber,
			LinesNumber:     size.LinesNumber,
			Width:           rin.Bounds().Max.X,
			Height:          rin.Bounds().Max.Y,
		}

		newFilename := exportType.OsFullPath(filePath, fmt.Sprintf("%.2d", indice)+".png")
		if err := Png(newFilename, rin); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create image (%s) error :%v\n", newFilename, err)
		}
		if err := SpriteTransform(rin, p, newSize, mode, newFilename, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create sprite image (%s) error %v\n", newFilename, err)
		}
		indice++
	}
	return nil
}
