// Package effect provides CPC screen effect utilities.
package effect

import (
	"image"
	"image/color"
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
	"github.com/jeromelesaux/martine/log"
)

// Egx generates an EGX screen from two input screen files.
// nolint: funlen, gocognit, gocyclo
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
		if cfg.ScrCfg.Type == config.FullscreenFormat {
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
		if err := ToEgx1(in0, in1, p, uint8(m1), "egx.scr", cfg); err != nil {
			return err
		}
		if cfg.ScrCfg.Type != config.FullscreenFormat {
			if err := ocpartstudio.EgxLoader("egx.scr", p, uint8(m1), uint8(m2), cfg); err != nil {
				return err
			}
		}
	} else if m1 == 1 && m2 == 2 || m2 == 1 && m1 == 2 {
		var f2, f1 string
		var mode2, mode1 uint8
		var err error
		mode1 = 1
		mode2 = 2
		if m1 == 1 {
			f1 = filepath1
			f2 = filepath2
		} else {
			f1 = filepath2
			f2 = filepath1
		}
		var in2, in1 *image.NRGBA
		if cfg.ScrCfg.Type.IsFullScreen() {
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
		if err := ToEgx2(in1, in2, p, uint8(m1), "egx.scr", cfg); err != nil {
			return err
		}
		if cfg.ScrCfg.Type != config.FullscreenFormat {
			if err := ocpartstudio.EgxLoader("egx.scr", p, uint8(m1), uint8(m2), cfg); err != nil {
				return err
			}
		}
	} else {
		return errors.ErrorFeatureNotImplemented
	}
	return nil
}

func prepareEgx(
	in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string,
) (*image.NRGBA, color.Palette, error) {
	var err error
	size := constants.Size{
		Width:  cfg.ScrCfg.Size.Width,
		Height: cfg.ScrCfg.Size.Height}

	im := ci.Resize(in, size, cfg.ScrCfg.Process.ResizingAlgo)
	loadedPalette := gfx.LoadPalette(cfg)
	var p color.Palette
	var downgraded *image.NRGBA

	if len(loadedPalette) > 0 {
		p, downgraded = ci.DowngradingWithPalette(im, loadedPalette)
	} else {
		p, downgraded, err = ci.DowngradingPalette(im, cfg.ScrCfg.Size, cfg.ScrCfg.IsPlus)
		if err != nil {
			log.GetLogger().Error("Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	p = constants.SortColorsByDistance(p)
	log.GetLogger().Info("Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := png.Png(filepath.Join(cfg.ScrCfg.OutputPath, filename+"_down.png"), downgraded); err != nil {
		return nil, nil, err
	}
	out, p := gfx.DoDithering(
		downgraded,
		p,
		cfg.ScrCfg.Process.Dithering.Algo,
		cfg.ScrCfg.Process.Dithering.Type,
		cfg.ScrCfg.Process.Dithering.WithQuantification,
		cfg.ScrCfg.Process.Dithering.Matrix,
		float32(cfg.ScrCfg.Process.Dithering.Multiplier),
		cfg.ScrCfg.IsPlus,
		cfg.ScrCfg.Size)
	return out, p, nil
}

func AutoEgx1(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string) error {
	downgraded, p, err := prepareEgx(in, cfg, filename, picturePath)
	if err != nil {
		return err
	}
	return ToEgx1(downgraded, downgraded, p, 0, picturePath, cfg)
}

func AutoEgx2(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string) error {
	downgraded, p, err := prepareEgx(in, cfg, filename, picturePath)
	if err != nil {
		return err
	}
	return ToEgx2(downgraded, downgraded, p, 1, picturePath, cfg)
}

func ToEgx1(inMode0, inMode1 *image.NRGBA, p color.Palette, firstLineMode uint8, picturePath string, cfg *config.MartineConfig) error {
	bw, p := ToEgx1Raw(inMode0, inMode1, p, firstLineMode, cfg)
	return export.Export(picturePath, bw, p, 1, cfg)
}

// ToEgx1Raw interleaves mode 0 and mode 1 lines into a raw EGX1 buffer.
// nolint: funlen
func ToEgx1Raw(inMode0, inMode1 *image.NRGBA, p color.Palette, firstLineMode uint8, cfg *config.MartineConfig) ([]byte, color.Palette) {
	var bw []byte
	if cfg.ScrCfg.Type == config.FullscreenFormat {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	mode0Line := 1
	mode1Line := 0
	if firstLineMode == 1 {
		mode0Line = 0
		mode1Line = 1
	}

	for y := inMode0.Bounds().Min.Y + mode0Line; y < inMode0.Bounds().Max.Y; y += 2 {
		for x := inMode0.Bounds().Min.X; x < inMode0.Bounds().Max.X; x += 2 {
			c1 := inMode0.At(x, y)
			pp1 := palette.PalettePositionOrDefault(c1, p)
			c2 := inMode0.At(x+1, y)
			pp2 := palette.PalettePositionOrDefault(c2, p)
			pix := pixel.PixelMode0(pp1, pp2)
			addr := address.CpcScreenAddress(0, x, y, 0, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}
	}
	for y := inMode1.Bounds().Min.Y + mode1Line; y < inMode1.Bounds().Max.Y; y += 2 {
		for x := inMode1.Bounds().Min.X; x < inMode1.Bounds().Max.X; x += 4 {
			c1 := inMode1.At(x, y)
			pp1 := palette.PalettePositionOrDefault(c1, p)
			c2 := inMode1.At(x+1, y)
			pp2 := palette.PalettePositionOrDefault(c2, p)
			c3 := inMode1.At(x+2, y)
			pp3 := palette.PalettePositionOrDefault(c3, p)
			c4 := inMode1.At(x+3, y)
			pp4 := palette.PalettePositionOrDefault(c4, p)

			pix := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			addr := address.CpcScreenAddress(0, x, y, 1, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
			addr = address.CpcScreenAddress(0, x+1, y, 1, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}
	}
	return bw, p
}

func ToEgx2(inMode1, inMode2 *image.NRGBA, p color.Palette, firstLineMode uint8, picturePath string, cfg *config.MartineConfig) error {
	bw, p := ToEgx2Raw(inMode1, inMode2, p, firstLineMode, cfg)
	return export.Export(picturePath, bw, p, 2, cfg)
}

// ToEgx2Raw interleaves mode 1 and mode 2 lines into a raw EGX2 buffer.
// nolint: funlen, gocognit
func ToEgx2Raw(inMode1, inMode2 *image.NRGBA, p color.Palette, firstLineMode uint8, cfg *config.MartineConfig) ([]byte, color.Palette) {
	var bw []byte
	if cfg.ScrCfg.Type.IsFullScreen() {
		bw = make([]byte, 0x8000)
	} else {
		bw = make([]byte, 0x4000)
	}
	mode1Line := 1
	mode2Line := 0
	if firstLineMode == 2 {
		mode1Line = 0
		mode2Line = 1
	}

	for y := inMode1.Bounds().Min.Y + mode1Line; y < inMode1.Bounds().Max.Y; y += 2 {
		for x := inMode1.Bounds().Min.X; x < inMode1.Bounds().Max.X; x += 4 {
			c1 := inMode1.At(x, y)
			pp1 := palette.PalettePositionOrDefault(c1, p)
			c2 := inMode1.At(x+1, y)
			pp2 := palette.PalettePositionOrDefault(c2, p)
			c3 := inMode1.At(x+2, y)
			pp3 := palette.PalettePositionOrDefault(c3, p)
			c4 := inMode1.At(x+3, y)
			pp4 := palette.PalettePositionOrDefault(c4, p)

			pix := pixel.PixelMode1(pp1, pp2, pp3, pp4)
			addr := address.CpcScreenAddress(0, x, y, 1, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}
	}
	for y := inMode2.Bounds().Min.Y + mode2Line; y < inMode2.Bounds().Max.Y; y += 2 {
		for x := inMode2.Bounds().Min.X; x < inMode2.Bounds().Max.X; x += 8 {
			c1 := inMode2.At(x, y)
			pp1 := palette.PalettePositionOrDefault(c1, p)
			c2 := inMode2.At(x+1, y)
			pp2 := palette.PalettePositionOrDefault(c2, p)
			c3 := inMode2.At(x+2, y)
			pp3 := palette.PalettePositionOrDefault(c3, p)
			c4 := inMode2.At(x+3, y)
			pp4 := palette.PalettePositionOrDefault(c4, p)
			c5 := inMode2.At(x+4, y)
			pp5 := palette.PalettePositionOrDefault(c5, p)
			c6 := inMode2.At(x+5, y)
			pp6 := palette.PalettePositionOrDefault(c6, p)
			c7 := inMode2.At(x+6, y)
			pp7 := palette.PalettePositionOrDefault(c7, p)
			c8 := inMode2.At(x+7, y)
			pp8 := palette.PalettePositionOrDefault(c8, p)
			pix := pixel.PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8)
			addr := address.CpcScreenAddress(0, x, y, 2, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
			addr = address.CpcScreenAddress(0, x+1, y, 2, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
			bw[addr] = pix
		}
	}
	return bw, p
}

// EgxRaw converts two images into a raw EGX screen buffer.
// nolint: funlen, gocognit
func EgxRaw(img1, img2 image.Image, p color.Palette, mode1, mode2 int, cfg *config.MartineConfig) (bw []byte, outPalette color.Palette, lineMode int, err error) {
	if mode1 == 0 && mode2 == 1 || mode2 == 0 && mode1 == 1 {
		var in0, in1 *image.NRGBA
		var firstLineMode uint8
		if mode1 == 0 {
			in0 = png.ImageToNRGBA(img1)
			in1 = png.ImageToNRGBA(img2)
			firstLineMode = 0
		} else {
			in0 = png.ImageToNRGBA(img2)
			in1 = png.ImageToNRGBA(img1)
			firstLineMode = 1
		}
		res, p := ToEgx1Raw(in0, in1, p, firstLineMode, cfg)
		return res, p, 1, nil
	} else if mode1 == 1 && mode2 == 2 || mode2 == 1 && mode1 == 2 {
		var in0, in1 *image.NRGBA
		var firstLineMode uint8
		if mode1 == 1 {
			in0 = png.ImageToNRGBA(img1)
			in1 = png.ImageToNRGBA(img2)
			firstLineMode = 1
		} else {
			in0 = png.ImageToNRGBA(img2)
			in1 = png.ImageToNRGBA(img1)
			firstLineMode = 2
		}
		res, p := ToEgx2Raw(in0, in1, p, firstLineMode, cfg)
		return res, p, 2, nil
	}
	return nil, p, 0, errors.ErrorFeatureNotImplemented
}
