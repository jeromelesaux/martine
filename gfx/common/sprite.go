package common

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func ToSpriteAndExport(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, ex *export.MartineConfig) error {

	data, firmwareColorUsed, lineSize, err := ToSprite(in, p, size, mode, ex)
	if err != nil {
		return err
	}
	fmt.Println(firmwareColorUsed)
	return ExportSprite(data, lineSize, p, size, mode, filename, dontImportDsk, ex)
}

func ToSpriteHardAndExport(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, ex *export.MartineConfig) error {

	data, firmwareColorUsed := ToSpriteHard(in, p, size, mode, ex)
	fmt.Println(firmwareColorUsed)
	return ExportSprite(data, 16, p, size, mode, filename, false, ex)
}

func ToSprite(in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	mode uint8,
	ex *export.MartineConfig) (data []byte, firmwareColorUsed map[int]int, lineSize int, err error) {

	firmwareColorUsed = make(map[int]int)
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X
	lineToAdd := 1

	if ex.OneLine {
		lineToAdd = 2
	}
	if mode == 0 {
		lineSize = int(math.Ceil(float64(size.Width) / 2.))
		data = make([]byte, size.Height*lineSize)
		offset := 0

		for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
			for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {
				c1 := in.At(x, y)
				pp1, err := PalettePosition(c1, p)
				if err != nil {
					//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
					pp1 = 0
				}
				pp1 = ex.SwapInk(pp1)
				firmwareColorUsed[pp1]++
				//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
				c2 := in.At(x+1, y)
				pp2, err := PalettePosition(c2, p)
				if err != nil {
					//fmt.Fprintf(os.Stdout, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
					pp2 = 0
				}
				pp2 = ex.SwapInk(pp2)
				firmwareColorUsed[pp2]++
				if ex.OneRow {
					pp2 = 0
				}
				pixel := PixelMode0(pp1, pp2)
				if ex.MaskAndOperation {
					pixel = pixel & ex.MaskSprite
				}
				if ex.MaskOrOperation {
					pixel = pixel | ex.MaskSprite
				}
				if len(ex.ScanlineSequence) > 0 {
					scanlineSize := len(ex.ScanlineSequence)
					scanlineIndex := y % scanlineSize
					scanlineValue := ex.ScanlineSequence[scanlineIndex]
					newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 2)
					data[newOffset] = pixel
				} else {
					data[offset] = pixel
				}
				offset++
			}
			if ex.OneLine {
				y++
				for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {
					pp := 0
					firmwareColorUsed[pp]++
					pixel := PixelMode0(pp, pp)
					if len(ex.ScanlineSequence) > 0 {
						scanlineSize := len(ex.ScanlineSequence)
						scanlineIndex := y % scanlineSize
						scanlineValue := ex.ScanlineSequence[scanlineIndex]
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
						//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
						pp1 = 0
					}
					pp1 = ex.SwapInk(pp1)
					firmwareColorUsed[pp1]++
					//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
					c2 := in.At(x+1, y)
					pp2, err := PalettePosition(c2, p)
					if err != nil {
						//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
						pp2 = 0
					}
					pp2 = ex.SwapInk(pp2)
					firmwareColorUsed[pp2]++
					c3 := in.At(x+2, y)
					pp3, err := PalettePosition(c3, p)
					if err != nil {
						//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
						pp3 = 0
					}
					pp3 = ex.SwapInk(pp3)
					firmwareColorUsed[pp3]++
					c4 := in.At(x+3, y)
					pp4, err := PalettePosition(c4, p)
					if err != nil {
						//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
						pp4 = 0
					}
					pp4 = ex.SwapInk(pp4)
					firmwareColorUsed[pp4]++
					if ex.OneRow {
						pp2 = 0
						pp4 = 0
					}
					pixel := PixelMode1(pp1, pp2, pp3, pp4)
					//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
					// MACRO PIXM0 COL2,COL1
					// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
					//	MEND
					if ex.MaskAndOperation {
						pixel = pixel & ex.MaskSprite
					}
					if ex.MaskOrOperation {
						pixel = pixel | ex.MaskSprite
					}
					if len(ex.ScanlineSequence) > 0 {
						scanlineSize := len(ex.ScanlineSequence)
						scanlineIndex := y % scanlineSize
						scanlineValue := ex.ScanlineSequence[scanlineIndex]
						newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 4)
						data[newOffset] = pixel
					} else {
						data[offset] = pixel
					}
					offset++
				}
				if ex.OneLine {
					y++
					for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {
						pp := 0
						firmwareColorUsed[pp]++
						pixel := PixelMode1(pp, pp, pp, pp)
						if len(ex.ScanlineSequence) > 0 {
							scanlineSize := len(ex.ScanlineSequence)
							scanlineIndex := y % scanlineSize
							scanlineValue := ex.ScanlineSequence[scanlineIndex]
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
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
							pp1 = 0
						}
						pp1 = ex.SwapInk(pp1)
						firmwareColorUsed[pp1]++
						//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
						c2 := in.At(x+1, y)
						pp2, err := PalettePosition(c2, p)
						if err != nil {
							//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
							pp2 = 0
						}
						pp2 = ex.SwapInk(pp2)
						firmwareColorUsed[pp2]++
						c3 := in.At(x+2, y)
						pp3, err := PalettePosition(c3, p)
						if err != nil {
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
							pp3 = 0
						}
						pp3 = ex.SwapInk(pp3)
						firmwareColorUsed[pp3]++
						c4 := in.At(x+3, y)
						pp4, err := PalettePosition(c4, p)
						if err != nil {
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
							pp4 = 0
						}
						pp4 = ex.SwapInk(pp4)
						firmwareColorUsed[pp4]++
						c5 := in.At(x+4, y)
						pp5, err := PalettePosition(c5, p)
						if err != nil {
							//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
							pp5 = 0
						}
						pp5 = ex.SwapInk(pp5)
						firmwareColorUsed[pp5]++
						//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
						c6 := in.At(x+5, y)
						pp6, err := PalettePosition(c6, p)
						if err != nil {
							//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
							pp6 = 0
						}
						pp6 = ex.SwapInk(pp6)
						firmwareColorUsed[pp6]++
						c7 := in.At(x+6, y)
						pp7, err := PalettePosition(c7, p)
						if err != nil {
							//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
							pp7 = 0
						}
						pp7 = ex.SwapInk(pp7)
						firmwareColorUsed[pp7]++
						c8 := in.At(x+7, y)
						pp8, err := PalettePosition(c8, p)
						if err != nil {
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+7, y)
							pp8 = 0
						}
						pp8 = ex.SwapInk(pp8)
						firmwareColorUsed[pp8]++
						if ex.OneRow {
							pp2 = 0
							pp4 = 0
							pp6 = 0
							pp8 = 0
						}
						pixel := PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
						//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
						// MACRO PIXM0 COL2,COL1
						// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
						//	MEND
						if len(ex.ScanlineSequence) > 0 {
							scanlineSize := len(ex.ScanlineSequence)
							scanlineIndex := y % scanlineSize
							scanlineValue := ex.ScanlineSequence[scanlineIndex]
							newOffset := (scanlineValue * ((y / scanlineSize) + 1) * lineSize) + (x / 8)
							data[newOffset] = pixel
						} else {
							data[offset] = pixel
						}
						offset++
					}
					if ex.OneLine {
						y++
						for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {
							pp := 0
							firmwareColorUsed[pp]++
							pixel := PixelMode2(pp, pp, pp, pp, pp, pp, pp, pp)
							if len(ex.ScanlineSequence) > 0 {
								scanlineSize := len(ex.ScanlineSequence)
								scanlineIndex := y % scanlineSize
								scanlineValue := ex.ScanlineSequence[scanlineIndex]
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
				return data, firmwareColorUsed, lineSize, errors.ErrorModeNotFound
			}
		}
	}
	return data, firmwareColorUsed, lineSize, nil
}

func ToSpriteHard(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, ex *export.MartineConfig) (data []byte, firmwareColorUsed map[int]int) {
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X
	firmwareColorUsed = make(map[int]int)
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
			data[offset] = byte(ex.SwapInk(pp))
			offset++
		}
	}
	return data, firmwareColorUsed
}
