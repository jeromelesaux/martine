package diskimage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/dsk"
	x "github.com/jeromelesaux/martine/export"
)

func ImportInDsk(filePath string, cont *x.MartineConfig) error {
	var suffix string
	if cont.EgxFormat == x.Egx1Mode {
		suffix += "-egx1"
	}
	if cont.EgxFormat == x.Egx2Mode {
		suffix += "-egx2"
	}
	if cont.CpcPlus {
		suffix += "-cpcplus"
	}
	if cont.Overscan {
		suffix += "-overscan"
	}
	if cont.Flash {
		suffix += "-flash"
	}
	if cont.CustomDimension || cont.SpriteHard {
		suffix += "-sprite"
	}
	if cont.DitheringAlgo != 0 {
		suffix += "-dithering"
	}
	if cont.SplitRaster {
		suffix += "-splitrasters"
	}

	dskFullpath := cont.AmsdosFullPath(filePath, suffix+".dsk")

	var floppy *dsk.DSK
	if cont.ExtendedDsk {
		floppy = dsk.FormatDsk(10, 80, 1, 1, dsk.DataFormat)
	} else {
		floppy = dsk.FormatDsk(9, 40, 1, 0, dsk.DataFormat)
	}

	dsk.WriteDsk(dskFullpath, floppy)
	for _, v := range cont.DskFiles {
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
