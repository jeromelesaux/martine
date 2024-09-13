package spritehard

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
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

func SpriteHardToImg(in string, p color.Palette) (*image.NRGBA, error) {

	data, err := RawSpriteHard(in)
	if err != nil {
		return nil, err
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 16, Y: 16}})
	offset := 0
	for x := out.Bounds().Min.X; x < out.Bounds().Max.X; x++ {
		for y := out.Bounds().Min.Y; y < out.Bounds().Max.Y; y++ {
			if offset > len(data) {
				return nil, errors.New("sprite hard exceed the size")
			}
			px := p[data[offset]]
			out.Set(x, y, px)
			offset++
		}
	}
	return out, nil

}

func RawSpriteHard(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	header := &cpc.CpcHead{}
	if err := binary.Read(f, binary.LittleEndian, header); err != nil {
		log.GetLogger().Error("Cannot read the RawScr Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return []byte{}, err
		}
	}
	if header.Checksum != header.ComputedChecksum16() {
		log.GetLogger().Error("Cannot read the RawScr Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return []byte{}, err
		}
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
