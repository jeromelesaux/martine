package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/compression"
)

// var amsdosFilenameOnce sync.Once
var (
	Egx1Mode = 1
	Egx2Mode = 2
)

var ErrorNotAllowed = errors.New("error not allowed")

type MartineConfig struct {
	InputPath                   string
	OutputPath                  string
	PalettePath                 string
	InkPath                     string
	KitPath                     string
	M4RemotePath                string
	M4Host                      string
	M4Autoexec                  bool
	Size                        constants.Size
	Compression                 compression.CompressionMethod
	NoAmsdosHeader              bool
	RotationMode                bool
	Rotation3DMode              bool
	Rotation3DX0                int
	Rotation3DY0                int
	Rotation3DType              int
	TileMode                    bool
	RollMode                    bool
	RollIteration               int
	TileIterationX              int
	TileIterationY              int
	M4                          bool
	Dsk                         bool
	Ink                         bool
	Kit                         bool
	Pal                         bool
	Scr                         bool
	Win                         bool
	Overscan                    bool
	Json                        bool
	Ascii                       bool
	CpcPlus                     bool
	CustomDimension             bool
	amsdosFilename              []byte
	DskFiles                    []string
	Tiles                       *export.JsonSlice
	DeltaMode                   bool
	ExtendedDsk                 bool
	ResizingAlgo                imaging.ResampleFilter
	DitheringAlgo               int
	DitheringMatrix             [][]float32
	DitheringMultiplier         float64
	DitheringWithQuantification bool
	DitheringType               constants.DitheringType
	RotationRraBit              int
	RotationRlaBit              int
	RotationSraBit              int
	RotationSlaBit              int
	RotationLosthighBit         int
	RotationLostlowBit          int
	RotationKeephighBit         int
	RotationKeeplowBit          int
	RotationIterations          int
	Flash                       bool
	FlashScreenFilepath1        string
	FlashScreenFilepath2        string
	FlashPaletteFilepath1       string
	FlashPaletteFilepath2       string
	EgxFormat                   int
	EgxMode1                    uint8
	EgxMode2                    uint8
	Sna                         bool
	SnaPath                     string
	SpriteHard                  bool
	SplitRaster                 bool
	ScanlineSequence            []int
	CustomScanlineSequence      bool
	MaskSprite                  uint8
	MaskOrOperation             bool
	MaskAndOperation            bool
	ZigZag                      bool
	Animate                     bool
	Reducer                     int
	OneLine                     bool
	OneRow                      bool
	InkSwapper                  map[int]int
	LineWidth                   int
	FilloutGif                  bool
	Saturation                  float64
	Brightness                  float64
	ExportAsGoFile              bool
	DoubleScreenAddress         bool
}

func MaskIsAllowed(mode uint8, value uint8) bool {
	values, err := ModeMaskSprite(mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with mode %d error :%v\n", mode, err)
		return false
	}
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func ModeMaskSprite(mode uint8) ([]uint8, error) {
	switch mode {
	case 0:
		return []uint8{0xAA, 0x55}, nil
	case 1:
		return []uint8{0x88, 0x44, 0x24, 0x11}, nil
	default:
		return make([]uint8, 0), ErrorNotAllowed
	}
}

func NewMartineConfig(input, output string) *MartineConfig {
	return &MartineConfig{
		Scr:            true,
		Pal:            true,
		Ink:            true,
		InputPath:      input,
		OutputPath:     output,
		amsdosFilename: make([]byte, 8),
		DskFiles:       make([]string, 0),
		Rotation3DX0:   -1,
		Rotation3DY0:   -1,
		Tiles:          export.NewJsonSlice(),
		InkSwapper:     make(map[int]int),
		LineWidth:      0x50,
	}
}

func (e *MartineConfig) AddFile(file string) {
	e.DskFiles = append(e.DskFiles, file)
}

func (e *MartineConfig) ImportInkSwap(s string) error {
	if s == "" {
		return nil
	}
	items := strings.Split(s, ",")
	for _, v := range items {
		var key, val int
		values := strings.Split(v, "=")
		if len(values) != 2 {
			return fmt.Errorf("expects two values parsed and gets %d, from [%s]",
				len(values),
				v)
		}
		key, err := strconv.Atoi(values[0])
		if err != nil {
			return err
		}
		val, err = strconv.Atoi(values[1])
		if err != nil {
			return err
		}
		e.InkSwapper[key] = val
	}

	return nil
}

func (e *MartineConfig) SwapInk(inkIndex int) int {
	if v, ok := e.InkSwapper[inkIndex]; ok {
		return v
	}
	return inkIndex
}

func RemoveUnsupportedChar(s string) string {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}

func (e *MartineConfig) AmsdosFilename() []byte {
	for i := 0; i < 8; i++ {
		e.amsdosFilename[i] = ' '
	}
	file := strings.ToUpper(filepath.Base(e.InputPath))
	filename := RemoveUnsupportedChar(strings.TrimSuffix(file, filepath.Ext(file)))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	copy(e.amsdosFilename[:], file[0:filenameSize])

	return e.amsdosFilename
}

func (e *MartineConfig) Filename() string {
	return string(e.AmsdosFilename())
}

func (e *MartineConfig) Fullpath(ext string) string {
	return filepath.Join(e.OutputPath, e.OsFilename(ext))
}

func (e *MartineConfig) TransformToAmsdosFile(filePath string) string {
	amsdosFile := make([]byte, 8)
	file := strings.ToUpper(filepath.Base(e.InputPath))
	filename := RemoveUnsupportedChar(strings.TrimSuffix(file, filepath.Ext(file)))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	copy(amsdosFile[:], file[0:filenameSize])
	return string(amsdosFile)
}

func AmsdosFilename(inputPath, ext string) string {
	file := strings.ToUpper(filepath.Base(inputPath))
	filename := RemoveUnsupportedChar(strings.TrimSuffix(file, filepath.Ext(file)))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	osFile := make([]byte, filenameSize)
	copy(osFile, filename[0:filenameSize])
	return string(osFile) + ext
}

func (e *MartineConfig) OsFilename(ext string) string {
	file := strings.ToUpper(filepath.Base(e.InputPath))
	filename := RemoveUnsupportedChar(strings.TrimSuffix(file, filepath.Ext(file)))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	osFile := make([]byte, filenameSize)
	copy(osFile, filename[0:filenameSize])
	return string(osFile) + ext
}

func (e *MartineConfig) GetAmsdosFilename(filePath string, ext string) string {
	file := strings.ToUpper(filepath.Base(filePath))
	filename := RemoveUnsupportedChar(strings.TrimSuffix(file, filepath.Ext(file)))
	filenameSize := 7
	end := len(filename)
	if len(filename) < 8 {
		filenameSize = len(filename) - 1
	}
	osFile := filename[0:filenameSize] + filename[end-1:end]
	return osFile + ext
}

func (e *MartineConfig) AmsdosFullPath(filePath string, newExtension string) string {
	filename := filepath.Base(filePath)
	file := RemoveUnsupportedChar(strings.TrimSuffix(filename, filepath.Ext(filename)))
	length := 6
	if len(file) < 8 {
		length = len(file)
	}

	newFilename := file[0:length] + newExtension
	return filepath.Join(e.OutputPath, strings.ToUpper(newFilename))
}

func (e *MartineConfig) OsFullPath(filePath string, newExtension string) string {
	filename := filepath.Base(filePath)
	file := RemoveUnsupportedChar(strings.TrimSuffix(filename, filepath.Ext(filename)))
	newFilename := file + newExtension
	return filepath.Join(e.OutputPath, newFilename)
}

func (e *MartineConfig) SetLineWith(i string) error {
	v, err := common.ParseHexadecimal8(i)
	if err != nil {
		return err
	}
	e.LineWidth = int(v)
	return nil
}
