package gfx

import (
	"encoding/binary"
	"fmt"
	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/rle"
	"image/color"
	"io"
	"os"
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
	G uint16
	R uint16
	B uint16
}

func NewRawCpcPlusColor(v uint16) *CpcPlusColor {
	c := &CpcPlusColor{}
	c.B = v & 0xf //1111
	c.R = v >> 4 & 0xf
	c.G = v >> 8 & 0xf
	return c
}

func (c *CpcPlusColor) ToString() string {
	return fmt.Sprintf("R:%.2b(%d),G:%.2b(%d),B:%.2b(%d)", c.R, c.R, c.G, c.G, c.B, c.B)
}

func (c *CpcPlusColor) Value() uint16 {
	v := c.B | c.R<<4 | c.G<<8
	fmt.Fprintf(os.Stderr, "value(%d,%d,%d)(%b,%b,%b) #%.4x (%.b): %d\n", c.R, c.G, c.B, c.R, c.G, c.B,
		v, v, c.B+(c.R*16)+c.G*256)
	return v
}
func (c *CpcPlusColor) Bytes() []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, c.Value())
	//	fmt.Fprintf(os.Stderr, "%b\n", buf)
	return buf
}

func NewCpcPlusColor(c color.Color) CpcPlusColor {
	r, g, b, _ := c.RGBA()
	//	fmt.Fprintf(os.Stderr,"original colors r:%d,g:%d,b:%d\n",r,g,b)
	return CpcPlusColor{G: uint16(g / 4096), R: uint16(r / 4096), B: uint16(b / 4096)}
}

type InkPalette struct {
	Colors [16]CpcPlusColor
}

func (i *InkPalette) ToString() string {
	var out string
	for _, v := range i.Colors {
		out += v.ToString() + "\n"
	}
	return out
}

func Overscan(filePath string, data []byte, p color.Palette, screenMode uint8, exportType *ExportType) error {
	o := make([]byte, 0x7e90-0x80)
	osFilepath := exportType.AmsdosFullPath(filePath, ".SCR")
	fmt.Fprintf(os.Stdout, "Saving overscan file (%s)\n", osFilepath)
	header := cpc.CpcHead{Type: 0, User: 0, Address: 0x170, Exec: 0x0,
		Size:        uint16(binary.Size(o)),
		Size2:       uint16(binary.Size(o)),
		LogicalSize: uint16(binary.Size(o))}

	cpcFilename := exportType.OsFilename(".SCR")
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", osFilepath, err)
		return err
	}
	if !exportType.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	copy(o, OverscanBoot[:])
	copy(o[0x200-0x170:], data[:])
	//o[(0x1ac-0x170)] = 0 // cpc old
	switch exportType.CpcPlus {
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
	if exportType.CpcPlus {
		offset := 0
		for i := 0; i < len(p); i++ {
			cp := NewCpcPlusColor(p[i])
			fmt.Fprintf(os.Stderr, "i:%d,r:%d,g:%d,b:%d\n", i, cp.R, cp.G, cp.B)
			v := cp.Bytes()
			copy(o[(0x801-0x170)+offset:], v[:])
			offset += 2
		}
	} else {
		for i := 0; i < len(p); i++ {
			v, err := constants.HardwareValues(p[i])
			if err == nil {
				o[(0x7f00-0x170)+i] = v[0]
			} else {
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
			}
		}
	}
	binary.Write(fw, binary.LittleEndian, o)
	fw.Close()
	return nil
}

func OpenInk(filePath string) (color.Palette, *InkPalette, error) {
	fmt.Fprintf(os.Stdout, "Opening (%s) file\n", filePath)
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return color.Palette{}, &InkPalette{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, 0)
	}

	inkPalette := &InkPalette{}
	buf := make([]uint16, 16)
	if err := binary.Read(fr, binary.LittleEndian, buf); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Palette from file (%s) error %v\n", filePath, err)
		return color.Palette{}, inkPalette, err
	}
	for i, v := range buf {
		c := NewRawCpcPlusColor(v)
		c.B *= 30
		c.R *= 30
		c.G *= 30
		inkPalette.Colors[i] = *c
	}

	p := color.Palette{}
	for _, v := range inkPalette.Colors {
		c := constants.CpcPlusPalette.Convert(color.RGBA{R: uint8(v.R), B: uint8(v.B), G: uint8(v.G), A: 0xFF})
		p = append(p, c)
	}
	return p, inkPalette, nil
}

func Ink(filePath string, p color.Palette, screenMode uint8, exportType *ExportType) error {
	osFilepath := exportType.AmsdosFullPath(filePath, ".INK")
	fmt.Fprintf(os.Stdout, "Saving INK file (%s)\n", osFilepath)
	data := [16]uint16{}

	for i := 0; i < len(p); i++ {
		cp := NewCpcPlusColor(p[i])
		data[i] = cp.Value()
	}
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x8809, Exec: 0x8809,
		Size:        uint16(binary.Size(data)),
		Size2:       uint16(binary.Size(data)),
		LogicalSize: uint16(binary.Size(data))}

	cpcFilename := exportType.OsFilename(".INK")
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", osFilepath, err)
		return err
	}
	if !exportType.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	return nil
}

func Scr(filePath string, data []byte, exportType *ExportType) error {
	osFilepath := exportType.AmsdosFullPath(filePath, ".SCR")
	fmt.Fprintf(os.Stdout, "Saving SCR file (%s)\n", osFilepath)
	if exportType.Compression != -1 {
		switch exportType.Compression {
		case 1:
			fmt.Fprintf(os.Stdout, "Using RLE compression\n")
			data = rle.Encode(data)
		case 2:
			fmt.Fprintf(os.Stdout, "Using RLE 16 bits compression\n")
			data = rle.Encode16(data)
		}
	}
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0xc000, Exec: 0xC7D0,
		Size:        uint16(binary.Size(data)),
		Size2:       uint16(binary.Size(data)),
		LogicalSize: uint16(binary.Size(data))}

	cpcFilename := exportType.OsFilename(".SCR")
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", osFilepath, err)
		return err
	}
	if !exportType.NoAmsdosHeader {
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

func (o *OcpPalette) ToString() string {
	out := fmt.Sprintf("Mode:(%d)\n", o.ScreenMode)
	out += fmt.Sprintf("Color Animation:(%d)\n", o.ColorAnimation)
	out += fmt.Sprintf("Color Animation delay :(%d)\n", o.ColorAnimationDelay)
	for index, v := range o.PaletteColors {
		out += fmt.Sprintf("Color (%d) : value (%d)(%.2x)\n", index, v[0], v[0])
	}
	for index, v := range o.BorderColor {
		out += fmt.Sprintf("Color border (%d) : value (%d)(%.2x)\n", index, v, v)
	}
	return out
}

func OpenPal(filePath string) (color.Palette, *OcpPalette, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return color.Palette{}, &OcpPalette{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, 0)
	}

	ocpPalette := &OcpPalette{}
	if err := binary.Read(fr, binary.LittleEndian, ocpPalette); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Palette from file (%s) error %v\n", filePath, err)
		return color.Palette{}, ocpPalette, err
	}

	p := color.Palette{}
	for _, v := range ocpPalette.PaletteColors {
		c, err := constants.ColorFromHardware(v[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hardware color value %.2x is not recognized error :%v\n", v[0], err)
		} else {
			p = append(p, c)
		}
	}

	return p, ocpPalette, nil
}

func Pal(filePath string, p color.Palette, screenMode uint8, exportType *ExportType) error {
	fmt.Fprintf(os.Stdout, "Saving PAL file (%s)\n", filePath)
	data := OcpPalette{ScreenMode: screenMode, ColorAnimation: 0, ColorAnimationDelay: 0}
	for i := 0; i < 16; i++ {
		for j := 0; j < 12; j++ {
			data.PaletteColors[i][j] = 54
		}
	}
	fmt.Fprintf(os.Stdout, "Palette size %d\n", len(p))
	for i := 0; i < len(p); i++ {
		v, err := constants.HardwareValues(p[i])
		if err == nil {
			for j := 0; j < 12; j++ {
				data.PaletteColors[i][j] = v[0]
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
		}
	}
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x8809, Exec: 0x8809,
		Size:        uint16(binary.Size(data)),
		Size2:       uint16(binary.Size(data)),
		LogicalSize: uint16(binary.Size(data))}

	cpcFilename := exportType.OsFilename(".PAL")
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header lenght %d\n", binary.Size(header))
	osFilepath := exportType.AmsdosFullPath(filePath, ".PAL")
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", cpcFilename, err)
		return err
	}
	if !exportType.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	exportType.AddFile(osFilepath)
	return nil
}

type OcpWinFooter struct {
	Unused2 byte
	Width   uint16
	Height  byte
	Unused  byte
}

func (o *OcpWinFooter) ToString() string {
	return fmt.Sprintf("Width:(%d)\nHeight:(%d)\n", o.Width, o.Height)
}

func OpenWin(filePath string) (*OcpWinFooter, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return &OcpWinFooter{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
	}

	ocpWinFooter := &OcpWinFooter{}
	_, err = fr.Seek(-5, io.SeekEnd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while seek in the file (%s) with error %v\n", filePath, err)
		return &OcpWinFooter{}, err
	}

	if err := binary.Read(fr, binary.LittleEndian, ocpWinFooter); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Win from file (%s) error %v\n", filePath, err)
		return ocpWinFooter, err
	}
	return ocpWinFooter, nil
}

func Win(filePath string, data []byte, screenMode uint8, width, height int, exportType *ExportType) error {
	osFilepath := exportType.AmsdosFullPath(filePath, ".WIN")
	fmt.Fprintf(os.Stdout, "Saving WIN file (%s), screen mode %d, (%d,%d)\n", osFilepath, screenMode, width, height)
	win := OcpWinFooter{Unused: 3, Height: byte(height), Unused2: 0, Width: uint16(width * 8)}
	if exportType.Compression != -1 {
		switch exportType.Compression {
		case 1:
			fmt.Fprintf(os.Stdout, "Using RLE compression\n")
			data = rle.Encode(data)
		case 2:
			fmt.Fprintf(os.Stdout, "Using RLE 16 bits compression\n")
			data = rle.Encode16(data)
		}
	}
	filesize := binary.Size(data) + binary.Size(win)
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x4000, Exec: 0x4000,
		Size:        uint16(filesize),
		Size2:       uint16(filesize),
		LogicalSize: uint16(filesize)}

	cpcFilename := exportType.OsFilename(".WIN")
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "filesize:%d,#%.2x\n", filesize, filesize)
	fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))
	fmt.Fprintf(os.Stderr, "Data length %d\n", binary.Size(data))
	fmt.Fprintf(os.Stderr, "Footer length %d\n", binary.Size(win))
	osFilename := exportType.Fullpath(".WIN")
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", osFilename, err)
		return err
	}
	fmt.Fprintf(os.Stdout, "%s, data size :%d\n", win.ToString(), len(data))
	if !exportType.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, data)
	binary.Write(fw, binary.LittleEndian, win)
	fw.Close()
	exportType.AddFile(osFilepath)
	return nil

}
