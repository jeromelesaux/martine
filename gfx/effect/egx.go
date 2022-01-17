package effect

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func Egx(filepath1, filepath2 string, p color.Palette, m1, m2 int, cont *export.MartineContext) error {
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
		if cont.Overscan {
			in0, err = common.OverscanToImg(f0, mode0, p)
			if err != nil {
				return err
			}
			in1, err = common.OverscanToImg(f1, mode1, p)
			if err != nil {
				return err
			}
		} else {
			in0, err = common.ScrToImg(f0, 0, p)
			if err != nil {
				return err
			}
			in1, err = common.ScrToImg(f1, 1, p)
			if err != nil {
				return err
			}
		}
		if err = ToEgx1(in0, in1, p, uint8(m1), "egx.scr", cont); err != nil {
			return err
		}
		if !cont.Overscan {
			if err = file.EgxLoader("egx.scr", p, uint8(m1), uint8(m2), cont); err != nil {
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
			if cont.Overscan {
				in1, err = common.OverscanToImg(f1, mode1, p)
				if err != nil {
					return err
				}
				in2, err = common.OverscanToImg(f2, mode2, p)
				if err != nil {
					return err
				}
			} else {
				in1, err = common.ScrToImg(f1, mode1, p)
				if err != nil {
					return err
				}
				in2, err = common.ScrToImg(f2, mode2, p)
				if err != nil {
					return err
				}
			}
			if err = ToEgx2(in1, in2, p, uint8(m1), "egx.scr", cont); err != nil {
				return err
			}
			if !cont.Overscan {
				if err = file.EgxLoader("egx.scr", p, uint8(m1), uint8(m2), cont); err != nil {
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
	cont *export.MartineContext,
	filename, picturePath string) error {
	var err error

	size := constants.Size{
		Width:  cont.Size.Width,
		Height: cont.Size.Height}

	im := convert.Resize(in, size, cont.ResizingAlgo)
	var palette color.Palette // palette de l'image
	var p color.Palette       // palette cpc de l'image
	var downgraded *image.NRGBA

	if cont.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cont.PalettePath)
		palette, _, err = file.OpenPal(cont.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if len(palette) > 0 {
		p, downgraded = convert.DowngradingWithPalette(im, palette)
	} else {
		p, downgraded, err = convert.DowngradingPalette(im, cont.Size, cont.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	p = constants.SortColorsByDistance(p)
	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := file.Png(filepath.Join(cont.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	downgraded, p = gfx.DoDithering(downgraded, p, cont)

	return ToEgx1(downgraded, downgraded, p, 0, picturePath, cont)
}

func AutoEgx2(in image.Image,
	cont *export.MartineContext,
	filename, picturePath string) error {
	var err error

	size := constants.Size{
		Width:  cont.Size.Width,
		Height: cont.Size.Height}

	im := convert.Resize(in, size, cont.ResizingAlgo)
	var palette color.Palette // palette de l'image
	var p color.Palette       // palette cpc de l'image
	var downgraded *image.NRGBA

	if cont.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cont.PalettePath)
		palette, _, err = file.OpenPal(cont.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	if len(palette) > 0 {
		p, downgraded = convert.DowngradingWithPalette(im, palette)
	} else {
		p, downgraded, err = convert.DowngradingPalette(im, cont.Size, cont.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	p = constants.SortColorsByDistance(p)
	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := file.Png(filepath.Join(cont.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	downgraded, p = gfx.DoDithering(downgraded, p, cont)

	return ToEgx2(downgraded, downgraded, p, 1, picturePath, cont)
}

func ToEgx1(inMode0, inMode1 *image.NRGBA, p color.Palette, firstLineMode uint8, picturePath string, cont *export.MartineContext) error {
	bw, p := ToEgx1Memory(inMode0, inMode1, p, firstLineMode, cont)
	return common.Export(picturePath, bw, p, 1, cont)
}

func ToEgx1Memory(inMode0, inMode1 *image.NRGBA, p color.Palette, firstLineMode uint8, cont *export.MartineContext) ([]byte, color.Palette) {
	var bw []byte
	if cont.Overscan {
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
			pp1, err := common.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode0.At(x+1, y)
			pp2, err := common.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			pixel := common.PixelMode0(pp1, pp2)
			addr := common.CpcScreenAddress(0, x, y, 0, cont.Overscan)
			bw[addr] = pixel
		}
	}
	for y := inMode1.Bounds().Min.Y + mode1Line; y < inMode1.Bounds().Max.Y; y += 2 {
		for x := inMode1.Bounds().Min.X; x < inMode1.Bounds().Max.X; x += 4 {
			c1 := inMode1.At(x, y)
			pp1, err := common.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode1.At(x+1, y)
			pp2, err := common.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := inMode1.At(x+2, y)
			pp3, err := common.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := inMode1.At(x+3, y)
			pp4, err := common.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++

			pixel := common.PixelMode1(pp1, pp2, pp3, pp4)
			addr := common.CpcScreenAddress(0, x, y, 1, cont.Overscan)
			bw[addr] = pixel
			addr = common.CpcScreenAddress(0, x+1, y, 1, cont.Overscan)
			bw[addr] = pixel
		}
	}
	return bw, p
}

func ToEgx2(inMode1, inMode2 *image.NRGBA, p color.Palette, firstLineMode uint8, picturePath string, cont *export.MartineContext) error {
	bw, p := ToEgx2Memory(inMode1, inMode2, p, firstLineMode, cont)
	return common.Export(picturePath, bw, p, 2, cont)
}

func ToEgx2Memory(inMode1, inMode2 *image.NRGBA, p color.Palette, firstLineMode uint8, cont *export.MartineContext) ([]byte, color.Palette) {
	var bw []byte
	if cont.Overscan {
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
			pp1, err := common.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode1.At(x+1, y)
			pp2, err := common.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := inMode1.At(x+2, y)
			pp3, err := common.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := inMode1.At(x+3, y)
			pp4, err := common.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++

			pixel := common.PixelMode1(pp1, pp2, pp3, pp4)
			addr := common.CpcScreenAddress(0, x, y, 1, cont.Overscan)
			bw[addr] = pixel
		}
	}
	for y := inMode2.Bounds().Min.Y + mode2Line; y < inMode2.Bounds().Max.Y; y += 2 {
		for x := inMode2.Bounds().Min.X; x < inMode2.Bounds().Max.X; x += 8 {
			c1 := inMode2.At(x, y)
			pp1, err := common.PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := inMode2.At(x+1, y)
			pp2, err := common.PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}
			firmwareColorUsed[pp2]++
			c3 := inMode2.At(x+2, y)
			pp3, err := common.PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := inMode2.At(x+3, y)
			pp4, err := common.PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++
			c5 := inMode2.At(x+4, y)
			pp5, err := common.PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
				pp5 = 0
			}
			firmwareColorUsed[pp5]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c6 := inMode2.At(x+5, y)
			pp6, err := common.PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
				pp6 = 0
			}
			firmwareColorUsed[pp6]++
			c7 := inMode2.At(x+6, y)
			pp7, err := common.PalettePosition(c7, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
				pp7 = 0
			}
			firmwareColorUsed[pp7]++
			c8 := inMode2.At(x+7, y)
			pp8, err := common.PalettePosition(c8, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+3, y)
				pp8 = 0
			}
			firmwareColorUsed[pp8]++
			pixel := common.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
			addr := common.CpcScreenAddress(0, x, y, 2, cont.Overscan)
			bw[addr] = pixel
			addr = common.CpcScreenAddress(0, x+1, y, 2, cont.Overscan)
			bw[addr] = pixel
		}
	}
	return bw, p
}
