package gfx

import (
	"fmt"
	"github.com/jeromelesaux/dsk"
	"os"
	"path/filepath"
)

func ImportInDsk(exportType *ExportType) error {
	dskFullpath := exportType.Fullpath(".dsk")
	floppy := dsk.FormatDsk(9, 40)
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
	//	}

	return dsk.WriteDsk(dskFullpath, floppy)
}
