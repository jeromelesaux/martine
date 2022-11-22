package overscan

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"os"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/common/errors"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/compression"
)

func OverscanPalette(filePath string) (color.Palette, uint8, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return color.Palette{}, 0xff, err
	}
	defer fr.Close()
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Overscan Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Overscan Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	palette := color.Palette{}
	b, err := io.ReadAll(fr)
	if err != nil {
		return palette, 0xff, err
	}
	fmt.Fprintf(os.Stdout, "Read (%X)\n", len(b))
	var mode uint8
	isPlus := false
	if b[(0x1ac-0x170)] == 1 {
		isPlus = true
	}
	pens := 0
	switch b[0x184-0x170] {
	case 0x0e:
		pens = 16
		mode = 0
	case 0x0f:
		pens = 4
		mode = 1
	case 0x10:
		pens = 2
		mode = 2
	}
	if isPlus {
		offset := 0
		for i := 0; i < pens; i++ {
			pc := binary.LittleEndian.Uint16(b[(0x801-0x170)+offset:])
			fmt.Fprintf(os.Stdout, "Read color %d\n", pc)
			if err == nil {
				c := constants.NewRawCpcPlusColor(pc)
				fmt.Fprintf(os.Stdout, "PEN(%d) R(%d) G(%d) B(%d)\n", i, c.R, c.G, c.B)
				col := color.RGBA{A: 0xff, B: uint8(c.B) << 4, G: uint8(c.G) << 4, R: uint8(c.R) << 4}
				palette = append(palette, col)
			} else {
				palette = append(palette, color.Black)
				fmt.Fprintf(os.Stderr, "Error while retreiving color from hardware value %X error %v\n", pc, err)
			}
			offset += 2
		}
	} else {
		for i := 0; i < pens; i++ {
			v := b[(0x7f00-0x170)+i]
			fmt.Fprintf(os.Stdout, "PEN(%d) Hardware value (%X)\n", i, v)
			c, err := constants.ColorFromHardware(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while retreiving color from hardware value %X error %v\n", v, err)
				palette = append(palette, color.Black)
			} else {

				palette = append(palette, c)
			}
		}
	}
	fmt.Fprintf(os.Stdout, "Overscan file (%s) palette length (%d) mode (%d)\n", filePath, len(palette), mode)
	return palette, mode, nil
}

func Overscan(filePath string, data []byte, p color.Palette, screenMode uint8, cfg *config.MartineConfig) error {
	o := make([]byte, 0x7e90-0x80)

	// remove first line to keep #38 address free
	var width int
	switch screenMode {
	case 0:
		width = cfg.Size.Width / 2
	case 1:
		width = cfg.Size.Width / 4
	case 2:
		width = cfg.Size.Width / 8
	}
	for i := 0; i < width; i++ {
		data[i] = 0
	}
	// end of the hack

	copy(o, OverscanBoot[:])
	copy(o[0x200-0x170:], data[:])
	//o[(0x1ac-0x170)] = 0 // cpc old
	switch cfg.CpcPlus {
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
	if cfg.CpcPlus {
		offset := 0
		for i := 0; i < len(p); i++ {
			cp := constants.NewCpcPlusColor(p[i])
			//fmt.Fprintf(os.Stderr, "i:%d,r:%d,g:%d,b:%d\n", i, cp.R, cp.G, cp.B)
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

	o, _ = compression.Compress(o, cfg.Compression)

	osFilepath := cfg.AmsdosFullPath(filePath, ".SCR")
	fmt.Fprintf(os.Stdout, "Saving overscan file (%s)\n", osFilepath)
	if !cfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".SCR", o, 0, 0, 0x170, 0); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(filePath, o); err != nil {
			return err
		}
	}

	cfg.AddFile(osFilepath)
	return nil
}

func RawOverscan(filePath string) ([]byte, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return []byte{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Overscan Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Overscan Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	bf, err := io.ReadAll(fr)
	if err != nil {
		return nil, err
	}
	data := make([]byte, 0x8000)
	copy(data, bf[0x200-0x170:])
	fmt.Fprintf(os.Stdout, "Raw overscan length #%X\n", len(data))
	if len(data) <= 0x4000 {
		return nil, errors.ErrorBadFileFormat
	}
	return data, nil
}

func EgxOverscan(filePath string, data []byte, p color.Palette, mode1, mode2 uint8, cfg *config.MartineConfig) error {
	o := make([]byte, 0x8000-0x80)
	osFilepath := cfg.AmsdosFullPath(filePath, ".SCR")
	fmt.Fprintf(os.Stdout, "Saving overscan file (%s)\n", osFilepath)

	//fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))

	var overscanTemplate []byte
	if cfg.CpcPlus {
		overscanTemplate = egxPlusOverscanTemplate
	} else {
		overscanTemplate = egxOverscanTemplate
	}
	copy(o[:], overscanTemplate[:])
	copy(o[0x200-0x170:], data[:]) //  - 0x170  to have the file offset
	//o[(0x1ac-0x170)] = 0 // cpc old
	switch cfg.CpcPlus {
	case true:
		o[(0x1ac - 0x170)] = 1
	case false:
		o[(0x1ac - 0x170)] = 0
	}

	screenMode := mode1
	if mode2 < mode1 {
		screenMode = mode2
	}
	switch screenMode {
	case 0:
		o[0x184-0x170] = 0x0e
	case 1:
		o[0x184-0x170] = 0x0f
	case 2:
		o[0x184-0x170] = 0x10
	}

	extraFlag := 0

	if mode1 == 0 && mode2 == 1 {
		extraFlag = 2
	}
	if mode2 == 0 && mode1 == 1 {
		extraFlag = 1
	}
	if mode1 == 1 && mode2 == 2 {
		extraFlag = 4
	}
	if mode1 == 2 && mode2 == 1 {
		extraFlag = 3
	}
	o[0x8f] = byte(extraFlag)

	// affectation de la palette CPC old
	if cfg.CpcPlus {
		offset := 0
		for i := 0; i < len(p); i++ {
			cp := constants.NewCpcPlusColor(p[i])
			//fmt.Fprintf(os.Stderr, "i:%d,r:%d,g:%d,b:%d\n", i, cp.R, cp.G, cp.B)
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
	if cfg.CpcPlus {
		copy(o[0x6b2:0x6c8], egxPlusOverscanTemplate[0x6b2:0x6c8])
		copy(o[0x7da0:], egxPlusOverscanTemplate[0x7da0:]) // copy egx routine
	} else {
		copy(o[0x7da0:], egxOverscanTemplate[0x7da0:]) // copy egx routine
	}
	if !cfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".SCR", o, 0, 0, 0x170, 0x170); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, o); err != nil {
			return err
		}
	}

	cfg.AddFile(osFilepath)
	return nil
}
