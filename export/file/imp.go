package file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	x "github.com/jeromelesaux/martine/export"
)

type ImpFooter struct {
	Width    byte
	Height   byte
	NbFrames byte
}

func Imp(sprites []byte, nbFrames, width, height, mode uint, filename string, export *x.MartineContext) error {
	w := width
	switch mode {
	case 0:
		w /= 2
	case 1:
		w /= 4
	case 2:
		w /= 8
	}
	impHeader := ImpFooter{
		Width:    byte(w),
		Height:   byte(height),
		NbFrames: byte(nbFrames),
	}
	output := make([]byte, 0)
	output = append(output, sprites...)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, impHeader); err != nil {
		fmt.Fprintf(os.Stderr, "Error while feeding imp header. error :%v\n", err)
	}
	output = append(output, buf.Bytes()...)

	impPath := filepath.Join(export.OutputPath, export.GetAmsdosFilename(filename, ".IMP"))

	if !export.NoAmsdosHeader {
		if err := SaveAmsdosFile(impPath, ".IMP", output, 0, 0, 0x4000, 0x0); err != nil {
			return err
		}
	} else {
		if err := SaveOSFile(impPath, output); err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "Imp-Catcher file exported in [%s]\n", impPath)
	export.AddFile(impPath)
	return nil
}

func TileMap(data []byte, filename string, export *x.MartineContext) error {

	output := make([]byte, 0x4000)
	copy(output[0:], data[:])

	impPath := filepath.Join(export.OutputPath, export.GetAmsdosFilename(filename, ".TIL"))

	if !export.NoAmsdosHeader {
		if err := SaveAmsdosFile(impPath, ".TIL", output, 0, 0, 0x4000, 0); err != nil {
			return err
		}
	} else {
		if err := SaveOSFile(impPath, output); err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "Imp-TileMap file exported in [%s]\n", impPath)
	export.AddFile(impPath)
	return nil
}
