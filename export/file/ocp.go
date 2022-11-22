package file

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/compression"
)

var codeScrStandard = []byte{ // Routine à mettre en #C7D0
	0x3A, 0xD0, 0xD7, //      LD      A,  (#D7D0)
	0xCD, 0x1C, 0xBD, //      CALL    #BD1C
	0x21, 0xD1, 0xD7, //      LD      HL, #D7D1
	0x46,             //      LD      B,  (HL)
	0x48,             //      LD      C,  B
	0xCD, 0x38, 0xBC, //      CALL    #BC38
	0xAF,             //      XOR     A
	0x21, 0xD1, 0xD7, //      LD      HL, #D7D1
	0x46,             // BCL: LD      B,  (HL)
	0x48,             //      LD      C,  B
	0xF5,             //      PUSH    AF
	0xE5,             //      PUSH    HL
	0xCD, 0x32, 0xBC, //      CALL    #BC32
	0xE1,       //      POP     HL
	0xF1,       //      POP     AF
	0x23,       //      INC     HL
	0x3C,       //      INC     A
	0xFE, 0x10, //      CP      #10
	0x20, 0xF1, //      JR      NZ,BCL
	0xC3, 0x18, 0xBB, //      JP      #BB18
}

var codeScrPlusP0 = []byte{
	0xF3,             //				DI
	0x01, 0x11, 0xBC, //				LD		BC,#BC11
	0x21, 0xD0, 0xDF, //				LD		HL,#DFD0
	0x7E,       //	BCL1:		LD		A,(HL)
	0xED, 0x79, //				OUT		(C),A
	0x23,       //				INC		HL
	0x0D,       //				DEC		C
	0x20, 0xF9, //				JR		NZ,BCL1
	0x01, 0xA0, 0x7F, //				LD		BC,#7FA0
	0x3A, 0xD0, 0xD7, //				LD		A,(#D7D0)
	0xED, 0x79, //				OUT		(C),A
	0xED, 0x49, //				OUT		(C),C
	0x01, 0xB8, 0x7F, //				LD		BC,#7FB8
	0xED, 0x49, //				OUT		(C),C
	0x21, 0xD1, 0xD7, //				LD		HL,#D7D1
	0x11, 0x00, 0x64, //				LD		DE,#6400
	0x01, 0x22, 0x00, //				LD		BC,#0022
	0xED, 0xB0, //				LDIR
	0xCD, 0xD0, 0xCF, //	BCL2:		CALL	WaitKey
	0x38, 0xFB, //				JR		C,BCL2
	0xFB, //				EI
	0xC9, //				RET
}

var codeScrPlusP1 = []byte{
	0x01, 0x0E, 0xF4, //	WaitKey:	LD		BC,#F40E
	0xED, 0x49, //				OUT		(C),C
	0x01, 0xC0, 0xF6, //				LD		BC,#F6C0
	0xED, 0x49, //				OUT		(C),C
	0xAF,       //				XOR		A
	0xED, 0x79, //				OUT		(C),A
	0x01, 0x92, 0xF7, //				LD		BC,#F792
	0xED, 0x49, //				OUT		(C),C
	0x01, 0x45, 0xF6, //				LD		BC,#F645
	0xED, 0x49, //				OUT		(C),C
	0x06, 0xF4, //				LD		B,#F4
	0xED, 0x78, //				IN		A,(C)
	0x01, 0x82, 0xF7, //				LD		BC,#F782
	0xED, 0x49, //				OUT		(C),C
	0x01, 0x00, 0xF6, //				LD		BC,#F600
	0xED, 0x49, //				OUT		(C),C
	0x17, //				RLA
	0xC9, //				RET
}

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

func Ink(filePath string, p color.Palette, screenMode uint8, dontImportDsk bool, cont *x.MartineContext) error {
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
	osFilepath := cont.AmsdosFullPath(filePath, ".INK")

	if !cont.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".INK", data, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, data); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cont.AddFile(osFilepath)
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

func Kit(filePath string, p color.Palette, screenMode uint8, dontImportDsk bool, cont *x.MartineContext) error {
	osFilepath := cont.AmsdosFullPath(filePath, ".KIT")
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
	if !cont.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".KIT", res, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, res); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cont.AddFile(osFilepath)
	}
	return nil
}

func DepackOCP(buf []byte) ([]byte, error) {
	var bmpCpc []byte
	var PosIn, PosOut int
	var LgOut, CntBlock int
	bmpCpc = make([]byte, 0x10000)
	copy(bmpCpc, buf)
	for PosOut < 0x4000 {
		if buf[PosIn] == 'M' && buf[PosIn+1] == 'J' && buf[PosIn+2] == 'H' {
			PosIn += 3
			LgOut = int(buf[PosIn])
			PosIn++
			LgOut += (int(buf[PosIn]) << 8)
			PosIn++
			CntBlock = 0
			for CntBlock < LgOut {
				if buf[PosIn] == 'M' && buf[PosIn+1] == 'J' && buf[PosIn+2] == 'H' {
					break
				}

				a := buf[PosIn]
				PosIn++
				if a == 1 { // MARKER_OCP
					var c int
					c = int(buf[PosIn])
					PosIn++
					a = buf[PosIn]
					PosIn++
					if c == 0 {
						c = 0x100
					}

					for i := 0; i < c && CntBlock < LgOut; i++ {
						bmpCpc[PosOut] = a
						PosOut++
						CntBlock++
					}
				} else {
					bmpCpc[PosOut] = a
					PosOut++
					CntBlock++
				}
			}
		} else {
			PosOut = 0x4000
		}
	}

	return bmpCpc, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RawScr(filePath string) ([]byte, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return []byte{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the RawScr Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the RawScr Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	bf, err := ioutil.ReadAll(fr)
	if err != nil {
		return nil, err
	}

	var sz int = min(0x4000, len(bf))

	var rawSrc []byte = make([]byte, sz)
	for i := 0; i < sz; i++ {
		rawSrc[i] = bf[i]
	}

	if rawSrc[0] == 'M' && rawSrc[1] == 'J' && rawSrc[2] == 'H' { // Compression OCP
		return DepackOCP(rawSrc)
	}
	return rawSrc, nil
}

func Scr(filePath string, data []byte, p color.Palette, screenMode uint8, cont *x.MartineContext) error {
	osFilepath := cont.AmsdosFullPath(filePath, ".SCR")
	fmt.Fprintf(os.Stdout, "Saving SCR file (%s)\n", osFilepath)
	var exec uint16
	if cont.CpcPlus {
		exec = 0x821
		switch screenMode {
		case 0:
			data[0x17d0] = 0
		case 1:
			data[0x17d0] = 1
		case 2:
			data[0x17d0] = 2
		}
		offset := 1
		for i := 0; i < len(p); i++ {
			cp := constants.NewCpcPlusColor(p[i])
			//fmt.Fprintf(os.Stderr, "i:%d,r:%d,g:%d,b:%d\n", i, cp.R, cp.G, cp.B)
			v := cp.Bytes()
			data[0x17d0+offset] = v[0]
			offset++
			data[0x17d0+offset] = v[1]
			offset++
		}
		copy(data[0x07d0:], codeScrPlusP0[:])
		copy(data[0x0fd0:], codeScrPlusP1[:])

	} else {
		exec = 0x811
		switch screenMode {
		case 0:
			data[0x17D0] = 0
		case 1:
			data[0x17D0] = 1
		case 2:
			data[0x17D0] = 2
		}
		for i := 0; i < len(p); i++ {
			v, err := constants.HardwareValues(p[i])
			if err == nil {
				data[(0x17D0)+i] = v[0]
			} else {
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
			}
		}
		copy(data[0x07d0:], codeScrStandard[:])
	}
	data, _ = compression.Compress(data, cont.Compression)

	if !cont.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".SCR", data, 2, 0, 0xc000, exec); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, data); err != nil {
			return err
		}
	}

	cont.AddFile(osFilepath)
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
		out += fmt.Sprintf("Color (%d) [%s] : value (%d)(#%.2x)\n", index, constants.CpcColorStringFromHardwareNumber(v[0]), v[0], v[0])
	}
	for index, v := range o.BorderColor {
		out += fmt.Sprintf("Color border (%d) [%s] : value (%d)(#%.2x)\n", index, constants.CpcColorStringFromHardwareNumber(v), v, v)
	}
	out += "Colors from Gatearray:\n"
	for _, v := range o.PaletteColors {
		out += fmt.Sprintf("#%.2x, ", v[0])
	}
	out += "\nColors from firmware:\n"

	for _, v := range o.PaletteColors {
		hcolor, err := constants.ColorFromHardware(v[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error color :%v\n", err)
			out += "00, "
		} else {
			fcolor, _ := constants.FirmwareNumber(hcolor)
			out += fmt.Sprintf("%.2d, ", fcolor)
		}
	}
	out += "\n"
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
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
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
			p = append(p, color.White)

		} else {
			p = append(p, c)
		}
	}

	return p, ocpPalette, nil
}

func SavePal(filePath string, p color.Palette, screenMode uint8, noAmsdosHeader bool) error {
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

	res, err := common.StructToBytes(data)
	if err != nil {
		return err
	}

	if !noAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(filePath, ".PAL", res, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(filePath, res); err != nil {
			return err
		}
	}

	return nil
}

func Pal(filePath string, p color.Palette, screenMode uint8, dontImportDsk bool, cont *x.MartineContext) error {
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
	osFilepath := cont.AmsdosFullPath(filePath, ".PAL")

	res, err := common.StructToBytes(data)
	if err != nil {
		return err
	}
	if !cont.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".PAL", res, 2, 0, 0x8809, 0x8809); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, res); err != nil {
			return err
		}
	}
	if !dontImportDsk {
		cont.AddFile(osFilepath)
	}
	return nil
}

type OcpWinFooter struct {
	Unused2 byte
	Width   uint16
	Height  byte
	Unused  byte
}

func (o *OcpWinFooter) ToString() string {
	return fmt.Sprintf("Width:(%d)\nHeight:(%d)\n", o.Width/8, o.Height)
}

func RawWin(filePath string) ([]byte, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return []byte{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Win Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Win Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	bf, err := ioutil.ReadAll(fr)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, len(bf)-5)
	copy(raw[:], bf[0:len(bf)-5])

	if raw[0] == 'M' && raw[1] == 'J' && raw[2] == 'H' { // Compression OCP
		return DepackOCP(raw)
	}

	return raw, nil
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
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	ocpWinFooter := &OcpWinFooter{}
	//_, err = fr.Seek(-5, io.SeekEnd)

	fmt.Fprintf(os.Stdout, "LogicalSize=%d\n", header.LogicalSize)
	_, err = fr.Seek(0x80+int64(header.LogicalSize)-5, io.SeekStart)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while seek in the file (%s) with error %v\n", filePath, err)
		return &OcpWinFooter{}, err
	}

	if err := binary.Read(fr, binary.LittleEndian, ocpWinFooter); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Win from file (%s) error %v\n", filePath, err)
		return ocpWinFooter, err
	}
	ocpWinFooter.Width = uint16(uint(ocpWinFooter.Width / 8))
	return ocpWinFooter, nil
}

func Win(filePath string, data []byte, screenMode uint8, width, height int, dontImportDsk bool, cont *x.MartineContext) error {
	osFilepath := cont.AmsdosFullPath(filePath, ".WIN")
	fmt.Fprintf(os.Stdout, "Saving WIN file (%s), screen mode %d, (%d,%d)\n", osFilepath, screenMode, width, height)
	win := OcpWinFooter{Unused: 3, Height: byte(height), Unused2: 0, Width: uint16(width * 8)}

	data, _ = compression.Compress(data, cont.Compression)

	//fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))
	fmt.Fprintf(os.Stderr, "Data length %d\n", binary.Size(data))
	fmt.Fprintf(os.Stderr, "Footer length %d\n", binary.Size(win))
	osFilename := cont.Fullpath(".WIN")

	body, err := common.StructToBytes(data)
	if err != nil {
		return err
	}
	footer, err := common.StructToBytes(win)
	if err != nil {
		return err
	}
	content := body
	content = append(content, footer...)

	fmt.Fprintf(os.Stdout, "%s, data size :%d\n", win.ToString(), len(data))
	if !cont.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilename, ".WIN", content, 2, 0, 0x4000, 0x4000); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilename, content); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cont.AddFile(osFilepath)
	}
	return nil
}
