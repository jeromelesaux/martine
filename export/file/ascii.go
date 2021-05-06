package file

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
	"runtime"
	"strings"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/rle"
	zx0 "github.com/jeromelesaux/zx0/encode"
)

// ByteToken is the token by default
var ByteToken = "db" // "BYTE"

func Ascii(filePath string, data []byte, p color.Palette, dontImportDsk bool, exportType *x.ExportType) error {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var out string
	if exportType.Compression != -1 {
		switch exportType.Compression {
		case 1:
			fmt.Fprintf(os.Stdout, "Using RLE compression\n")
			data = rle.Encode(data)
		case 2:
			fmt.Fprintf(os.Stdout, "Using RLE 16 bits compression\n")
			data = rle.Encode16(data)
		case 5:
			fmt.Fprintf(os.Stdout, "Using Zx0 cruncher")
			data = zx0.Encode(data)
		}

	}
	cpcFilename := string(exportType.AmsdosFilename()) + ".TXT"
	osFilepath := exportType.AmsdosFullPath(filePath, ".TXT")
	fmt.Fprintf(os.Stdout, "Writing ascii file (%s) data length (%d)\n", osFilepath, len(data))
	sizeInfos := fmt.Sprintf("; width %d height %d %s", exportType.Size.Width, exportType.Size.Height, eol)
	out += "; Screen " + cpcFilename + eol + ".screen:" + eol + sizeInfos
	out += FormatAssemblyDatabyte(data, eol)
	out += "; Palette " + cpcFilename + eol + ".palette:" + eol + ByteToken + " "

	if exportType.CpcPlus {
		out += FormatAssemblyCPCPlusPalette(p, eol)
	} else {
		out += FormatAssemblyCPCPalette(p, eol)
		out += eol + "; Basic Palette " + cpcFilename + eol + ".basic_palette:" + eol + ByteToken + " "
		out += FormatAssemblyBasicPalette(p, eol)
		out += eol
	}
	//fmt.Fprintf(os.Stdout,"%s",out)
	header := cpc.CpcHead{Type: 0, User: 0, Address: 0x0, Exec: 0x0,
		Size:        uint16(len(out)),
		Size2:       uint16(len(out)),
		LogicalSize: uint16(len(out))}

	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", osFilepath, err)
		return err
	}
	if !exportType.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, []byte(out))
	fw.Close()
	if !dontImportDsk {
		exportType.AddFile(osFilepath)
	}
	if exportType.Json {
		palette := make([]string, len(p))
		for i := 0; i < len(p); i++ {
			v, err := constants.FirmwareNumber(p[i])
			if err == nil {
				palette[i] = fmt.Sprintf("%.2d", v)
			} else {
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
			}
		}
		hardwarepalette := make([]string, len(p))
		for i := 0; i < len(p); i++ {
			fcolor, _ := constants.FirmwareNumber(p[i])
			hardwarepalette[i] = fmt.Sprintf("0x%.2x", fcolor)
		}
		screen := make([]string, len(data))
		for i := 0; i < len(data); i++ {
			screen[i] = fmt.Sprintf("0x%.2x", data[i])
		}
		j := x.NewJson(exportType.Filename(), exportType.Size.Width, exportType.Size.Height, screen, palette, hardwarepalette)
		fmt.Fprintf(os.Stdout, "Filepath:%s\n", filePath)
		if exportType.TileMode {
			exportType.Tiles.Sprites = append(exportType.Tiles.Sprites, j)
			return nil
		}
		return j.Save(exportType.OsFullPath(filePath, ".json"))
	}
	return nil
}

func AsciiByColumn(filePath string, data []byte, p color.Palette, dontImportDsk bool, mode uint8, exportType *x.ExportType) error {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var out string
	var i int
	var jsonData []string

	cpcFilename := string(exportType.AmsdosFilename()) + "C.TXT"
	osFilepath := exportType.AmsdosFullPath(filePath, "C.TXT")
	fmt.Fprintf(os.Stdout, "Writing ascii file (%s) values by columns data length (%d)\n", osFilepath, len(data))
	sizeInfos := fmt.Sprintf("; width %d height %d %s", exportType.Size.Width, exportType.Size.Height, eol)
	out += "; Screen by column " + cpcFilename + eol + ".screen:" + eol + sizeInfos
	var adjustMode int
	switch mode {
	case 0:
		adjustMode = 2
	case 1:
		adjustMode = 4
	case 2:
		adjustMode = 8
	}
	pas := exportType.Size.Width / adjustMode
	h := 0
	nbValues := 1
	octetsRead := 0
	end := 17
	if (exportType.Size.Width + 1) < end {
		end = (exportType.Size.Width + 1)
	}
	for {

		if nbValues == 1 {
			out += fmt.Sprintf("%s ", ByteToken)
		}
		out += fmt.Sprintf("#%0.2x", data[i])
		jsonData = append(jsonData, fmt.Sprintf("0x%.2x", data[i]))
		nbValues++

		i += pas
		octetsRead++
		if nbValues < end && octetsRead != len(data) {
			out += " ,"
		}
		if octetsRead == len(data) {
			break
		}

		if i >= len(data) {
			h++
			i = h
		}

		if nbValues == end {
			out += eol
			nbValues = 1
		}
	}
	out += eol
	out += "; Palette " + cpcFilename + eol + ".palette:" + eol + ByteToken + " "

	if exportType.CpcPlus {
		out += FormatAssemblyCPCPlusPalette(p, eol)
	} else {
		out += FormatAssemblyCPCPalette(p, eol)
		out += eol + "; Basic Palette " + cpcFilename + eol + ".basic_palette:" + eol + ByteToken + " "
		out += FormatAssemblyBasicPalette(p, eol)
		out += eol
	}
	//fmt.Fprintf(os.Stdout,"%s",out)
	header := cpc.CpcHead{Type: 0, User: 0, Address: 0x0, Exec: 0x0,
		Size:        uint16(len(out)),
		Size2:       uint16(len(out)),
		LogicalSize: uint16(len(out))}

	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))
	fw, err := os.Create(osFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", osFilepath, err)
		return err
	}
	if !exportType.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, []byte(out))
	fw.Close()
	if !dontImportDsk {
		exportType.AddFile(osFilepath)
	}

	if exportType.Json {
		palette := make([]string, len(p))
		for i := 0; i < len(p); i++ {
			v, err := constants.FirmwareNumber(p[i])
			if err == nil {
				palette[i] = fmt.Sprintf("%.2d", v)
			} else {
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
			}
		}
		hardwarepalette := make([]string, len(p))
		for i := 0; i < len(p); i++ {
			fcolor, _ := constants.FirmwareNumber(p[i])
			hardwarepalette[i] = fmt.Sprintf("0x%.2x", fcolor)
		}

		j := x.NewJson(exportType.Filename(), exportType.Size.Width, exportType.Size.Height, jsonData, palette, hardwarepalette)
		fmt.Fprintf(os.Stdout, "Filepath:%s\n", filePath)
		if exportType.TileMode {
			exportType.Tiles.Sprites = append(exportType.Tiles.Sprites, j)
			return nil
		}
		return j.Save(exportType.OsFullPath(filePath, "_column.json"))
	}
	return nil
}

func FormatAssemblyString(data []string, eol string) string {
	var out string
	for i := 0; i < len(data); i += 8 {
		out += fmt.Sprintf("%s ", ByteToken)
		if i < len(data) {
			out += data[i]
		}
		if i+1 < len(data) {
			out += fmt.Sprintf(", %s", data[i+1])
		}
		if i+2 < len(data) {
			out += fmt.Sprintf(", %s", data[i+2])
		}
		if i+3 < len(data) {
			out += fmt.Sprintf(", %s", data[i+3])
		}
		if i+4 < len(data) {
			out += fmt.Sprintf(", %s", data[i+4])
		}
		if i+5 < len(data) {
			out += fmt.Sprintf(", %s", data[i+5])
		}
		if i+6 < len(data) {
			out += fmt.Sprintf(", %s", data[i+6])
		}
		if i+7 < len(data) {
			out += fmt.Sprintf(", %s", data[i+7])
		}
		out += eol
	}
	return out
}

func FormatAssemblyDatabyte(data []byte, eol string) string {
	var out string
	for i := 0; i < len(data); i += 8 {
		out += fmt.Sprintf("%s ", ByteToken)
		if i < len(data) {
			out += fmt.Sprintf("#%0.2x", data[i])
		}
		if i+1 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+1])
		}
		if i+2 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+2])
		}
		if i+3 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+3])
		}
		if i+4 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+4])
		}
		if i+5 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+5])
		}
		if i+6 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+6])
		}
		if i+7 < len(data) {
			out += fmt.Sprintf(", #%0.2x", data[i+7])
		}
		out += eol
	}
	return out
}

func FormatAssemblyCPCPalette(p color.Palette, eol string) string {
	var out string
	for i := 0; i < len(p); i++ {
		v, err := constants.HardwareValues(p[i])
		if err == nil {
			out += fmt.Sprintf("#%0.2x", v[0])
			if (i+1)%8 == 0 && i+1 < len(p) {
				out += eol + ByteToken + " "
			} else {
				if i+1 < len(p) {
					out += ", "
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
		}
	}
	return out
}

func FormatAssemblyBasicPalette(p color.Palette, eol string) string {
	var out string
	for i := 0; i < len(p); i++ {
		v, err := constants.FirmwareNumber(p[i])
		if err == nil {
			out += fmt.Sprintf("%0.2d", v)
			if (i+1)%8 == 0 && i+1 < len(p) {
				out += eol + ByteToken + " "
			} else {
				if i+1 < len(p) {
					out += ", "
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%v\n", p[0], err)
		}
	}
	return out
}

func FormatAssemblyCPCPlusPalette(p color.Palette, eol string) string {
	var out string
	for i := 0; i < len(p); i++ {
		cp := constants.NewCpcPlusColor(p[i])
		v := cp.Value()
		out += fmt.Sprintf("#%.2x, #%.2x", byte(v), byte(v>>8))
		if (i+1)%8 == 0 && i+1 < len(p) {
			out += eol + ByteToken + " "
		} else {
			if i+1 < len(p) {
				out += ", "
			}
		}
	}
	return out
}
