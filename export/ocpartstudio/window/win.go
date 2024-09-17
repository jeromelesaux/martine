package window

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/log"
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
		log.GetLogger().Error("Error while opening file (%s) error %v\n", filePath, err)
		return []byte{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		log.GetLogger().Error("Cannot read the Ocp Win Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		_, err = fr.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}
	if header.Checksum != header.ComputedChecksum16() {
		log.GetLogger().Error("Cannot read the Ocp Win Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		_, err = fr.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	bf, err := io.ReadAll(fr)
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
		log.GetLogger().Error("Error while opening file (%s) error %v\n", filePath, err)
		return &OcpWinFooter{}, err
	}
	header := &cpc.CpcHead{}
	if err := binary.Read(fr, binary.LittleEndian, header); err != nil {
		log.GetLogger().Error("Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		_, err := fr.Seek(0, io.SeekStart)
		if err != nil {
			return &OcpWinFooter{}, err
		}
	}
	if header.Checksum != header.ComputedChecksum16() {
		log.GetLogger().Error("Cannot read the Ocp Amsdos header (%s) with error :%v, trying to skip it\n", filePath, err)
		_, err := fr.Seek(0, io.SeekStart)
		if err != nil {
			return &OcpWinFooter{}, err
		}
	}

	ocpWinFooter := &OcpWinFooter{}

	log.GetLogger().Info("LogicalSize=%d\n", header.LogicalSize)
	_, err = fr.Seek(0x80+int64(header.LogicalSize)-5, io.SeekStart)

	if err != nil {
		log.GetLogger().Error("Error while seek in the file (%s) with error %v\n", filePath, err)
		return &OcpWinFooter{}, err
	}

	if err := binary.Read(fr, binary.LittleEndian, ocpWinFooter); err != nil {
		log.GetLogger().Error("Error while reading Ocp Win from file (%s) error %v\n", filePath, err)
		return ocpWinFooter, err
	}
	ocpWinFooter.Width = uint16(uint(ocpWinFooter.Width / 8))
	return ocpWinFooter, nil
}

func Win(filePath string, data []byte, screenMode uint8, width, height int, dontImportDsk bool, cfg *config.MartineConfig) error {
	osFilepath := cfg.AmsdosFullPath(filePath, ".WIN")
	log.GetLogger().Info("Saving WIN file (%s), screen mode %d, (%d,%d)\n", osFilepath, screenMode, width, height)
	win := OcpWinFooter{Unused: 3, Height: byte(height), Unused2: 0, Width: uint16(width * 8)}

	data, _ = compression.Compress(data, cfg.ScrCfg.Compression)

	// log.GetLogger().Error( "Header length %d\n", binary.Size(header))
	log.GetLogger().Error("Data length %d\n", binary.Size(data))
	log.GetLogger().Error("Footer length %d\n", binary.Size(win))

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

	log.GetLogger().Info("%s, data size :%d\n", win.ToString(), len(data))
	if !cfg.ScrCfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(osFilepath, ".WIN", content, 2, 0, 0x4000, 0x4000); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(osFilepath, content); err != nil {
			return err
		}
	}

	if !dontImportDsk {
		cfg.AddFile(osFilepath)
	}
	return nil
}

func WinInformation(filePath string) {
	log.GetLogger().Info("Input window to open : (%s)\n", filePath)
	win, err := OpenWin(filePath)
	if err != nil {
		log.GetLogger().Error("Window in file (%s) can not be read skipped\n", filePath)
	} else {
		log.GetLogger().Info("Window from file %s\n\n%s", filePath, win.ToString())
	}
}
