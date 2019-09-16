package gfx

import (
	"fmt"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Egx1(in image.Image,
	exportType *export.ExportType,
	filename, picturePath string) error {
	var err error

	size := constants.Size{
		Width:  exportType.Size.Width,
		Height: exportType.Size.Height}

	im := convert.Resize(in, size, exportType.ResizingAlgo)
	bw := make([]byte, 0x4000) // ecran 17ko
	firmwareColorUsed := make(map[int]int, 0)
	var palette color.Palette // palette de l'image
	var p color.Palette       // palette cpc de l'image
	var downgraded *image.NRGBA

	if exportType.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", exportType.PalettePath)
		palette, _, err = file.OpenPal(exportType.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", exportType.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if len(palette) > 0 {
		p, downgraded = convert.DowngradingWithPalette(im, palette)
	} else {
		p, downgraded, err = convert.DowngradingPalette(im, exportType.Size, exportType.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := file.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_down.png", downgraded); err != nil {
		os.Exit(-2)
	}
	for y := downgraded.Bounds().Min.Y + 1; y < downgraded.Bounds().Max.Y; y += 2 {
		for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x += 2 {
			c1 := downgraded.At(x, y)
			pp1, err := PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := downgraded.At(x+1, y)
			pp2, err := PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			pixel := pixelMode0(pp1, pp2)
			addr := CpcScreenAddress(0, x, y, 0, exportType.Overscan)
			bw[addr] = pixel
		}
	}
	for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y += 2 {
		for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x += 4 {
			c1 := downgraded.At(x, y)
			pp1, err := PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := downgraded.At(x+1, y)
			pp2, err := PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := downgraded.At(x+2, y)
			pp3, err := PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := downgraded.At(x+3, y)
			pp4, err := PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++

			pixel := pixelMode1(pp1, pp2, pp3, pp4)
			addr := CpcScreenAddress(0, x, y, 1, exportType.Overscan)
			bw[addr] = pixel
		}
	}
	return Export(picturePath, bw, p, 1, exportType)
}

func Flash(filepath1, filepath2, palpath1, palpath2 string, m1, m2 int, exportType *export.ExportType) error {
	if filepath2 == "" && filepath1 != "" {
		filename := filepath.Base(filepath1)
		f, err := os.Open(filepath1)
		if err != nil {
			return err
		}
		defer f.Close()
		in, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		return AutoFlash(in, exportType, filename, filepath1, m1, uint8(m1))
	}
	filename1 := filepath.Base(filepath1)
	filename2 := filepath.Base(filepath2)
	p1, _, err := file.OpenPal(palpath1)
	if err != nil {
		return err
	}
	p2, _, err := file.OpenPal(palpath2)
	if err != nil {
		return err
	}
	exportType.AddFile(filepath1)
	exportType.AddFile(filepath2)
	exportType.AddFile(palpath1)
	exportType.AddFile(palpath2)
	return file.FlashLoader(filename1, filename2, p1, p2, uint8(m1), uint8(m2), exportType)
}

func AutoFlash(in image.Image,
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

	exportType.PalettePath = flashPalettePath1

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
	return file.FlashLoader(filenameLeft, filenameRigth, p1, p2, uint8(mode), uint8(flashMode), exportType)
}
