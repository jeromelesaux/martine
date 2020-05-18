package gfx

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
)

var (
	ErrorColorNotFound     = errors.New("Color not found in palette.")
	ErrorNotYetImplemented = errors.New("Function is not yet implemented.")
	ErrorModeNotFound      = errors.New("Mode not found or not implemented.")
	ErrorBadSize           = errors.New("Width height does not correspond to data size.")
)

func Transform(in *image.NRGBA, p color.Palette, size constants.Size, filepath string, exportType *x.ExportType) error {
	switch size {
	case constants.Mode0:
		return TransformMode0(in, p, size, filepath, exportType)
	case constants.Mode1:
		return TransformMode1(in, p, size, filepath, exportType)
	case constants.Mode2:
		return TransformMode2(in, p, size, filepath, exportType)
	case constants.OverscanMode0:
		return TransformMode0(in, p, size, filepath, exportType)
	case constants.OverscanMode1:
		return TransformMode1(in, p, size, filepath, exportType)
	case constants.OverscanMode2:
		return TransformMode2(in, p, size, filepath, exportType)
	default:
		return ErrorNotYetImplemented
	}
}

func SpriteHardTransform(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, exportType *x.ExportType) error {
	var data []byte
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X
	firmwareColorUsed := make(map[int]int, 0)
	offset := 0
	data = make([]byte, 256)
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			pp, err := PalettePosition(c, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c, x, y)
				pp = 0
			}
			firmwareColorUsed[pp]++
			data[offset] = byte(pp)
			offset++
		}
	}
	fmt.Println(firmwareColorUsed)
	if err := file.Win(filename, data, mode, 16, size.Height, false, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if !exportType.CpcPlus {
		if err := file.Pal(filename, p, mode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettepal.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Ink(filename, p, mode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath = exportType.OsFullPath(filename, "_paletteink.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := file.Kit(filename, p, mode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettekit.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if err := file.Ascii(filename, data, p, false, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return file.AsciiByColumn(filename, data, p, false, mode, exportType)
}

func SpriteTransform(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, exportType *x.ExportType) error {
	var data []byte
	firmwareColorUsed := make(map[int]int, 0)
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X
	var lineSize int
	lineToAdd := 1
	if exportType.OneLine {
		lineToAdd = 2
	}
	fmt.Fprintf(os.Stderr, "%v\n", size)
	if mode == 0 {
		lineSize = int(math.Ceil(float64(size.Width) / 2.))
		data = make([]byte, size.Height*lineSize)
		offset := 0

		for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
			for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {
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
					fmt.Fprintf(os.Stdout, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
					pp2 = 0
				}
				firmwareColorUsed[pp2]++
				pixel := pixelMode0(pp1, pp2)
				if exportType.MaskAndOperation {
					pixel = pixel & exportType.MaskSprite
				}
				if exportType.MaskOrOperation {
					pixel = pixel | exportType.MaskSprite
				}
				if len(exportType.ScanlineSequence) > 0 {
					scanlineSize := len(exportType.ScanlineSequence)
					scanlineIndex := y % scanlineSize
					scanlineValue := exportType.ScanlineSequence[scanlineIndex]
					newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 2)
					data[newOffset] = pixel
				} else {
					data[offset] = pixel
				}
				offset++
			}
			if exportType.OneLine {
				y++
				for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {
					pp := 0
					firmwareColorUsed[pp]++
					pixel := pixelMode0(pp, pp)
					if len(exportType.ScanlineSequence) > 0 {
						scanlineSize := len(exportType.ScanlineSequence)
						scanlineIndex := y % scanlineSize
						scanlineValue := exportType.ScanlineSequence[scanlineIndex]
						newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 2)
						data[newOffset] = pixel
					} else {
						data[offset] = pixel
					}
					offset++
				}
			}
		}
	} else {
		if mode == 1 {
			lineSize = int(math.Ceil(float64(size.Width) / 4.))
			data = make([]byte, size.Height*lineSize)
			offset := 0
			for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
				for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {

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
					if exportType.MaskAndOperation {
						pixel = pixel & exportType.MaskSprite
					}
					if exportType.MaskOrOperation {
						pixel = pixel | exportType.MaskSprite
					}
					if len(exportType.ScanlineSequence) > 0 {
						scanlineSize := len(exportType.ScanlineSequence)
						scanlineIndex := y % scanlineSize
						scanlineValue := exportType.ScanlineSequence[scanlineIndex]
						newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 4)
						data[newOffset] = pixel
					} else {
						data[offset] = pixel
					}
					offset++
				}
				if exportType.OneLine {
					y++
					for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {
						pp := 0
						firmwareColorUsed[pp]++
						pixel := pixelMode1(pp, pp, pp, pp)
						if len(exportType.ScanlineSequence) > 0 {
							scanlineSize := len(exportType.ScanlineSequence)
							scanlineIndex := y % scanlineSize
							scanlineValue := exportType.ScanlineSequence[scanlineIndex]
							newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 2)
							data[newOffset] = pixel
						} else {
							data[offset] = pixel
						}
						offset++
					}
				}
			}

		} else {
			if mode == 2 {
				lineSize = int(math.Ceil(float64(size.Width) / 8.))
				data = make([]byte, size.Height*lineSize)
				offset := 0
				for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
					for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {

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
						if len(exportType.ScanlineSequence) > 0 {
							scanlineSize := len(exportType.ScanlineSequence)
							scanlineIndex := y % scanlineSize
							scanlineValue := exportType.ScanlineSequence[scanlineIndex]
							newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 8)
							data[newOffset] = pixel
						} else {
							data[offset] = pixel
						}
						offset++
					}
					if exportType.OneLine {
						y++
						for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {
							pp := 0
							firmwareColorUsed[pp]++
							pixel := pixelMode2(pp, pp, pp, pp, pp, pp, pp, pp)
							if len(exportType.ScanlineSequence) > 0 {
								scanlineSize := len(exportType.ScanlineSequence)
								scanlineIndex := y % scanlineSize
								scanlineValue := exportType.ScanlineSequence[scanlineIndex]
								newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 2)
								data[newOffset] = pixel
							} else {
								data[offset] = pixel
							}
							offset++
						}
					}
				}
			} else {
				return ErrorModeNotFound
			}
		}
	}
	fmt.Println(firmwareColorUsed)
	if err := file.Win(filename, data, mode, lineSize, size.Height, dontImportDsk, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if !exportType.CpcPlus {
		if err := file.Pal(filename, p, mode, dontImportDsk, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettepal.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Ink(filename, p, 2, dontImportDsk, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath = exportType.OsFullPath(filename, "_paletteink.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := file.Kit(filename, p, mode, dontImportDsk, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettekit.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if err := file.Ascii(filename, data, p, dontImportDsk, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return file.AsciiByColumn(filename, data, p, dontImportDsk, mode, exportType)
}

func PalettePosition(c color.Color, p color.Palette) (int, error) {
	r, g, b, a := c.RGBA()
	for index, cp := range p {
		//fmt.Fprintf(os.Stdout,"index(%d), c:%v,cp:%v\n",index,c,cp)
		rp, gp, bp, ap := cp.RGBA()
		if r == rp && g == gp && b == bp && a == ap {
			//fmt.Fprintf(os.Stdout,"Position found")
			return index, nil
		}
	}
	return -1, ErrorColorNotFound
}

func pixelMode0(pp1, pp2 int) byte {
	var pixel byte
	//fmt.Fprintf(os.Stderr,"1:(%.8b)2:(%.8b)4:(%.8b)8:(%.8b)\n",1,2,4,8)
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&1:%.8b\n",uint8(pp1)&1)
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&2:%.8b\n",uint8(pp1)&2)
	if uint8(pp1)&2 == 2 {
		pixel += 8
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&4:%.8b\n",uint8(pp1)&4)
	if uint8(pp1)&4 == 4 {
		pixel += 32
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&8:%.8b\n",uint8(pp1)&8)
	if uint8(pp1)&8 == 8 {
		pixel += 2
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp2)&2 == 2 {
		pixel += 4
	}
	if uint8(pp2)&4 == 4 {
		pixel += 16
	}
	if uint8(pp2)&8 == 8 {
		pixel++
	}
	return pixel
}

func pixelMode1(pp1, pp2, pp3, pp4 int) byte {
	var pixel byte
	//fmt.Fprintf(os.Stderr,"1:(%.8b)2:(%.8b)4:(%.8b)8:(%.8b)\n",1,2,4,8)
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&1:%.8b\n",uint8(pp1)&1)
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&2:%.8b\n",uint8(pp1)&2)
	if uint8(pp1)&2 == 2 {
		pixel += 8
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&4:%.8b\n",uint8(pp1)&4)
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&8:%.8b\n",uint8(pp1)&8)
	if uint8(pp2)&2 == 2 {
		pixel += 4
	}
	if uint8(pp3)&1 == 1 {
		pixel += 32
	}
	if uint8(pp3)&2 == 2 {
		pixel += 2
	}
	if uint8(pp4)&1 == 1 {
		pixel += 16
	}
	if uint8(pp4)&2 == 2 {
		pixel++
	}
	return pixel
}

func pixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 int) byte {
	var pixel byte
	//fmt.Fprintf(os.Stderr,"1:(%.8b)2:(%.8b)4:(%.8b)8:(%.8b)\n",1,2,4,8)
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&1:%.8b\n",uint8(pp1)&1)
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&2:%.8b\n",uint8(pp1)&2)
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&4:%.8b\n",uint8(pp1)&4)
	if uint8(pp3)&1 == 1 {
		pixel += 32
	}
	//fmt.Fprintf(os.Stderr,"uint8(pp1)&8:%.8b\n",uint8(pp1)&8)
	if uint8(pp4)&1 == 1 {
		pixel += 16
	}
	if uint8(pp5)&1 == 1 {
		pixel += 8
	}
	if uint8(pp6)&1 == 1 {
		pixel += 4
	}
	if uint8(pp7)&1 == 1 {
		pixel += 2
	}
	if uint8(pp8)&1 == 1 {
		pixel++
	}
	return pixel
}

func rawPixelMode2(b byte) (pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 int) {
	val := int(b)
	if val-128 >= 0 {
		pp1 = 1
		val -= 128
	}
	if val-64 >= 0 {
		pp2 = 1
		val -= 64
	}
	if val-32 >= 0 {
		pp3 = 1
		val -= 32
	}
	if val-16 >= 0 {
		pp4 = 1
		val -= 16
	}
	if val-8 >= 0 {
		pp5 = 1
		val -= 8
	}
	if val-4 >= 0 {
		pp6 = 1
		val -= 4
	}
	if val-2 >= 0 {
		pp7 = 1
		val -= 2
	}
	if val-1 >= 0 {
		pp8 = 1
	}
	return
}

func rawPixelMode1(b byte) (pp1, pp2, pp3, pp4 int) {
	val := int(b)
	if val-128 >= 0 {
		pp1 |= 1
		val -= 128
	}
	if val-64 >= 0 {
		pp2 |= 1
		val -= 64
	}
	if val-32 >= 0 {
		pp3 |= 1
		val -= 32
	}
	if val-16 >= 0 {
		pp4 |= 1
		val -= 16
	}
	if val-8 >= 0 {
		pp1 |= 2
		val -= 8
	}
	if val-4 >= 0 {
		pp2 |= 2
		val -= 4
	}
	if val-2 >= 0 {
		pp3 |= 2
		val -= 2
	}
	if val-1 >= 0 {
		pp4 |= 2
	}

	return
}

func rawPixelMode0(b byte) (pp1, pp2 int) {
	val := int(b)
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-128 >= 0 {
		pp1 |= 1
		val -= 128
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-64 >= 0 {
		pp2 |= 1
		val -= 64
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-32 >= 0 {
		pp1 |= 4
		val -= 32
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-16 >= 0 {
		pp2 |= 4
		val -= 16
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-8 >= 0 {
		pp1 |= 2
		val -= 8
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-4 >= 0 {
		pp2 |= 2
		val -= 4
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-2 >= 0 {
		pp1 |= 8
		val -= 2
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-1 >= 0 {
		pp2 |= 8
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	return
}

func ToMode0(in *image.NRGBA, p color.Palette, exportType *x.ExportType) []byte {
	var bw []byte

	lineToAdd := 1
	if exportType.OneLine {
		lineToAdd = 2
	}
	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {

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
		}
	}

	fmt.Println(firmwareColorUsed)
	return bw
}

func TransformMode0(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *x.ExportType) error {
	bw := ToMode0(in, p, exportType)
	return Export(filePath, bw, p, 0, exportType)
}

func Export(filePath string, bw []byte, p color.Palette, screenMode uint8, exportType *x.ExportType) error {
	if exportType.Overscan {
		if exportType.EgxFormat == 0 {
			if err := file.Overscan(filePath, bw, p, screenMode, exportType); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
				return err
			}
		} else {
			if err := file.EgxOverscan(filePath, bw, p, exportType.EgxMode1, exportType.EgxMode2, exportType); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
				return err
			}
		}

	} else {
		if err := file.Scr(filePath, bw, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Loader(filePath, p, screenMode, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if !exportType.CpcPlus {
		if err := file.Pal(filePath, p, screenMode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := exportType.OsFullPath(filePath, "_palettepal.png")
		if err := file.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
		if err := file.Ink(filePath, p, screenMode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 = exportType.OsFullPath(filePath, "_paletteink.png")
		if err := file.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	} else {
		if err := file.Kit(filePath, p, screenMode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := exportType.OsFullPath(filePath, "_palettekit.png")
		if err := file.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	}
	return file.Ascii(filePath, bw, p, false, exportType)
}

func CpcScreenAddress(intialeAddresse int, x, y int, mode uint8, isOverscan bool) int {
	var addr int
	var adjustMode int
	switch mode {
	case 0:
		adjustMode = 2
	case 1:
		adjustMode = 4
	case 2:
		adjustMode = 8
	}
	if isOverscan {
		if y > 127 {
			addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / adjustMode) + (0x3800)
		} else {
			addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / adjustMode)
		}
	} else {
		addr = (0x800 * (y % 8)) + (0x50 * (y / 8)) + ((x + 1) / adjustMode)
	}
	if intialeAddresse == 0 {
		return addr
	}
	return intialeAddresse + addr
}

func CpcScreenAddressOffset(line int) int {
	return int(math.Floor(float64(line)/8)*80) + (line%8)*2048
}

func ToMode1(in *image.NRGBA, p color.Palette, exportType *x.ExportType) []byte {
	var bw []byte

	lineToAdd := 1
	if exportType.OneLine {
		lineToAdd = 2
	}
	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}

	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {

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
		}
	}
	return bw
}

func TransformMode1(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *x.ExportType) error {
	bw := ToMode1(in, p, exportType)
	return Export(filePath, bw, p, 1, exportType)
}

func ToMode2(in *image.NRGBA, p color.Palette, exportType *x.ExportType) []byte {
	var bw []byte

	lineToAdd := 1
	if exportType.OneLine {
		lineToAdd = 2
	}
	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {

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
		}

	}

	fmt.Println(firmwareColorUsed)
	return bw
}

func TransformMode2(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *x.ExportType) error {
	bw := ToMode2(in, p, exportType)
	return Export(filePath, bw, p, 2, exportType)
}

func revertColor(rawColor uint8, index int, isPlus bool) color.Color {
	var newColor color.Color
	var err error
	if !isPlus {
		newColor, err = constants.ColorFromHardware(rawColor)
		if err != nil {
			fmt.Fprintf(os.Stderr, "No color found in data at index %d\n", index)
			return constants.White.Color
		}
	} else {
		plusColor := constants.NewRawCpcPlusColor(uint16(rawColor))
		c := color.RGBA{A: 0xFF, R: uint8(plusColor.R), G: uint8(plusColor.G), B: uint8(plusColor.B)}
		newColor = constants.CpcPlusPalette.Convert(c)
	}
	return newColor
}

func TransformRawCpcData(data, palette []int, width, height int, mode int, isPlus bool) (*image.NRGBA, error) {

	in := image.NewNRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: width, Y: height}})
	x := 0
	y := 0
	for index, val := range data {

		switch mode {
		case 0:
			p1, p2 := rawPixelMode0(byte(val))
			c1 := palette[p1]
			newColor := revertColor(uint8(c1), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c2 := palette[p2]
			newColor = revertColor(uint8(c2), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
		case 1:
			p1, p2, p3, p4 := rawPixelMode1(byte(val))
			c1 := palette[p1]
			newColor := revertColor(uint8(c1), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c2 := palette[p2]
			newColor = revertColor(uint8(c2), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c3 := palette[p3]
			newColor = revertColor(uint8(c3), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c4 := palette[p4]
			newColor = revertColor(uint8(c4), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
		case 2:
			p1, p2, p3, p4, p5, p6, p7, p8 := rawPixelMode2(byte(val))
			c1 := palette[p1]
			newColor := revertColor(uint8(c1), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c2 := palette[p2]
			newColor = revertColor(uint8(c2), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c3 := palette[p3]
			newColor = revertColor(uint8(c3), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c4 := palette[p4]
			newColor = revertColor(uint8(c4), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c5 := palette[p5]
			newColor = revertColor(uint8(c5), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c6 := palette[p6]
			newColor = revertColor(uint8(c6), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c7 := palette[p7]
			newColor = revertColor(uint8(c7), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c8 := palette[p8]
			newColor = revertColor(uint8(c8), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
		}
	}
	return in, nil
}
