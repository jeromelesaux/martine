package gfx

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jeromelesaux/martine/convert"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/oliamb/cutter"
)

func TileMode(exportType *x.ExportType, mode uint8, iterationX, iterationY int) error {
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

	factorX := width/iterationX + 1
	factorY := height/iterationY + 1

	if factorX != factorY {
		fmt.Fprintf(os.Stdout, "factor x (%d) differs from factor y (%d)\n", factorX, factorY)
	}
	if factorY == 0 {
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

			resized := convert.Resize(cropped, exportType.Size, exportType.ResizingAlgo)
			ext := "_resized_" + strconv.Itoa(index) + ".png"
			filePath := filepath.Join(exportType.OutputPath, exportType.OsFilename(ext))
			if err := file.Png(filePath, resized); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not resized image, error %v\n", err)
			}
			p, downgraded, err := convert.DowngradingPalette(resized, exportType.Size, exportType.CpcPlus)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot downgrad the palette, error :%v\n", err)
			}

			if exportType.ZigZag {
				zizagImg := image.NewNRGBA(image.Rectangle{
					image.Point{X: 0, Y: 0},
					image.Point{X: downgraded.Bounds().Max.X, Y: downgraded.Bounds().Max.Y}})
				for y := 1; y < downgraded.Bounds().Max.Y; y += 2 {
					xZigZag := 0
					for x := downgraded.Bounds().Max.X - 1; x >= 0; x-- {
						zizagImg.Set(xZigZag, y, downgraded.At(x, y))
						xZigZag++
					}
				}
				for y := 0; y < downgraded.Bounds().Max.Y; y += 2 {
					for x := 0; x < downgraded.Bounds().Max.X; x++ {
						zizagImg.Set(x, y, downgraded.At(x, y))
					}
				}
				downgraded = zizagImg
			}

			ext = "_downgraded_" + strconv.Itoa(index) + ".png"
			filePath = filepath.Join(exportType.OutputPath, exportType.OsFilename(ext))
			if err := file.Png(filePath, downgraded); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not downgrad image, error %v\n", err)
			}
			ext = strconv.Itoa(index) + ".png"
			exportType.Size.Width = resized.Bounds().Max.X
			exportType.Size.Height = resized.Bounds().Max.Y
			if err := SpriteTransform(downgraded, p, exportType.Size, mode, exportType.OsFilename(ext), exportType); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create sprite from image, error :%v\n", err)
			}
			index++
		}
	}

	return exportType.Tiles.Save(exportType.Fullpath(".json"))
}
