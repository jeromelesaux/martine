package gfx

import (
	"errors"
	"github.com/jeromelesaux/martine/constants"
	xl "github.com/jeromelesaux/martine/export/file"
	"image"
	"image/color"
)

var (
	UNDEFINED_MODE = errors.New("Undefined mode")
)

func SpriteToPng(winPath string, output string, p color.Palette) error {
	footer, err := xl.OpenWin(winPath)
	if err != nil {
		return err
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(footer.Width), Y: int(footer.Height)}})
	d, err := xl.RawWin(winPath)
	if err != nil {
		return err
	}
	for x := 0; x < int(footer.Height); x++ {
		for y := 0; y < int(footer.Width); y++ {
			i := d[x+y]
			c := p[i]
			out.Set(x, y, c)
		}
	}
	return xl.Png(output, out)
}

func ScrToPng(scrPath string, output string, mode uint8, p color.Palette) error {
	var m constants.Size
	switch mode {
	case 0:
		m = constants.Mode0
	case 1:
		m = constants.Mode1
	case 2:
		m = constants.Mode2
	default:
		return UNDEFINED_MODE
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Height), Y: int(m.Width)}})

	d, err := xl.RawScr(scrPath)
	if err != nil {
		return err
	}
	cpcRow := 0
	for y := 0; y < m.Height; y++ {
		cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
		for x := 0; x < m.Width; x += 2 {
			val := d[cpcLine+cpcRow]
			pp1, pp2 := rawPixelMode0(val)
			c1 := p[pp1]
			c2 := p[pp2]

			out.Set(x, y, c1)
			out.Set(x+1, y, c2)
			cpcRow++
		}
		cpcRow = 0
	}

	return xl.Png(output, out)
}
