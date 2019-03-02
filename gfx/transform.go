package gfx

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jeromelesaux/m4client/cpc"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

var (
	ColorNotFound     = errors.New("Color not found in palette")
	NotYetImplemented = errors.New("Function is not yet implemented")
)

func Transform(in *image.NRGBA, p color.Palette, size Size, filepath string) error {
	switch size {
	case Mode0:
		return TransformMode0(in, p, size, filepath)
	default:
		return NotYetImplemented
	}
	return nil
}

func PalettePosition(c color.Color, p color.Palette) (int, error) {
	r, g, b, a := c.RGBA()
	for index, cp := range p {
		//fmt.Fprintf(os.Stdout,"index(%d), c:%v,cp:%v\n",index,c,cp)
		rp, gp, bp, ap := cp.RGBA()
		if r == rp && g == gp && b == bp && a == ap {
			//fmt.Fprintf(os.Stdout,"Position found")
			return index, nil
		}
	}
	return -1, ColorNotFound
}

func TransformMode0(in *image.NRGBA, p color.Palette, size Size, filePath string) error {
	bw := make([]byte, 0X4000)
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 2 {

			c1 := in.At(x, y)
			pp1, err := PalettePosition(c1, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, x, y)
				pp1 = 0
			}
			firmwareColorUsed[pp1]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c2 := in.At(x+1, y)
			pp2, err := PalettePosition(c2, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, x+1, y)
				pp2 = 0
			}

			firmwareColorUsed[pp2]++

			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x+1, y+j, c1, pp2)
			var pixel byte
			//fmt.Fprintf(os.Stderr,"1:(%.8b)2:(%.8b)4:(%.8b)8:(%.8b)\n",1,2,4,8)
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&1:%.8b\n",uint8(pp1)&1)
			if uint8(pp1)&1 == 1 {
				pixel += 128
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&2:%.8b\n",uint8(pp1)&2)
			if uint8(pp1)&2 == 2 {
				pixel += 8
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&4:%.8b\n",uint8(pp1)&4)
			if uint8(pp1)&4 == 4 {
				pixel += 32
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&8:%.8b\n",uint8(pp1)&8)
			if uint8(pp1)&8 == 8 {
				pixel += 2
			}
			if uint8(pp2)&1 == 1 {
				pixel += 64
			}
			if uint8(pp2)&2 == 2 {
				pixel += 4
			}
			if uint8(pp2)&4 == 4 {
				pixel += 16
			}
			if uint8(pp2)&8 == 8 {
				pixel++
			}
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			//pixel = (uint8(pp1)&8)/8 | ((uint8(pp1)&4)*4) | ((uint8(pp1)&2)*2) | ((uint8(pp1)&1)*64) | ((uint8(pp2)&8)/4) | ((uint8(pp2)&4)*8) | ((uint8(pp2)&2)*4) | ((uint8(pp2)&1)*128)
			//pixel = (uint8(pp2) & 128)>>7  + (uint8(pp1) & 32)>>4  + (uint8(pp2) & 8)>>1 + (uint8(pp1) & 2)<<2 +
			// (uint8(pp2) & 64 )>>6 + (uint8(pp1) & 16)>>3  + (uint8(pp2) & 4) + (uint8(pp1) & 1)<<3
			bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/2)] = pixel
			//bw = append(bw, pixel)
		}

	}

	fmt.Println(firmwareColorUsed)
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0xc000, Exec: 0xC7D0,
		Size: uint16(binary.Size(bw)), Size2: uint16(binary.Size(bw)), LogicalSize: uint16(binary.Size(bw))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename, extension, ".SCR", -1))
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	binary.Write(fw, binary.LittleEndian, header)
	binary.Write(fw, binary.LittleEndian, bw)
	fw.Close()

	return nil
}
