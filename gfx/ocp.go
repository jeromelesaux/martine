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
