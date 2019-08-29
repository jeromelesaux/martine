package gfx

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
)

func Rotate3d(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filePath string, resizeAlgo imaging.ResampleFilter, exportType *x.ExportType) error {
	if exportType.RollIteration == -1 {
		return ErrorMissingNumberOfImageToGenerate
	}

	var indice int
	angle := 360. / float64(exportType.RollIteration)
	//targetSize := in.Bounds().Max.X
	//if in.Bounds().Max.Y > in.Bounds().Max.X {
	//	targetSize = in.Bounds().Max.Y
	//}

	for i := 0.; i < 360.; i += angle {
		background := image.NewNRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: size.Width, Y: size.Height}})
		draw.Draw(background, background.Bounds(), &image.Uniform{p[0]}, image.ZP, draw.Src)
		rin := rotateImage(in, background, i, exportType)
		_, rin = convert.DowngradingWithPalette(rin, p)

		newFilename := exportType.OsFullPath(filePath, fmt.Sprintf("%.2d", indice)+".png")
		if err := file.Png(newFilename, rin); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create image (%s) error :%v\n", newFilename, err)
		}
		if err := SpriteTransform(rin, p, constants.Size{Width: size.Width, Height: size.Height}, mode, newFilename, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create sprite image (%s) error %v\n", newFilename, err)
		}
		indice++
	}

	return nil
}

func rotateImage(in, out *image.NRGBA, angle float64, exportType *x.ExportType) *image.NRGBA {
	var xc, yc int
	if exportType.Rotation3DX0 != -1 {
		xc = exportType.Rotation3DX0
	} else {
		xc = out.Bounds().Max.X / 2
	}
	if exportType.Rotation3DY0 != -1 {
		yc = exportType.Rotation3DY0
	} else {
		yc = out.Bounds().Max.Y / 2
	}
	for x := 0; x < in.Bounds().Max.X; x++ {
		for y := 0; y < in.Bounds().Max.Y; y++ {
			c := in.At(x, y)
			var x3d, y3d int
			switch exportType.Rotation3DType {
			case 1:
				x3d, y3d = rotateXAxisCoordinates(x, y, xc, yc, angle)
			case 2:
				x3d, y3d = rotateYAxisCoordinates(x, y, xc, yc, angle)
			case 3:
				x3d, y3d = rotateToReverseXAxisCoordinates(x, y, xc, yc, angle)
			case 4:
				x3d, y3d = rotateLeftToRightYAxisCoordinates(x, y, xc, yc, angle)
			case 5:
				x3d, y3d = rotateDiagonalXAxisCoordinates(x, y, xc, yc, angle)
			case 6:
				x3d, y3d = rotateDiagonalYAxisCoordinates(x, y, xc, yc, angle)
			default:
				x3d, y3d = rotateXAxisCoordinates(x, y, xc, yc, angle)
			}

			out.Set(x3d, y3d, c)
		}
	}
	return out
}

func rotateCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	x3d := (float64(x-xc) * cosTheta) - (float64(y-yc) * sinTheta) + float64(xc)
	y3d := (float64(y-yc) * cosTheta) + (float64(x-xc) * sinTheta) + float64(yc)
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}

// source : https://slideplayer.com/slide/9723655/
func rotateYAxisCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	cosTheta := math.Cos(theta)
	x3d := (float64(x-xc) * cosTheta) + float64(xc)
	y3d := (float64(y))
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}

func rotateXAxisCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	cosTheta := math.Cos(theta)
	x3d := (float64(x))
	y3d := (float64(y-yc) * cosTheta) + float64(yc)
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}

func rotateLeftToRightYAxisCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	sinTheta := math.Sin(theta)
	x3d := (float64(x-xc) * sinTheta) + float64(xc)
	y3d := (float64(y))
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}

func rotateToReverseXAxisCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	sinTheta := math.Sin(theta)
	x3d := (float64(x))
	y3d := (float64(y-yc) * sinTheta) + float64(yc)
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}

func rotateDiagonalXAxisCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	x3d := (float64(x-xc) * cosTheta) - (float64(y-yc) * sinTheta) + float64(xc)
	y3d := (float64(y))
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}

func rotateDiagonalYAxisCoordinates(x, y, xc, yc int, angle float64) (int, int) {
	theta := angle * math.Pi / 180.
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	x3d := (float64(x))
	y3d := (float64(y-yc) * cosTheta) + (float64(x-xc) * sinTheta) + float64(yc)
	return int(math.Floor(x3d)), int(math.Floor(y3d))
}
