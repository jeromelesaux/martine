package file

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func OpenImp(filePath string, mode int) (*ImpFooter, error) {

	footer := &ImpFooter{}
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return footer, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	fmt.Fprintf(os.Stdout, "LogicalSize=%d\n", header.LogicalSize)
	_, err = fr.Seek(0x80+int64(header.LogicalSize)-5, io.SeekStart)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while seek in the file (%s) with error %v\n", filePath, err)
		return footer, err
	}

	if err := binary.Read(fr, binary.LittleEndian, footer); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Win from file (%s) error %v\n", filePath, err)
		return footer, err
	}
	switch mode {
	case 0:
		footer.Width *= 2
	case 1:
		footer.Width *= 4
	case 2:
		footer.Width = 8
	}

	if footer.Width == 0 || footer.Height == 0 {
		return footer, errors.New("empty footer")
	}
	return footer, nil
}

func RawImp(filePath string) ([]byte, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return []byte{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Win Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}
	if header.Checksum != header.ComputedChecksum16() {
		fmt.Fprintf(os.Stderr, "Cannot read the Ocp Win Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		fr.Seek(0, io.SeekStart)
	}

	bf, err := ioutil.ReadAll(fr)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, len(bf)-3)
	copy(raw[:], bf[0:len(bf)-3])

	return raw, nil
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
