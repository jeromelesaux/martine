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
	NbFrames uint
	Height   uint
	Width    uint
}

func Imp(sprites []byte, width, height uint, filename string, export *x.ExportType) error {
	impHeader := &ImpFooter{
		NbFrames: uint(len(sprites)),
		Height:   height,
		Width:    width,
	}
	output := make([]byte, 0)
	for _, v := range sprites {
		output = append(output, v)
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, impHeader)
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
