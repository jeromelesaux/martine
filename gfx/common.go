package gfx

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/common"
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
				p = convert.ExtractPalette(in, isCpcPlus, 27)
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
		palette, _, err = ocpartstudio.OpenPal(cont.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if cont.InkPath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", cont.InkPath)
		palette, _, err = impPalette.OpenInk(cont.InkPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.InkPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if cont.KitPath != "" {
		fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", cont.KitPath)
		palette, _, err = ocpartstudio.OpenKit(cont.KitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", cont.KitPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}

	out := convert.Resize(in, cont.Size, cont.ResizingAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
	if err := png.Png(filepath.Join(cont.OutputPath, filename+"_resized.png"), out); err != nil {
		os.Exit(-2)
	}

	if cont.Reducer > 0 {
		out = convert.Reducer(out, cont.Reducer)
		if err := png.Png(filepath.Join(cont.OutputPath, filename+"_resized.png"), out); err != nil {
			os.Exit(-2)
		}
	}

	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, cont.Size, cont.CpcPlus)
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

	out, _ = DoDithering(out, newPalette, cont.DitheringAlgo, cont.DitheringType, cont.DitheringWithQuantification, cont.DitheringMatrix, float32(cont.DitheringMultiplier), cont.CpcPlus, cont.Size)
	if cont.Saturation > 0 || cont.Brightness > 0 {
		palette = convert.EnhanceBrightness(newPalette, cont.Brightness, cont.Saturation)
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
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
	if err := png.Png(filepath.Join(cont.OutputPath, filename+"_down.png"), downgraded); err != nil {
		os.Exit(-2)
	}

	images, err := DoTransformation(downgraded, newPalette,
		screenMode, cont.RollMode, cont.RotationMode, cont.Rotation3DMode,
		cont.RotationRlaBit, cont.RotationSlaBit, cont.RotationRraBit, cont.RotationSraBit,
		cont.RotationKeephighBit, cont.RotationLosthighBit,
		cont.RotationKeeplowBit, cont.RotationLostlowBit, cont.RotationIterations,
		cont.RollIteration, cont.Rotation3DX0, cont.Rotation3DY0, cont.Rotation3DType, cont.ResizingAlgo, cont.Size)
	if err != nil {
		os.Exit(-2)
	} else {

		for indice := 0; indice < cont.RollIteration; indice++ {
			img := images[indice]
			newFilename := cont.OsFullPath(filename, fmt.Sprintf("%.2d", indice)+".png")
			if err := png.Png(newFilename, img); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create image (%s) error :%v\n", newFilename, err)
			}
			if err := common.ToSpriteAndExport(img, newPalette, constants.Size{Width: cont.Size.Width, Height: cont.Size.Height}, screenMode, newFilename, false, cont); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create sprite image (%s) error %v\n", newFilename, err)
			}
		}
	}

	if !cont.CustomDimension && !cont.SpriteHard {
		Transform(downgraded, newPalette, cont.Size, picturePath, cont)
	} else {
		if cont.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cont.SpriteHard {
			//fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			common.ToSpriteAndExport(downgraded, newPalette, cont.Size, screenMode, filename, false, cont)
		} else {
			//fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			common.ToSpriteHardAndExport(downgraded, newPalette, cont.Size, screenMode, filename, cont)
		}
	}
	return err
}

func ApplyOneImage(in image.Image,
	cont *export.MartineContext,
	mode int,
	palette color.Palette,
	screenMode uint8) ([]byte, *image.NRGBA, color.Palette, int, error) {

	var newPalette color.Palette
	var downgraded *image.NRGBA
	var err error

	out := convert.Resize(in, cont.Size, cont.ResizingAlgo)

	if cont.Reducer > -1 {
		out = convert.Reducer(out, cont.Reducer)
	}

	if len(palette) > 0 {
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = convert.DowngradingPalette(out, cont.Size, cont.CpcPlus)
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
	out, _ = DoDithering(out, newPalette, cont.DitheringAlgo, cont.DitheringType, cont.DitheringWithQuantification, cont.DitheringMatrix, float32(cont.DitheringMultiplier), cont.CpcPlus, cont.Size)

	if cont.Saturation > 0 || cont.Brightness > 0 {
		palette = convert.EnhanceBrightness(newPalette, cont.Brightness, cont.Saturation)
		newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
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
	if !cont.CustomDimension && !cont.SpriteHard {
		data = InternalTransform(downgraded, newPalette, cont.Size, cont)
		lineSize = cont.Size.Width
	} else {
		if cont.ZigZag {
			// prepare zigzag transformation
			downgraded = transformation.Zigzag(downgraded)
		}
		if !cont.SpriteHard {
			//fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
			data, _, lineSize, err = common.ToSprite(downgraded, newPalette, cont.Size, screenMode, cont)
		} else {
			//	fmt.Fprintf(os.Stdout, "Transform image in sprite hard.\n")
			data, _ = common.ToSpriteHard(downgraded, newPalette, cont.Size, screenMode, cont)
			lineSize = 16
		}
	}
	if cont.OneRow {
		for y := 0; y < downgraded.Bounds().Max.Y; y += 2 {
			for x := 0; x < downgraded.Bounds().Max.X; x++ {
				downgraded.Set(x, y, newPalette[0])
			}
		}
	}
	if cont.OneLine {
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
