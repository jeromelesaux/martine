package export

import (
	"image/color"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/jeromelesaux/m4client/cpc"
)

func Scr(filePath string, data []byte) error {
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0xc000, Exec: 0xC7D0,
		Size: uint16(binary.Size(data)), 
		Size2: uint16(binary.Size(data)), 
		LogicalSize: uint16(binary.Size(data))}
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
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	return nil
}

type OcpPalette struct {
	ScreenMode uint8
	ColorAnimation uint8
	ColorAnimationDelay uint8
	PaletteColor0 [12]uint8
	PaletteColor1 [12]uint8
	PaletteColor2 [12]uint8
	PaletteColor3 [12]uint8
	PaletteColor4 [12]uint8
	PaletteColor5 [12]uint8
	PaletteColor6 [12]uint8
	PaletteColor7 [12]uint8
	PaletteColor8 [12]uint8
	PaletteColor9 [12]uint8
	PaletteColor10 [12]uint8
	PaletteColor11 [12]uint8
	PaletteColor12 [12]uint8
	PaletteColor13 [12]uint8
	PaletteColor14 [12]uint8
	PaletteColor15 [12]uint8
	BorderColor [12]uint8
	Excluded [16]uint8
	Protected [16]uint8

}

func Pal(filePath string, p color.Palette, screenMode uint8) error {
	data := OcpPalette{ScreenMode:screenMode, ColorAnimation:0, ColorAnimationDelay:0}
	header := cpc.CpcHead{Type: 2, User: 0, Address: 0x8809, Exec: 0x8809,
		Size: uint16(binary.Size(data)), 
		Size2: uint16(binary.Size(data)), 
		LogicalSize: uint16(binary.Size(data))}
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
	binary.Write(fw, binary.LittleEndian, data)
	fw.Close()
	return nil
}