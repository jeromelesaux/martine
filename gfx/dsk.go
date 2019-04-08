package gfx

import (
	"github.com/jeromelesaux/dsk"
	"path/filepath"
)

func ImportInDsk(exportType *ExportType) error {
	dskFullpath := exportType.Fullpath(".dsk")
	floppy := dsk.FormatDsk(9, 40)
	dsk.WriteDsk(dskFullpath, floppy)
	if exportType.Kit {
		floppy.PutFile(exportType.Fullpath(".KIT"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
	}
	if exportType.Ink {
		floppy.PutFile(exportType.Fullpath(".INK"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
	}
	if exportType.Pal {
		floppy.PutFile(exportType.Fullpath(".PAL"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
	}
	if exportType.Scr {
		floppy.PutFile(exportType.Fullpath(".SCR"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
		floppy.PutFile(exportType.Fullpath(".BAS"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
	}
	if exportType.Overscan {
		floppy.PutFile(exportType.Fullpath(".SCR"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
	}
	if exportType.Win {
		floppy.PutFile(exportType.Fullpath(".WIN"), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
	}
	if exportType.Ascii {
		floppy.PutFile(exportType.Fullpath(".TXT"), dsk.MODE_ASCII, 0, 0, 0, false, false)
	}
	if exportType.RollMode || exportType.TileMode {
		for _, v := range exportType.DskFiles {
			if filepath.Ext(v) == ".TXT" {
				floppy.PutFile(v, dsk.MODE_ASCII, 0, 0, 0, false, false)
			} else {
				floppy.PutFile(v, dsk.MODE_BINAIRE, 0, 0, 0, false, false)
			}
		}
	}

	return dsk.WriteDsk(dskFullpath, floppy)
}
