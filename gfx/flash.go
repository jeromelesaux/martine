package gfx

import (
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"image"
	"path"
	"path/filepath"
	"strings"
)

func Flash(in image.Image,
	exportType *export.ExportType,
	filename, picturePath string,
	mode int,
	screenMode uint8) error {

	var err error
	size := constants.Size{
		Width:  exportType.Size.Width * 2,
		Height: exportType.Size.Height * 2}
	im := convert.Resize(in, size, exportType.ResizingAlgo)
	leftIm := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{exportType.Size.Width, exportType.Size.Height}})
	rigthIm := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{exportType.Size.Width, exportType.Size.Height}})
	indexExtFilename := strings.LastIndex(filename, ".")
	indexExtPath := strings.LastIndex(picturePath, ".")

	bFilename := make([]byte, indexExtFilename)
	bPath := make([]byte, indexExtPath)
	formerExt := filename[indexExtFilename:len(filename)]
	copy(bFilename, filename[0:indexExtFilename])
	copy(bPath, picturePath[0:indexExtPath])

	filenameLeft := string(bFilename) + "1" + formerExt
	filenameRigth := string(bFilename) + "2" + formerExt
	filepathLeft := string(bPath) + "1" + formerExt
	filepathRigth := string(bPath) + "2" + formerExt
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
	ext := path.Ext(filename)
	name := strings.Replace(filename, ext, "", 1)
	namesize := len(name)
	if namesize > 8 {
		namesize = 7
	}
	flashPaletteFilename1 := strings.ToUpper(name)[0:namesize] + "1.PAL"
	flashPalettePath1 := exportType.OutputPath + string(filepath.Separator) + flashPaletteFilename1

	err = ApplyOneImage(leftIm,
		exportType,
		filenameLeft, filepathLeft,
		mode,
		screenMode)
	if err != nil {
		return err
	}
	p1, _, err := file.OpenPal(flashPalettePath1)
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

	flashMode := mode - 1
	if flashMode < 0 {
		flashMode = 2
	}

	switch flashMode {
	case 0:
		exportType.Size = constants.Mode0
	case 1:
		exportType.Size = constants.Mode1
	case 2:
		exportType.Size = constants.Mode2
	}

	flashPaletteFilename2 := strings.ToUpper(name)[0:namesize] + "2.PAL"
	flashPalettePath2 := exportType.OutputPath + string(filepath.Separator) + flashPaletteFilename2

	//exportType.PalettePath = flashPalettePath1

	err = ApplyOneImage(rigthIm,
		exportType,
		filenameRigth, filepathRigth,
		flashMode,
		uint8(flashMode))
	if err != nil {
		return err
	}
	p2, _, err := file.OpenPal(flashPalettePath2)
	if err != nil {
		return err
	}
	return file.FlashLoader(filenameLeft, filenameRigth, p1, p2,uint8(mode),uint8(flashMode), exportType)
}
