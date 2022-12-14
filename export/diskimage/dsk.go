package diskimage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/dsk"
	"github.com/jeromelesaux/martine/config"
)

func ImportInDsk(filePath string, cfg *config.MartineConfig) error {
	var suffix string
	if cfg.EgxFormat == config.Egx1Mode {
		suffix += "-egx1"
	}
	if cfg.EgxFormat == config.Egx2Mode {
		suffix += "-egx2"
	}
	if cfg.CpcPlus {
		suffix += "-cpcplus"
	}
	if cfg.Overscan {
		suffix += "-overscan"
	}
	if cfg.Flash {
		suffix += "-flash"
	}
	if cfg.CustomDimension || cfg.SpriteHard {
		suffix += "-sprite"
	}
	if cfg.DitheringAlgo != 0 {
		suffix += "-dithering"
	}
	if cfg.SplitRaster {
		suffix += "-splitrasters"
	}

	dskFullpath := cfg.AmsdosFullPath(filePath, suffix+".dsk")

	var floppy *dsk.DSK
	if cfg.ExtendedDsk {
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
				fmt.Fprintf(os.Stderr, "Error while insert (%s) in dsk (%s) error :%v\n", v, dskFullpath, err)
			}
		} else {
			if err := floppy.PutFile(v, dsk.MODE_BINAIRE, 0, 0, 0, false, false); err != nil {
				fmt.Fprintf(os.Stderr, "Error while insert (%s) in dsk (%s) error :%v\n", v, dskFullpath, err)
			}
		}
	}
	fmt.Fprintf(os.Stdout, "Saving final dsk in path {%s}\n", dskFullpath)
	return dsk.WriteDsk(dskFullpath, floppy)
}
