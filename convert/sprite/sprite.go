// Package sprite provides sprite conversion and export utilities.
package sprite

import (
	"image"
	"image/color"
	"math"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/export"
	pal "github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/log"
)

func ToImg(d [][]byte, p color.Palette) *image.NRGBA {
	o := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: len(d[0]), Y: len(d)},
	})

	for y := 0; y < len(d); y++ {
		for x := 0; x < len(d[0]); x++ {
			o.Set(x, y, p[d[y][x]])
		}
	}
	return o
}

// RawSpriteToImg converts raw sprite bytes into an NRGBA image.
// nolint:funlen
func RawSpriteToImg(data []byte, height, width, mode uint8, p color.Palette) *image.NRGBA {
	var out *image.NRGBA
	switch mode {
	case 0:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 2), Y: int(height)},
		})
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
			Max: image.Point{X: int(width * 4), Y: int(height)},
		})
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
			Max: image.Point{X: int(width * 8), Y: int(height)},
		})
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

// SpriteToImg loads an OCP win filepath into image.NRGBA.
// using the mode and palette as arguments
func SpriteToImg(winPath string, mode uint8, p color.Palette) (*image.NRGBA, constants.Size, error) {
	var s constants.Size
	footer, err := window.OpenWin(winPath)
	if err != nil {
		return nil, s, err
	}
	d, err := window.RawWin(winPath)
	if err != nil {
		return nil, s, err
	}
	pixelsPerByte := []int{2, 4, 8}
	s.Width = int(footer.Width) * pixelsPerByte[mode]
	s.Height = int(footer.Height)
	out := RawSpriteToImg(d, footer.Height, uint8(footer.Width), mode, p)
	return out, s, nil
}

// SpriteToPng saves the sprite OCP win filepath as a png file.
// output will be the png filpath, using the mode and palette as arguments
func SpriteToPng(winPath, output string, mode uint8, p color.Palette) error {
	out, _, err := SpriteToImg(winPath, mode, p)
	if err != nil {
		return err
	}
	return png.Png(output+".png", out)
}

func ExportSprite(data []byte, lineSize int, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, cfg *config.MartineConfig) error {
	if err := window.Win(filename, data, mode, lineSize, size.Height, dontImportDsk, cfg); err != nil {
		log.GetLogger().Error("Error while saving file %s error :%v", filename, err)
		return err
	}
	if err := export.ExportPalette(filename, data, p, mode, dontImportDsk, cfg); err != nil {
		return err
	}
	if err := ascii.ASCII(filename, data, p, dontImportDsk, cfg); err != nil {
		log.GetLogger().Error("Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return ascii.ASCIIByColumn(filename, data, p, dontImportDsk, mode, cfg)
}

func ToSpriteAndExport(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, ex *config.MartineConfig) error {
	data, _, lineSize, err := ToSprite(in, p, size, mode, ex)
	if err != nil {
		return err
	}

	return ExportSprite(data, lineSize, p, size, mode, filename, dontImportDsk, ex)
}

func toSpriteMode(
	in *image.NRGBA,
	_ color.Palette,
	size constants.Size,
	xStep int,
	ex *config.MartineConfig,
	readPixel func(x, y int) byte,
	blankPixel func() byte,
) (data []byte, count int) {
	lineSize := int(math.Ceil(float64(size.Width) / float64(xStep)))
	data = make([]byte, size.Height*lineSize)
	offset := 0
	lineToAdd := 1
	if ex.ScrCfg.Process.OneLine {
		lineToAdd = 2
	}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += xStep {
			px := readPixel(x, y)
			if ex.ScrCfg.Process.Mask.AndOperation {
				px &= ex.ScrCfg.Process.Mask.Sprite
			}
			if ex.ScrCfg.Process.Mask.OrOperation {
				px |= ex.ScrCfg.Process.Mask.Sprite
			}
			if len(ex.ScanlineSequence) > 0 {
				scanlineSize := len(ex.ScanlineSequence)
				scanlineIndex := y % scanlineSize
				scanlineValue := ex.ScanlineSequence[scanlineIndex]
				newOffset := (scanlineValue*((y/scanlineSize)+1)*lineSize + (x / xStep))
				data[newOffset] = px
			} else {
				data[offset] = px
			}
			offset++
		}
		if ex.ScrCfg.Process.OneLine {
			y++
			for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += xStep {
				px := blankPixel()
				if len(ex.ScanlineSequence) > 0 {
					scanlineSize := len(ex.ScanlineSequence)
					scanlineIndex := y % scanlineSize
					scanlineValue := ex.ScanlineSequence[scanlineIndex]
					newOffset := (scanlineValue*((y/scanlineSize)+1)*lineSize + (x / xStep))
					data[newOffset] = px
				} else {
					data[offset] = px
				}
				offset++
			}
		}
	}
	return data, lineSize
}

// ToSprite converts an NRGBA image into CPC sprite bytes for the given mode.
// nolint:funlen, gocognit
func ToSprite(in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	mode uint8,
	ex *config.MartineConfig,
) (data []byte, firmwareColorUsed map[int]int, lineSize int, err error) {
	firmwareColorUsed = make(map[int]int)
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X

	switch mode {
	case 0:
		data, lineSize = toSpriteMode(in, p, size, 2, ex,
			func(x, y int) byte {
				pp1 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x, y), p))
				pp2 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+1, y), p))
				if ex.ScrCfg.Process.OneRow {
					pp2 = 0
				}
				return pixel.PixelMode0(pp1, pp2)
			},
			func() byte { return pixel.PixelMode0(0, 0) },
		)
	case 1:
		data, lineSize = toSpriteMode(in, p, size, 4, ex,
			func(x, y int) byte {
				pp1 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x, y), p))
				pp2 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+1, y), p))
				pp3 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+2, y), p))
				pp4 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+3, y), p))
				if ex.ScrCfg.Process.OneRow {
					pp2, pp4 = 0, 0
				}
				return pixel.PixelMode1(pp1, pp2, pp3, pp4)
			},
			func() byte { return pixel.PixelMode1(0, 0, 0, 0) },
		)
	case 2:
		data, lineSize = toSpriteMode(in, p, size, 8, ex,
			func(x, y int) byte {
				pp1 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x, y), p))
				pp2 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+1, y), p))
				pp3 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+2, y), p))
				pp4 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+3, y), p))
				pp5 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+4, y), p))
				pp6 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+5, y), p))
				pp7 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+6, y), p))
				pp8 := ex.SwapInk(pal.PalettePositionOrDefault(in.At(x+7, y), p))
				if ex.ScrCfg.Process.OneRow {
					pp2, pp4, pp6, pp8 = 0, 0, 0, 0
				}
				return pixel.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
			},
			func() byte { return pixel.PixelMode2(0, 0, 0, 0, 0, 0, 0, 0) },
		)
	default:
		return data, firmwareColorUsed, lineSize, errors.ErrorModeNotFound
	}
	return data, firmwareColorUsed, lineSize, nil
}
