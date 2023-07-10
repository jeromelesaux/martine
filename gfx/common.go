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
	"github.com/jeromelesaux/martine/log"
)

func DoDithering(in *image.NRGBA,
	p color.Palette,
	ditheringAlgo int,
	ditheringType constants.DitheringType,
	ditheringWithQuantification bool,
	ditheringMatrix [][]float32,
	ditheringMultiplier float32,
	isCpcPlus bool,
	size constants.Size,
) (*image.NRGBA, color.Palette) {
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
	size constants.Size,
) ([]*image.NRGBA, error) {
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
			log.GetLogger().Error("Error while perform rotation on image error :%v\n", err)
		}
	}
	if rotation3DMode {
		if images, err = transformation.Rotate3d(in, p, size, screenMode, resizingAlgo, rollIterations, rotation3DX0, rotation3DY0, rotation3DType); err != nil {
			log.GetLogger().Error("Error while perform rotation on image error :%v\n", err)
		}
	}

	return images, err
}

func ApplyOneImageAndExport(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string,
	mode int,
	screenMode uint8,
) error {
	var palette color.Palette
	var newPalette color.Palette
	var downgraded, out *image.NRGBA
	var err error

	if cfg.PalettePath != "" {
		log.GetLogger().Info("Input palette to apply : (%s)\n", cfg.PalettePath)
		palette, _, err = ocpartstudio.OpenPal(cfg.PalettePath)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", cfg.PalettePath)
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}
	if cfg.InkPath != "" {
		log.GetLogger().Info("Input palette to apply : (%s)\n", cfg.InkPath)
		palette, _, err = impPalette.OpenInk(cfg.InkPath)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", cfg.InkPath)
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}
	if cfg.KitPath != "" {
		log.GetLogger().Info("Input plus palette to apply : (%s)\n", cfg.KitPath)
		palette, _, err = impPalette.OpenKit(cfg.KitPath)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", cfg.KitPath)
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}

	// if cfg.UseKmeans {
	// 	downgraded, err = ci.Kmeans(cfg.Size.ColorsAvailable, cfg.KmeansInterations, in)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	out = ci.Resize(downgraded, cfg.Size, cfg.ResizingAlgo)
	// } else {
	out = ci.Resize(in, cfg.Size, cfg.ResizingAlgo)
	// }
	log.GetLogger().Info("Saving resized image into (%s)\n", filename+"_resized.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}

	if cfg.UseKmeans {
		out, err = ci.Kmeans(cfg.Size.ColorsAvailable, cfg.KmeansInterations, out)
		if err != nil {
			log.GetLogger().Info("error while applying kmeans with iterations [%d] (%v)\n", cfg.KmeansInterations, err)
			return err
		}
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
			log.GetLogger().Error("Cannot downgrade colors palette for this image %s\n", picturePath)
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

	log.GetLogger().Info("Saving downgraded image into (%s)\n", filename+"_down.png")
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
				log.GetLogger().Error("Cannot create image (%s) error :%v\n", newFilename, err)
			}
			if err := sprite.ToSpriteAndExport(img, newPalette, constants.Size{Width: cfg.Size.Width, Height: cfg.Size.Height}, screenMode, newFilename, false, cfg); err != nil {
				log.GetLogger().Error("Cannot create sprite image (%s) error %v\n", newFilename, err)
			}
		}
	}

	if !cfg.CustomDimension && !cfg.SpriteHard {
		err = Transform(downgraded, newPalette, cfg.Size, picturePath, cfg)
		if err != nil {
			return err
		}
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cfg.SpriteHard {
			// log.GetLogger().Info( "Transform image in sprite.\n")
			err = sprite.ToSpriteAndExport(downgraded, newPalette, cfg.Size, screenMode, filename, false, cfg)
			if err != nil {
				return err
			}
		} else {
			// log.GetLogger().Info( "Transform image in sprite hard.\n")
			err = spritehard.ToSpriteHardAndExport(downgraded, newPalette, cfg.Size, screenMode, filename, cfg)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func ApplyImages(
	in []*image.NRGBA,
	cfg *config.MartineConfig,
	mode int,
	palette color.Palette,
	screenMode uint8,
) ([][]byte, []*image.NRGBA, color.Palette, error) {
	var gErr error
	raw := make([][]byte, 0)
	images := make([]*image.NRGBA, 0)
	for _, img := range in {
		v, r, _, _, err := ApplyOneImage(img, cfg, mode, palette, screenMode)
		if err != nil {
			gErr = err
		}
		raw = append(raw, v)
		images = append(images, r)
	}
	return raw, images, palette, gErr
}

func ApplyOneImage(in image.Image,
	cfg *config.MartineConfig,
	mode int,
	palette color.Palette,
	screenMode uint8,
) ([]byte, *image.NRGBA, color.Palette, int, error) {
	var newPalette color.Palette
	var downgraded, out *image.NRGBA
	var err error

	// if cfg.UseKmeans {
	// 	downgraded, err = ci.Kmeans(cfg.Size.ColorsAvailable, cfg.KmeansInterations, in)
	// 	if err != nil {
	// 		return []byte{}, downgraded, palette, 0, err
	// 	}
	out = ci.Resize(in, cfg.Size, cfg.ResizingAlgo)
	// } else {
	// 	out = ci.Resize(in, cfg.Size, cfg.ResizingAlgo)
	// }

	if cfg.UseKmeans {
		out, err = ci.Kmeans(cfg.Size.ColorsAvailable, cfg.KmeansInterations, out)
		if err != nil {
			log.GetLogger().Info("error while applying kmeans with iterations [%d] (%v)\n", cfg.KmeansInterations, err)
			return []byte{}, out, palette, 0, err
		}
	}
	if cfg.Reducer > -1 {
		out = ci.Reducer(out, cfg.Reducer)
	}

	if len(palette) > 0 {
		newPalette, downgraded = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = ci.DowngradingPalette(out, cfg.Size, cfg.CpcPlus)
		if err != nil {
			log.GetLogger().Error("Cannot downgrade colors palette for this image")
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
			// log.GetLogger().Info( "Transform image in sprite.\n")
			data, _, lineSize, err = sprite.ToSprite(downgraded, newPalette, cfg.Size, screenMode, cfg)
		} else {
			//	log.GetLogger().Info( "Transform image in sprite hard.\n")
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
