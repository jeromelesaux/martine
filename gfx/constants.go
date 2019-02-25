package gfx

import (
	"image/color"
)

type Size struct {
	Width         int
	High          int
	LinesNumber   int
	ColumnsNumber int
}

type CpcColor struct {
	HardwareNumber int
	HardwareValues []int16
	FirmwareNumber int
	Color          color.RGBA
}

var (
	Mode0    = Size{Width: 160, High: 200, LinesNumber: 200, ColumnsNumber: 80}
	Mode1    = Size{Width: 320, High: 200, LinesNumber: 200, ColumnsNumber: 80}
	Mode2    = Size{Width: 640, High: 200, LinesNumber: 200, ColumnsNumber: 80}
	Overscan = Size{Width: 640, High: 400, LinesNumber: 272, ColumnsNumber: 96}
)

var (
	White         = CpcColor{HardwareNumber: 0, FirmwareNumber: 13, HardwareValues: []int16{0x40}, Color: color.RGBA{A: 1, R: 50, G: 50, B: 50}}
	SeaGreen      = CpcColor{HardwareNumber: 2, FirmwareNumber: 19, HardwareValues: []int16{0x42, 0x51}, Color: color.RGBA{A: 1, R: 0, G: 100, B: 50}}
	PastelYellow  = CpcColor{HardwareNumber: 3, FirmwareNumber: 25, HardwareValues: []int16{0x43, 0x49}, Color: color.RGBA{A: 1, R: 100, G: 100, B: 50}}
	Blue          = CpcColor{HardwareNumber: 4, FirmwareNumber: 1, HardwareValues: []int16{0x44, 0x50}, Color: color.RGBA{A: 1, R: 0, G: 0, B: 50}}
	Purple        = CpcColor{HardwareNumber: 5, FirmwareNumber: 7, HardwareValues: []int16{0x45, 0x48}, Color: color.RGBA{A: 1, R: 100, G: 0, B: 50}}
	Cyan          = CpcColor{HardwareNumber: 6, FirmwareNumber: 10, HardwareValues: []int16{0x46}, Color: color.RGBA{A: 1, R: 0, G: 50, B: 50}}
	Pink          = CpcColor{HardwareNumber: 7, FirmwareNumber: 16, HardwareValues: []int16{0x40}, Color: color.RGBA{A: 1, R: 100, G: 50, B: 50}}
	BrightYellow  = CpcColor{HardwareNumber: 10, FirmwareNumber: 24, HardwareValues: []int16{0x4A}, Color: color.RGBA{A: 1, R: 100, G: 100, B: 0}}
	BrightWhite   = CpcColor{HardwareNumber: 11, FirmwareNumber: 26, HardwareValues: []int16{0x4B}, Color: color.RGBA{A: 1, R: 100, G: 100, B: 100}}
	BrightRed     = CpcColor{HardwareNumber: 12, FirmwareNumber: 6, HardwareValues: []int16{0x4C}, Color: color.RGBA{A: 1, R: 100, G: 0, B: 0}}
	BrightMagenta = CpcColor{HardwareNumber: 13, FirmwareNumber: 8, HardwareValues: []int16{0x4D}, Color: color.RGBA{A: 1, R: 100, G: 0, B: 100}}
	Orange        = CpcColor{HardwareNumber: 14, FirmwareNumber: 15, HardwareValues: []int16{0x4E}, Color: color.RGBA{A: 1, R: 100, G: 50, B: 0}}
	PastelMagenta = CpcColor{HardwareNumber: 15, FirmwareNumber: 17, HardwareValues: []int16{0x4F}, Color: color.RGBA{A: 1, R: 100, G: 50, B: 100}}
	BrightGreen   = CpcColor{HardwareNumber: 18, FirmwareNumber: 18, HardwareValues: []int16{0x52}, Color: color.RGBA{A: 1, R: 0, G: 100, B: 0}}
	BrightCyan    = CpcColor{HardwareNumber: 19, FirmwareNumber: 20, HardwareValues: []int16{0x53}, Color: color.RGBA{A: 1, R: 0, G: 100, B: 100}}
	Black         = CpcColor{HardwareNumber: 20, FirmwareNumber: 0, HardwareValues: []int16{0x54}, Color: color.RGBA{A: 1, R: 0, G: 0, B: 0}}
	BrightBlue    = CpcColor{HardwareNumber: 21, FirmwareNumber: 2, HardwareValues: []int16{0x55}, Color: color.RGBA{A: 1, R: 0, G: 0, B: 100}}
	Green         = CpcColor{HardwareNumber: 22, FirmwareNumber: 9, HardwareValues: []int16{0x56}, Color: color.RGBA{A: 1, R: 0, G: 50, B: 0}}
	SkyBlue       = CpcColor{HardwareNumber: 23, FirmwareNumber: 11, HardwareValues: []int16{0x57}, Color: color.RGBA{A: 1, R: 0, G: 50, B: 100}}
	Magenta       = CpcColor{HardwareNumber: 24, FirmwareNumber: 4, HardwareValues: []int16{0x58}, Color: color.RGBA{A: 1, R: 50, G: 0, B: 50}}
	PastelGreen   = CpcColor{HardwareNumber: 25, FirmwareNumber: 22, HardwareValues: []int16{0x59}, Color: color.RGBA{A: 1, R: 50, G: 100, B: 50}}
	Lime          = CpcColor{HardwareNumber: 26, FirmwareNumber: 21, HardwareValues: []int16{0x5A}, Color: color.RGBA{A: 1, R: 50, G: 100, B: 0}}
	PastelCyan    = CpcColor{HardwareNumber: 27, FirmwareNumber: 23, HardwareValues: []int16{0x5B}, Color: color.RGBA{A: 1, R: 50, G: 100, B: 100}}
	Red           = CpcColor{HardwareNumber: 28, FirmwareNumber: 3, HardwareValues: []int16{0x5C}, Color: color.RGBA{A: 1, R: 50, G: 0, B: 0}}
	Mauve         = CpcColor{HardwareNumber: 29, FirmwareNumber: 5, HardwareValues: []int16{0x5D}, Color: color.RGBA{A: 1, R: 50, G: 0, B: 100}}
	Yellow        = CpcColor{HardwareNumber: 30, FirmwareNumber: 12, HardwareValues: []int16{0x5E}, Color: color.RGBA{A: 1, R: 50, G: 50, B: 0}}
	PastelBlue    = CpcColor{HardwareNumber: 31, FirmwareNumber: 14, HardwareValues: []int16{0x5F}, Color: color.RGBA{A: 1, R: 50, G: 50, B: 100}}
)

var CpcOldPalette = []color.Color{White.Color,
	SeaGreen.Color,
	PastelYellow.Color,
	Blue.Color,
	Purple.Color,
	Cyan.Color,
	Pink.Color,
	BrightYellow.Color,
	BrightWhite.Color,
	BrightRed.Color,
	BrightMagenta.Color,
	Orange.Color,
	PastelMagenta.Color,
	BrightGreen.Color,
	BrightCyan.Color,
	Black.Color,
	BrightBlue.Color,
	Green.Color,
	SkyBlue.Color,
	Magenta.Color,
	PastelGreen.Color,
	Lime.Color,
	PastelCyan.Color,
	Red.Color,
	Mauve.Color,
	Yellow.Color,
	PastelBlue.Color,
}