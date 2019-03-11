package gfx

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
)

var (
	ErrorColorNotFound     = errors.New("Color not found in palette.")
	ErrorNotYetImplemented = errors.New("Function is not yet implemented.")
	ErrorModeNotFound      = errors.New("Mode not found or not implemented.")
)

func Transform(in *image.NRGBA, p color.Palette, size Size, filepath, dirpath string, noAmsdosHeader, isCpcPlus bool) error {
	switch size {
	case Mode0:
		return TransformMode0(in, p, size, filepath, dirpath, false, noAmsdosHeader, isCpcPlus)
	case Mode1:
		return TransformMode1(in, p, size, filepath, dirpath, false, noAmsdosHeader, isCpcPlus)
	case Mode2:
		return TransformMode2(in, p, size, filepath, dirpath, false, noAmsdosHeader, isCpcPlus)
	case OverscanMode0:
		return TransformMode0(in, p, size, filepath, dirpath, true, noAmsdosHeader, isCpcPlus)
	case OverscanMode1:
		return TransformMode1(in, p, size, filepath, dirpath, true, noAmsdosHeader, isCpcPlus)
	case OverscanMode2:
		return TransformMode2(in, p, size, filepath, dirpath, true, noAmsdosHeader, isCpcPlus)
	default:
		return ErrorNotYetImplemented
	}
}

func SpriteTransform(in *image.NRGBA, p color.Palette, size Size, mode uint8, filePath, dirPath string, noAmsdosHeader, isCpcPlus bool) error {
	var data []byte
	firmwareColorUsed := make(map[int]int, 0)
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X

	if mode == 0 {
		data = make([]byte, (size.Height * (size.Width/2)))
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
					fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
					pp2 = 0
				}

				firmwareColorUsed[pp2]++

				pixel := pixelMode0(pp1, pp2)
				//fmt.Fprintf(os.Stderr,"(%d,%d)[#%.2x]:#%.2x\n",y,x,offset,pixel)
				data[offset] = pixel
				offset++
			}
		}
	} else {
		if mode == 1 {
			data = make([]byte, (size.Height * size.Width))
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
					offset+=4
				}
			}
		} else {
			if mode == 2 {
				data = make([]byte, (size.Height * size.Width))
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
						offset+=8
					}
				}
			} else {
				return ErrorModeNotFound
			}
		}
	}
	fmt.Println(firmwareColorUsed)
	if err := Win(filePath, dirPath, data, mode, size.Width, size.Height, noAmsdosHeader); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
		return err
	}
	if err := Pal(filePath, dirPath, p, mode, noAmsdosHeader); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
		return err
	}
	return Ascii(filePath, dirPath, data, p, noAmsdosHeader, isCpcPlus)
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

func TransformMode0(in *image.NRGBA, p color.Palette, size Size, filePath, dirPath string, overscan, noAmsdosHeader, isCpcPlus bool) error {
	var bw []byte
	if overscan {
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
			if overscan {
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
	if overscan {
		if err := Overscan(filePath, dirPath, bw, p, 0, noAmsdosHeader, isCpcPlus); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Scr(filePath, dirPath, bw, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if isCpcPlus {
		if err := Ink(filePath, dirPath, p, 0, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Pal(filePath, dirPath, p, 0, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}

	return Ascii(filePath, dirPath, bw, p, noAmsdosHeader, isCpcPlus)
}

func TransformMode1(in *image.NRGBA, p color.Palette, size Size, filePath, dirPath string, overscan, noAmsdosHeader, isCpcPlus bool) error {
	var bw []byte
	if overscan {
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
			if overscan {
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
	if overscan {
		if err := Overscan(filePath, dirPath, bw, p, 1, noAmsdosHeader, isCpcPlus); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Scr(filePath, dirPath, bw, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if isCpcPlus {
		if err := Ink(filePath, dirPath, p, 0, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Pal(filePath, dirPath, p, 1, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}

	return Ascii(filePath, dirPath, bw, p, noAmsdosHeader, isCpcPlus)
}

func TransformMode2(in *image.NRGBA, p color.Palette, size Size, filePath, dirPath string, overscan, noAmsdosHeader, isCpcPlus bool) error {
	var bw []byte

	if overscan {
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
			if overscan {
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
	if overscan {
		if err := Overscan(filePath, dirPath, bw, p, 2, noAmsdosHeader, isCpcPlus); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Scr(filePath, dirPath, bw, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if isCpcPlus {
		if err := Ink(filePath, dirPath, p, 0, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := Pal(filePath, dirPath, p, 2, noAmsdosHeader); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}

	return Ascii(filePath, dirPath, bw, p, noAmsdosHeader, isCpcPlus)
}
