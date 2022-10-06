package file

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
	x "github.com/jeromelesaux/martine/export"
)

type SpriteHard struct {
	Data [256]byte
}

type SprImpdraw struct {
	Data []SpriteHard
}

func (s *SprImpdraw) Images(pal color.Palette) []*image.NRGBA {
	imgs := make([]*image.NRGBA, 0)
	for _, v := range s.Data {
		imgs = append(imgs, v.Image(pal))
	}
	return imgs
}

func (s *SpriteHard) Image(pal color.Palette) *image.NRGBA {
	img := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{16, 16}})
	var index int
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			var c color.Color = color.Black
			if index < len(s.Data) && int(s.Data[index]) < len(pal) {
				c = pal[int(s.Data[index])]
			}
			img.Set(x, y, c)
			index++
		}
	}
	return img
}

func OpenSpr(filePath string) (*SprImpdraw, error) {
	spr := SprImpdraw{Data: make([]SpriteHard, 0)}
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return &spr, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Sprite Hard Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Sprite Hard Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	for {
		spriteHard := SpriteHard{}
		if err = binary.Read(fr, binary.LittleEndian, &spriteHard); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
				return &spr, err
			}
		}
		spr.Data = append(spr.Data, spriteHard)
	}
	return &spr, nil
}

func Spr(filePath string, spr SprImpdraw, cont *x.MartineContext) error {
	osFilename := cont.AmsdosFullPath(filePath, ".SPR")
	fmt.Fprintf(os.Stdout, "Saving SPR file (%s)\n", osFilename)
	content := make([]byte, 0)
	for _, v := range spr.Data {
		content = append(content, v.Data[:]...)
	}
	content, _ = Compress(content, cont.Compression)
	ext := ".SPR"
	if cont.Compression != -1 {
		ext = ".SPR.zxo"
	}
	if !cont.NoAmsdosHeader {
		if err := SaveAmsdosFile(osFilename, ext, content, 2, 0, 0x0, 0x4000); err != nil {
			return err
		}
	} else {
		if err := SaveOSFile(osFilename, content); err != nil {
			return err
		}
	}

	return nil
}
