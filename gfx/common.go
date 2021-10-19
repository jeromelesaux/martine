package gfx

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
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/gfx/filter"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

func DoDithering(in *image.NRGBA, p color.Palette, cont *export.MartineContext) (*image.NRGBA, color.Palette) {
	if cont.DitheringAlgo != -1 {
		switch cont.DitheringType {
		case constants.ErrorDiffusionDither:
			if cont.DitheringWithQuantification {
				in = filter.QuantizeWithDither(in, cont.DitheringMatrix, cont.Size.ColorsAvailable, p)
			} else {
				in = filter.Dithering(in, cont.DitheringMatrix, float32(cont.DitheringMultiplier))
			}
		case constants.OrderedDither:
			if cont.CpcPlus {
				p = convert.ExtractPalette(in, cont.CpcPlus, 27)
				in = filter.BayerDiphering(in, cont.DitheringMatrix, p)
			} else {
				in = filter.BayerDiphering(in, cont.DitheringMatrix, constants.CpcOldPalette)
			}
		}
	}
	return in, p
}

func DoTransformation(in *image.NRGBA, p color.Palette, filename, picturePath string, screenMode uint8, mode int, cont *export.MartineContext) error {
	var err error
	if cont.RollMode {
		if cont.RotationRlaBit != -1 || cont.RotationSlaBit != -1 {
			transformation.RollLeft(cont.RotationRlaBit, cont.RotationSlaBit, cont.RotationIterations, screenMode, cont.Size, in, p, filename, cont)
		} else {
			if cont.RotationRraBit != -1 || cont.RotationSraBit != -1 {
				transformation.RollRight(cont.RotationRraBit, cont.RotationSraBit, cont.RotationIterations, screenMode, cont.Size, in, p, filename, cont)
			}
		}
		if cont.RotationKeephighBit != -1 || cont.RotationLosthighBit != -1 {
			transformation.RollUp(cont.RotationKeephighBit, cont.RotationLosthighBit, cont.RotationIterations, screenMode, cont.Size, in, p, filename, cont)
		} else {
			if cont.RotationKeeplowBit != -1 || cont.RotationLostlowBit != -1 {
				transformation.RollLow(cont.RotationKeeplowBit, cont.RotationLostlowBit, cont.RotationIterations, screenMode, cont.Size, in, p, filename, cont)
			}
		}
	}
	if cont.RotationMode {
		if err = transformation.Rotate(in, p, cont.Size, uint8(mode), picturePath, cont.ResizingAlgo, cont); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if cont.Rotation3DMode {
		if err = transformation.Rotate3d(in, p, cont.Size, uint8(mode), picturePath, cont.ResizingAlgo, cont); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	return err
}

func ApplyOneImageAndExport(in image.Image,
	cont *export.MartineContext,
	filename, picturePath string,
	mode int,
	screenMode uint8) error {

	var palette color.Palette
	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	if cont.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cont.PalettePath)
		palette, _, err = file.OpenPal(cont.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if cont.InkPath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cont.InkPath)
		palette, _, err = file.OpenInk(cont.InkPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.InkPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if cont.KitPath != "" {
		fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", cont.KitPath)
		palette, _, err = file.OpenKit(cont.KitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.KitPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	out := convert.Resize(in, cont.Size, cont.ResizingAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := file.Png(filepath.Join(cont.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}

	if cont.Reducer > -1 {
		out = convert.Reducer(out, cont.Reducer)
		if err := file.Png(filepath.Join(cont.OutputPath, filename+"_resized.png"), out); err != nil {
			os.Exit(-2)
		}
	}

	if newPalette == nil { // in case of dithering without input palette
		if cont.CpcPlus {
			newPalette = constants.CpcPlusPalette
		} else {
			newPalette = constants.CpcOldPalette
		}
	}
	out, _ = DoDithering(out, newPalette, cont)

	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, cont.Size, cont.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}
	newPalette = constants.SortColorsByDistance(newPalette)
	if cont.Saturation > 0 || cont.Brightness > 0 {
		palette = convert.EnhanceBrightness(newPalette, cont.Brightness, cont.Saturation)
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
		newPalette = constants.SortColorsByDistance(newPalette)
	}

	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := file.Png(filepath.Join(cont.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	DoTransformation(downgraded, newPalette, filename, picturePath, screenMode, mode, cont)

	if !cont.CustomDimension && !cont.SpriteHard {
		Transform(downgraded, newPalette, cont.Size, picturePath, cont)
	} else {
		if cont.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cont.SpriteHard {
			fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			common.ToSpriteAndExport(downgraded, newPalette, cont.Size, screenMode, filename, false, cont)
		} else {
			fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			common.ToSpriteHardAndExport(downgraded, newPalette, cont.Size, screenMode, filename, cont)
		}
	}
	return err
}

func ApplyOneImage(in image.Image,
	cont *export.MartineContext,
	mode int,
	palette color.Palette,
	screenMode uint8) ([]byte, color.Palette, int, error) {

	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	out := convert.Resize(in, cont.Size, cont.ResizingAlgo)

	if cont.Reducer > -1 {
		out = convert.Reducer(out, cont.Reducer)
	}

	if newPalette == nil { // in case of dithering without input palette
		if cont.CpcPlus {
			newPalette = constants.CpcPlusPalette
		} else {
			newPalette = constants.CpcOldPalette
		}
	}
	out, _ = DoDithering(out, newPalette, cont)

	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, cont.Size, cont.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image")
		}
	}

	newPalette = constants.SortColorsByDistance(newPalette)
	if cont.Saturation > 0 || cont.Brightness > 0 {
		palette = convert.EnhanceBrightness(newPalette, cont.Brightness, cont.Saturation)
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
		newPalette = constants.SortColorsByDistance(newPalette)
	}
	var data []byte
	var lineSize int
	if !cont.CustomDimension && !cont.SpriteHard {
		data = InternalTransform(downgraded, newPalette, cont.Size, cont)
		lineSize = cont.Size.Width
	} else {
		if cont.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cont.SpriteHard {
			fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			data, _, lineSize, err = common.ToSprite(downgraded, newPalette, cont.Size, screenMode, cont)
		} else {
			fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			data, _ = common.ToSpriteHard(downgraded, newPalette, cont.Size, screenMode, cont)
			lineSize = 16
		}
	}
	return data, newPalette, lineSize, err
}
