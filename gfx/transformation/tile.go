package transformation

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"

	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/oliamb/cutter"
)

// nolint: funlen
func TileMode(ex *config.MartineConfig, mode uint8, iterationX, iterationY int) error {
	fr, err := os.Open(ex.ScrCfg.InputPath)
	if err != nil {
		log.GetLogger().Error("Cannot open (%s),error :%v\n", ex.ScrCfg.InputPath, err)
		return err
	}
	defer fr.Close()

	in, _, err := image.Decode(fr)
	if err != nil {
		log.GetLogger().Error("Cannot decode image (%s) error :%v\n", ex.ScrCfg.InputPath, err)
		return err
	}

	factorX := in.Bounds().Max.X/iterationX + 1
	factorY := in.Bounds().Max.Y/iterationY + 1

	if factorX != factorY {
		log.GetLogger().Info("factor x (%d) differs from factor y (%d)\n", factorX, factorY)
	}
	if factorY == 0 {
		factorY = in.Bounds().Max.Y
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
				log.GetLogger().Error("Cannot crop image for (%d,%d), error :%v\n", i, y, err)
				return err
			}

			resized := ci.Resize(cropped, ex.ScrCfg.Size, ex.ScrCfg.Treatment.ResizingAlgo)
			ext := "_resized_" + strconv.Itoa(index) + ".png"
			filePath := filepath.Join(ex.ScrCfg.OutputPath, ex.OsFilename(ext))
			if err := png.Png(filePath, resized); err != nil {
				log.GetLogger().Error("Cannot resized image, error %v\n", err)
			}
			p, downgraded, err := ci.DowngradingPalette(resized, ex.ScrCfg.Size, ex.ScrCfg.IsPlus)
			if err != nil {
				log.GetLogger().Error("Cannot downgrad the palette, error :%v\n", err)
			}

			if ex.ZigZag {
				downgraded = Zigzag(downgraded)
			}

			ext = "_downgraded_" + strconv.Itoa(index) + ".png"
			filePath = filepath.Join(ex.ScrCfg.OutputPath, ex.OsFilename(ext))
			if err := png.Png(filePath, downgraded); err != nil {
				log.GetLogger().Error("Cannot downgrade image, error %v\n", err)
			}
			ext = strconv.Itoa(index) + ".png"
			ex.ScrCfg.Size.Width = resized.Bounds().Max.X
			ex.ScrCfg.Size.Height = resized.Bounds().Max.Y
			if err := sprite.ToSpriteAndExport(downgraded, p, ex.ScrCfg.Size, mode, ex.OsFilename(ext), false, ex); err != nil {
				log.GetLogger().Error("Cannot create sprite from image, error :%v\n", err)
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
