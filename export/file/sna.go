package file

import (
	"encoding/binary"
	"fmt"
	"github.com/jeromelesaux/dsk"
	m "github.com/jeromelesaux/m4client/cpc"
	"os"
)

func ImportInSna(filePath, snaPath string, screenMode uint8) error {
	sna := &dsk.SNA{Data: make([]byte, 0xFFFF), Header: dsk.NewSnaHeader()}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	header := &m.CpcHead{}
	if err := binary.Read(f, binary.LittleEndian, header); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Import file %s at address:#%4x size:%4x\n", filePath, header.Address, header.Size)
	buff := make([]byte, 0xFFFF)
	_, err = f.Read(buff)
	if err != nil {
		return err
	}
	if err := sna.Put(buff, header.Address, header.Size); err != nil {
		return err
	}
	sna.Header.RegisterPCHigh = 0x1
	sna.Header.RegisterPCLow = 0xad
	switch screenMode {
	case 0:
		sna.Header.GAMultiConfiguration = 0x8c
	case 1:
		sna.Header.GAMultiConfiguration = 0x8d
	case 2:
		sna.Header.GAMultiConfiguration = 0x8e
	}
	w, err := os.Create(snaPath)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := binary.Write(w, binary.LittleEndian, sna.Header); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, sna.Data); err != nil {
		return err
	}
	return nil
}
