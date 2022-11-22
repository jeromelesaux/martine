package screen

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/address"
	"github.com/jeromelesaux/martine/convert/export"
	"github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
)

func ToMode2(in *image.NRGBA, p color.Palette, ex *config.MartineConfig) []byte {
	var bw []byte

	lineToAdd := 1

	if ex.OneRow {
		lineToAdd = 2
	}

	if ex.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int)
	//fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	//fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {

			c1 := in.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			pp1 = ex.SwapInk(pp1)
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			pp2 = ex.SwapInk(pp2)
			firmwareColorUsed[pp2]++
			c3 := in.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			pp3 = ex.SwapInk(pp3)
			firmwareColorUsed[pp3]++
			c4 := in.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			pp4 = ex.SwapInk(pp4)
			firmwareColorUsed[pp4]++
			c5 := in.At(x+4, y)
			pp5, err := palette.PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
				pp5 = 0
			}
			pp5 = ex.SwapInk(pp5)
			firmwareColorUsed[pp5]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c6 := in.At(x+5, y)
			pp6, err := palette.PalettePosition(c6, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
				pp6 = 0
			}
			pp6 = ex.SwapInk(pp6)
			firmwareColorUsed[pp6]++
			c7 := in.At(x+6, y)
			pp7, err := palette.PalettePosition(c7, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
				pp7 = 0
			}
			pp7 = ex.SwapInk(pp7)
			firmwareColorUsed[pp7]++
			c8 := in.At(x+7, y)
			pp8, err := palette.PalettePosition(c8, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+7, y)
				pp8 = 0
			}
			pp8 = ex.SwapInk(pp8)
			firmwareColorUsed[pp8]++
			if ex.OneLine {
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
			addr := address.CpcScreenAddress(0, x, y, 2, ex.Overscan)
			bw[addr] = pixel
		}

	}

	//fmt.Println(firmwareColorUsed)
	return bw
}

func ToMode1(in *image.NRGBA, p color.Palette, ex *config.MartineConfig) []byte {
	var bw []byte

	lineToAdd := 1

	if ex.OneRow {
		lineToAdd = 2
	}
	if ex.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}

	firmwareColorUsed := make(map[int]int)
	//fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	//fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {

			c1 := in.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			pp1 = ex.SwapInk(pp1)
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			pp2 = ex.SwapInk(pp2)
			firmwareColorUsed[pp2]++
			c3 := in.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			pp3 = ex.SwapInk(pp3)
			firmwareColorUsed[pp3]++
			c4 := in.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			pp4 = ex.SwapInk(pp4)
			firmwareColorUsed[pp4]++
			if ex.OneLine {
				pp4 = 0
				pp2 = 0
			}
			pixel := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcScreenAddress(0, x, y, 1, ex.Overscan)
			bw[addr] = pixel
		}
	}
	return bw
}

func ToMode0(in *image.NRGBA, p color.Palette, ex *config.MartineConfig) []byte {
	var bw []byte

	lineToAdd := 1
	if ex.OneRow {
		lineToAdd = 2
	}
	if ex.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int)
	//fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	//fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y += lineToAdd {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {

			c1 := in.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			pp1 = ex.SwapInk(pp1)
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			pp2 = ex.SwapInk(pp2)
			firmwareColorUsed[pp2]++
			if ex.OneLine {
				pp2 = 0
			}
			pixel := pixel.PixelMode0(pp1, pp2)
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcScreenAddress(0, x, y, 0, ex.Overscan)
			bw[addr] = pixel
		}
	}

	//	fmt.Println(firmwareColorUsed)
	return bw
}

func ToMode0AndExport(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, cfg *config.MartineConfig) error {
	bw := ToMode0(in, p, cfg)
	return export.Export(filePath, bw, p, 0, cfg)
}

func ToMode1AndExport(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, cfg *config.MartineConfig) error {
	bw := ToMode1(in, p, cfg)
	return export.Export(filePath, bw, p, 1, cfg)
}

func ToMode2AndExport(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, cfg *config.MartineConfig) error {
	bw := ToMode2(in, p, cfg)
	return export.Export(filePath, bw, p, 2, cfg)
}

// scrRawToImg will convert the classical OCP screen slice of bytes  into image.NRGBA structure
// using the mode and the palette as arguments
func ScrRawToImg(d []byte, mode uint8, p color.Palette) (*image.NRGBA, error) {
	m := constants.NewSizeMode(mode, false)

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
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
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
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
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
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

// SrcToImg load the amstrad classical 17ko  screen image to image.NRBGA
// using the mode and palette as arguments
func ScrToImg(scrPath string, mode uint8, p color.Palette) (*image.NRGBA, error) {
	m := constants.NewSizeMode(mode, false)

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	d, err := ocpartstudio.RawScr(scrPath)
	if err != nil {
		return nil, err
	}
	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
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
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
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
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
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

func ScrToPng(scrPath string, output string, mode uint8, p color.Palette) error {

	out, err := ScrToImg(scrPath, mode, p)
	if err != nil {
		return err
	}
	return png.Png(output, out)
}
