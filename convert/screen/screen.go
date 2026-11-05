// Package screen provides CPC screen conversion utilities.
package screen

import (
	"errors"
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/address"

	"github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/png"
)

// ToMode2 converts an NRGBA image and palette into CPC screen mode 2 bytes.
// nolint: funlen, gocognit
func ToMode2(in *image.NRGBA, p color.Palette, cfg *config.MartineConfig) []byte {
	var bw []byte

	lineToAdd := 1

	if cfg.ScrCfg.Process.OneRow {
		lineToAdd = 2
	}

	if cfg.ScrCfg.Type.IsFullScreen() {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {
			pp1 := palette.PalettePositionOrDefault(in.At(x, y), p)
			pp1 = cfg.SwapInk(pp1)
			pp2 := palette.PalettePositionOrDefault(in.At(x+1, y), p)
			pp2 = cfg.SwapInk(pp2)
			pp3 := palette.PalettePositionOrDefault(in.At(x+2, y), p)
			pp3 = cfg.SwapInk(pp3)
			pp4 := palette.PalettePositionOrDefault(in.At(x+3, y), p)
			pp4 = cfg.SwapInk(pp4)
			pp5 := palette.PalettePositionOrDefault(in.At(x+4, y), p)
			pp5 = cfg.SwapInk(pp5)
			pp6 := palette.PalettePositionOrDefault(in.At(x+5, y), p)
			pp6 = cfg.SwapInk(pp6)
			pp7 := palette.PalettePositionOrDefault(in.At(x+6, y), p)
			pp7 = cfg.SwapInk(pp7)
			pp8 := palette.PalettePositionOrDefault(in.At(x+7, y), p)
			pp8 = cfg.SwapInk(pp8)
			if cfg.ScrCfg.Process.OneLine {
				pp2 = 0
				pp4 = 0
				pp6 = 0
				pp8 = 0
			}
			pix := pixel.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcScreenAddress(0, x, y, 2, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}

	}

	return bw
}

// ToMode1 converts an NRGBA image and palette into CPC screen mode 1 bytes.
// nolint: funlen
func ToMode1(in *image.NRGBA, p color.Palette, cfg *config.MartineConfig) []byte {
	var bw []byte

	lineToAdd := 1

	if cfg.ScrCfg.Process.OneRow {
		lineToAdd = 2
	}
	if cfg.ScrCfg.Type.IsFullScreen() {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {
			pp1 := palette.PalettePositionOrDefault(in.At(x, y), p)
			pp1 = cfg.SwapInk(pp1)
			pp2 := palette.PalettePositionOrDefault(in.At(x+1, y), p)
			pp2 = cfg.SwapInk(pp2)
			pp3 := palette.PalettePositionOrDefault(in.At(x+2, y), p)
			pp3 = cfg.SwapInk(pp3)
			pp4 := palette.PalettePositionOrDefault(in.At(x+3, y), p)
			pp4 = cfg.SwapInk(pp4)
			if cfg.ScrCfg.Process.OneLine {
				pp4 = 0
				pp2 = 0
			}
			pix := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcScreenAddress(0, x, y, 1, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}
	}
	return bw
}

func ToMode0(in *image.NRGBA, p color.Palette, cfg *config.MartineConfig) []byte {
	var bw []byte

	lineToAdd := 1
	if cfg.ScrCfg.Process.OneRow {
		lineToAdd = 2
	}
	if cfg.ScrCfg.Type.IsFullScreen() {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {
			pp1 := palette.PalettePositionOrDefault(in.At(x, y), p)
			pp1 = cfg.SwapInk(pp1)
			pp2 := palette.PalettePositionOrDefault(in.At(x+1, y), p)
			pp2 = cfg.SwapInk(pp2)
			if cfg.ScrCfg.Process.OneLine {
				pp2 = 0
			}
			pix := pixel.PixelMode0(pp1, pp2)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcScreenAddress(0, x, y, 0, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}
	}

	return bw
}

// ScrRawToImg converts the classical OCP screen slice of bytes into image.NRGBA.
// using the mode and the palette as arguments
func ScrRawToImg(d []byte, mode uint8, p color.Palette) (*image.NRGBA, error) {
	m := constants.NewSizeMode(mode, false)

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: m.Width, Y: m.Height}})

	switch mode {
	case 0:
		out = setImageMode0(out, p, d, m.Width, m.Height)
	case 1:
		out = setImageMode1(out, p, d, m.Width, m.Height)
	case 2:
		out = setImageMode2(out, p, d, m.Width, m.Height)

	}
	return out, nil
}

func setImage(out *image.NRGBA, p color.Palette, d []byte, width, height, pixelsPerByte int, decode func(byte) []int) *image.NRGBA {
	cpcRow := 0
	for y := 0; y < height; y++ {
		cpcLine := (y/0x8)*0x50 + (y%0x8)*0x800
		for x := 0; x < width; x += pixelsPerByte {
			indices := decode(d[cpcLine+cpcRow])
			for i, idx := range indices {
				out.Set(x+i, y, p[idx])
			}
			cpcRow++
		}
		cpcRow = 0
	}
	return out
}

func setImageMode0(out *image.NRGBA, p color.Palette, d []byte, width, height int) *image.NRGBA {
	return setImage(out, p, d, width, height, 2, func(val byte) []int {
		pp1, pp2 := pixel.RawPixelMode0(val)
		return []int{pp1, pp2}
	})
}

func setImageMode1(out *image.NRGBA, p color.Palette, d []byte, width, height int) *image.NRGBA {
	return setImage(out, p, d, width, height, 4, func(val byte) []int {
		pp1, pp2, pp3, pp4 := pixel.RawPixelMode1(val)
		return []int{pp1, pp2, pp3, pp4}
	})
}

func setImageMode2(out *image.NRGBA, p color.Palette, d []byte, width, height int) *image.NRGBA {
	return setImage(out, p, d, width, height, 8, func(val byte) []int {
		pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := pixel.RawPixelMode2(val)
		return []int{pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8}
	})
}

func WinToImg(path string, mode uint8, p color.Palette) (*image.NRGBA, error) {
	footer, err := window.OpenWin(path)
	if err != nil {
		return nil, err
	}
	data, err := window.RawWin(path)
	if err != nil {
		return nil, err
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(footer.Width), Y: int(footer.Height)}})
	for w := 0; w < int(footer.Width); w++ {
		for h := 0; h < int(footer.Height); h++ {
			d := data[h+(w*int(footer.Width))]
			out.Set(w, h, p[d])
		}
	}
	return out, nil
}

// ScrToImg loads the amstrad classical 17ko screen image into image.NRGBA.
// using the mode and palette as arguments
func ScrToImg(scrPath string, mode uint8, p color.Palette) (*image.NRGBA, error) {
	m := constants.NewSizeMode(mode, false)

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: m.Width, Y: m.Height}})

	d, err := ocpartstudio.RawScr(scrPath)
	if err != nil {
		return nil, err
	}
	if len(d) < (m.Width * m.Height) {
		return nil, errors.New("bad format")
	}
	switch mode {
	case 0:
		out = setImageMode0(out, p, d, m.Width, m.Height)
	case 1:
		out = setImageMode1(out, p, d, m.Width, m.Height)
	case 2:
		out = setImageMode2(out, p, d, m.Width, m.Height)
	}
	return out, nil
}

func ScrToPng(scrPath, output string, mode uint8, p color.Palette) error {

	out, err := ScrToImg(scrPath, mode, p)
	if err != nil {
		return err
	}
	return png.Png(output, out)
}
