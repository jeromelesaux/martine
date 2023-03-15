package gfx

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/export"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/convert/screen"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/log"
)

func Transform(in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	filepath string,
	cfg *config.MartineConfig) error {
	switch size {
	case constants.Mode0:
		return export.ToMode0AndExport(in, p, size, filepath, cfg)
	case constants.Mode1:
		return export.ToMode1AndExport(in, p, size, filepath, cfg)
	case constants.Mode2:
		return export.ToMode2AndExport(in, p, size, filepath, cfg)
	case constants.OverscanMode0:
		return export.ToMode0AndExport(in, p, size, filepath, cfg)
	case constants.OverscanMode1:
		return export.ToMode1AndExport(in, p, size, filepath, cfg)
	case constants.OverscanMode2:
		return export.ToMode2AndExport(in, p, size, filepath, cfg)
	default:
		return errors.ErrorNotYetImplemented
	}
}

func InternalTransform(
	in *image.NRGBA,
	p color.Palette,
	size constants.Size,
	cfg *config.MartineConfig) []byte {

	switch size {
	case constants.Mode0:
		return screen.ToMode0(in, p, cfg)
	case constants.Mode1:
		return screen.ToMode1(in, p, cfg)
	case constants.Mode2:
		return screen.ToMode2(in, p, cfg)
	case constants.OverscanMode0:
		return screen.ToMode0(in, p, cfg)
	case constants.OverscanMode1:
		return screen.ToMode1(in, p, cfg)
	case constants.OverscanMode2:
		return screen.ToMode2(in, p, cfg)
	default:
		return []byte{}
	}
}

func revertColor(rawColor uint8, index int, isPlus bool) color.Color {
	var newColor color.Color
	var err error
	if !isPlus {
		newColor, err = constants.ColorFromHardware(rawColor)
		if err != nil {
			log.GetLogger().Error("No color found in data at index %d\n", index)
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
			p1, p2 := pixel.RawPixelMode0(byte(val))
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
			p1, p2, p3, p4 := pixel.RawPixelMode1(byte(val))
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
			p1, p2, p3, p4, p5, p6, p7, p8 := pixel.RawPixelMode2(byte(val))
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
