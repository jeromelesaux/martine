package overscan

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/address"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
)

func ToGo(data []byte, screenMode uint8, p color.Palette, isCpcPlus bool) ([]byte, []byte, error) {
	orig, err := OverscanRawToImg(data, screenMode, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while converting into image  error :%v", err)
		return nil, nil, err
	}

	imgUp, imgDown, err := ci.SplitImage(orig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while splitting image  error :%v", err)
		return nil, nil, err
	}
	config := config.NewMartineConfig("", "")
	config.Size = constants.Size{Width: imgUp.Bounds().Max.X, Height: imgUp.Bounds().Max.Y}
	config.Overscan = true
	config.CpcPlus = isCpcPlus
	var dataUp, dataDown []byte
	switch screenMode {
	case 0:
		dataUp = ToMode0(imgUp, p, config)
		dataDown = ToMode0(imgDown, p, config)

	case 1:
		dataUp = ToMode1(imgUp, p, config)
		dataDown = ToMode1(imgDown, p, config)
	case 2:
		dataUp = ToMode2(imgUp, p, config)
		dataDown = ToMode2(imgDown, p, config)
	}
	return dataUp, dataDown, nil
}

func ToMode2(in *image.NRGBA, p color.Palette, cfg *config.MartineConfig) []byte {
	bw := make([]byte, 0x4000)

	lineToAdd := 1

	if cfg.OneRow {
		lineToAdd = 2
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
			pp1 = cfg.SwapInk(pp1)
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			pp2 = cfg.SwapInk(pp2)
			firmwareColorUsed[pp2]++
			c3 := in.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			pp3 = cfg.SwapInk(pp3)
			firmwareColorUsed[pp3]++
			c4 := in.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			pp4 = cfg.SwapInk(pp4)
			firmwareColorUsed[pp4]++
			c5 := in.At(x+4, y)
			pp5, err := palette.PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
				pp5 = 0
			}
			pp5 = cfg.SwapInk(pp5)
			firmwareColorUsed[pp5]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c6 := in.At(x+5, y)
			pp6, err := palette.PalettePosition(c6, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
				pp6 = 0
			}
			pp6 = cfg.SwapInk(pp6)
			firmwareColorUsed[pp6]++
			c7 := in.At(x+6, y)
			pp7, err := palette.PalettePosition(c7, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
				pp7 = 0
			}
			pp7 = cfg.SwapInk(pp7)
			firmwareColorUsed[pp7]++
			c8 := in.At(x+7, y)
			pp8, err := palette.PalettePosition(c8, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+7, y)
				pp8 = 0
			}
			pp8 = cfg.SwapInk(pp8)
			firmwareColorUsed[pp8]++
			if cfg.OneLine {
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
			addr := address.CpcOverscanSplitScreenAddress(0, x, y, 2, cfg.Overscan)
			bw[addr] = pixel
		}

	}

	//fmt.Println(firmwareColorUsed)
	return bw
}

func ToMode1(in *image.NRGBA, p color.Palette, cfg *config.MartineConfig) []byte {
	bw := make([]byte, 0x4000)

	lineToAdd := 1

	if cfg.OneRow {
		lineToAdd = 2
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
			pp1 = cfg.SwapInk(pp1)
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			pp2 = cfg.SwapInk(pp2)
			firmwareColorUsed[pp2]++
			c3 := in.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			pp3 = cfg.SwapInk(pp3)
			firmwareColorUsed[pp3]++
			c4 := in.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			pp4 = cfg.SwapInk(pp4)
			firmwareColorUsed[pp4]++
			if cfg.OneLine {
				pp4 = 0
				pp2 = 0
			}
			pixel := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcOverscanSplitScreenAddress(0, x, y, 1, cfg.Overscan)
			bw[addr] = pixel
		}
	}
	return bw
}

func ToMode0(in *image.NRGBA, p color.Palette, cfg *config.MartineConfig) []byte {
	bw := make([]byte, 0x4000)

	lineToAdd := 1
	if cfg.OneRow {
		lineToAdd = 2
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
			pp1 = cfg.SwapInk(pp1)
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			pp2 = cfg.SwapInk(pp2)
			firmwareColorUsed[pp2]++
			if cfg.OneLine {
				pp2 = 0
			}
			pixel := pixel.PixelMode0(pp1, pp2)
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			addr := address.CpcOverscanSplitScreenAddress(0, x, y, 0, cfg.Overscan)
			bw[addr] = pixel
		}
	}

	//	fmt.Println(firmwareColorUsed)
	return bw
}
