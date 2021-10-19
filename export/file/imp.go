package file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/m4client/cpc"
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
	header := cpc.CpcHead{Type: 0, User: 0, Address: 0x4000, Exec: 0x0,
		Size:        uint16(binary.Size(output)),
		Size2:       uint16(binary.Size(output)),
		LogicalSize: uint16(binary.Size(output))}
	copy(header.Filename[:], export.GetAmsdosFilename(filename, ".IMP"))
	header.Checksum = uint16(header.ComputedChecksum16())
	impPath := filepath.Join(export.OutputPath, export.GetAmsdosFilename(filename, ".IMP"))
	fw, err := os.Create(impPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", impPath, err)
		return err
	}
	if !export.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, output)
	fw.Close()
	fmt.Fprintf(os.Stdout, "Imp-Catcher file exported in [%s]\n", impPath)
	export.AddFile(impPath)
	return nil
}

func TileMap(data []byte, filename string, export *x.MartineContext) error {

	output := make([]byte, 0x4000)
	copy(output[0:], data[:])
	header := cpc.CpcHead{Type: 0, User: 0, Address: 0x4000, Exec: 0x0,
		Size:        uint16(binary.Size(output)),
		Size2:       uint16(binary.Size(output)),
		LogicalSize: uint16(binary.Size(output))}
	copy(header.Filename[:], export.GetAmsdosFilename(filename, ".TIL"))
	header.Checksum = uint16(header.ComputedChecksum16())
	impPath := filepath.Join(export.OutputPath, export.GetAmsdosFilename(filename, ".TIL"))
	fw, err := os.Create(impPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file (%s) error :%s\n", impPath, err)
		return err
	}
	if !export.NoAmsdosHeader {
		binary.Write(fw, binary.LittleEndian, header)
	}
	binary.Write(fw, binary.LittleEndian, output)
	fw.Close()
	fmt.Fprintf(os.Stdout, "Imp-TileMap file exported in [%s]\n", impPath)
	export.AddFile(impPath)
	return nil
}
