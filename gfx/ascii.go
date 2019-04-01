package gfx

import (
	"encoding/binary"
	"fmt"
	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/constants"
	"image/color"
	"os"
	"runtime"
	"strings"
)

var ByteToken = "BYTE"

func Ascii(filePath string, data []byte, p color.Palette, exportType *ExportType) error {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var out string
	var i int

	cpcFilename := string(exportType.AmsdosFilename()) + ".TXT"
	osFilepath := exportType.OsFullPath(filePath, ".TXT")
	fmt.Fprintf(os.Stdout, "Writing ascii file (%s) data length (%d)\n", osFilepath, len(data))
	out += "; Screen " + cpcFilename + eol + ".screen:" + eol
	for i = 0; i < len(data); i += 8 {
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
		out += fmt.Sprintf("%s", eol)
	}
	out += "; Palette " + cpcFilename + eol + ".palette:" + eol + ByteToken + " "

	if exportType.CpcPlus {
		for i := 0; i < len(p); i++ {
			cp := NewCpcPlusColor(p[i])
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
	} else {
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
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%d\n", p[0], err)
			}
		}
		out += eol + "; Basic Palette " + cpcFilename + eol + ".basic_palette:" + eol + ByteToken + " "
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
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%d\n", p[0], err)
			}
		}
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

	if exportType.Json {
		palette := make([]string, len(p))
		for i := 0; i < len(p); i++ {
			v, err := constants.FirmwareNumber(p[i])
			if err == nil {
				palette[i] = fmt.Sprintf("%.2d", v)
			} else {
				fmt.Fprintf(os.Stderr, "Error while getting the hardware values for color %v, error :%d\n", p[0], err)
			}
		}
		screen := make([]string, len(data))
		for i := 0; i < len(data); i++ {
			screen[i] = fmt.Sprintf("0x%.2x", data[i])
		}
		j := NewJson(exportType.Filename(), exportType.Width, exportType.Height, screen, palette)
		return j.Save(exportType.OsFullPath(filePath, ".json"))

	}
	return nil
}
