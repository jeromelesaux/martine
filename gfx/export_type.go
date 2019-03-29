package gfx

import (
	"path/filepath"
	"strings"
	"sync"
)

var amsdosFilenameOnce sync.Once

type ExportType struct {
	InputPath      string
	OutputPath     string
	NoAmsdosHeader bool
	RollMode       bool
	RollIteration  int
	Dsk            bool
	Ink            bool
	Pal            bool
	Scr            bool
	Win            bool
	Overscan       bool
	Json           bool
	Ascii          bool
	CpcPlus        bool
	Width          int
	Height         int
	amsdosFilename []byte
}

func NewExportType(input, output string) *ExportType {
	return &ExportType{Json: true, Ascii: true, Scr: true, Pal: true, InputPath: input, OutputPath: output, amsdosFilename: make([]byte, 8)}
}

func (e *ExportType) AmsdosFilename() []byte {
	amsdosFilenameOnce.Do(
		func() {
			for i:=0; i < 8 ; i++ {
				e.amsdosFilename[i] =' '
			}
			file := strings.ToUpper(filepath.Base(e.InputPath))
			filename := strings.TrimSuffix(file, filepath.Ext(file))
			filenameSize := len(filename)
			if filenameSize > 8 {
				filenameSize = 8
			}
			copy(e.amsdosFilename[:], file[0:filenameSize])
		})
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
