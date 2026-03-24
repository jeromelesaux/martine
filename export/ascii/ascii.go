package ascii

import (
	"fmt"
	"image/color"
	"runtime"
	"strings"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/log"
)

// ByteToken is the token by default
var ByteToken = "db" // "BYTE"

// nolint: funlen
func Ascii(filePath string, data []byte, p color.Palette, dontImportDsk bool, cgf *config.MartineConfig) error {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var out string

	data, _ = compression.Compress(data, cgf.ScrCfg.Compression)

	cpcFilename := string(cgf.AmsdosFilename()) + ".TXT"
	osFilepath := cgf.AmsdosFullPath(filePath, ".TXT")
	log.GetLogger().Info("Writing ascii file (%s) data length (%d)\n", osFilepath, len(data))
	sizeInfos := fmt.Sprintf("; width %d height %d %s", cgf.ScrCfg.Size.Width, cgf.ScrCfg.Size.Height, eol)
	out += "; Screen " + cpcFilename + eol + ".screen:" + eol + sizeInfos
	out += FormatAssemblyDatabyte(data, eol)
	out += "; Palette " + cpcFilename + eol + ".palette:" + eol + ByteToken + " "

	if cgf.ScrCfg.IsPlus {
		out += FormatAssemblyCPCPlusPalette(p, eol)
	} else {
		out += FormatAssemblyCPCPalette(p, eol) + eol + "; Basic Palette " + cpcFilename + eol + ".basic_palette:" + eol + ByteToken + " "
		out += FormatAssemblyBasicPalette(p, eol) + eol
	}
	if !cgf.ScrCfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".TXT", []byte(out), 0, 0, 0, 0); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, []byte(out)); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cgf.AddFile(osFilepath)
	}
	if cgf.ScrCfg.IsExport(config.JsonExport) {
		palette := make([]string, len(p))
		for i := range p {
			v, err := constants.FirmwareNumber(p[i])
			if err == nil {
				palette[i] = fmt.Sprintf("%.2d", v)
			}
		}
		hardwarepalette := make([]string, len(p))
		for i := range p {
			fcolor, _ := constants.FirmwareNumber(p[i])
			hardwarepalette[i] = fmt.Sprintf("0x%.2x", fcolor)
		}
		screen := make([]string, len(data))
		for i := 0; i < len(data); i++ {
			screen[i] = fmt.Sprintf("0x%.2x", data[i])
		}
		j := x.NewJson(cgf.Filename(), cgf.ScrCfg.Size.Width, cgf.ScrCfg.Size.Height, screen, palette, hardwarepalette)
		log.GetLogger().Info("Filepath:%s\n", filePath)
		if cgf.TileMode {
			cgf.Tiles.Sprites = append(cgf.Tiles.Sprites, j)
			return nil
		}
		return j.Save(cgf.OsFullPath(filePath, ".json"))
	}
	return nil
}

// nolint: funlen, gocognit
func AsciiByColumn(filePath string, data []byte, p color.Palette, dontImportDsk bool, mode uint8, cfg *config.MartineConfig) error {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var out string
	var i int
	var jsonData []string

	cpcFilename := string(cfg.AmsdosFilename()) + "C.TXT"
	osFilepath := cfg.AmsdosFullPath(filePath, "C.TXT")
	log.GetLogger().Info("Writing ascii file (%s) values by columns data length (%d)\n", osFilepath, len(data))
	sizeInfos := fmt.Sprintf("; width %d height %d %s", cfg.ScrCfg.Size.Width, cfg.ScrCfg.Size.Height, eol)
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
	pas := cfg.ScrCfg.Size.Width / adjustMode
	h := 0
	nbValues := 1
	octetsRead := 0
	end := min((cfg.ScrCfg.Size.Width + 1), 17)
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

	if cfg.ScrCfg.IsPlus {
		out += FormatAssemblyCPCPlusPalette(p, eol)
	} else {
		out += FormatAssemblyCPCPalette(p, eol)
		out += eol + "; Basic Palette " + cpcFilename + eol + ".basic_palette:" + eol + ByteToken + " "
		out += FormatAssemblyBasicPalette(p, eol)
		out += eol
	}

	if !cfg.ScrCfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".TXT", []byte(out), 0, 0, 0, 0); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, []byte(out)); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cfg.AddFile(osFilepath)
	}

	if cfg.ScrCfg.IsExport(config.JsonExport) {
		palette := make([]string, len(p))
		for i := 0; i < len(p); i++ {
			v, err := constants.FirmwareNumber(p[i])
			if err == nil {
				palette[i] = fmt.Sprintf("%.2d", v)
			}
		}
		hardwarepalette := make([]string, len(p))
		for i := range p {
			fcolor, _ := constants.FirmwareNumber(p[i])
			hardwarepalette[i] = fmt.Sprintf("0x%.2x", fcolor)
		}

		j := x.NewJson(cfg.Filename(), cfg.ScrCfg.Size.Width, cfg.ScrCfg.Size.Height, jsonData, palette, hardwarepalette)
		log.GetLogger().Info("Filepath:%s\n", filePath)
		if cfg.TileMode {
			cfg.Tiles.Sprites = append(cfg.Tiles.Sprites, j)
			return nil
		}
		return j.Save(cfg.OsFullPath(filePath, "_column.json"))
	}
	return nil
}

func FormatAssemblyString(data []string, eol string) string {
	var out strings.Builder
	for i := 0; i < len(data); i += 8 {
		fmt.Fprintf(&out, "%s ", ByteToken)
		if i < len(data) {
			out.WriteString(data[i])
		}
		for j := 1; j < 8; j++ {
			if i+j < len(data) {
				fmt.Fprintf(&out, ", %s", data[i+j])
			}
		}
		out.WriteString(eol)
	}
	return out.String()
}

func SpritesHardText(data [][]byte, compressionType compression.CompressionMethod) string {
	var out strings.Builder
	for i, v := range data {
		fmt.Fprintf(&out, "Sprite_%02d\n", i)
		compressed, err := compression.Compress(v, compressionType)
		if err == nil {
			out.WriteString(FormatAssemblyDatabyte(compressed, "\n"))
		} else {
			out.WriteString(err.Error())
		}

	}

	return out.String()
}

func FormatAssemblyDatabyte(data []byte, eol string) string {
	var out strings.Builder
	for i := 0; i < len(data); i += 8 {
		fmt.Fprintf(&out, "%s ", ByteToken)
		if i < len(data) {
			fmt.Fprintf(&out, "#%0.2x", data[i])
		}
		for j := 1; j < 8; j++ {
			if i+j < len(data) {
				fmt.Fprintf(&out, ", #%0.2x", data[i+j])
			}
		}

		out.WriteString(eol)
	}
	return out.String()
}

func FormatAssemblyCPCPalette(p color.Palette, eol string) string {
	var out strings.Builder
	for i := range p {
		v, err := constants.HardwareValues(p[i])
		if err == nil {
			fmt.Fprintf(&out, "#%0.2x", v[0])
			if (i+1)%8 == 0 && i+1 < len(p) {
				out.WriteString(eol + ByteToken + " ")
			} else if i+1 < len(p) {
				out.WriteString(", ")
			}
		}
	}
	return out.String()
}

func FormatAssemblyBasicPalette(p color.Palette, eol string) string {
	var out strings.Builder
	for i := range p {
		v, err := constants.FirmwareNumber(p[i])
		if err == nil {
			fmt.Fprintf(&out, "%0.2d", v)
			if (i+1)%8 == 0 && i+1 < len(p) {
				out.WriteString(eol + ByteToken + " ")
			} else if i+1 < len(p) {
				out.WriteString(", ")
			}
		}
	}
	return out.String()
}

func FormatAssemblyCPCPlusPalette(p color.Palette, eol string) string {
	var out strings.Builder
	out.WriteString(ByteToken)
	for i := range p {
		cp := constants.NewCpcPlusColor(p[i])
		v := cp.Value()
		fmt.Fprintf(&out, "#%.2x, #%.2x", byte(v), byte(v>>8))
		if (i+1)%8 == 0 && i+1 < len(p) {
			out.WriteString(eol + ByteToken + " ")
		} else if i+1 < len(p) {
			out.WriteString(", ")
		}

	}
	return out.String()
}
