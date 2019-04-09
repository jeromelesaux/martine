package gfx

import (
	"errors"
	"fmt"
	"github.com/jeromelesaux/martine/constants"
	"image"
	"image/color"
	"os"
	"math"
)

var (
	ErrorColorNotFound     = errors.New("Color not found in palette.")
	ErrorNotYetImplemented = errors.New("Function is not yet implemented.")
	ErrorModeNotFound      = errors.New("Mode not found or not implemented.")
)

func Transform(in *image.NRGBA, p color.Palette, size constants.Size, filepath string, exportType *ExportType) error {
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

func SpriteTransform(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, exportType *ExportType) error {
	var data []byte
	firmwareColorUsed := make(map[int]int, 0)
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X
	var lineSize int
	fmt.Fprintf(os.Stderr, "%v\n", size)
	if mode == 0 {
		lineSize = int(math.Ceil(float64(size.Width) / 2.))
		data = make([]byte, size.Height*lineSize)
		offset := 0
		for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
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

				data[offset] = pixel
				offset++
			}
		}
	} else {
		if mode == 1 {
			lineSize = int(math.Ceil(float64(size.Width) / 4.))
			data = make([]byte, size.Height*lineSize)
			offset := 0
			for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
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

					data[offset] = pixel
					offset++
				}
			}
		} else {
			if mode == 2 {
				lineSize = int(math.Ceil(float64(size.Width) / 8.))
				data = make([]byte, size.Height*lineSize)
				offset := 0
				for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
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
						data[offset] = pixel
						offset++
					}
				}
			} else {
				return ErrorModeNotFound
			}
		}
	}
	fmt.Println(firmwareColorUsed)
	if err := Win(filename, data, mode, lineSize, size.Height, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if err := Pal(filename, p, mode, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if err := Ink(filename, p, 2, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if err := Ascii(filename, data, p, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return AsciiByColumn(filename, data, p, exportType)
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

func TransformMode0(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *ExportType) error {
	var bw []byte
	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
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
			if exportType.Overscan {
				var addr int
				if y > 127 {
					addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / 2) + (0x3800)
				} else {
					addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / 2)

				}
				bw[addr] = pixel
			} else {
				bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/2)] = pixel
			}
		}
	}

	fmt.Println(firmwareColorUsed)
	if exportType.Overscan {
		if err := Overscan(filePath, bw, p, 0, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Scr(filePath, bw, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := Loader(filePath, p, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if exportType.CpcPlus {
		if err := Kit(filePath, p, 0, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Pal(filePath, p, 0, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := Ink(filePath, p, 2, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}

	return Ascii(filePath, bw, p, exportType)
}

func TransformMode1(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *ExportType) error {
	var bw []byte
	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}

	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
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
			if exportType.Overscan {
				var addr int
				if y > 127 {
					addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / 4) + (0x3800)
				} else {
					addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / 4)

				}
				bw[addr] = pixel
			} else {
				bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/4)] = pixel
			}
		}
	}

	fmt.Println(firmwareColorUsed)
	if exportType.Overscan {
		if err := Overscan(filePath, bw, p, 1, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Scr(filePath, bw, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := Loader(filePath, p, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if exportType.CpcPlus {
		if err := Kit(filePath, p, 0, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Pal(filePath, p, 1, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := Ink(filePath, p, 2, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}

	return Ascii(filePath, bw, p, exportType)
}

func TransformMode2(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *ExportType) error {
	var bw []byte

	if exportType.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
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
			if exportType.Overscan {
				var addr int
				if y > 127 {
					addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / 8) + (0x3800)
				} else {
					addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / 8)

				}
				bw[addr] = pixel
			} else {
				bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/8)] = pixel
			}
			//bw = append(bw, pixel)
		}

	}

	fmt.Println(firmwareColorUsed)
	if exportType.Overscan {
		if err := Overscan(filePath, bw, p, 2, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Scr(filePath, bw, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := Loader(filePath, p, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if exportType.CpcPlus {
		if err := Kit(filePath, p, 0, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Pal(filePath, p, 2, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := Ink(filePath, p, 2, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}

	return Ascii(filePath, bw, p, exportType)
}
