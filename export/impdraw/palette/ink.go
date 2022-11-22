package palette

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/amsdos"
)

type InkPalette struct {
	Colors [16]constants.CpcColor
}

func (i *InkPalette) ToString() string {
	var out string
	for _, v := range i.Colors {
		out += v.ToString() + "\n"
	}
	return out
}

func Ink(filePath string, p color.Palette, screenMode uint8, dontImportDsk bool, cfg *config.MartineConfig) error {
	fmt.Fprintf(os.Stdout, "Saving INK file (%s)\n", filePath)
	data := make([]uint8, 16)
	fmt.Fprintf(os.Stdout, "Palette size %d\n", len(p))
	for i := 0; i < len(p); i++ {
		v, err := constants.HardwareNumber(p[i])
		if err == nil {
			for j := 0; j < 12; j++ {
				data[i] = uint8(v)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
		}
	}

	//fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))
	osFilepath := cfg.AmsdosFullPath(filePath, ".INK")

	if !cfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".INK", data, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, data); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cfg.AddFile(osFilepath)
	}
	return nil
}

func OpenInk(filePath string) (color.Palette, *InkPalette, error) {
	fmt.Fprintf(os.Stdout, "Opening (%s) file\n", filePath)
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return color.Palette{}, &InkPalette{}, err
	}
	defer fr.Close()
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ink Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Ink Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	inkPalette := &InkPalette{}
	buf := make([]uint8, 16)
	if err := binary.Read(fr, binary.LittleEndian, buf); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Palette from file (%s) error %v\n", filePath, err)
		return color.Palette{}, &InkPalette{}, err
	}
	for i, v := range buf {
		c, err := constants.CpcColorFromHardwareNumber(int(v))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Color error :%v\n", err)
		} else {
			inkPalette.Colors[i] = c
		}
	}

	p := color.Palette{}
	for _, v := range inkPalette.Colors {
		c, err := constants.ColorFromHardware(uint8(v.HardwareNumber))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Color error :%v\n", err)
		} else {
			p = append(p, c)
		}
	}
	return p, inkPalette, nil
}

func InkInformation(filePath string) {
	fmt.Fprintf(os.Stdout, "Input kit palette to open : (%s)\n", filePath)
	_, palette, err := OpenInk(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", filePath)
	} else {
		fmt.Fprintf(os.Stdout, "Palette from file %s\n\n%s", filePath, palette.ToString())
	}
}
