package palette

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/amsdos"
)

type KitPalette struct {
	Colors [16]constants.CpcPlusColor
}

func (i *KitPalette) ToString() string {
	var out string
	for _, v := range i.Colors {
		out += v.ToString() + "\n"
	}

	out += "db "
	for index, v := range i.Colors {
		out += fmt.Sprintf("#%0.2X, #%0.2X", v.Bytes()[0], v.Bytes()[1])
		if index == (len(i.Colors) - 1) {
			out += ""
		} else {
			if (index+1)%8 == 0 {
				out += "\ndb "
			} else {
				out += ", "
			}
		}
	}
	out += "\n"
	return out
}

func OpenKit(filePath string) (color.Palette, *KitPalette, error) {
	fmt.Fprintf(os.Stdout, "Opening (%s) file\n", filePath)
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return color.Palette{}, &KitPalette{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Kit Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Kit Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	KitPalette := &KitPalette{}
	buf := make([]uint16, 0)
	for {
		var b uint16
		if err := binary.Read(fr, binary.LittleEndian, &b); err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error while reading Ocp Palette from file (%s) error %v\n", filePath, err)
				return color.Palette{}, KitPalette, err
			}
			break
		}
		buf = append(buf, b)
	}
	p := color.Palette{}
	for i, v := range buf {
		c := constants.NewRawCpcPlusColor(v)
		KitPalette.Colors[i] = *c
		pp := constants.NewColorCpcPlusColor(*c)
		p = append(p, pp)
	}
	return p, KitPalette, nil
}

func SaveKit(filePath string, p color.Palette, noAmsdosHeader bool) error {

	fmt.Fprintf(os.Stdout, "Saving Kit file (%s)\n", filePath)
	data := [16]uint16{}
	paletteSize := len(p)
	if len(p) > 16 {
		paletteSize = 16
	}
	for i := 0; i < paletteSize; i++ {
		cp := constants.NewCpcPlusColor(p[i])
		data[i] = cp.Value()
	}

	//fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))

	v, err := common.StructToBytes(data)
	if err != nil {
		return err
	}
	if !noAmsdosHeader {
		if err = amsdos.SaveAmsdosFile(filePath, ".KIT", v, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(filePath, v); err != nil {
			return err
		}
	}

	return nil
}

func Kit(filePath string, p color.Palette, screenMode uint8, dontImportDsk bool, cfg *config.MartineConfig) error {
	osFilepath := cfg.AmsdosFullPath(filePath, ".KIT")
	fmt.Fprintf(os.Stdout, "Saving Kit file (%s)\n", osFilepath)
	data := [16]uint16{}
	paletteSize := len(p)
	if len(p) > 16 {
		paletteSize = 16
	}
	for i := 0; i < paletteSize; i++ {
		cp := constants.NewCpcPlusColor(p[i])
		data[i] = cp.Value()
	}

	res, err := common.StructToBytes(data)
	if err != nil {
		return err
	}
	if !cfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".KIT", res, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, res); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cfg.AddFile(osFilepath)
	}
	return nil
}

func KitInformation(filePath string) {
	fmt.Fprintf(os.Stdout, "Input kit palette to open : (%s)\n", filePath)
	_, palette, err := OpenKit(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", filePath)
	} else {
		fmt.Fprintf(os.Stdout, "Palette from file %s\n\n%s", filePath, palette.ToString())
	}
}
