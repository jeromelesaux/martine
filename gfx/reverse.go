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
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	d, err := xl.RawScr(scrPath)
	if err != nil {
		return err
	}
	cpcRow := 0
	switch mode {
	case 0:
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
	case 1:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := rawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				cpcRow++
			}
			cpcRow = 0
		}
	case 2:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := rawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				out.Set(x+4, y, c5)
				out.Set(x+5, y, c6)
				out.Set(x+6, y, c7)
				out.Set(x+7, y, c8)
				cpcRow++
			}
			cpcRow = 0
		}
	}

	return xl.Png(output, out)
}

func OverscanToPng(scrPath string, output string, mode uint8, p color.Palette) error {
	var m constants.Size
	switch mode {
	case 0:
		m = constants.OverscanMode0
	case 1:
		m = constants.OverscanMode1
	case 2:
		m = constants.OverscanMode2
	default:
		return UNDEFINED_MODE
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	d, err := xl.RawOverscan(scrPath) // RawOverscan data commence en 0x30
	if err != nil {
		return err
	}
	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 135 {
				cpcLine += (0x3800)
			}
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
	case 1:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 135 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := rawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				cpcRow++
			}
			cpcRow = 0
		}
	case 2:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if cpcLine > 135 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := rawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				out.Set(x+4, y, c5)
				out.Set(x+5, y, c6)
				out.Set(x+6, y, c7)
				out.Set(x+7, y, c8)
				cpcRow++
			}
			cpcRow = 0
		}
	}

	return xl.Png(output, out)
}
