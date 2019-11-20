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

func DoSpliteRaster(in image.Image, screenMode uint8, filename string, exportType *export.ExportType) error {

	var p color.Palette
	var bw []byte
	var rasters []*constants.SplitRaster
	var err error
	if !exportType.Overscan {
		return ErrorNotYetImplemented
	}
	switch exportType.CpcPlus {
	case false:
		p, bw, rasters, err = ToSplitRasterCPCOld(in, screenMode, filename, exportType)
		if err != nil {
			return err
		}
	default:
		fmt.Fprintf(os.Stderr, "Not yet implemented.")
		return ErrorNotYetImplemented
	}
	// export des données
	if err := Export(filename, bw, p, screenMode, exportType); err != nil {
		return err
	}
	return file.ExportSplitRaster(filename, p, rasters, exportType)
}

func ToSplitRasterCPCOld(in image.Image, screenMode uint8, filename string, exportType *export.ExportType) (color.Palette, []byte, []*constants.SplitRaster, error) {
	rasters := make([]*constants.SplitRaster, 0)
	var bw []byte
	out := convert.Resize(in, exportType.Size, exportType.ResizingAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := file.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_resized.png", out); err != nil {
		return nil, bw, rasters, err
	}
	p, newIm, err := convert.DowngradingPalette(out, exportType.Size, exportType.CpcPlus)
	if err != nil {
		return p, bw, rasters, err
	}
	if err := file.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_downgraded.png", newIm); err != nil {
		return nil, bw, rasters, err
	}

	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), newIm.Bounds().Max.X, newIm.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	var occurence int
	notSplitRaster := true
	backgroundColor := p[0]
	for y := 0; y < exportType.Size.Height; y++ {
		for x := 0; x < exportType.Size.Width; {
			if x%16 == 0 {
				if isSplitRaster(newIm, backgroundColor, x, y, 16) {
					occurence++
					notSplitRaster = false
					pp, _ := PalettePosition(backgroundColor, p)
					fmt.Fprintf(os.Stdout, "X{%d,%d},Y{%d} might be a splitraster\n", x, (x + 16), y)
					switch screenMode {
					case 0:
						for i := 0; i < 16; {
							pixel := pixelMode0(pp, pp)
							addr := CpcScreenAddress(0, x+i, y, 0, exportType.Overscan)
							bw[addr] = pixel
							i += 2
							firmwareColorUsed[pp] += 2
						}
						addr := CpcScreenAddress(0, x, y, 0, exportType.Overscan)
						rasters = append(rasters, constants.NewSpliteRaster(uint16(addr), 16, occurence, pp))
					case 1:
						for i := 0; i < 16; {
							pixel := pixelMode1(pp, pp, pp, pp)
							addr := CpcScreenAddress(0, x+i, y, 1, exportType.Overscan)
							bw[addr] = pixel
							i += 4
							firmwareColorUsed[pp] += 4
						}
						addr := CpcScreenAddress(0, x, y, 1, exportType.Overscan)
						rasters = append(rasters, constants.NewSpliteRaster(uint16(addr), 16, occurence, pp))
					case 2:
						for i := 0; i < 16; {
							pixel := pixelMode2(pp, pp, pp, pp, pp, pp, pp, pp)
							addr := CpcScreenAddress(0, x+i, y, 2, exportType.Overscan)
							bw[addr] = pixel
							i += 8
							firmwareColorUsed[pp] += 8
						}
						addr := CpcScreenAddress(0, x, y, 2, exportType.Overscan)
						rasters = append(rasters, constants.NewSpliteRaster(uint16(addr), 16, occurence, pp))
					}
					// ajout d'un split raster
					// modification de l'image destination pour utiliser celle du background
					// gestion des modes à faire
					x += 16
				} else {
					notSplitRaster = true
				}
			}
			if notSplitRaster {
				// traitement normal des pixels
				switch screenMode {
				case 0:
					bw, firmwareColorUsed = setPixelMode0(newIm, p, x, y, bw, firmwareColorUsed, exportType)
					x += 2
				case 1:
					bw, firmwareColorUsed = setPixelMode1(newIm, p, x, y, bw, firmwareColorUsed, exportType)
					x += 4
				case 2:
					bw, firmwareColorUsed = setPixelMode2(newIm, p, x, y, bw, firmwareColorUsed, exportType)
					x += 4
				}
			}

		}
	}
	fmt.Println(firmwareColorUsed)
	return p, bw, rasters, nil
}

func isSplitRaster(in *image.NRGBA, backgroundColor color.Color, pos, y, length int) bool {
	occ := 0
	for x := pos; x < pos+length || x < in.Bounds().Max.X; x++ {
		c := in.At(x, y)
		if !constants.ColorsAreEquals(c, backgroundColor) {
			return false
		}
		occ++
	}
	if occ < (length - 1) {
		return false
	}
	return true
}

/*
func extractPixelMode0(in *image.NRGBA, p color.Palette, x, y int, exportType *export.ExportType) (pixel byte, addr int) {
	c1 := in.At(x, y)
	pp1, err := PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	pp2, err := PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}
	pixel = pixelMode0(pp1, pp2)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr = CpcScreenAddress(0, x, y, 0, exportType.Overscan)
	return
}
*/

func setPixelMode0(in *image.NRGBA, p color.Palette, x, y int, bw []byte, firmwareColorUsed map[int]int, exportType *export.ExportType) ([]byte, map[int]int) {
	c1 := in.At(x, y)
	pp1, err := PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	firmwareColorUsed[pp1]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	pp2, err := PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}

	firmwareColorUsed[pp2]++

	pixel := pixelMode0(pp1, pp2)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr := CpcScreenAddress(0, x, y, 0, exportType.Overscan)
	bw[addr] = pixel
	return bw, firmwareColorUsed
}

func setPixelMode1(in *image.NRGBA, p color.Palette, x, y int, bw []byte, firmwareColorUsed map[int]int, exportType *export.ExportType) ([]byte, map[int]int) {
	c1 := in.At(x, y)
	pp1, err := PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	firmwareColorUsed[pp1]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	pp2, err := PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}
	firmwareColorUsed[pp2]++
	c3 := in.At(x+2, y)
	pp3, err := PalettePosition(c3, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
		pp3 = 0
	}
	firmwareColorUsed[pp3]++
	c4 := in.At(x+3, y)
	pp4, err := PalettePosition(c4, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
		pp4 = 0
	}
	firmwareColorUsed[pp4]++

	pixel := pixelMode1(pp1, pp2, pp3, pp4)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr := CpcScreenAddress(0, x, y, 1, exportType.Overscan)
	bw[addr] = pixel
	return bw, firmwareColorUsed
}

func setPixelMode2(in *image.NRGBA, p color.Palette, x, y int, bw []byte, firmwareColorUsed map[int]int, exportType *export.ExportType) ([]byte, map[int]int) {
	c1 := in.At(x, y)
	pp1, err := PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	firmwareColorUsed[pp1]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	pp2, err := PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}
	firmwareColorUsed[pp2]++
	c3 := in.At(x+2, y)
	pp3, err := PalettePosition(c3, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
		pp3 = 0
	}
	firmwareColorUsed[pp3]++
	c4 := in.At(x+3, y)
	pp4, err := PalettePosition(c4, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
		pp4 = 0
	}
	firmwareColorUsed[pp4]++
	c5 := in.At(x+4, y)
	pp5, err := PalettePosition(c5, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
		pp5 = 0
	}
	firmwareColorUsed[pp5]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c6 := in.At(x+5, y)
	pp6, err := PalettePosition(c6, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
		pp6 = 0
	}
	firmwareColorUsed[pp6]++
	c7 := in.At(x+6, y)
	pp7, err := PalettePosition(c7, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
		pp3 = 0
	}
	firmwareColorUsed[pp7]++
	c8 := in.At(x+7, y)
	pp8, err := PalettePosition(c8, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+7, y)
		pp8 = 0
	}
	firmwareColorUsed[pp8]++

	pixel := pixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr := CpcScreenAddress(0, x, y, 2, exportType.Overscan)
	bw[addr] = pixel
	return bw, firmwareColorUsed
}
