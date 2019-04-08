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
)

var (
	ErrorMissingNumberOfImageToGenerate = errors.New("Iteration is not set, cannot define the number of images to generate.")
	ErrorSizeMismatch                   = errors.New("Error width and height mismatch cannot perform rotation.")
)

func Rotate(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filePath string, resizeAlgo imaging.ResampleFilter,exportType *ExportType) error {
	if exportType.RollIteration == -1 {
		return ErrorMissingNumberOfImageToGenerate
	}
	maxSize := constants.Size{}
	if size.Width != size.Height {
		return ErrorSizeMismatch
	}
	var indice int
	angle := 360. / float64(exportType.RollIteration)
	for i := 0. ; i < 360.; i += angle {
		rin := imaging.Rotate(in, float64(i), color.Transparent)
		if maxSize.Width < rin.Bounds().Max.X || maxSize.Height < rin.Bounds().Max.Y {
			maxSize.Width = rin.Bounds().Max.X 
			maxSize.Height = rin.Bounds().Max.Y
		}
	}

	fmt.Fprintf(os.Stdout,"initiale size (%d,%d) maxsize (%d,%d)\n",size.Width,size.Height,maxSize.Width,maxSize.Height )

	for i := 0.; i < 360.; i += angle {
		rin := imaging.Rotate(in, float64(i), color.Transparent)
		rin = convert.Resize(rin,size,resizeAlgo)
		newFilename := exportType.OsFullPath(filePath,fmt.Sprintf("%.2d",indice)+ ".png")
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
