package gfx

import (
	"github.com/jeromelesaux/dsk"
	"strconv"
)

func ImportInDsk(exportType *ExportType) error {
	dskFullpath := exportType.Fullpath(".dsk")
	floppy := dsk.FormatDsk(9, 40)
	dsk.WriteDsk(dskFullpath, floppy)
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
	if exportType.RollMode {
		for i := 0; i < exportType.RollIteration; i++ {
			ext := strconv.Itoa(i) + ".WIN"
			floppy.PutFile(exportType.Fullpath(ext), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
			ext = strconv.Itoa(i) + ".PAL"
			floppy.PutFile(exportType.Fullpath(ext), dsk.MODE_BINAIRE, 0, 0, 0, false, false)
		}
	}

	return dsk.WriteDsk(dskFullpath, floppy)
}
