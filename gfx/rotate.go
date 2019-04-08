package gfx

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"image"
	"image/color"
	"os"
	"strconv"
)

var (
	ErrorMissingNumberOfImageToGenerate = errors.New("Iteration is not set, cannot define the number of images to generate.")
	ErrorSizeMismatch                   = errors.New("Error width and height mismatch cannot perform rotation.")
)

func Rotate(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filePath string, resizeAlgo imaging.ResampleFilter,exportType *ExportType) error {
	if exportType.RollIteration == -1 {
		return ErrorMissingNumberOfImageToGenerate
	}
	if size.Width != size.Height {
		return ErrorSizeMismatch
	}
	var indice int
	angle := 360 / exportType.RollIteration
	for i := 0; i < 360; i += angle {
		rin := imaging.Rotate(in, float64(i), color.Black)
		rin = convert.Resize(rin,size,resizeAlgo)
		newFilename := exportType.OsFullPath(filePath,strconv.Itoa(indice)+ ".png")
		if err := Png(newFilename, rin); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create image (%s) error :%v\n", newFilename, err)
		}
		if err := SpriteTransform(rin, p, size, mode, newFilename, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create sprite image (%s) error %v\n", newFilename, err)
		}
		indice++
	}
	return nil
}
