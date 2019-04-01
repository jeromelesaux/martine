package gfx

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/convert"
	"github.com/oliamb/cutter"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
)

func TileMode(exportType *ExportType, mode uint8, iteration int, algo imaging.ResampleFilter) error {
	fr, err := os.Open(exportType.InputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open (%s),error :%v\n", exportType.InputPath, err)
		return err
	}
	defer fr.Close()

	in, _, err := image.Decode(fr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot decode image (%s) error :%v\n", exportType.InputPath, err)
		return err
	}

	width := in.Bounds().Max.X
	height := in.Bounds().Max.Y

	factorX := width / iteration
	factorY := height / iteration

	if factorX != factorY {
		fmt.Fprintf(os.Stdout, "factor x (%d) differs from factor y (%d)\n", factorX, factorY)
		factorY = height
	}
	index := 0
	for i := 0; i < in.Bounds().Max.X; i += factorX {
		for y := 0; y < in.Bounds().Max.Y; y += factorY {
			cropped, err := cutter.Crop(in, cutter.Config{
				Width:  factorX,
				Height: factorY,
				Anchor: image.Point{i, y},
				Mode:   cutter.TopLeft,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not crop image for (%d,%d), error :%v\n", i, y, err)
				return err
			}

			resized := convert.Resize(cropped, exportType.Size, algo)
			ext := "_resized_" + strconv.Itoa(index) + ".png"
			filePath := exportType.OutputPath + string(filepath.Separator) + exportType.OsFilename(ext)
			if err := Png(filePath, resized); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not resized image, error %v\n", err)
			}
			p, downgraded, err := convert.DowngradingPalette(resized, exportType.Size, exportType.CpcPlus)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot downgrad the palette, error :%v\n", err)
			}
			ext = "_downgraded_" + strconv.Itoa(index) + ".png"
			filePath = exportType.OutputPath + string(filepath.Separator) + exportType.OsFilename(ext)
			if err := Png(filePath, downgraded); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not downgrad image, error %v\n", err)
			}
			ext = strconv.Itoa(index) + ".png"
			if err := SpriteTransform(downgraded, p, exportType.Size, mode, exportType.OsFilename(ext), exportType); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create sprite from image, error :%v\n", err)
			}
			index++
		}
	}

	return nil
}
