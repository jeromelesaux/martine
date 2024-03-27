package gfx

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common/iface"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/convert/spritehard"
	palettepath "github.com/jeromelesaux/martine/export/palette_path"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/filter"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
)

func DoDithering(in *image.NRGBA,
	p color.Palette,
	d constants.Dithering,
	isCpcPlus bool,
	size constants.Size,
) (*image.NRGBA, color.Palette) {
	if d.DitheringAlgo != -1 {
		switch d.DitheringType {
		case constants.ErrorDiffusionDither:
			if d.DitheringWithQuantification {
				in = filter.QuantizeWithDither(in, d.DitheringMatrix, size.ColorsAvailable, p)
			} else {
				in = filter.Dithering(in, d.DitheringMatrix, float32(d.DitheringMultiplier))
			}
		case constants.OrderedDither:
			if isCpcPlus {
				p = ci.ExtractPalette(in, isCpcPlus, 27)
				in = filter.BayerDiphering(in, d.DitheringMatrix, p)
			} else {
				in = filter.BayerDiphering(in, d.DitheringMatrix, constants.CpcOldPalette)
			}
		}
	}
	return in, p
}

func DoTransformation(in *image.NRGBA,
	p color.Palette,
	screenMode uint8,
	tr constants.Transformation,
	resizingAlgo imaging.ResampleFilter,
	size constants.Size,
) ([]*image.NRGBA, error) {
	var err error

	var images []*image.NRGBA
	if tr.RollMode {
		if tr.RotationRlaBit != -1 || tr.RotationSlaBit != -1 {
			images = transformation.RollLeft(tr.RotationRlaBit, tr.RotationSlaBit, tr.RotationIterations, screenMode, size, in, p)
		} else {
			if tr.RotationRraBit != -1 || tr.RotationSraBit != -1 {
				images = transformation.RollRight(tr.RotationRraBit, tr.RotationSraBit, tr.RotationIterations, screenMode, size, in, p)
			}
		}
		if tr.RotationKeephighBit != -1 || tr.RotationLosthighBit != -1 {
			images = transformation.RollUp(tr.RotationKeephighBit, tr.RotationLosthighBit, tr.RotationIterations, screenMode, size, in, p)
		} else {
			if tr.RotationKeeplowBit != -1 || tr.RotationLostlowBit != -1 {
				images = transformation.RollLow(tr.RotationKeeplowBit, tr.RotationLostlowBit, tr.RotationIterations, screenMode, size, in, p)
			}
		}
	}
	if tr.RotationMode {
		if images, err = transformation.Rotate(in, p, size, screenMode, tr.RollIteration, resizingAlgo); err != nil {
			log.GetLogger().Error("Error while perform rotation on image error :%v\n", err)
		}
	}
	if tr.Rotation3DMode {
		if images, err = transformation.Rotate3d(in, p, size, screenMode, resizingAlgo, tr.RollIteration, tr.Rotation3DX0, tr.Rotation3DY0, tr.Rotation3DType); err != nil {
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

	out = ci.Resize(in, cfg.Size, cfg.ResizingAlgo)
	if cfg.UseKmeans {
		log.GetLogger().Info("kmeans with %f threshold", cfg.KmeansThreshold)
		out, err = ci.Kmeans(cfg.Size.ColorsAvailable, cfg.KmeansThreshold, out)
		if err != nil {
			return err
		}
	}

	log.GetLogger().Info("Saving resized image into (%s)\n", filename+"_resized.png")
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}

	palette, err = palettepath.Open(cfg.PalettePath)
	if err != nil {
		log.GetLogger().Error("Palette in file (%v) can not be read skipped\n", cfg.PalettePath)
	}

	if cfg.Reducer > 0 {
		out = ci.Reducer(out, cfg.Reducer)
		if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_resized.png"), out); err != nil {
			os.Exit(-2)
		}
	}

	if len(palette) > 0 {
		newPalette, out = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, out, err = ci.DowngradingPalette(out, cfg.Size, cfg.CpcPlus)
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

	out, _ = DoDithering(out, newPalette, cfg.Dithering, cfg.CpcPlus, cfg.Size)
	if cfg.Saturation > 0 || cfg.Brightness > 0 {
		palette = ci.EnhanceBrightness(newPalette, cfg.Brightness, cfg.Saturation)
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
	if err := png.Png(filepath.Join(cfg.OutputPath, filename+"_down.png"), out); err != nil {
		os.Exit(-2)
	}

	images, err := DoTransformation(out, newPalette,
		screenMode, cfg.Transformation, cfg.ResizingAlgo, cfg.Size)
	if err != nil {
		os.Exit(-2)
	} else {
		for indice := 0; indice < cfg.Transformation.RollIteration; indice++ {
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
		err = Transform(out, newPalette, cfg.Size, picturePath, cfg)
		if err != nil {
			return err
		}
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			out = transformation.Zigzag(out)
		}
		if !cfg.SpriteHard {
			// log.GetLogger().Info( "Transform image in sprite.\n")
			err = sprite.ToSpriteAndExport(out, newPalette, cfg.Size, screenMode, filename, false, cfg)
			if err != nil {
				return err
			}
		} else {
			// log.GetLogger().Info( "Transform image in sprite hard.\n")
			err = spritehard.ToSpriteHardAndExport(out, newPalette, cfg.Size, screenMode, filename, cfg)
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

	out = ci.Resize(in, cfg.Size, cfg.ResizingAlgo)
	if cfg.UseKmeans {
		log.GetLogger().Info("kmeans with %f threshold", cfg.KmeansThreshold)
		out, err = ci.Kmeans(cfg.Size.ColorsAvailable, cfg.KmeansThreshold, out)
		if err != nil {
			return []byte{}, out, palette, 0, err
		}
	}

	if cfg.Reducer > -1 {
		out = ci.Reducer(out, cfg.Reducer)
	}

	if len(palette) > 0 {
		newPalette, out = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, out, err = ci.DowngradingPalette(out, cfg.Size, cfg.CpcPlus)
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
	out, _ = DoDithering(out, newPalette, cfg.Dithering, cfg.CpcPlus, cfg.Size)

	if cfg.Saturation > 0 || cfg.Brightness > 0 {
		palette = ci.EnhanceBrightness(newPalette, cfg.Brightness, cfg.Saturation)
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
	if !cfg.CustomDimension && !cfg.SpriteHard {
		data = InternalTransform(out, newPalette, cfg.Size, cfg)
		lineSize = cfg.Size.Width
	} else {
		if cfg.ZigZag {
			// prepare zigzag transformation
			out = transformation.Zigzag(out)
		}
		if !cfg.SpriteHard {
			// log.GetLogger().Info( "Transform image in sprite.\n")
			data, _, lineSize, err = sprite.ToSprite(out, newPalette, cfg.Size, screenMode, cfg)
		} else {
			//	log.GetLogger().Info( "Transform image in sprite hard.\n")
			data, _ = spritehard.ToSpriteHard(out, newPalette, cfg.Size, screenMode, cfg)
			lineSize = 16
		}
	}
	if cfg.OneRow {
		for y := 0; y < out.Bounds().Max.Y; y += 2 {
			for x := 0; x < out.Bounds().Max.X; x++ {
				out.Set(x, y, newPalette[0])
			}
		}
	}
	if cfg.OneLine {
		for y := 0; y < out.Bounds().Max.Y; y++ {
			for x := 0; x < out.Bounds().Max.X; x += 2 {
				out.Set(x, y, newPalette[0])
			}
		}
	}
	return data, out, newPalette, lineSize, err
}

func ApplyKmeans(cfg iface.ImageIface) error {
	out := cfg.Img()
	log.GetLogger().Info("kmeans with %f threshold", cfg.KmeansThreshold())
	out, err := ci.Kmeans(cfg.Size().ColorsAvailable, cfg.KmeansThreshold(), out)
	if err != nil {
		return err
	}
	cfg.SetImg(out)
	return nil
}

func ApplyResize(cfg iface.ImageIface) error {
	f, err := os.Open(cfg.ImagePath())
	if err != nil {
		return err
	}
	in, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	cfg.SetImg(ci.Resize(in, cfg.Size(), cfg.ResizingAlgo()))
	return nil
}

func ApplySortPalette(cfg iface.ImageIface) error {
	p := cfg.Palette()
	var paletteToSort color.Palette
	switch cfg.Mode() {
	case 1:
		paletteToSort = p[0:4]
	case 2:
		paletteToSort = p[0:2]
	default:
		paletteToSort = p
	}
	paletteToSort = constants.FillColorPalette(paletteToSort)
	p = constants.SortColorsByDistance(paletteToSort)
	cfg.SetPalette(p)
	return nil
}

func ApplyColorReducer(cfg iface.ImageIface) error {
	out := cfg.Img()
	out = ci.Reducer(out, cfg.Reducer())
	cfg.SetImg(out)
	return nil
}

func ApplyDowngradingPalette(cfg iface.ImageIface) error {
	p := cfg.Palette()
	out := cfg.Img()
	if len(p) > 0 {
		p, out = ci.DowngradingWithPalette(out, p)
	} else {
		var err error
		p, out, err = ci.DowngradingPalette(out, cfg.Size(), cfg.CpcPlus())
		if err != nil {
			log.GetLogger().Error("Cannot downgrade colors palette for this image")
			return err
		}
	}
	cfg.SetPalette(p)
	cfg.SetImg(out)
	return nil
}

func ApplySaturationBrightness(cfg iface.ImageIface) error {
	p := cfg.Palette()
	out := cfg.Img()
	p = ci.EnhanceBrightness(p, cfg.Brightness(), cfg.Saturation())
	p, out = ci.DowngradingWithPalette(out, p)
	var paletteToSort color.Palette
	switch cfg.Mode() {
	case 1:
		paletteToSort = p[0:4]
	case 2:
		paletteToSort = p[0:2]
	default:
		paletteToSort = p
	}

	p = constants.SortColorsByDistance(
		constants.FillColorPalette(paletteToSort))
	cfg.SetPalette(p)
	cfg.SetImg(out)
	return nil
}

func ApplyZigZag(cfg iface.ImageIface) error {
	out := cfg.Img()
	// prepare zigzag transformation
	out = transformation.Zigzag(out)
	cfg.SetImg(out)
	return nil
}

func ApplyOneRow(cfg iface.ImageIface) error {
	out := cfg.Img()
	palette := cfg.Palette()
	for y := 0; y < out.Bounds().Max.Y; y += 2 {
		for x := 0; x < out.Bounds().Max.X; x++ {
			out.Set(x, y, palette[0])
		}
	}
	cfg.SetImg(out)
	return nil
}

func ApplyOneLine(cfg iface.ImageIface) error {
	out := cfg.Img()
	palette := cfg.Palette()
	for y := 0; y < out.Bounds().Max.Y; y++ {
		for x := 0; x < out.Bounds().Max.X; x += 2 {
			out.Set(x, y, palette[0])
		}
	}
	cfg.SetImg(out)
	return nil
}
