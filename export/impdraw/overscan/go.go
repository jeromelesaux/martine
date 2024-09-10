package overscan

import (
	"image/color"

	"github.com/jeromelesaux/martine/common/errors"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/amsdos"
)

func SaveGo(filePath string, data []byte, p color.Palette, screenMode uint8, cfg *config.MartineConfig) error {
	data1, data2, err := OverscanToGo(data)
	if err != nil {
		return err
	}
	go1Filename := cfg.AmsdosFullPath(filePath, ".GO1")
	go2Filename := cfg.AmsdosFullPath(filePath, ".GO2")
	if !cfg.NoAmsdosHeader {
		if err := amsdos.SaveAmsdosFile(go1Filename, ".GO1", data1, 2, 0, 0x20, 0); err != nil {
			return err
		}
		if err := amsdos.SaveAmsdosFile(go2Filename, ".GO2", data2, 2, 0, 0x4000, 0); err != nil {
			return err
		}
	} else {
		if err := amsdos.SaveOSFile(go1Filename, data1); err != nil {
			return err
		}
		if err := amsdos.SaveOSFile(go2Filename, data2); err != nil {
			return err
		}
	}

	cfg.AddFile(go1Filename)
	cfg.AddFile(go2Filename)

	return nil
}

func OverscanToGo(data []byte) (go1, go2 []byte, err error) {
	if len(data) < 0x8000 {
		return go1, go2, errors.ErrorBadFileFormat
	}
	go1 = make([]byte, 0x4000)
	go2 = make([]byte, 0x4000)
	copy(go1, data[0:0x4000])
	copy(go2, data[0x3FE0:0x8000])
	return go1, go2, nil
}
