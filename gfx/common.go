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
	rotation3DY0 int,
	rotation3DType config.Rotation3d,
	resizingAlgo imaging.ResampleFilter,
	size constants.Size,
) ([]*image.NRGBA, error) {
	var err error

	var images []*image.NRGBA
	if rollMode {
		if rotationRlaBit != -1 || rotationSlaBit != -1 {
			images = transformation.RollLeft(
				rotationRlaBit,
				rotationSlaBit,
				rotationIterations,
				screenMode,
				size,
				in,
				p)
		} else {
			if rotationRraBit != -1 || rotationSraBit != -1 {
				images = transformation.RollRight(
					rotationRraBit,
					rotationSraBit,
					rotationIterations,
					screenMode,
					size,
					in,
					p)
			}
		}
		if rotationKeephighBit != -1 || rotationLosthighBit != -1 {
			images = transformation.RollUp(
				rotationKeephighBit,
				rotationLosthighBit,
				rotationIterations,
				screenMode,
				size,
				in,
				p)
		} else {
			if rotationKeeplowBit != -1 || rotationLostlowBit != -1 {
				images = transformation.RollLow(
					rotationKeeplowBit,
					rotationLostlowBit,
					rotationIterations,
					screenMode,
					size,
					in,
					p)
			}
		}
	}
	if rotationMode {
		if images, err = transformation.Rotate(
			in,
			p,
			size,
			screenMode,
			rollIterations,
			resizingAlgo); err != nil {
			log.GetLogger().Error("Error while perform rotation on image error :%v\n", err)
		}
	}
	if rotation3DMode {
		if images, err = transformation.Rotate3d(
			in,
			p,
			size,
			screenMode,
			resizingAlgo,
			rollIterations,
			rotation3DX0,
			rotation3DY0,
			rotation3DType); err != nil {
			log.GetLogger().Error("Error while perform rotation on image error :%v\n", err)
		}
	}

	return images, err
}

// nolint:funlen, gocognit
func ApplyOneImageAndExport(in image.Image,
	cfg *config.MartineConfig,
	filename, picturePath string,
	mode int,
	screenMode uint8,
) error {
	var palette color.Palette
	var newPalette color.Palette
	var out *image.NRGBA
	var err error

	out = ci.Resize(in, cfg.ScrCfg.Size, cfg.ScrCfg.Process.ResizingAlgo)
	if cfg.ScrCfg.Process.Kmeans.Used {
		log.GetLogger().Info("kmeans with %f threshold", cfg.ScrCfg.Process.Kmeans.Threshold)
		out, err = ci.Kmeans(cfg.ScrCfg.Size.ColorsAvailable, cfg.ScrCfg.Process.Kmeans.Threshold, out)
		if err != nil {
			return err
		}
	}

	log.GetLogger().Info("Saving resized image into (%s)\n", filename+"_resized.png")
	if err := png.Png(filepath.Join(cfg.ScrCfg.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}
	switch cfg.PalCfg.Type {
	case config.PalPalette:
		log.GetLogger().Info("Input palette to apply : (%s)\n", cfg.PalCfg.Path)
		palette, _, err = ocpartstudio.OpenPal(cfg.PalCfg.Path)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", cfg.PalCfg.Path)
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	case config.InkPalette:
		log.GetLogger().Info("Input palette to apply : (%s)\n", cfg.PalCfg.Path)
		palette, _, err = impPalette.OpenInk(cfg.PalCfg.Path)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", cfg.PalCfg.Path)
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	case config.KitPalette:
		log.GetLogger().Info("Input plus palette to apply : (%s)\n", cfg.PalCfg.Path)
		palette, _, err = impPalette.OpenKit(cfg.PalCfg.Path)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", cfg.PalCfg.Path)
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}

	if cfg.ScrCfg.Process.Reducer > 0 {
		out = ci.Reducer(out, cfg.ScrCfg.Process.Reducer)
		if err := png.Png(filepath.Join(cfg.ScrCfg.OutputPath, filename+"_resized.png"), out); err != nil {
			os.Exit(-2)
		}
	}

	if len(palette) > 0 {
		newPalette, out = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, out, err = ci.DowngradingPalette(out, cfg.ScrCfg.Size, cfg.ScrCfg.IsPlus)
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
	paletteToSort = constants.FillColorPalette(paletteToSort)
	newPalette = constants.SortColorsByDistance(paletteToSort)

	out, _ = DoDithering(
		out,
		newPalette,
		cfg.ScrCfg.Process.Dithering.Algo,
		cfg.ScrCfg.Process.Dithering.Type,
		cfg.ScrCfg.Process.Dithering.WithQuantification,
		cfg.ScrCfg.Process.Dithering.Matrix,
		float32(cfg.ScrCfg.Process.Dithering.Multiplier),
		cfg.ScrCfg.IsPlus,
		cfg.ScrCfg.Size)
	if cfg.ScrCfg.Process.Saturation > 0 || cfg.ScrCfg.Process.Brightness > 0 {
		palette = ci.EnhanceBrightness(newPalette, cfg.ScrCfg.Process.Brightness, cfg.ScrCfg.Process.Saturation)
		newPalette, out = ci.DowngradingWithPalette(out, palette)
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
		paletteToSort = constants.FillColorPalette(paletteToSort)
		newPalette = constants.SortColorsByDistance(paletteToSort)
	}

	log.GetLogger().Info("Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := png.Png(filepath.Join(cfg.ScrCfg.OutputPath, filename+"_down.png"), out); err != nil {
		os.Exit(-2)
	}

	images, err := DoTransformation(out, newPalette,
		screenMode, cfg.RotateCfg.RollMode, cfg.RotateCfg.RotationMode, cfg.RotateCfg.Rotation3DMode,
		cfg.RotateCfg.RotationRlaBit, cfg.RotateCfg.RotationSlaBit, cfg.RotateCfg.RotationRraBit, cfg.RotateCfg.RotationSraBit,
		cfg.RotateCfg.RotationKeephighBit, cfg.RotateCfg.RotationLosthighBit,
		cfg.RotateCfg.RotationKeeplowBit, cfg.RotateCfg.RotationLostlowBit, cfg.RotateCfg.RotationIterations,
		cfg.RotateCfg.RollIteration, cfg.RotateCfg.Rotation3DX0, cfg.RotateCfg.Rotation3DY0, cfg.RotateCfg.Rotation3DType, cfg.ScrCfg.Process.ResizingAlgo, cfg.ScrCfg.Size)
	if err != nil {
		os.Exit(-2)
	} else {
		for indice := 0; indice < cfg.RotateCfg.RollIteration; indice++ {
			img := images[indice]
			newFilename := cfg.OsFullPath(filename, fmt.Sprintf("%.2d", indice)+".png")
			if err := png.Png(newFilename, img); err != nil {
				log.GetLogger().Error("Cannot create image (%s) error :%v\n", newFilename, err)
			}
			if err := sprite.ToSpriteAndExport(img, newPalette, constants.Size{Width: cfg.ScrCfg.Size.Width, Height: cfg.ScrCfg.Size.Height}, screenMode, newFilename, false, cfg); err != nil {
				log.GetLogger().Error("Cannot create sprite image (%s) error %v\n", newFilename, err)
			}
		}
	}

	if !cfg.CustomDimension && cfg.ScrCfg.Type != config.SpriteHardFormat {
		err = Transform(out, newPalette, cfg.ScrCfg.Size, picturePath, cfg)
		if err != nil {
			return err
		}
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			out = transformation.Zigzag(out)
		}
		if cfg.ScrCfg.Type != config.SpriteHardFormat {
			// log.GetLogger().Info( "Transform image in sprite.\n")
			err = sprite.ToSpriteAndExport(out, newPalette, cfg.ScrCfg.Size, screenMode, filename, false, cfg)
			if err != nil {
				return err
			}
		} else {
			// log.GetLogger().Info( "Transform image in sprite hard.\n")
			err = spritehard.ToSpriteHardAndExport(out, newPalette, cfg.ScrCfg.Size, screenMode, filename, cfg)
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

// nolint: funlen, gocognit
func ApplyOneImage(in image.Image,
	cfg *config.MartineConfig,
	mode int,
	palette color.Palette,
	screenMode uint8,
) ([]byte, *image.NRGBA, color.Palette, int, error) {
	var newPalette color.Palette
	var out *image.NRGBA
	var err error

	out = ci.Resize(in, cfg.ScrCfg.Size, cfg.ScrCfg.Process.ResizingAlgo)
	if cfg.ScrCfg.Process.Kmeans.Used {
		log.GetLogger().Info("kmeans with %f threshold", cfg.ScrCfg.Process.Kmeans.Threshold)
		out, err = ci.Kmeans(cfg.ScrCfg.Size.ColorsAvailable, cfg.ScrCfg.Process.Kmeans.Threshold, out)
		if err != nil {
			return []byte{}, out, palette, 0, err
		}
	}

	if cfg.ScrCfg.Process.Reducer > -1 {
		out = ci.Reducer(out, cfg.ScrCfg.Process.Reducer)
	}

	if len(palette) > 0 {
		newPalette, out = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, out, err = ci.DowngradingPalette(out, cfg.ScrCfg.Size, cfg.ScrCfg.IsPlus)
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
	paletteToSort = constants.FillColorPalette(paletteToSort)
	newPalette = constants.SortColorsByDistance(paletteToSort)
	out, _ = DoDithering(
		out,
		newPalette,
		cfg.ScrCfg.Process.Dithering.Algo,
		cfg.ScrCfg.Process.Dithering.Type,
		cfg.ScrCfg.Process.Dithering.WithQuantification,
		cfg.ScrCfg.Process.Dithering.Matrix,
		float32(cfg.ScrCfg.Process.Dithering.Multiplier),
		cfg.ScrCfg.IsPlus,
		cfg.ScrCfg.Size)

	if cfg.ScrCfg.Process.Saturation > 0 || cfg.ScrCfg.Process.Brightness > 0 {
		palette = ci.EnhanceBrightness(newPalette, cfg.ScrCfg.Process.Brightness, cfg.ScrCfg.Process.Saturation)
		newPalette, out = ci.DowngradingWithPalette(out, palette)
		var paletteToSort color.Palette
		switch mode {
		case 1:
			paletteToSort = newPalette[0:4]
		case 2:
			paletteToSort = newPalette[0:2]
		default:
			paletteToSort = newPalette
		}
		paletteToSort = constants.FillColorPalette(paletteToSort)
		newPalette = constants.SortColorsByDistance(paletteToSort)
	}
	var data []byte
	var lineSize int
	if !cfg.CustomDimension && cfg.ScrCfg.Type != config.SpriteHardFormat {
		data = InternalTransform(out, newPalette, cfg.ScrCfg.Size, cfg)
		lineSize = cfg.ScrCfg.Size.Width
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			out = transformation.Zigzag(out)
		}
		if cfg.ScrCfg.Type != config.SpriteHardFormat {
			// log.GetLogger().Info( "Transform image in sprite.\n")
			data, _, lineSize, err = sprite.ToSprite(out, newPalette, cfg.ScrCfg.Size, screenMode, cfg)
		} else {
			//	log.GetLogger().Info( "Transform image in sprite hard.\n")
			data, _ = spritehard.ToSpriteHard(out, newPalette, cfg.ScrCfg.Size, screenMode, cfg)
			lineSize = 16
		}
	}
	if cfg.ScrCfg.Process.OneRow {
		for y := 0; y < out.Bounds().Max.Y; y += 2 {
			for x := 0; x < out.Bounds().Max.X; x++ {
				out.Set(x, y, newPalette[0])
			}
		}
	}
	if cfg.ScrCfg.Process.OneLine {
		for y := 0; y < out.Bounds().Max.Y; y++ {
			for x := 0; x < out.Bounds().Max.X; x += 2 {
				out.Set(x, y, newPalette[0])
			}
		}
	}
	return data, out, newPalette, lineSize, err
}

func ExportRawImage(in image.Image,
	palette color.Palette,
	cfg *config.MartineConfig,
	filename, picturePath string,
	screenMode uint8) error {
	var err error
	//check the palette and sort it
	palette = constants.FillColorPalette(palette)
	palette = constants.SortColorsByDistance(palette)

	// at least resize the image to be sure
	out := ci.Resize(in, cfg.ScrCfg.Size, cfg.ScrCfg.Process.ResizingAlgo)
	if cfg.ScrCfg.Process.Kmeans.Used {
		log.GetLogger().Info("kmeans with %f threshold", cfg.ScrCfg.Process.Kmeans.Threshold)
		out, err = ci.Kmeans(cfg.ScrCfg.Size.ColorsAvailable, cfg.ScrCfg.Process.Kmeans.Threshold, out)
		if err != nil {
			return err
		}
	}
	if !cfg.CustomDimension && cfg.ScrCfg.Type != config.SpriteHardFormat {
		err = Transform(out, palette, cfg.ScrCfg.Size, picturePath, cfg)
		if err != nil {
			return err
		}
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			out = transformation.Zigzag(out)
		}
		if cfg.ScrCfg.Type != config.SpriteHardFormat {
			// log.GetLogger().Info( "Transform image in sprite.\n")
			err = sprite.ToSpriteAndExport(out, palette, cfg.ScrCfg.Size, screenMode, filename, false, cfg)
			if err != nil {
				return err
			}
		} else {
			// log.GetLogger().Info( "Transform image in sprite hard.\n")
			err = spritehard.ToSpriteHardAndExport(out, palette, cfg.ScrCfg.Size, screenMode, filename, cfg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
