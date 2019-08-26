package gfx

import (
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"image"
	"strings"
)

func Flash(in image.Image,
	resizeAlgo imaging.ResampleFilter,
	exportType *export.ExportType,
	filename, picturePath, palettePath, inkPath, kitPath string,
	mode, ditheringAlgo, rla, sla, rra, sra, keephigh, keeplow, losthigh, lostlow, iterations int,
	screenMode uint8,
	ditheringMultiplier float64,
	ditheringMatrix [][]float32,
	ditherType DitheringType,
	customDimension, withQuantization bool) error {

	var err error
	size := constants.Size{
		Width:  exportType.Size.Width * 2,
		Height: exportType.Size.Height * 2}
	im := convert.Resize(in, size, resizeAlgo)
	leftIm := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{exportType.Size.Width, exportType.Size.Height}})
	rigthIm := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{exportType.Size.Width, exportType.Size.Height}})
	filenameLeft := strings.Replace(filename, ".", "1.", 1)
	filenameRigth := strings.Replace(filename, ".", "2.", 1)
	filepathLeft := strings.Replace(picturePath, ".", "1.", 1)
	filepathRigth := strings.Replace(picturePath, ".", "2.", 1)
	x := 0
	for i := 0; i < size.Width; i += 2 {
		y := 0
		for j := 0; j < size.Height; j += 2 {
			leftIm.Set(x, y, im.At(i, j))
			y++
		}
		x++
		y = 0
	}
	err = ApplyOneImage(leftIm,
		resizeAlgo,
		exportType,
		filenameLeft, filepathLeft, palettePath, inkPath, kitPath,
		mode, ditheringAlgo, rla, sla, rra, sra, keephigh, keeplow, losthigh, lostlow, iterations,
		screenMode,
		ditheringMultiplier,
		ditheringMatrix,
		ditherType,
		customDimension, withQuantization)
	if err != nil {
		return err
	}

	x = 0
	for i := 1; i < size.Width; i += 2 {
		y := 0
		for j := 1; j < size.Height; j += 2 {
			rigthIm.Set(x, y, im.At(i, j))
			y++
		}
		x++
		y = 0
	}

	err = ApplyOneImage(rigthIm,
		resizeAlgo,
		exportType,
		filenameRigth, filepathRigth, palettePath, inkPath, kitPath,
		mode, ditheringAlgo, rla, sla, rra, sra, keephigh, keeplow, losthigh, lostlow, iterations,
		screenMode,
		ditheringMultiplier,
		ditheringMatrix,
		ditherType,
		customDimension, withQuantization)
	return err

}
