package window

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
)

type OcpWinFooter struct {
	Unused2 byte
	Width   uint16
	Height  byte
	Unused  byte
}

func (o *OcpWinFooter) ToString() string {
	return fmt.Sprintf("Width:(%d)\nHeight:(%d)\n", o.Width/8, o.Height)
}

func RawWin(filePath string) ([]byte, error) {
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
	raw := make([]byte, len(bf)-5)
	copy(raw[:], bf[0:len(bf)-5])

	if raw[0] == 'M' && raw[1] == 'J' && raw[2] == 'H' { // Compression OCP
		return ocpartstudio.DepackOCP(raw)
	}

	return raw, nil
}

func OpenWin(filePath string) (*OcpWinFooter, error) {
	fr, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", filePath, err)
		return &OcpWinFooter{}, err
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

	ocpWinFooter := &OcpWinFooter{}
	//_, err = fr.Seek(-5, io.SeekEnd)

	fmt.Fprintf(os.Stdout, "LogicalSize=%d\n", header.LogicalSize)
	_, err = fr.Seek(0x80+int64(header.LogicalSize)-5, io.SeekStart)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while seek in the file (%s) with error %v\n", filePath, err)
		return &OcpWinFooter{}, err
	}

	if err := binary.Read(fr, binary.LittleEndian, ocpWinFooter); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading Ocp Win from file (%s) error %v\n", filePath, err)
		return ocpWinFooter, err
	}
	ocpWinFooter.Width = uint16(uint(ocpWinFooter.Width / 8))
	return ocpWinFooter, nil
}

func Win(filePath string, data []byte, screenMode uint8, width, height int, dontImportDsk bool, cont *export.MartineConfig) error {
	osFilepath := cont.AmsdosFullPath(filePath, ".WIN")
	fmt.Fprintf(os.Stdout, "Saving WIN file (%s), screen mode %d, (%d,%d)\n", osFilepath, screenMode, width, height)
	win := OcpWinFooter{Unused: 3, Height: byte(height), Unused2: 0, Width: uint16(width * 8)}

	data, _ = compression.Compress(data, cont.Compression)

	//fmt.Fprintf(os.Stderr, "Header length %d\n", binary.Size(header))
	fmt.Fprintf(os.Stderr, "Data length %d\n", binary.Size(data))
	fmt.Fprintf(os.Stderr, "Footer length %d\n", binary.Size(win))
	osFilename := cont.Fullpath(".WIN")

	body, err := common.StructToBytes(data)
	if err != nil {
		return err
	}
	footer, err := common.StructToBytes(win)
	if err != nil {
		return err
	}
	content := body
	content = append(content, footer...)

	fmt.Fprintf(os.Stdout, "%s, data size :%d\n", win.ToString(), len(data))
	if !cont.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilename, ".WIN", content, 2, 0, 0x4000, 0x4000); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilename, content); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cont.AddFile(osFilepath)
	}
	return nil
}

func WinInformation(filePath string) {
	fmt.Fprintf(os.Stdout, "Input window to open : (%s)\n", filePath)
	win, err := OpenWin(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Window in file (%s) can not be read skipped\n", filePath)
	} else {
		fmt.Fprintf(os.Stdout, "Window from file %s\n\n%s", filePath, win.ToString())
	}
}
