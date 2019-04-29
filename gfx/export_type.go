package gfx

import (
	"github.com/jeromelesaux/martine/constants"
	"path/filepath"
	"strings"
	"sync"
)

var amsdosFilenameOnce sync.Once

type ExportType struct {
	InputPath      string
	OutputPath     string
	M4RemotePath string
	M4Host string
	M4Autoexec bool
	Size           constants.Size
	Compression    int
	NoAmsdosHeader bool
	RotationMode   bool
	Rotation3DMode bool
	TileMode       bool
	RollMode       bool
	RollIteration  int
	TileIterationX int
	TileIterationY int
	M4 bool
	Dsk            bool
	Ink            bool
	Kit            bool
	Pal            bool
	Scr            bool
	Win            bool
	Overscan       bool
	Json           bool
	Ascii          bool
	CpcPlus        bool
	amsdosFilename []byte
	DskFiles       []string
	Tiles          *JsonSlice
}

func NewExportType(input, output string) *ExportType {
	return &ExportType{
		Json:           true,
		Ascii:          true,
		Scr:            true,
		Pal:            true,
		Ink:            true,
		InputPath:      input,
		OutputPath:     output,
		amsdosFilename: make([]byte, 8),
		DskFiles:       make([]string, 0),
		Tiles:          NewJsonSlice()}
}

func (e *ExportType) AddFile(file string) {
	e.DskFiles = append(e.DskFiles, file)
}

func (e *ExportType) AmsdosFilename() []byte {
	for i := 0; i < 8; i++ {
		e.amsdosFilename[i] = ' '
	}
	file := strings.ToUpper(filepath.Base(e.InputPath))
	filename := strings.TrimSuffix(file, filepath.Ext(file))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	copy(e.amsdosFilename[:], file[0:filenameSize])

	return e.amsdosFilename
}

func (e *ExportType) Filename() string {
	return string(e.AmsdosFilename())
}

func (e *ExportType) Fullpath(ext string) string {
	return e.OutputPath + string(filepath.Separator) + e.OsFilename(ext)
}

func (e *ExportType) TransformToAmsdosFile(filePath string) string {
	amsdosFile := make([]byte, 8)
	file := strings.ToUpper(filepath.Base(e.InputPath))
	filename := strings.TrimSuffix(file, filepath.Ext(file))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	copy(amsdosFile[:], file[0:filenameSize])
	return string(amsdosFile)

}

func (e *ExportType) OsFilename(ext string) string {
	file := strings.ToUpper(filepath.Base(e.InputPath))
	filename := strings.TrimSuffix(file, filepath.Ext(file))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	osFile := make([]byte, filenameSize)
	copy(osFile, filename[0:filenameSize])
	return string(osFile) + ext
}

func (e *ExportType) AmsdosFullPath(filePath string, newExtension string) string {
	filename := filepath.Base(filePath)
	file := strings.TrimSuffix(filename, filepath.Ext(filename))
	length := 8
	if len(file) < 8 {
		length = len(file)
	}
	newFilename := file[0:length] + newExtension
	return e.OutputPath + string(filepath.Separator) + strings.ToUpper(newFilename)
}

func (e *ExportType) OsFullPath(filePath string, newExtension string) string {
	filename := filepath.Base(filePath)
	file := strings.TrimSuffix(filename, filepath.Ext(filename))
	newFilename := file + newExtension
	return e.OutputPath + string(filepath.Separator) + newFilename
}
