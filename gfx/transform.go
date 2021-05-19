package gfx

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func Transform(in *image.NRGBA, p color.Palette, size constants.Size, filepath string, exportType *x.ExportType) error {
	switch size {
	case constants.Mode0:
		return TransformMode0(in, p, size, filepath, exportType)
	case constants.Mode1:
		return TransformMode1(in, p, size, filepath, exportType)
	case constants.Mode2:
		return TransformMode2(in, p, size, filepath, exportType)
	case constants.OverscanMode0:
		return TransformMode0(in, p, size, filepath, exportType)
	case constants.OverscanMode1:
		return TransformMode1(in, p, size, filepath, exportType)
	case constants.OverscanMode2:
		return TransformMode2(in, p, size, filepath, exportType)
	default:
		return errors.ErrorNotYetImplemented
	}
}

func InternalTransform(in *image.NRGBA, p color.Palette, size constants.Size, exportType *x.ExportType) []byte {
	switch size {
	case constants.Mode0:
		return common.ToMode0(in, p, exportType)
	case constants.Mode1:
		return common.ToMode1(in, p, exportType)
	case constants.Mode2:
		return common.ToMode2(in, p, exportType)
	case constants.OverscanMode0:
		return common.ToMode0(in, p, exportType)
	case constants.OverscanMode1:
		return common.ToMode1(in, p, exportType)
	case constants.OverscanMode2:
		return common.ToMode2(in, p, exportType)
	default:
		return []byte{}
	}
}

func SpriteHardTransformAndSave(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, exportType *x.ExportType) error {

	data, firmwareColorUsed := common.ToSpriteHard(in, p, size, mode, exportType)
	fmt.Println(firmwareColorUsed)
	if err := file.Win(filename, data, mode, 16, size.Height, false, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if !exportType.CpcPlus {
		if err := file.Pal(filename, p, mode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettepal.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Ink(filename, p, mode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath = exportType.OsFullPath(filename, "_paletteink.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := file.Kit(filename, p, mode, false, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettekit.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if err := file.Ascii(filename, data, p, false, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return file.AsciiByColumn(filename, data, p, false, mode, exportType)
}

func SpriteTransformAndSave(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, exportType *x.ExportType) error {

	data, firmwareColorUsed, lineSize, err := common.ToSprite(in, p, size, mode, exportType)
	if err != nil {
		return err
	}
	fmt.Println(firmwareColorUsed)
	if err := file.Win(filename, data, mode, lineSize, size.Height, dontImportDsk, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if !exportType.CpcPlus {
		if err := file.Pal(filename, p, mode, dontImportDsk, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettepal.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Ink(filename, p, 2, dontImportDsk, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath = exportType.OsFullPath(filename, "_paletteink.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := file.Kit(filename, p, mode, dontImportDsk, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := exportType.OsFullPath(filename, "_palettekit.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if err := file.Ascii(filename, data, p, dontImportDsk, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return file.AsciiByColumn(filename, data, p, dontImportDsk, mode, exportType)
}

func TransformMode0(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *x.ExportType) error {
	bw := common.ToMode0(in, p, exportType)
	return common.Export(filePath, bw, p, 0, exportType)
}

func TransformMode1(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *x.ExportType) error {
	bw := common.ToMode1(in, p, exportType)
	return common.Export(filePath, bw, p, 1, exportType)
}

func TransformMode2(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, exportType *x.ExportType) error {
	bw := common.ToMode2(in, p, exportType)
	return common.Export(filePath, bw, p, 2, exportType)
}

func revertColor(rawColor uint8, index int, isPlus bool) color.Color {
	var newColor color.Color
	var err error
	if !isPlus {
		newColor, err = constants.ColorFromHardware(rawColor)
		if err != nil {
			fmt.Fprintf(os.Stderr, "No color found in data at index %d\n", index)
			return constants.White.Color
		}
	} else {
		plusColor := constants.NewRawCpcPlusColor(uint16(rawColor))
		c := color.RGBA{A: 0xFF, R: uint8(plusColor.R), G: uint8(plusColor.G), B: uint8(plusColor.B)}
		newColor = constants.CpcPlusPalette.Convert(c)
	}
	return newColor
}

func TransformRawCpcData(data, palette []int, width, height int, mode int, isPlus bool) (*image.NRGBA, error) {

	in := image.NewNRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: width, Y: height}})
	x := 0
	y := 0
	for index, val := range data {

		switch mode {
		case 0:
			p1, p2 := common.RawPixelMode0(byte(val))
			c1 := palette[p1]
			newColor := revertColor(uint8(c1), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c2 := palette[p2]
			newColor = revertColor(uint8(c2), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
		case 1:
			p1, p2, p3, p4 := common.RawPixelMode1(byte(val))
			c1 := palette[p1]
			newColor := revertColor(uint8(c1), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c2 := palette[p2]
			newColor = revertColor(uint8(c2), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c3 := palette[p3]
			newColor = revertColor(uint8(c3), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c4 := palette[p4]
			newColor = revertColor(uint8(c4), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
		case 2:
			p1, p2, p3, p4, p5, p6, p7, p8 := common.RawPixelMode2(byte(val))
			c1 := palette[p1]
			newColor := revertColor(uint8(c1), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c2 := palette[p2]
			newColor = revertColor(uint8(c2), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c3 := palette[p3]
			newColor = revertColor(uint8(c3), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c4 := palette[p4]
			newColor = revertColor(uint8(c4), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c5 := palette[p5]
			newColor = revertColor(uint8(c5), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c6 := palette[p6]
			newColor = revertColor(uint8(c6), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c7 := palette[p7]
			newColor = revertColor(uint8(c7), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
			c8 := palette[p8]
			newColor = revertColor(uint8(c8), index, isPlus)
			in.Set(x, y, newColor)
			x++
			if (x % width) == 0 {
				x = 0
				y++
			}
		}
	}
	return in, nil
}
