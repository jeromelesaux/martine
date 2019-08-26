package gfx

import (
	"fmt"
	"github.com/disintegration/imaging"
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
	resizeAlgo imaging.ResampleFilter,
	exportType *export.ExportType,
	filename, picturePath, palettePath, inkPath, kitPath string,
	mode, ditheringAlgo, rla, sla, rra, sra, keephigh, keeplow, losthigh, lostlow, iterations int,
	screenMode uint8,
	ditheringMultiplier float64,
	ditheringMatrix [][]float32,
	ditherType DitheringType,
	customDimension, withQuantization bool) error {

	var palette color.Palette
	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	if palettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", palettePath)
		palette, _, err = file.OpenPal(palettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", palettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if inkPath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", inkPath)
		palette, _, err = file.OpenInk(inkPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", inkPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if kitPath != "" {
		fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", kitPath)
		palette, _, err = file.OpenKit(kitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", palettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	out := convert.Resize(in, exportType.Size, resizeAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := file.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_resized.png", out); err != nil {
		os.Exit(-2)
	}

	if ditheringAlgo != -1 {
		switch ditherType {
		case ErrorDiffusionDither:
			if withQuantization {
				out = QuantizeWithDither(out, ditheringMatrix, exportType.Size.ColorsAvailable, newPalette)
			} else {
				out = Dithering(out, ditheringMatrix, float32(ditheringMultiplier))
			}
		case OrderedDither:
			//newPalette = convert.PaletteUsed(out,exportType.CpcPlus)
			if exportType.CpcPlus {
				newPalette = convert.ExtractPalette(out, exportType.CpcPlus, 27)
				out = BayerDiphering(out, ditheringMatrix, newPalette)
			} else {
				out = BayerDiphering(out, ditheringMatrix, constants.CpcOldPalette)
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
		if rla != -1 || sla != -1 {
			RollLeft(rla, sla, iterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
		} else {
			if rra != -1 || sra != -1 {
				RollRight(rra, sra, iterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
			}
		}
		if keephigh != -1 || losthigh != -1 {
			RollUp(keephigh, losthigh, iterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
		} else {
			if keeplow != -1 || lostlow != -1 {
				RollLow(keeplow, lostlow, iterations, screenMode, exportType.Size, downgraded, newPalette, filename, exportType)
			}
		}
	}
	if exportType.RotationMode {
		if err := Rotate(downgraded, newPalette, exportType.Size, uint8(mode), picturePath, resizeAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if exportType.Rotation3DMode {
		if err := Rotate3d(downgraded, newPalette, exportType.Size, uint8(mode), picturePath, resizeAlgo, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", picturePath, err)
		}
	}
	if !customDimension {
		Transform(downgraded, newPalette, exportType.Size, picturePath, exportType)
	} else {
		fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
		SpriteTransform(downgraded, newPalette, exportType.Size, screenMode, filename, exportType)
	}
	return err
}
