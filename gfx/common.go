package gfx

import (
	"fmt"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

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
	if err := file.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_resized.png", out); err != nil {
		os.Exit(-2)
	}

	if exportType.DitheringAlgo != -1 {
		switch exportType.DitheringType {
		case constants.ErrorDiffusionDither:
			if exportType.DitheringWithQuantification {
				out = QuantizeWithDither(out, exportType.DitheringMatrix, exportType.Size.ColorsAvailable, newPalette)
			} else {
				out = Dithering(out, exportType.DitheringMatrix, float32(exportType.DitheringMultiplier))
			}
		case constants.OrderedDither:
			//newPalette = convert.PaletteUsed(out,exportType.CpcPlus)
			if exportType.CpcPlus {
				newPalette = convert.ExtractPalette(out, exportType.CpcPlus, 27)
				out = BayerDiphering(out, exportType.DitheringMatrix, newPalette)
			} else {
				out = BayerDiphering(out, exportType.DitheringMatrix, constants.CpcOldPalette)
			}
		}
	}
	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, exportType.Size, exportType.CpcPlus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", picturePath)
		}
	}

	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := file.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_down.png", downgraded); err != nil {
		os.Exit(-2)
	}

	if exportType.RollMode {
		if exportType.RotationRlaBit != -1 || exportType.RotationSlaBit != -1 {
			RollLeft(exportType.RotationRlaBit, exportType.RotationSlaBit, exportType.RotationIterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
		} else {
			if exportType.RotationRraBit != -1 || exportType.RotationSraBit != -1 {
				RollRight(exportType.RotationRraBit, exportType.RotationSraBit, exportType.RotationIterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
			}
		}
		if exportType.RotationKeephighBit != -1 || exportType.RotationLosthighBit != -1 {
			RollUp(exportType.RotationKeephighBit, exportType.RotationLosthighBit, exportType.RotationIterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
		} else {
			if exportType.RotationKeeplowBit != -1 || exportType.RotationLostlowBit != -1 {
				RollLow(exportType.RotationKeeplowBit, exportType.RotationLostlowBit, exportType.RotationIterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
			}
		}
	}
	if exportType.RotationMode {
		if err := Rotate(downgraded, newPalette, exportType.Size, uint8(mode), picturePath, exportType.ResizingAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if exportType.Rotation3DMode {
		if err := Rotate3d(downgraded, newPalette, exportType.Size, uint8(mode), picturePath, exportType.ResizingAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if !exportType.CustomDimension {
		Transform(downgraded, newPalette, exportType.Size, picturePath, exportType)
	} else {
		fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
		SpriteTransform(downgraded, newPalette, exportType.Size, screenMode, filename, exportType)
	}
	return err
}
