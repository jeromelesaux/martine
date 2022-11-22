package gfx

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/convert/spritehard"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/filter"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

func DoDithering(in *image.NRGBA,
	p color.Palette,
	ditheringAlgo int,
	ditheringType constants.DitheringType,
	ditheringWithQuantification bool,
	ditheringMatrix [][]float32,
	ditheringMultiplier float32,
	isCpcPlus bool,
	size constants.Size) (*image.NRGBA, color.Palette) {
	if ditheringAlgo != -1 {
		switch ditheringType {
		case constants.ErrorDiffusionDither:
			if ditheringWithQuantification {
				in = filter.QuantizeWithDither(in, ditheringMatrix, size.ColorsAvailable, p)
			} else {
				in = filter.Dithering(in, ditheringMatrix, ditheringMultiplier)
			}
		case constants.OrderedDither:
			if isCpcPlus {
				p = ci.ExtractPalette(in, isCpcPlus, 27)
				in = filter.BayerDiphering(in, ditheringMatrix, p)
			} else {
				in = filter.BayerDiphering(in, ditheringMatrix, constants.CpcOldPalette)
			}
		}
	}
	return in, p
}

func DoTransformation(in *image.NRGBA,
	p color.Palette,
	screenMode uint8,
	rollMode,
	rotationMode,
	rotation3DMode bool,
	rotationRlaBit,
	rotationSlaBit,
	rotationRraBit,
	rotationSraBit,
	rotationKeephighBit,
	rotationLosthighBit,
	rotationKeeplowBit,
	rotationLostlowBit,
	rotationIterations,
	rollIterations int,
	rotation3DX0,
	rotation3DY0,
	rotation3DType int,
	resizingAlgo imaging.ResampleFilter,
	size constants.Size) ([]*image.NRGBA, error) {
	var err error

	var images []*image.NRGBA
	if rollMode {
		if rotationRlaBit != -1 || rotationSlaBit != -1 {
			images = transformation.RollLeft(rotationRlaBit, rotationSlaBit, rotationIterations, screenMode, size, in, p)
		} else {
			if rotationRraBit != -1 || rotationSraBit != -1 {
				images = transformation.RollRight(rotationRraBit, rotationSraBit, rotationIterations, screenMode, size, in, p)
			}
		}
		if rotationKeephighBit != -1 || rotationLosthighBit != -1 {
			images = transformation.RollUp(rotationKeephighBit, rotationLosthighBit, rotationIterations, screenMode, size, in, p)
		} else {
			if rotationKeeplowBit != -1 || rotationLostlowBit != -1 {
				images = transformation.RollLow(rotationKeeplowBit, rotationLostlowBit, rotationIterations, screenMode, size, in, p)
			}
		}
	}
	if rotationMode {
		if images, err = transformation.Rotate(in, p, size, screenMode, rollIterations, resizingAlgo); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image error :%v\n", err)
		}
	}
	if rotation3DMode {
		if images, err = transformation.Rotate3d(in, p, size, screenMode, resizingAlgo, rollIterations, rotation3DX0, rotation3DY0, rotation3DType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image error :%v\n", err)
		}
	}

	return images, err
}

func ApplyOneImageAndExport(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string,
	mode int,
	screenMode uint8) error {

	var palette color.Palette
	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	if cfg.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cfg.PalettePath)
		palette, _, err = ocpartstudio.OpenPal(cfg.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cfg.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if cfg.InkPath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cfg.InkPath)
		palette, _, err = impPalette.OpenInk(cfg.InkPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cfg.InkPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if cfg.KitPath != "" {
		fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", cfg.KitPath)
		palette, _, err = impPalette.OpenKit(cfg.KitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cfg.KitPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	out := ci.Resize(in, cfg.Size, cfg.ResizingAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}

	if cfg.Reducer > 0 {
		out = ci.Reducer(out, cfg.Reducer)
		if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_resized.png"), out); err != nil {
			os.Exit(-2)
		}
	}

	if len(palette) > 0 {
		newPalette, downgraded = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = ci.DowngradingPalette(out, cfg.Size, cfg.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	var paletteToSort color.Palette
	switch mode {
	case 1:
		paletteToSort = newPalette[0:4]
	case 2:
		paletteToSort = newPalette[0:2]
	default:
		paletteToSort = newPalette
	}
	paletteToSort = fillColorPalette(paletteToSort)
	newPalette = constants.SortColorsByDistance(paletteToSort)

	out, _ = DoDithering(out, newPalette, cfg.DitheringAlgo, cfg.DitheringType, cfg.DitheringWithQuantification, cfg.DitheringMatrix, float32(cfg.DitheringMultiplier), cfg.CpcPlus, cfg.Size)
	if cfg.Saturation > 0 || cfg.Brightness > 0 {
		palette = ci.EnhanceBrightness(newPalette, cfg.Brightness, cfg.Saturation)
		newPalette, downgraded = ci.DowngradingWithPalette(out, palette)
		var paletteToSort color.Palette
		switch mode {
		case 1:
			end := len(newPalette)
			if len(newPalette) >= 4 {
				end = 4
			}
			paletteToSort = newPalette[0:end]
		case 2:
			end := len(newPalette)
			if len(newPalette) >= 2 {
				end = 2
			}
			paletteToSort = newPalette[0:end]
		default:
			paletteToSort = newPalette
		}
		paletteToSort = fillColorPalette(paletteToSort)
		newPalette = constants.SortColorsByDistance(paletteToSort)
	}

	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	images, err := DoTransformation(downgraded, newPalette,
		screenMode, cfg.RollMode, cfg.RotationMode, cfg.Rotation3DMode,
		cfg.RotationRlaBit, cfg.RotationSlaBit, cfg.RotationRraBit, cfg.RotationSraBit,
		cfg.RotationKeephighBit, cfg.RotationLosthighBit,
		cfg.RotationKeeplowBit, cfg.RotationLostlowBit, cfg.RotationIterations,
		cfg.RollIteration, cfg.Rotation3DX0, cfg.Rotation3DY0, cfg.Rotation3DType, cfg.ResizingAlgo, cfg.Size)
	if err != nil {
		os.Exit(-2)
	} else {

		for indice := 0; indice < cfg.RollIteration; indice++ {
			img := images[indice]
			newFilename := cfg.OsFullPath(filename, fmt.Sprintf("%.2d", indice)+".png")
			if err := png.Png(newFilename, img); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create image (%s) error :%v\n", newFilename, err)
			}
			if err := sprite.ToSpriteAndExport(img, newPalette, constants.Size{Width: cfg.Size.Width, Height: cfg.Size.Height}, screenMode, newFilename, false, cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create sprite image (%s) error %v\n", newFilename, err)
			}
		}
	}

	if !cfg.CustomDimension && !cfg.SpriteHard {
		Transform(downgraded, newPalette, cfg.Size, picturePath, cfg)
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cfg.SpriteHard {
			//fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			sprite.ToSpriteAndExport(downgraded, newPalette, cfg.Size, screenMode, filename, false, cfg)
		} else {
			//fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			spritehard.ToSpriteHardAndExport(downgraded, newPalette, cfg.Size, screenMode, filename, cfg)
		}
	}
	return err
}

func ApplyOneImage(in image.Image,
	cfg *config.MartineConfig,
	mode int,
	palette color.Palette,
	screenMode uint8) ([]byte, *image.NRGBA, color.Palette, int, error) {

	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	out := ci.Resize(in, cfg.Size, cfg.ResizingAlgo)

	if cfg.Reducer > -1 {
		out = ci.Reducer(out, cfg.Reducer)
	}

	if len(palette) > 0 {
		newPalette, downgraded = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = ci.DowngradingPalette(out, cfg.Size, cfg.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image")
		}
	}

	var paletteToSort color.Palette
	switch mode {
	case 1:
		end := len(newPalette)
		if len(newPalette) >= 4 {
			end = 4
		}
		paletteToSort = newPalette[0:end]
	case 2:
		end := len(newPalette)
		if len(newPalette) >= 2 {
			end = 2
		}
		paletteToSort = newPalette[0:end]
	default:
		paletteToSort = newPalette
	}
	paletteToSort = fillColorPalette(paletteToSort)
	newPalette = constants.SortColorsByDistance(paletteToSort)
	out, _ = DoDithering(out, newPalette, cfg.DitheringAlgo, cfg.DitheringType, cfg.DitheringWithQuantification, cfg.DitheringMatrix, float32(cfg.DitheringMultiplier), cfg.CpcPlus, cfg.Size)

	if cfg.Saturation > 0 || cfg.Brightness > 0 {
		palette = ci.EnhanceBrightness(newPalette, cfg.Brightness, cfg.Saturation)
		newPalette, downgraded = ci.DowngradingWithPalette(out, palette)
		var paletteToSort color.Palette
		switch mode {
		case 1:
			paletteToSort = newPalette[0:4]
		case 2:
			paletteToSort = newPalette[0:2]
		default:
			paletteToSort = newPalette
		}
		paletteToSort = fillColorPalette(paletteToSort)
		newPalette = constants.SortColorsByDistance(paletteToSort)
	}
	var data []byte
	var lineSize int
	if !cfg.CustomDimension && !cfg.SpriteHard {
		data = InternalTransform(downgraded, newPalette, cfg.Size, cfg)
		lineSize = cfg.Size.Width
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cfg.SpriteHard {
			//fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			data, _, lineSize, err = sprite.ToSprite(downgraded, newPalette, cfg.Size, screenMode, cfg)
		} else {
			//	fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			data, _ = spritehard.ToSpriteHard(downgraded, newPalette, cfg.Size, screenMode, cfg)
			lineSize = 16
		}
	}
	if cfg.OneRow {
		for y := 0; y < downgraded.Bounds().Max.Y; y += 2 {
			for x := 0; x < downgraded.Bounds().Max.X; x++ {
				downgraded.Set(x, y, newPalette[0])
			}
		}
	}
	if cfg.OneLine {
		for y := 0; y < downgraded.Bounds().Max.Y; y++ {
			for x := 0; x < downgraded.Bounds().Max.X; x += 2 {
				downgraded.Set(x, y, newPalette[0])
			}
		}
	}
	return data, downgraded, newPalette, lineSize, err
}

func fillColorPalette(p color.Palette) color.Palette {
	for i, v := range p {
		if v == nil {
			p[i] = color.Black
		}
	}
	return p
}
