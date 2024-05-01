package dialog

import (
	fyne "fyne.io/fyne/v2"
	dl "fyne.io/fyne/v2/dialog"
)

type dialogIface interface {
	Show()
}

func CheckAmsdosHeaderExport(inDsk, addAmsdosHeader bool, d dialogIface, w fyne.Window) {
	if inDsk && !addAmsdosHeader {
		dl.NewConfirm("Warning",
			"You are about to export files in DSK without Amsdos header, continue ? ",
			func(b bool) {
				if b {
					d.Show()
				} else {
					return
				}
			},
			w).Show()

	} else {
		d.Show()
	}
}
