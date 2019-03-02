package gfx

import (
	"errors"
	"fmt"
	"image/color"
)

type Size struct {
	Width           int
	Height          int
	LinesNumber     int
	ColumnsNumber   int
	ColorsAvailable int
}

type CpcColor struct {
	HardwareNumber int
	HardwareValues []int16
	FirmwareNumber int
	Color          color.RGBA
}

func (s *Size) ToString() string {
	return fmt.Sprintf("Size:\nWidth (%d) pixels\nHigh (%d) pixels\nNumber of lines (%d)\nNumber of columns (%d)\nColors available in this mode (%d)\n",
		s.Width,
		s.Height,
		s.LinesNumber,
		s.ColumnsNumber,
		s.ColorsAvailable)
}

var (
	Mode0    = Size{Width: 160, Height: 200, LinesNumber: 200, ColumnsNumber: 20, ColorsAvailable: 16}
	Mode1    = Size{Width: 320, Height: 200, LinesNumber: 200, ColumnsNumber: 40, ColorsAvailable: 4}
	Mode2    = Size{Width: 640, Height: 200, LinesNumber: 200, ColumnsNumber: 80, ColorsAvailable: 2}
	OverscanMode0 = Size{Width: 192, Height: 232, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 16}
	OverscanMode1 = Size{Width: 384, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 4}
	OverscanMode2 = Size{Width: 768, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 2}
	SelfMode = Size{}
)
var (
	CpcColorNotFound = errors.New("Cpc color not found")
)

// values 50% RGB = 0x7F
// values 100% RGB = 0xFF
var (
	White         = CpcColor{HardwareNumber: 0, FirmwareNumber: 13, HardwareValues: []int16{0x40}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0x7F, B: 0x7F}}
	SeaGreen      = CpcColor{HardwareNumber: 2, FirmwareNumber: 19, HardwareValues: []int16{0x42, 0x51}, Color: color.RGBA{A: 0xFF, R: 0, G: 0xFF, B: 0x7F}}
	PastelYellow  = CpcColor{HardwareNumber: 3, FirmwareNumber: 25, HardwareValues: []int16{0x43, 0x49}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0x7F}}
	Blue          = CpcColor{HardwareNumber: 4, FirmwareNumber: 1, HardwareValues: []int16{0x44, 0x50}, Color: color.RGBA{A: 0xFF, R: 0, G: 0, B: 0x7F}}
	Purple        = CpcColor{HardwareNumber: 5, FirmwareNumber: 7, HardwareValues: []int16{0x45, 0x48}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0, B: 0x7F}}
	Cyan          = CpcColor{HardwareNumber: 6, FirmwareNumber: 10, HardwareValues: []int16{0x46}, Color: color.RGBA{A: 0xFF, R: 0, G: 0x7F, B: 0x7F}}
	Pink          = CpcColor{HardwareNumber: 7, FirmwareNumber: 16, HardwareValues: []int16{0x40}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0x7F}}
	BrightYellow  = CpcColor{HardwareNumber: 10, FirmwareNumber: 24, HardwareValues: []int16{0x4A}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0}}
	BrightWhite   = CpcColor{HardwareNumber: 11, FirmwareNumber: 26, HardwareValues: []int16{0x4B}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0xFF}}
	BrightRed     = CpcColor{HardwareNumber: 12, FirmwareNumber: 6, HardwareValues: []int16{0x4C}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0, B: 0}}
	BrightMagenta = CpcColor{HardwareNumber: 13, FirmwareNumber: 8, HardwareValues: []int16{0x4D}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0, B: 0xFF}}
	Orange        = CpcColor{HardwareNumber: 14, FirmwareNumber: 15, HardwareValues: []int16{0x4E}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0}}
	PastelMagenta = CpcColor{HardwareNumber: 15, FirmwareNumber: 17, HardwareValues: []int16{0x4F}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0xFF}}
	BrightGreen   = CpcColor{HardwareNumber: 18, FirmwareNumber: 18, HardwareValues: []int16{0x52}, Color: color.RGBA{A: 0xFF, R: 0, G: 0xFF, B: 0}}
	BrightCyan    = CpcColor{HardwareNumber: 19, FirmwareNumber: 20, HardwareValues: []int16{0x53}, Color: color.RGBA{A: 0xFF, R: 0, G: 0xFF, B: 0xFF}}
	Black         = CpcColor{HardwareNumber: 20, FirmwareNumber: 0, HardwareValues: []int16{0x54}, Color: color.RGBA{A: 0xFF, R: 0, G: 0, B: 0}}
	BrightBlue    = CpcColor{HardwareNumber: 21, FirmwareNumber: 2, HardwareValues: []int16{0x55}, Color: color.RGBA{A: 0xFF, R: 0, G: 0, B: 0xFF}}
	Green         = CpcColor{HardwareNumber: 22, FirmwareNumber: 9, HardwareValues: []int16{0x56}, Color: color.RGBA{A: 0xFF, R: 0, G: 0x7F, B: 0}}
	SkyBlue       = CpcColor{HardwareNumber: 23, FirmwareNumber: 11, HardwareValues: []int16{0x57}, Color: color.RGBA{A: 0xFF, R: 0, G: 0x7F, B: 0xFF}}
	Magenta       = CpcColor{HardwareNumber: 24, FirmwareNumber: 4, HardwareValues: []int16{0x58}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0, B: 0x7F}}
	PastelGreen   = CpcColor{HardwareNumber: 25, FirmwareNumber: 22, HardwareValues: []int16{0x59}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0xFF, B: 0x7F}}
	Lime          = CpcColor{HardwareNumber: 26, FirmwareNumber: 21, HardwareValues: []int16{0x5A}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0xFF, B: 0}}
	PastelCyan    = CpcColor{HardwareNumber: 27, FirmwareNumber: 23, HardwareValues: []int16{0x5B}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0xFF, B: 0xFF}}
	Red           = CpcColor{HardwareNumber: 28, FirmwareNumber: 3, HardwareValues: []int16{0x5C}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0, B: 0}}
	Mauve         = CpcColor{HardwareNumber: 29, FirmwareNumber: 5, HardwareValues: []int16{0x5D}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0, B: 0xFF}}
	Yellow        = CpcColor{HardwareNumber: 30, FirmwareNumber: 12, HardwareValues: []int16{0x5E}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0x7F, B: 0}}
	PastelBlue    = CpcColor{HardwareNumber: 31, FirmwareNumber: 14, HardwareValues: []int16{0x5F}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0x7F, B: 0xFF}}
)

var CpcOldPalette = color.Palette{White.Color,
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

func ColorsAreEquals(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	if r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2 {
		return true
	}
	return false
}

func FirmwareNumber(c color.Color) (int, error) {
	if ColorsAreEquals(White.Color, c) {
		return White.FirmwareNumber, nil
	}
	if ColorsAreEquals(SeaGreen.Color, c) {
		return SeaGreen.FirmwareNumber, nil
	}
	if ColorsAreEquals(PastelYellow.Color, c) {
		return PastelYellow.FirmwareNumber, nil
	}
	if ColorsAreEquals(Blue.Color, c) {
		return Blue.FirmwareNumber, nil
	}
	if ColorsAreEquals(Purple.Color, c) {
		return Purple.FirmwareNumber, nil
	}
	if ColorsAreEquals(Cyan.Color, c) {
		return Cyan.FirmwareNumber, nil
	}
	if ColorsAreEquals(Pink.Color, c) {
		return Pink.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightYellow.Color, c) {
		return BrightYellow.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightWhite.Color, c) {
		return BrightWhite.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightRed.Color, c) {
		return BrightRed.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightMagenta.Color, c) {
		return BrightMagenta.FirmwareNumber, nil
	}
	if ColorsAreEquals(Orange.Color, c) {
		return Orange.FirmwareNumber, nil
	}
	if ColorsAreEquals(PastelMagenta.Color, c) {
		return PastelMagenta.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightGreen.Color, c) {
		return BrightGreen.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightCyan.Color, c) {
		return BrightCyan.FirmwareNumber, nil
	}
	if ColorsAreEquals(Black.Color, c) {
		return Black.FirmwareNumber, nil
	}
	if ColorsAreEquals(BrightBlue.Color, c) {
		return BrightBlue.FirmwareNumber, nil
	}
	if ColorsAreEquals(Green.Color, c) {
		return Green.FirmwareNumber, nil
	}
	if ColorsAreEquals(SkyBlue.Color, c) {
		return SkyBlue.FirmwareNumber, nil
	}
	if ColorsAreEquals(Magenta.Color, c) {
		return Magenta.FirmwareNumber, nil
	}
	if ColorsAreEquals(PastelGreen.Color, c) {
		return PastelGreen.FirmwareNumber, nil
	}
	if ColorsAreEquals(Lime.Color, c) {
		return Lime.FirmwareNumber, nil
	}
	if ColorsAreEquals(PastelCyan.Color, c) {
		return PastelCyan.FirmwareNumber, nil
	}
	if ColorsAreEquals(Red.Color, c) {
		return Red.FirmwareNumber, nil
	}
	if ColorsAreEquals(Mauve.Color, c) {
		return Mauve.FirmwareNumber, nil
	}
	if ColorsAreEquals(Yellow.Color, c) {
		return Yellow.FirmwareNumber, nil
	}
	if ColorsAreEquals(PastelBlue.Color, c) {
		return PastelBlue.FirmwareNumber, nil
	}
	return -1, CpcColorNotFound

}
