package overscan

import (
	"image/color"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/amsdos"
)

func SaveGo(filePath string, dataUp, dataDown []byte, p color.Palette, screenMode uint8, cfg *config.MartineConfig) error {
	data1 := make([]byte, 0x4000)
	data2 := make([]byte, 0x4000)
	copy(data1, dataUp)
	copy(data2, dataDown)
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
