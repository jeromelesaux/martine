package transformation

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
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/oliamb/cutter"
)

func TileMode(ex *x.MartineContext, mode uint8, iterationX, iterationY int) error {
	fr, err := os.Open(ex.InputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open (%s),error :%v\n", ex.InputPath, err)
		return err
	}
	defer fr.Close()

	in, _, err := image.Decode(fr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot decode image (%s) error :%v\n", ex.InputPath, err)
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
				fmt.Fprintf(os.Stderr, "Cannot crop image for (%d,%d), error :%v\n", i, y, err)
				return err
			}

			resized := convert.Resize(cropped, ex.Size, ex.ResizingAlgo)
			ext := "_resized_" + strconv.Itoa(index) + ".png"
			filePath := filepath.Join(ex.OutputPath, ex.OsFilename(ext))
			if err := file.Png(filePath, resized); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot resized image, error %v\n", err)
			}
			p, downgraded, err := convert.DowngradingPalette(resized, ex.Size, ex.CpcPlus)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot downgrad the palette, error :%v\n", err)
			}

			if ex.ZigZag {
				downgraded = Zigzag(downgraded)
			}

			ext = "_downgraded_" + strconv.Itoa(index) + ".png"
			filePath = filepath.Join(ex.OutputPath, ex.OsFilename(ext))
			if err := file.Png(filePath, downgraded); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot downgrade image, error %v\n", err)
			}
			ext = strconv.Itoa(index) + ".png"
			ex.Size.Width = resized.Bounds().Max.X
			ex.Size.Height = resized.Bounds().Max.Y
			if err := common.ToSpriteAndExport(downgraded, p, ex.Size, mode, ex.OsFilename(ext), false, ex); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create sprite from image, error :%v\n", err)
			}
			index++
		}
	}

	return ex.Tiles.Save(ex.Fullpath(".json"))
}

func Zigzag(in *image.NRGBA) *image.NRGBA {
	zizagImg := image.NewNRGBA(image.Rectangle{
		image.Point{X: 0, Y: 0},
		image.Point{X: in.Bounds().Max.X, Y: in.Bounds().Max.Y}})
	for y := 1; y < in.Bounds().Max.Y; y += 2 {
		xZigZag := 0
		for x := in.Bounds().Max.X - 1; x >= 0; x-- {
			zizagImg.Set(xZigZag, y, in.At(x, y))
			xZigZag++
		}
	}
	for y := 0; y < in.Bounds().Max.Y; y += 2 {
		for x := 0; x < in.Bounds().Max.X; x++ {
			zizagImg.Set(x, y, in.At(x, y))
		}
	}
	in = zizagImg
	return in
}
