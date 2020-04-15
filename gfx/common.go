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
)

func DoDithering(in *image.NRGBA, p color.Palette, exportType *export.ExportType) (*image.NRGBA, color.Palette) {
	if exportType.DitheringAlgo != -1 {
		switch exportType.DitheringType {
		case constants.ErrorDiffusionDither:
			if exportType.DitheringWithQuantification {
				in = QuantizeWithDither(in, exportType.DitheringMatrix, exportType.Size.ColorsAvailable, p)
			} else {
				in = Dithering(in, exportType.DitheringMatrix, float32(exportType.DitheringMultiplier))
			}
		case constants.OrderedDither:
			//newPalette = convert.PaletteUsed(out,exportType.CpcPlus)
			if exportType.CpcPlus {
				p = convert.ExtractPalette(in, exportType.CpcPlus, 27)
				in = BayerDiphering(in, exportType.DitheringMatrix, p)
			} else {
				in = BayerDiphering(in, exportType.DitheringMatrix, constants.CpcOldPalette)
			}
		}
	}
	return in, p
}

func DoTransformation(in *image.NRGBA, p color.Palette, filename, picturePath string, screenMode uint8, mode int, exportType *export.ExportType) error {
	var err error
	if exportType.RollMode {
		if exportType.RotationRlaBit != -1 || exportType.RotationSlaBit != -1 {
			RollLeft(exportType.RotationRlaBit, exportType.RotationSlaBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
		} else {
			if exportType.RotationRraBit != -1 || exportType.RotationSraBit != -1 {
				RollRight(exportType.RotationRraBit, exportType.RotationSraBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
			}
		}
		if exportType.RotationKeephighBit != -1 || exportType.RotationLosthighBit != -1 {
			RollUp(exportType.RotationKeephighBit, exportType.RotationLosthighBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
		} else {
			if exportType.RotationKeeplowBit != -1 || exportType.RotationLostlowBit != -1 {
				RollLow(exportType.RotationKeeplowBit, exportType.RotationLostlowBit, exportType.RotationIterations, screenMode, exportType.Size, in, p, filename, exportType)
			}
		}
	}
	if exportType.RotationMode {
		if err = Rotate(in, p, exportType.Size, uint8(mode), picturePath, exportType.ResizingAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if exportType.Rotation3DMode {
		if err = Rotate3d(in, p, exportType.Size, uint8(mode), picturePath, exportType.ResizingAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	return err
}

func ApplyOneImage(in image.Image,
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
	out, newPalette = DoDithering(out, newPalette, exportType)

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
		if !exportType.SpriteHard {
			fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			SpriteTransform(downgraded, newPalette, exportType.Size, screenMode, filename, false, exportType)
		} else {
			fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			SpriteHardTransform(downgraded, newPalette, exportType.Size, screenMode, filename, exportType)
		}
	}
	return err
}
