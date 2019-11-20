package file

import (
	"fmt"
	"github.com/jeromelesaux/dsk"
	x "github.com/jeromelesaux/martine/export"
	"os"
	"path/filepath"
)

func ImportInDsk(filePath string, exportType *x.ExportType) error {
	var suffix string
	if exportType.EgxFormat == x.Egx1Mode {
		suffix += "-egx1"
	}
	if exportType.EgxFormat == x.Egx2Mode {
		suffix += "-egx2"
	}
	if exportType.CpcPlus {
		suffix += "-cpcplus"
	}
	if exportType.Overscan {
		suffix += "-overscan"
	}
	if exportType.Flash {
		suffix += "-flash"
	}
	if exportType.CustomDimension || exportType.SpriteHard {
		suffix += "-sprite"
	}
	if exportType.DitheringAlgo != 0 {
		suffix += "-dithering"
	}
	if exportType.SplitRaster {
		suffix += "-splitrasters"
	}

	dskFullpath := exportType.AmsdosFullPath(filePath, suffix+".dsk")

	var floppy *dsk.DSK
	if exportType.ExtendedDsk {
		floppy = dsk.FormatDsk(10, 80, 1, 1)
	} else {
		floppy = dsk.FormatDsk(9, 40, 1, 0)
	}

	dsk.WriteDsk(dskFullpath, floppy)
	for _, v := range exportType.DskFiles {
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
	return dsk.WriteDsk(dskFullpath, floppy)
}
