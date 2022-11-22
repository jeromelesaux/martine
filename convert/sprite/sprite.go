package sprite

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	pal "github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func RawSpriteToImg(data []byte, height, width, mode uint8, p color.Palette) *image.NRGBA {
	var out *image.NRGBA
	switch mode {
	case 0:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 2), Y: int(height)}})
		index := 0

		for y := 0; y < int(height); y++ {
			indexX := 0
			for x := 0; x < int(width); x++ {
				val := data[index]
				pp1, pp2 := pixel.RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				index++
			}
		}
	case 1:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 4), Y: int(height)}})
		index := 0
		for y := 0; y < int(height); y++ {
			indexX := 0
			for x := 0; x < int(width); x++ {
				val := data[index]
				pp1, pp2, pp3, pp4 := pixel.RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				index++
			}
		}
	case 2:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 8), Y: int(height)}})
		index := 0
		for y := 0; y < int(width); y++ {
			indexX := 0
			for x := 0; x < int(width/8); x++ {
				val := data[index]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := pixel.RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				out.Set(indexX, y, c5)
				indexX++
				out.Set(indexX, y, c6)
				indexX++
				out.Set(indexX, y, c7)
				indexX++
				out.Set(indexX, y, c8)
				indexX++
				index++
			}
		}
	}
	return out
}

// spriteToImg load a OCP win filepath to image.NRGBA
// using the mode and palette as arguments
func SpriteToImg(winPath string, mode uint8, p color.Palette) (*image.NRGBA, constants.Size, error) {
	var s constants.Size
	footer, err := window.OpenWin(winPath)
	if err != nil {
		return nil, s, err
	}
	var out *image.NRGBA

	d, err := window.RawWin(winPath)
	if err != nil {
		return nil, s, err
	}
	switch mode {
	case 0:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(footer.Width * 2), Y: int(footer.Height)}})
		index := 0
		s.Width = int(footer.Width * 2)
		s.Height = int(footer.Height)
		for y := 0; y < int(footer.Height); y++ {
			indexX := 0
			for x := 0; x < int(footer.Width); x++ {
				val := d[index]
				pp1, pp2 := pixel.RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				index++
			}
		}
	case 1:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(footer.Width * 4), Y: int(footer.Height)}})
		index := 0
		s.Width = int(footer.Width * 4)
		s.Height = int(footer.Height)
		for y := 0; y < int(footer.Height); y++ {
			indexX := 0
			for x := 0; x < int(footer.Width); x++ {
				val := d[index]
				pp1, pp2, pp3, pp4 := pixel.RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				index++
			}
		}
	case 2:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(footer.Width * 8), Y: int(footer.Height)}})
		index := 0
		s.Width = int(footer.Width * 8)
		s.Height = int(footer.Height)
		for y := 0; y < int(footer.Width); y++ {
			indexX := 0
			for x := 0; x < int(footer.Width/8); x++ {
				val := d[index]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := pixel.RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				out.Set(indexX, y, c5)
				indexX++
				out.Set(indexX, y, c6)
				indexX++
				out.Set(indexX, y, c7)
				indexX++
				out.Set(indexX, y, c8)
				indexX++
				index++
			}
		}
	}
	return out, s, err
}

// spriteToPng will save the sprite OCP win filepath in a png file
// output will be the png filpath, using the mode and palette as arguments
func SpriteToPng(winPath string, output string, mode uint8, p color.Palette) error {
	out, _, err := SpriteToImg(winPath, mode, p)
	if err != nil {
		return err
	}
	return png.Png(output+".png", out)
}

func ExportSprite(data []byte, lineSize int, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, cfg *config.MartineConfig) error {
	if err := window.Win(filename, data, mode, lineSize, size.Height, dontImportDsk, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if !cfg.CpcPlus {
		if err := ocpartstudio.Pal(filename, p, mode, dontImportDsk, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := cfg.OsFullPath(filename, "_palettepal.png")
		if err := png.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := palette.Ink(filename, p, 2, dontImportDsk, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath = cfg.OsFullPath(filename, "_paletteink.png")
		if err := png.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := palette.Kit(filename, p, mode, dontImportDsk, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := cfg.OsFullPath(filename, "_palettekit.png")
		if err := png.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if err := ascii.Ascii(filename, data, p, dontImportDsk, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return ascii.AsciiByColumn(filename, data, p, dontImportDsk, mode, cfg)
}

func ToSpriteAndExport(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, ex *config.MartineConfig) error {

	data, firmwareColorUsed, lineSize, err := ToSprite(in, p, size, mode, ex)
	if err != nil {
		return err
	}
	fmt.Println(firmwareColorUsed)
	return ExportSprite(data, lineSize, p, size, mode, filename, dontImportDsk, ex)
}

func ToSprite(in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	mode uint8,
	ex *config.MartineConfig) (data []byte, firmwareColorUsed map[int]int, lineSize int, err error) {

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
				pp1, err := pal.PalettePosition(c1, p)
				if err != nil {
					//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
					pp1 = 0
				}
				pp1 = ex.SwapInk(pp1)
				firmwareColorUsed[pp1]++
				//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
				c2 := in.At(x+1, y)
				pp2, err := pal.PalettePosition(c2, p)
				if err != nil {
					//fmt.Fprintf(os.Stdout, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
					pp2 = 0
				}
				pp2 = ex.SwapInk(pp2)
				firmwareColorUsed[pp2]++
				if ex.OneRow {
					pp2 = 0
				}
				pixel := pixel.PixelMode0(pp1, pp2)
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
					pixel := pixel.PixelMode0(pp, pp)
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
					pp1, err := pal.PalettePosition(c1, p)
					if err != nil {
						//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
						pp1 = 0
					}
					pp1 = ex.SwapInk(pp1)
					firmwareColorUsed[pp1]++
					//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
					c2 := in.At(x+1, y)
					pp2, err := pal.PalettePosition(c2, p)
					if err != nil {
						//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
						pp2 = 0
					}
					pp2 = ex.SwapInk(pp2)
					firmwareColorUsed[pp2]++
					c3 := in.At(x+2, y)
					pp3, err := pal.PalettePosition(c3, p)
					if err != nil {
						//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
						pp3 = 0
					}
					pp3 = ex.SwapInk(pp3)
					firmwareColorUsed[pp3]++
					c4 := in.At(x+3, y)
					pp4, err := pal.PalettePosition(c4, p)
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
					pixel := pixel.PixelMode1(pp1, pp2, pp3, pp4)
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
						pixel := pixel.PixelMode1(pp, pp, pp, pp)
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
						pp1, err := pal.PalettePosition(c1, p)
						if err != nil {
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
							pp1 = 0
						}
						pp1 = ex.SwapInk(pp1)
						firmwareColorUsed[pp1]++
						//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
						c2 := in.At(x+1, y)
						pp2, err := pal.PalettePosition(c2, p)
						if err != nil {
							//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
							pp2 = 0
						}
						pp2 = ex.SwapInk(pp2)
						firmwareColorUsed[pp2]++
						c3 := in.At(x+2, y)
						pp3, err := pal.PalettePosition(c3, p)
						if err != nil {
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
							pp3 = 0
						}
						pp3 = ex.SwapInk(pp3)
						firmwareColorUsed[pp3]++
						c4 := in.At(x+3, y)
						pp4, err := pal.PalettePosition(c4, p)
						if err != nil {
							//		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
							pp4 = 0
						}
						pp4 = ex.SwapInk(pp4)
						firmwareColorUsed[pp4]++
						c5 := in.At(x+4, y)
						pp5, err := pal.PalettePosition(c5, p)
						if err != nil {
							//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
							pp5 = 0
						}
						pp5 = ex.SwapInk(pp5)
						firmwareColorUsed[pp5]++
						//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
						c6 := in.At(x+5, y)
						pp6, err := pal.PalettePosition(c6, p)
						if err != nil {
							//	fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
							pp6 = 0
						}
						pp6 = ex.SwapInk(pp6)
						firmwareColorUsed[pp6]++
						c7 := in.At(x+6, y)
						pp7, err := pal.PalettePosition(c7, p)
						if err != nil {
							//fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
							pp7 = 0
						}
						pp7 = ex.SwapInk(pp7)
						firmwareColorUsed[pp7]++
						c8 := in.At(x+7, y)
						pp8, err := pal.PalettePosition(c8, p)
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
						pixel := pixel.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
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
							pixel := pixel.PixelMode2(pp, pp, pp, pp, pp, pp, pp, pp)
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
