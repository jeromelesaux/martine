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

func DoDithering(in *image.NRGBA, p color.Palette, exportType *export.ExportType) (*image.NRGBA, color.Palette) {
	if exportType.DitheringAlgo != -1 {
		switch exportType.DitheringType {
		case constants.ErrorDiffusionDither:
			if exportType.DitheringWithQuantification {
				in = filter.QuantizeWithDither(in, exportType.DitheringMatrix, exportType.Size.ColorsAvailable, p)
			} else {
				in = filter.Dithering(in, exportType.DitheringMatrix, float32(exportType.DitheringMultiplier))
			}
		case constants.OrderedDither:
			//newPalette = convert.PaletteUsed(out,exportType.CpcPlus)
			if exportType.CpcPlus {
				p = convert.ExtractPalette(in, exportType.CpcPlus, 27)
				in = filter.BayerDiphering(in, exportType.DitheringMatrix, p)
			} else {
				in = filter.BayerDiphering(in, exportType.DitheringMatrix, constants.CpcOldPalette)
			}
		}
	}
	return in, p
}

func DoTransformation(in *image.NRGBA, p color.Palette, filename, picturePath string, screenMode uint8, mode int, exportType *export.ExportType) error {
	var err error
	if exportType.RollMode {
		if exportType.RotationRlaBit != -1 || exportType.RotationSlaBit != -1 {
			transformation.RollLeft(exportType.RotationRlaBit, exportType.RotationSlaBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
		} else {
			if exportType.RotationRraBit != -1 || exportType.RotationSraBit != -1 {
				transformation.RollRight(exportType.RotationRraBit, exportType.RotationSraBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
			}
		}
		if exportType.RotationKeephighBit != -1 || exportType.RotationLosthighBit != -1 {
			transformation.RollUp(exportType.RotationKeephighBit, exportType.RotationLosthighBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
		} else {
			if exportType.RotationKeeplowBit != -1 || exportType.RotationLostlowBit != -1 {
				transformation.RollLow(exportType.RotationKeeplowBit, exportType.RotationLostlowBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
			}
		}
	}
	if exportType.RotationMode {
		if err = transformation.Rotate(in, p, exportType.Size, uint8(mode), picturePath, exportType.ResizingAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if exportType.Rotation3DMode {
		if err = transformation.Rotate3d(in, p, exportType.Size, uint8(mode), picturePath, exportType.ResizingAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	return err
}

func ApplyOneImageAndExport(in image.Image,
	exportType *export.ExportType,
	filename, picturePath string,
	mode int,
	screenMode uint8) error {

	var palette color.Palette
	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	if exportType.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", exportType.PalettePath)
		palette, _, err = file.OpenPal(exportType.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", exportType.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if exportType.InkPath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", exportType.InkPath)
		palette, _, err = file.OpenInk(exportType.InkPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", exportType.InkPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if exportType.KitPath != "" {
		fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", exportType.KitPath)
		palette, _, err = file.OpenKit(exportType.KitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", exportType.KitPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	out := convert.Resize(in, exportType.Size, exportType.ResizingAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := file.Png(filepath.Join(exportType.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}

	if exportType.Reducer > -1 {
		out = convert.Reducer(out, exportType.Reducer)
		if err := file.Png(filepath.Join(exportType.OutputPath, filename+"_resized.png"), out); err != nil {
			os.Exit(-2)
		}
	}

	if newPalette == nil { // in case of dithering without input palette
		if exportType.CpcPlus {
			newPalette = constants.CpcPlusPalette
		} else {
			newPalette = constants.CpcOldPalette
		}
	}
	out, _ = DoDithering(out, newPalette, exportType)

	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, exportType.Size, exportType.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}

	newPalette = constants.SortColorsByDistance(newPalette)

	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := file.Png(filepath.Join(exportType.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	DoTransformation(downgraded, newPalette, filename, picturePath, screenMode, mode, exportType)

	if !exportType.CustomDimension && !exportType.SpriteHard {
		Transform(downgraded, newPalette, exportType.Size, picturePath, exportType)
	} else {
		if exportType.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !exportType.SpriteHard {
			fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			common.ToSpriteAndExport(downgraded, newPalette, exportType.Size, screenMode, filename, false, exportType)
		} else {
			fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			common.ToSpriteHardAndExport(downgraded, newPalette, exportType.Size, screenMode, filename, exportType)
		}
	}
	return err
}

func ApplyOneImage(in image.Image,
	exportType *export.ExportType,
	mode int,
	palette color.Palette,
	screenMode uint8) ([]byte, color.Palette, int, error) {

	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	out := convert.Resize(in, exportType.Size, exportType.ResizingAlgo)

	if exportType.Reducer > -1 {
		out = convert.Reducer(out, exportType.Reducer)
	}

	if newPalette == nil { // in case of dithering without input palette
		if exportType.CpcPlus {
			newPalette = constants.CpcPlusPalette
		} else {
			newPalette = constants.CpcOldPalette
		}
	}
	out, _ = DoDithering(out, newPalette, exportType)

	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, exportType.Size, exportType.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image")
		}
	}

	newPalette = constants.SortColorsByDistance(newPalette)

	var data []byte
	var lineSize int
	if !exportType.CustomDimension && !exportType.SpriteHard {
		data = InternalTransform(downgraded, newPalette, exportType.Size, exportType)
		lineSize = exportType.Size.Width
	} else {
		if exportType.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !exportType.SpriteHard {
			fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			data, _, lineSize, err = common.ToSprite(downgraded, newPalette, exportType.Size, screenMode, exportType)
		} else {
			fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			data, _ = common.ToSpriteHard(downgraded, newPalette, exportType.Size, screenMode, exportType)
			lineSize = 16
		}
	}
	return data, newPalette, lineSize, err
}
