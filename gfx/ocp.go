package gfx

import (
	"encoding/binary"
	"fmt"
	"github.com/jeromelesaux/m4client/cpc"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

var OverscanBoot = [...]byte{0x0e, 0x00, 0x0a, 0x00, 0x01, 0xc0, 0x20, 0x69, 0x4d, 0x50, 0x20, 0x76, 0x32, 0x00, 0x0d,
	0x00, 0x14, 0x00, 0xad, 0x20, 0x0e, 0x01, 0x83, 0x1c, 0xad, 0x01, 0x00, 0x00, 0x00, 0x01, 0x30,
	0x02, 0x32, 0x06, 0x22, 0x07, 0x23, 0x0c, 0x0d, 0xd0, 0x00, 0x00, 0x3f, 0xff, 0x00, 0xff, 0x77,
	0xb3, 0x51, 0xa8, 0xd4, 0x62, 0x39, 0x9c, 0x46, 0x2b, 0x15, 0x8a, 0xcd, 0xee, 0x00, 0xf3, 0x21,
	0x8d, 0x01, 0x3e, 0x06, 0x01, 0xbe, 0xbd, 0xed, 0xa3, 0x41, 0xed, 0xa3, 0x3d, 0x20, 0xf8, 0x01,
	0x00, 0x7f, 0x1e, 0x10, 0x0a, 0xed, 0x49, 0xed, 0x79, 0x0c, 0x1d, 0x20, 0xf7, 0x3a, 0xac, 0x01,
	0xfe, 0x01, 0x20, 0x22, 0x21, 0x9b, 0x01, 0x01, 0x11, 0xbc, 0x7e, 0xed, 0x79, 0x2c, 0x0d, 0x20,
	0xf9, 0x01, 0xb8, 0x7f, 0xed, 0x49, 0x21, 0x01, 0x08, 0x11, 0x00, 0x64, 0x01, 0x20, 0x00, 0xed,
	0xb0, 0x01, 0xa0, 0x7f, 0xed, 0x49, 0x21, 0xf9, 0xb7, 0xcd, 0xdd, 0xbc, 0xfb, 0xc3, 0x18, 0xbb,
	0x00, 0xc3, 0xc6, 0xba, 0xc3, 0xc1, 0xb9, 0x00, 0x00, 0xc3, 0x35, 0xba, 0x00, 0xed, 0x49, 0xd9,
	0xfb, 0xc3, 0x00, 0xbe, 0x2b, 0x00, 0x71, 0x18, 0x08, 0xc3, 0x41, 0xb9, 0xc9, 0x00, 0x00, 0x00}

type CpcPlusColor struct {
	G      byte
	R      byte
	B      byte
	Unused byte
}

func (c *CpcPlusColor) Value() uint16 {
	return uint16(c.B) + uint16(c.R<<4) + uint16(c.G<<16)
}

func NewCpcPlusColor(c color.Color) CpcPlusColor {
	r, g, b, _ := c.RGBA()
	return CpcPlusColor{G: byte(g / 64), R: byte(r / 64), B: byte(b / 64)}
}

type InkPalette struct {
	Colors [16]CpcPlusColor
}

func Overscan(filePath, dirPath string, data []byte, p color.Palette, screenMode uint8, noAmsdosHeader, isCpcPlus bool) error {
	o := make([]byte, 0x7e90-0x80)
	fmt.Fprintf(os.Stdout, "Saving overscan file (%s)\n", filePath)
	header := cpc.CpcHead{Type: 0, User: 0, Address: 0x170, Exec: 0x0,
		Size:        uint16(binary.Size(o)),
		Size2:       uint16(binary.Size(o)),
		LogicalSize: uint16(binary.Size(o))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename, extension, ".SCR", -1))
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(dirPath + string(filepath.Separator) + cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	if !noAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	copy(o, OverscanBoot[:])
	copy(o[0x200-0x170:], data[:])
	//o[(0x1ac-0x170)] = 0 // cpc old
	switch isCpcPlus {
	case true:
		o[(0x1ac - 0x170)] = 1
	case false:
		o[(0x1ac - 0x170)] = 0
	}
	switch screenMode {
	case 0:
		o[0x184-0x170] = 0x0e
	case 1:
		o[0x184-0x170] = 0x0f
	case 2:
		o[0x184-0x170] = 0x10
	}
	// affectation de la palette CPC old
	if isCpcPlus {
		offset := 0
		for i := 0; i < len(p); i++ {
			cp := NewCpcPlusColor(p[i])
			fmt.Fprintf(os.Stderr, "i:%d,r:%d,g:%d,b:%d\n", i, cp.R, cp.G, cp.B)
			v := cp.Value()
			o[(0x800-0x170)+offset] = byte(v >> 8)
			offset++
			o[(0x800-0x170)+offset] = byte(v)
			offset++
		}
	} else {
		for i := 0; i < len(p); i++ {
			v, err := HardwareValues(p[i])
			if err == nil {
				o[(0x7f00-0x170)+i] = v[0]
			} else {
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%d\n", p[0], err)
			}
		}
	}
	binary.Write(fw, binary.LittleEndian, o)
	fw.Close()
	return nil
}

func Ink(filePath, dirPath string, p color.Palette, screenMode uint8, noAmsdosHeader bool) error {
	fmt.Fprintf(os.Stdout, "Saving INK file (%s)\n", filePath)
	data := [16]uint16{}
	for i := 0; i < len(p); i++ {
		cp := NewCpcPlusColor(p[i])
		data[i] = cp.Value()
	}
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x8809, Exec: 0x8809,
		Size:        uint16(binary.Size(data)),
		Size2:       uint16(binary.Size(data)),
		LogicalSize: uint16(binary.Size(data))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename, extension, ".INK", -1))
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(dirPath + string(filepath.Separator) + cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	if !noAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	return nil
}

func Scr(filePath, dirPath string, data []byte, noAmsdosHeader bool) error {
	fmt.Fprintf(os.Stdout, "Saving SCR file (%s)\n", filePath)
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0xc000, Exec: 0xC7D0,
		Size:        uint16(binary.Size(data)),
		Size2:       uint16(binary.Size(data)),
		LogicalSize: uint16(binary.Size(data))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename, extension, ".SCR", -1))
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(dirPath + string(filepath.Separator) + cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	if !noAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	return nil
}

type OcpPalette struct {
	ScreenMode          uint8
	ColorAnimation      uint8
	ColorAnimationDelay uint8
	PaletteColors       [16][12]uint8
	BorderColor         [12]uint8
	Excluded            [16]uint8
	Protected           [16]uint8
}

func Pal(filePath, dirPath string, p color.Palette, screenMode uint8, noAmsdosHeader bool) error {
	fmt.Fprintf(os.Stdout, "Saving PAL file (%s)\n", filePath)
	data := OcpPalette{ScreenMode: screenMode, ColorAnimation: 0, ColorAnimationDelay: 0}
	for i := 0; i < 16 ; i++ {
		for j := 0; j< 12; j++ {
			data.PaletteColors[i][j] = 54
		}
	}
	for i := 0; i < len(p); i++ {
		v, err := HardwareValues(p[i])
		if err == nil {
			for j := 0; j < 12; j++ {
				data.PaletteColors[i][j] = v[0]
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%d\n", p[0], err)
		}
	} 
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x8809, Exec: 0x8809,
		Size:        uint16(binary.Size(data)),
		Size2:       uint16(binary.Size(data)),
		LogicalSize: uint16(binary.Size(data))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename, extension, ".PAL", -1))
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(dirPath + string(filepath.Separator) + cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	if !noAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	return nil
}

type OcpWin struct {
	Data   []byte
	Width  uint8
	Height byte
}

func Win(filePath, dirPath string, data []byte, screenMode uint8, width, height int, noAmsdosHeader bool) error {
	fmt.Fprintf(os.Stdout, "Saving WIN file (%s)\n", filePath)
	win := OcpWin{Data: data}
	switch screenMode {
	case 0:
		win.Width = uint8(width / 4)
		win.Height = byte(height / 2)
	case 1:
		win.Width = uint8(width / 2)
		win.Height = uint8(height / 2)
	case 2:
		win.Width = uint8(width)
		win.Height = uint8(height / 2)
	default:
		fmt.Fprintf(os.Stderr, "Win export screen mode not supported. %d\n", screenMode)
	}
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x4000, Exec: 0x4000,
		Size:        uint16(binary.Size(win)),
		Size2:       uint16(binary.Size(win)),
		LogicalSize: uint16(binary.Size(win))}
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.ToUpper(strings.Replace(filename, extension, ".WIN", -1))
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(dirPath + string(filepath.Separator) + cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	if !noAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, win)
	fw.Close()
	return nil

}
