package diskimage

import (
	"path/filepath"

	"github.com/jeromelesaux/dsk"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"
)

func ImportInDsk(filePath string, cfg *config.MartineConfig) error {
	var suffix string
	if cfg.ScrCfg.Type == config.Egx1Format {
		suffix += "-egx1"
	}
	if cfg.ScrCfg.Type == config.Egx2Format {
		suffix += "-egx2"
	}
	if cfg.ScrCfg.IsPlus {
		suffix += "-cpcplus"
	}
	if cfg.ScrCfg.Type == config.FullscreenFormat {
		suffix += "-overscan"
	}
	if cfg.Flash {
		suffix += "-flash"
	}
	if cfg.CustomDimension || cfg.ScrCfg.Type == config.SpriteHardFormat {
		suffix += "-sprite"
	}
	if cfg.ScrCfg.Process.DitheringAlgo != 0 {
		suffix += "-dithering"
	}
	if cfg.SplitRaster {
		suffix += "-splitrasters"
	}

	dskFullpath := cfg.AmsdosFullPath(filePath, suffix+".dsk")

	var floppy *dsk.DSK
	if cfg.HasContainerExport(config.ExtendedDskContainer) {
		floppy = dsk.FormatDsk(10, 80, 1, 1, dsk.DataFormat)
	} else {
		floppy = dsk.FormatDsk(9, 40, 1, 0, dsk.DataFormat)
	}

	err := dsk.WriteDsk(dskFullpath, floppy)
	if err != nil {
		return err
	}
	for _, v := range cfg.DskFiles {
		if filepath.Ext(v) == ".TXT" {
			if err := floppy.PutFile(v, dsk.MODE_ASCII, 0, 0, 0, false, false); err != nil {
				log.GetLogger().Error("Error while insert (%s) in dsk (%s) error :%v\n", v, dskFullpath, err)
			}
		} else {
			if err := floppy.PutFile(v, dsk.MODE_BINAIRE, 0, 0, 0, false, false); err != nil {
				log.GetLogger().Error("Error while insert (%s) in dsk (%s) error :%v\n", v, dskFullpath, err)
			}
		}
	}
	log.GetLogger().Info("Saving final dsk in path {%s}\n", dskFullpath)
	return dsk.WriteDsk(dskFullpath, floppy)
}
