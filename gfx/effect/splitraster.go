package effect

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export/impdraw/splitraster"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func DoSpliteRaster(in image.Image, screenMode uint8, filename string, cfg *config.MartineConfig) error {

	var p color.Palette
	var bw []byte
	var rasters *constants.SplitRasterScreen
	var err error
	if !cfg.Overscan {
		return errors.ErrorNotYetImplemented
	}
	switch cfg.CpcPlus {
	case false:
		p, bw, rasters, err = ToSplitRasterCPCOld(in, screenMode, filename, cfg)
		if err != nil {
			return err
		}
	default:
		fmt.Fprintf(os.Stderr, "Not yet implemented.")
		return errors.ErrorNotYetImplemented
	}
	// export des données
	if err := common.Export(filename, bw, p, screenMode, cfg); err != nil {
		return err
	}
	return splitraster.ExportSplitRaster(filename, p, rasters, cfg)
}

func ToSplitRasterCPCOld(in image.Image, screenMode uint8, filename string, cfg *config.MartineConfig) (color.Palette, []byte, *constants.SplitRasterScreen, error) {

	var bw []byte
	srs := constants.NewSplitRasterScreen()
	out := convert.Resize(in, cfg.Size, cfg.ResizingAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_resized.png"), out); err != nil {
		return nil, bw, srs, err
	}
	p, newIm, err := convert.DowngradingPalette(out, cfg.Size, cfg.CpcPlus)
	if err != nil {
		return p, bw, srs, err
	}
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_downgraded.png"), newIm); err != nil {
		return nil, bw, srs, err
	}

	srIm := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: out.Bounds().Min.X, Y: out.Bounds().Min.Y},
		Max: image.Point{X: out.Bounds().Max.X, Y: out.Bounds().Max.Y}})

	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), newIm.Bounds().Max.X, newIm.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	if cfg.Overscan {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	firmwareColorUsed := make(map[int]int)
	backgroundColor := p[0]
	notSplitRaster := true
	for y := 0; y < cfg.Size.Height; y++ {
		for x := 0; x < cfg.Size.Width; {
			if x%16 == 0 {
				if !srs.IsFull() {
					notSplitRaster = false
					if isSplitRaster(newIm, x, y, 16) {
						srs = SetCpcOldSplitRaster(out, srIm, constants.CpcOldPalette, srs, x, y, 16)
					}
					pp, _ := common.PalettePosition(backgroundColor, p)
					fmt.Fprintf(os.Stdout, "X{%d,%d},Y{%d} might be a splitraster\n", x, (x + 16), y)
					switch screenMode {
					case 0:
						for i := 0; i < 16; {
							pixel := common.PixelMode0(pp, pp)
							addr := common.CpcScreenAddress(0, x+i, y, 0, cfg.Overscan)
							bw[addr] = pixel
							i += 2
							firmwareColorUsed[pp] += 2
						}
					case 1:
						for i := 0; i < 16; {
							pixel := common.PixelMode1(pp, pp, pp, pp)
							addr := common.CpcScreenAddress(0, x+i, y, 1, cfg.Overscan)
							bw[addr] = pixel
							i += 4
							firmwareColorUsed[pp] += 4
						}
					case 2:
						for i := 0; i < 16; {
							pixel := common.PixelMode2(pp, pp, pp, pp, pp, pp, pp, pp)
							addr := common.CpcScreenAddress(0, x+i, y, 2, cfg.Overscan)
							bw[addr] = pixel
							i += 8
							firmwareColorUsed[pp] += 8
						}
					}
					// ajout d'un split raster
					// modification de l'image destination pour utiliser celle du background
					// gestion des modes à faire
					x += 16
				} else {
					notSplitRaster = true
				}

			} else {
				notSplitRaster = true
			}
			if notSplitRaster {

				// traitement normal des pixels
				switch screenMode {
				case 0:
					bw, firmwareColorUsed = setPixelMode0(newIm, srIm, p, x, y, bw, firmwareColorUsed, cfg)
					x += 2
				case 1:
					bw, firmwareColorUsed = setPixelMode1(newIm, srIm, p, x, y, bw, firmwareColorUsed, cfg)
					x += 4
				case 2:
					bw, firmwareColorUsed = setPixelMode2(newIm, srIm, p, x, y, bw, firmwareColorUsed, cfg)
					x += 4
				}
			}
		}
	}
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_splitraster.png"), srIm); err != nil {
		return nil, bw, srs, err
	}
	fmt.Println(firmwareColorUsed)
	return p, bw, srs, nil
}

func SetCpcOldSplitRaster(in *image.NRGBA, out *image.NRGBA, p color.Palette, s *constants.SplitRasterScreen, pos, y, length int) *constants.SplitRasterScreen {
	occ := 0
	if !s.Add(constants.NewSpliteRaster(uint16(pos), length, occ)) {
		return s
	}
	for x := pos; x < pos+length && x < in.Bounds().Max.X; x++ {
		c := in.At(x, y)
		c2 := p.Convert(c) // find the color in palette old cpc
		out.Set(x, y, c2)
		hds, err := constants.HardwareValues(c2)
		if err != nil {
			r, g, b, _ := c.RGBA()
			fmt.Fprintf(os.Stderr, "not hardware value for color (%d,%d,%d)\n", r, g, b)
			continue
		}
		if !s.Values[len(s.Values)-1].Add(0, int(hds[0])) {
			continue
		}
		occ++
	}
	if occ < (length - 1) {
		return s
	}
	return s
}

func isSplitRaster(in *image.NRGBA, pos, y, length int) bool {
	occ := 0
	c := in.At(pos, y)
	for x := pos + 1; x < pos+length && x < in.Bounds().Max.X; x++ {
		c2 := in.At(x, y)
		if !constants.ColorsAreEquals(c2, c) {
			return false
		}
		occ++
	}
	return occ >= (length - 1)
}

func setPixelMode0(in *image.NRGBA, out *image.NRGBA, p color.Palette, x, y int, bw []byte, firmwareColorUsed map[int]int, cfg *config.MartineConfig) ([]byte, map[int]int) {
	c1 := in.At(x, y)
	out.Set(x, y, c1)
	pp1, err := common.PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	firmwareColorUsed[pp1]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	out.Set(x+1, y, c2)
	pp2, err := common.PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}

	firmwareColorUsed[pp2]++

	pixel := common.PixelMode0(pp1, pp2)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr := common.CpcScreenAddress(0, x, y, 0, cfg.Overscan)
	bw[addr] = pixel
	return bw, firmwareColorUsed
}

func setPixelMode1(in *image.NRGBA, out *image.NRGBA, p color.Palette, x, y int, bw []byte, firmwareColorUsed map[int]int, cfg *config.MartineConfig) ([]byte, map[int]int) {
	c1 := in.At(x, y)
	out.Set(x, y, c1)
	pp1, err := common.PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	firmwareColorUsed[pp1]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	out.Set(x+1, y, c2)
	pp2, err := common.PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}
	firmwareColorUsed[pp2]++
	c3 := in.At(x+2, y)
	out.Set(x+2, y, c3)
	pp3, err := common.PalettePosition(c3, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
		pp3 = 0
	}
	firmwareColorUsed[pp3]++
	c4 := in.At(x+3, y)
	out.Set(x+3, y, c4)
	pp4, err := common.PalettePosition(c4, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
		pp4 = 0
	}
	firmwareColorUsed[pp4]++

	pixel := common.PixelMode1(pp1, pp2, pp3, pp4)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr := common.CpcScreenAddress(0, x, y, 1, cfg.Overscan)
	bw[addr] = pixel
	return bw, firmwareColorUsed
}

func setPixelMode2(in *image.NRGBA, out *image.NRGBA, p color.Palette, x, y int, bw []byte, firmwareColorUsed map[int]int, cfg *config.MartineConfig) ([]byte, map[int]int) {
	c1 := in.At(x, y)
	out.Set(x, y, c1)
	pp1, err := common.PalettePosition(c1, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
		pp1 = 0
	}
	firmwareColorUsed[pp1]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c2 := in.At(x+1, y)
	out.Set(x+1, y, c2)
	pp2, err := common.PalettePosition(c2, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
		pp2 = 0
	}
	firmwareColorUsed[pp2]++
	c3 := in.At(x+2, y)
	out.Set(x+2, y, c3)
	pp3, err := common.PalettePosition(c3, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
		pp3 = 0
	}
	firmwareColorUsed[pp3]++
	c4 := in.At(x+3, y)
	out.Set(x+3, y, c4)
	pp4, err := common.PalettePosition(c4, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
		pp4 = 0
	}
	firmwareColorUsed[pp4]++
	c5 := in.At(x+4, y)
	out.Set(x+4, y, c5)
	pp5, err := common.PalettePosition(c5, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
		pp5 = 0
	}
	firmwareColorUsed[pp5]++
	//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
	c6 := in.At(x+5, y)
	out.Set(x+5, y, c6)
	pp6, err := common.PalettePosition(c6, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
		pp6 = 0
	}
	firmwareColorUsed[pp6]++
	c7 := in.At(x+6, y)
	out.Set(x+6, y, c7)
	pp7, err := common.PalettePosition(c7, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
		pp3 = 0
	}
	firmwareColorUsed[pp7]++
	c8 := in.At(x+7, y)
	out.Set(x+7, y, c8)
	pp8, err := common.PalettePosition(c8, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+7, y)
		pp8 = 0
	}
	firmwareColorUsed[pp8]++

	pixel := common.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
	//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
	// MACRO PIXM0 COL2,COL1
	// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
	//	MEND
	addr := common.CpcScreenAddress(0, x, y, 2, cfg.Overscan)
	bw[addr] = pixel
	return bw, firmwareColorUsed
}
