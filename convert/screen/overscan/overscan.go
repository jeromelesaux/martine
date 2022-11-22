package overscan

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/png"
)

// overscanRawToImg will convert fullscreen amstrad screen slice of bytes in image.NRGBA
// using the  screen mode  and the palette as arguments
func OverscanRawToImg(d []byte, mode uint8, p color.Palette) (*image.NRGBA, error) {
	m := constants.NewSizeMode(mode, true)

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 2 {
				val := d[cpcLine+cpcRow]
				pp1, pp2 := pixel.RawPixelMode0(val)
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
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := pixel.RawPixelMode1(val)
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
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := pixel.RawPixelMode2(val)
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
	return out, nil
}

// overscanToImg will convert fullscreen amstrad screen filepath 'srcPath' in image.NRGBA
// using the  screen mode  and the palette as arguments
func OverscanToImg(scrPath string, mode uint8, p color.Palette) (*image.NRGBA, error) {
	m := constants.NewSizeMode(mode, true)
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	d, err := overscan.RawOverscan(scrPath) // RawOverscan data commence en 0x30
	if err != nil {
		return nil, err
	}

	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 2 {
				val := d[cpcLine+cpcRow]
				pp1, pp2 := pixel.RawPixelMode0(val)
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
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := pixel.RawPixelMode1(val)
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
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := pixel.RawPixelMode2(val)
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
	return out, nil
}

// overscanToPng will convert fullscreen amstrad screen filepath in png file 'output'
// using the  screen mode  and the palette as arguments
func OverscanToPng(scrPath string, output string, mode uint8, p color.Palette) error {
	out, err := OverscanToImg(scrPath, mode, p)
	if err != nil {
		return err
	}
	return png.Png(output, out)
}
