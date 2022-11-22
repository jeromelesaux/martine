package file

import (
	"fmt"
	"os"

	"github.com/jeromelesaux/martine/export/ocpartstudio"
)

func PalInformation(filePath string) {
	fmt.Fprintf(os.Stdout, "Input palette to open : (%s)\n", filePath)
	_, palette, err := ocpartstudio.OpenPal(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", filePath)
	} else {
		fmt.Fprintf(os.Stdout, "Palette from file %s\n\n%s", filePath, palette.ToString())
	}
}

func WinInformation(filePath string) {
	fmt.Fprintf(os.Stdout, "Input window to open : (%s)\n", filePath)
	win, err := ocpartstudio.OpenWin(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Window in file (%s) can not be read skipped\n", filePath)
	} else {
		fmt.Fprintf(os.Stdout, "Window from file %s\n\n%s", filePath, win.ToString())
	}
}

func KitInformation(filePath string) {
	fmt.Fprintf(os.Stdout, "Input kit palette to open : (%s)\n", filePath)
	_, palette, err := ocpartstudio.OpenKit(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", filePath)
	} else {
		fmt.Fprintf(os.Stdout, "Palette from file %s\n\n%s", filePath, palette.ToString())
	}
}
