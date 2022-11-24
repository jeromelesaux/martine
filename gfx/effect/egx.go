package effect

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/address"
	"github.com/jeromelesaux/martine/convert/export"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/convert/screen"
	"github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func Egx(filepath1, filepath2 string, p color.Palette, m1, m2 int, cfg *config.MartineConfig) error {
	p = constants.SortColorsByDistance(p)
	if m1 == 0 && m2 == 1 || m2 == 0 && m1 == 1 {
		var f0, f1 string
		var mode0, mode1 uint8
		var err error
		mode0 = 0
		mode1 = 1
		if m1 == 0 {
			f0 = filepath1
			f1 = filepath2
		} else {
			f0 = filepath2
			f1 = filepath1
		}
		var in0, in1 *image.NRGBA
		if cfg.Overscan {
			in0, err = overscan.OverscanToImg(f0, mode0, p)
			if err != nil {
				return err
			}
			in1, err = overscan.OverscanToImg(f1, mode1, p)
			if err != nil {
				return err
			}
		} else {
			in0, err = screen.ScrToImg(f0, 0, p)
			if err != nil {
				return err
			}
			in1, err = screen.ScrToImg(f1, 1, p)
			if err != nil {
				return err
			}
		}
		if err = ToEgx1(in0, in1, p, uint8(m1), "egx.scr", cfg); err != nil {
			return err
		}
		if !cfg.Overscan {
			if err = ocpartstudio.EgxLoader("egx.scr", p, uint8(m1), uint8(m2), cfg); err != nil {
				return err
			}
		}
		return nil
	} else {
		if m1 == 1 && m2 == 2 || m2 == 1 && m1 == 2 {
			var f2, f1 string
			var mode2, mode1 uint8
			var err error
			mode1 = 1
			mode2 = 2
			if m1 == 1 {
				//filename := filepath.Base(filepath1)
				//filePath = exportType.OutputPath + string(filepath.Separator) + filename

				f1 = filepath1
				f2 = filepath2
			} else {
				//filename := filepath.Base(filepath2)
				//filePath = exportType.OutputPath + string(filepath.Separator) + filename
				f1 = filepath2
				f2 = filepath1
			}
			var in2, in1 *image.NRGBA
			if cfg.Overscan {
				in1, err = overscan.OverscanToImg(f1, mode1, p)
				if err != nil {
					return err
				}
				in2, err = overscan.OverscanToImg(f2, mode2, p)
				if err != nil {
					return err
				}
			} else {
				in1, err = screen.ScrToImg(f1, mode1, p)
				if err != nil {
					return err
				}
				in2, err = screen.ScrToImg(f2, mode2, p)
				if err != nil {
					return err
				}
			}
			if err = ToEgx2(in1, in2, p, uint8(m1), "egx.scr", cfg); err != nil {
				return err
			}
			if !cfg.Overscan {
				if err = ocpartstudio.EgxLoader("egx.scr", p, uint8(m1), uint8(m2), cfg); err != nil {
					return err
				}
			}
		} else {
			return errors.ErrorFeatureNotImplemented
		}
	}

	return nil
}

func AutoEgx1(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string) error {
	var err error

	size := constants.Size{
		Width:  cfg.Size.Width,
		Height: cfg.Size.Height}

	im := ci.Resize(in, size, cfg.ResizingAlgo)
	var palette color.Palette // palette de l'image
	var p color.Palette       // palette cpc de l'image
	var downgraded *image.NRGBA

	if cfg.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cfg.PalettePath)
		palette, _, err = ocpartstudio.OpenPal(cfg.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cfg.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if len(palette) > 0 {
		p, downgraded = ci.DowngradingWithPalette(im, palette)
	} else {
		p, downgraded, err = ci.DowngradingPalette(im, cfg.Size, cfg.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	p = constants.SortColorsByDistance(p)
	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	downgraded, p = gfx.DoDithering(downgraded, p, cfg.DitheringAlgo, cfg.DitheringType, cfg.DitheringWithQuantification, cfg.DitheringMatrix, float32(cfg.DitheringMultiplier), cfg.CpcPlus, cfg.Size)

	return ToEgx1(downgraded, downgraded, p, 0, picturePath, cfg)
}

func AutoEgx2(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string) error {
	var err error

	size := constants.Size{
		Width:  cfg.Size.Width,
		Height: cfg.Size.Height}

	im := ci.Resize(in, size, cfg.ResizingAlgo)
	var palette color.Palette // palette de l'image
	var p color.Palette       // palette cpc de l'image
	var downgraded *image.NRGBA

	if cfg.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cfg.PalettePath)
		palette, _, err = ocpartstudio.OpenPal(cfg.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cfg.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	if len(palette) > 0 {
		p, downgraded = ci.DowngradingWithPalette(im, palette)
	} else {
		p, downgraded, err = ci.DowngradingPalette(im, cfg.Size, cfg.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	p = constants.SortColorsByDistance(p)
	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	downgraded, p = gfx.DoDithering(downgraded, p, cfg.DitheringAlgo, cfg.DitheringType, cfg.DitheringWithQuantification, cfg.DitheringMatrix, float32(cfg.DitheringMultiplier), cfg.CpcPlus, cfg.Size)

	return ToEgx2(downgraded, downgraded, p, 1, picturePath, cfg)
}

func ToEgx1(inMode0, inMode1 *image.NRGBA, p color.Palette, firstLineMode uint8, picturePath string, cfg *config.MartineConfig) error {
	bw, p := ToEgx1Raw(inMode0, inMode1, p, firstLineMode, cfg)
	return export.Export(picturePath, bw, p, 1, cfg)
}

func ToEgx1Raw(inMode0, inMode1 *image.NRGBA, p color.Palette, firstLineMode uint8, cfg *config.MartineConfig) ([]byte, color.Palette) {
	var bw []byte
	if cfg.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int)
	mode0Line := 1
	mode1Line := 0
	if firstLineMode == 1 {
		mode0Line = 0
		mode1Line = 1
	}

	for y := inMode0.Bounds().Min.Y + mode0Line; y < inMode0.Bounds().Max.Y; y += 2 {
		for x := inMode0.Bounds().Min.X; x < inMode0.Bounds().Max.X; x += 2 {
			c1 := inMode0.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode0.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			pixel := pixel.PixelMode0(pp1, pp2)
			addr := address.CpcScreenAddress(0, x, y, 0, cfg.Overscan, cfg.DoubleScreenAddress)
			bw[addr] = pixel
		}
	}
	for y := inMode1.Bounds().Min.Y + mode1Line; y < inMode1.Bounds().Max.Y; y += 2 {
		for x := inMode1.Bounds().Min.X; x < inMode1.Bounds().Max.X; x += 4 {
			c1 := inMode1.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode1.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := inMode1.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := inMode1.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++

			pixel := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			addr := address.CpcScreenAddress(0, x, y, 1, cfg.Overscan, cfg.DoubleScreenAddress)
			bw[addr] = pixel
			addr = address.CpcScreenAddress(0, x+1, y, 1, cfg.Overscan, cfg.DoubleScreenAddress)
			bw[addr] = pixel
		}
	}
	return bw, p
}

func ToEgx2(inMode1, inMode2 *image.NRGBA, p color.Palette, firstLineMode uint8, picturePath string, cfg *config.MartineConfig) error {
	bw, p := ToEgx2Raw(inMode1, inMode2, p, firstLineMode, cfg)
	return export.Export(picturePath, bw, p, 2, cfg)
}

func ToEgx2Raw(inMode1, inMode2 *image.NRGBA, p color.Palette, firstLineMode uint8, cfg *config.MartineConfig) ([]byte, color.Palette) {
	var bw []byte
	if cfg.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int)
	mode1Line := 1
	mode2Line := 0
	if firstLineMode == 2 {
		mode1Line = 0
		mode2Line = 1
	}

	for y := inMode1.Bounds().Min.Y + mode1Line; y < inMode1.Bounds().Max.Y; y += 2 {
		for x := inMode1.Bounds().Min.X; x < inMode1.Bounds().Max.X; x += 4 {
			c1 := inMode1.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode1.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := inMode1.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := inMode1.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++

			pixel := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			addr := address.CpcScreenAddress(0, x, y, 1, cfg.Overscan, cfg.DoubleScreenAddress)
			bw[addr] = pixel
		}
	}
	for y := inMode2.Bounds().Min.Y + mode2Line; y < inMode2.Bounds().Max.Y; y += 2 {
		for x := inMode2.Bounds().Min.X; x < inMode2.Bounds().Max.X; x += 8 {
			c1 := inMode2.At(x, y)
			pp1, err := palette.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode2.At(x+1, y)
			pp2, err := palette.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := inMode2.At(x+2, y)
			pp3, err := palette.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := inMode2.At(x+3, y)
			pp4, err := palette.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++
			c5 := inMode2.At(x+4, y)
			pp5, err := palette.PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
				pp5 = 0
			}
			firmwareColorUsed[pp5]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c6 := inMode2.At(x+5, y)
			pp6, err := palette.PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
				pp6 = 0
			}
			firmwareColorUsed[pp6]++
			c7 := inMode2.At(x+6, y)
			pp7, err := palette.PalettePosition(c7, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
				pp7 = 0
			}
			firmwareColorUsed[pp7]++
			c8 := inMode2.At(x+7, y)
			pp8, err := palette.PalettePosition(c8, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+3, y)
				pp8 = 0
			}
			firmwareColorUsed[pp8]++
			pixel := pixel.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
			addr := address.CpcScreenAddress(0, x, y, 2, cfg.Overscan, cfg.DoubleScreenAddress)
			bw[addr] = pixel
			addr = address.CpcScreenAddress(0, x+1, y, 2, cfg.Overscan, cfg.DoubleScreenAddress)
			bw[addr] = pixel
		}
	}
	return bw, p
}

func EgxRaw(img1, img2 []byte, p color.Palette, mode1, mode2 int, cfg *config.MartineConfig) ([]byte, color.Palette, int, error) {
	p = constants.SortColorsByDistance(p)
	if mode1 == 0 && mode2 == 1 || mode2 == 0 && mode1 == 1 {
		var f0, f1 []byte
		var mode0, mode1 uint8
		var err error
		mode0 = 0
		mode1 = 1
		if mode1 == 0 {
			f0 = img1
			f1 = img2
		} else {
			f0 = img2
			f1 = img1
		}
		var in0, in1 *image.NRGBA
		if cfg.Overscan {
			in0, err = overscan.OverscanRawToImg(f0, mode0, p)
			if err != nil {
				return nil, p, 0, err
			}
			in1, err = overscan.OverscanRawToImg(f1, mode1, p)
			if err != nil {
				return nil, p, 0, err
			}
		} else {
			in0, err = screen.ScrRawToImg(f0, 0, p)
			if err != nil {
				return nil, p, 0, err
			}
			in1, err = screen.ScrRawToImg(f1, 1, p)
			if err != nil {
				return nil, p, 0, err
			}
		}
		res, p := ToEgx1Raw(in0, in1, p, uint8(mode1), cfg)
		return res, p, 1, nil
	} else {
		if mode1 == 1 && mode2 == 2 || mode2 == 1 && mode1 == 2 {
			var f2, f1 []byte
			var mode2, mode1 uint8
			var err error
			mode1 = 1
			mode2 = 2
			if mode1 == 1 {
				//filename := filepath.Base(filepath1)
				//filePath = exportType.OutputPath + string(filepath.Separator) + filename

				f1 = img1
				f2 = img2
			} else {
				//filename := filepath.Base(filepath2)
				//filePath = exportType.OutputPath + string(filepath.Separator) + filename
				f1 = img2
				f2 = img1
			}
			var in2, in1 *image.NRGBA
			if cfg.Overscan {
				in1, err = overscan.OverscanRawToImg(f1, mode1, p)
				if err != nil {
					return nil, p, 0, err
				}
				in2, err = overscan.OverscanRawToImg(f2, mode2, p)
				if err != nil {
					return nil, p, 0, err
				}
			} else {
				in1, err = screen.ScrRawToImg(f1, mode1, p)
				if err != nil {
					return nil, p, 0, err
				}
				in2, err = screen.ScrRawToImg(f2, mode2, p)
				if err != nil {
					return nil, p, 0, err
				}
			}
			res, p := ToEgx2Raw(in1, in2, p, uint8(mode1), cfg)
			return res, p, 2, err

		} else {
			return nil, p, 0, errors.ErrorFeatureNotImplemented
		}
	}
}
