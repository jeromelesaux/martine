package spritehard

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	pal "github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/log"
)

func ToSpriteHardAndExport(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, filename string, ex *config.MartineConfig) error {

	data, firmwareColorUsed := ToSpriteHard(in, p, size, mode, ex)
	log.GetLogger().Infoln(firmwareColorUsed)
	return sprite.ExportSprite(data, 16, p, size, mode, filename, false, ex)
}

func ToSpriteHard(in *image.NRGBA, p color.Palette, size constants.Size, mode uint8, ex *config.MartineConfig) (data []byte, firmwareColorUsed map[int]int) {
	size.Height = in.Bounds().Max.Y
	size.Width = in.Bounds().Max.X
	firmwareColorUsed = make(map[int]int)
	offset := 0
	data = make([]byte, 256)
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			pp, err := pal.PalettePosition(c, p)
			if err != nil {
				// log.GetLogger().Error("%v pixel position(%d,%d) not found in palette\n", c, x, y)
				pp = 0
			}
			firmwareColorUsed[pp]++
			data[offset] = byte(ex.SwapInk(pp))
			offset++
		}
	}
	return data, firmwareColorUsed
}
