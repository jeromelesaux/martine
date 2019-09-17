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
	"path/filepath"
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

	downgraded, p = DoDithering(downgraded, p, exportType)

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
			addr := CpcScreenAddress(0, x, y, 0, exportType.Overscan)
			bw[addr] = pixel
			addr = CpcScreenAddress(0, x+1, y, 0, exportType.Overscan)
			bw[addr] = pixel
		}
	}
	return Export(picturePath, bw, p, 1, exportType)
}
