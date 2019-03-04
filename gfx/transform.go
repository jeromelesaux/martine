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
		return TransformMode0(in, p, size, filepath, false)
	case Mode1:
		return TransformMode1(in, p, size, filepath, false)
	case Mode2:
		return TransformMode2(in, p, size, filepath, false)
	case OverscanMode0:
		return TransformMode0(in, p, size, filepath, true)
	case OverscanMode1:
		return TransformMode1(in, p, size, filepath, true)
	case OverscanMode2:
		return TransformMode2(in, p, size, filepath, true)
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

func TransformMode0(in *image.NRGBA, p color.Palette, size Size, filePath string, overscan bool) error {
	bw := make([]byte, 0x4000)
	if overscan {
		bw = make([]byte, 0x8000)
	}
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
			if overscan {
				bw[(0x800*(y%8))+(0x60*(y/8))+((x+1)/2)] = pixel
			} else {
				bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/2)] = pixel
			}
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

func TransformMode1(in *image.NRGBA, p color.Palette, size Size, filePath string, overscan bool) error {
	bw := make([]byte, 0x4000)
	if overscan {
		bw = make([]byte, 0x8000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 4 {

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
			c3 := in.At(x+2, y)
			pp3, err := PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := in.At(x+3, y)
			pp4, err := PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++

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
			if uint8(pp2)&4 == 4 {
				pixel += 32
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&8:%.8b\n",uint8(pp1)&8)
			if uint8(pp2)&8 == 8 {
				pixel += 2
			}
			if uint8(pp3)&1 == 1 {
				pixel += 64
			}
			if uint8(pp3)&2 == 2 {
				pixel += 4
			}
			if uint8(pp4)&4 == 4 {
				pixel += 16
			}
			if uint8(pp4)&8 == 8 {
				pixel++
			}
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			if overscan {
				bw[(0x800*(y%8))+(0x60*(y/8))+((x+1)/4)] = pixel
			} else {
				bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/4)] = pixel
			}
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

func TransformMode2(in *image.NRGBA, p color.Palette, size Size, filePath string, overscan bool) error {
	bw := make([]byte, 0x4000)
	if overscan {
		bw = make([]byte, 0x8000)
	}
	firmwareColorUsed := make(map[int]int, 0)
	fmt.Fprintf(os.Stdout, "Informations palette (%d) for image (%d,%d)\n", len(p), in.Bounds().Max.X, in.Bounds().Max.Y)
	fmt.Println(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x += 8 {

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
			c3 := in.At(x+2, y)
			pp3, err := PalettePosition(c3, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, x+2, y)
				pp3 = 0
			}
			firmwareColorUsed[pp3]++
			c4 := in.At(x+3, y)
			pp4, err := PalettePosition(c4, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c4, x+3, y)
				pp4 = 0
			}
			firmwareColorUsed[pp4]++
			c5 := in.At(x+4, y)
			pp5, err := PalettePosition(c5, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, x+4, y)
				pp5 = 0
			}
			firmwareColorUsed[pp5]++
			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x, y+j, c1, pp1)
			c6 := in.At(x+5, y)
			pp6, err := PalettePosition(c6, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, x+5, y)
				pp6 = 0
			}
			firmwareColorUsed[pp6]++
			c7 := in.At(x+6, y)
			pp7, err := PalettePosition(c7, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, x+6, y)
				pp3 = 0
			}
			firmwareColorUsed[pp7]++
			c8 := in.At(x+7, y)
			pp8, err := PalettePosition(c8, p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, x+7, y)
				pp8 = 0
			}
			firmwareColorUsed[pp8]++

			//fmt.Fprintf(os.Stdout, "(%d,%d), %v, position palette %d\n", x+1, y+j, c1, pp2)
			var pixel byte
			//fmt.Fprintf(os.Stderr,"1:(%.8b)2:(%.8b)4:(%.8b)8:(%.8b)\n",1,2,4,8)
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&1:%.8b\n",uint8(pp1)&1)
			if uint8(pp1)&1 == 1 {
				pixel += 128
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&2:%.8b\n",uint8(pp1)&2)
			if uint8(pp2)&2 == 2 {
				pixel += 64
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&4:%.8b\n",uint8(pp1)&4)
			if uint8(pp3)&4 == 4 {
				pixel += 32
			}
			//fmt.Fprintf(os.Stderr,"uint8(pp1)&8:%.8b\n",uint8(pp1)&8)
			if uint8(pp4)&8 == 8 {
				pixel += 16
			}
			if uint8(pp5)&16 == 16 {
				pixel += 8
			}
			if uint8(pp6)&32 == 32 {
				pixel += 4
			}
			if uint8(pp7)&64 == 64 {
				pixel += 2
			}
			if uint8(pp8)&128 == 128 {
				pixel++
			}
			//fmt.Fprintf(os.Stdout, "x(%d), y(%d), pp1(%.8b), pp2(%.8b) pixel(%.8b)(%d)(&%.2x)\n", x, y, pp1, pp2, pixel, pixel, pixel)
			// MACRO PIXM0 COL2,COL1
			// ({COL1}&8)/8 | (({COL1}&4)*4) | (({COL1}&2)*2) | (({COL1}&1)*64) | (({COL2}&8)/4) | (({COL2}&4)*8) | (({COL2}&2)*4) | (({COL2}&1)*128)
			//	MEND
			if overscan {
				bw[(0x800*(y%8))+(0x60*(y/8))+((x+1)/8)] = pixel
			} else {
				bw[(0x800*(y%8))+(0x50*(y/8))+((x+1)/8)] = pixel
			}
			//bw = append(bw, pixel)
		}

	}

	fmt.Println(firmwareColorUsed)
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0xc000, Exec: 0xC7D0,
		Size: uint16(binary.Size(bw)), Size2: uint16(binary.Size(bw)), LogicalSize: uint16(binary.Size(bw))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename[:8], extension, ".SCR", -1))
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
