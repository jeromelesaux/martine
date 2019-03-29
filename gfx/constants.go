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
	HardwareValues []uint8
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
	Mode0         = Size{Width: 160, Height: 200, LinesNumber: 200, ColumnsNumber: 20, ColorsAvailable: 16}
	Mode1         = Size{Width: 320, Height: 200, LinesNumber: 200, ColumnsNumber: 40, ColorsAvailable: 4}
	Mode2         = Size{Width: 640, Height: 200, LinesNumber: 200, ColumnsNumber: 80, ColorsAvailable: 2}
	OverscanMode0 = Size{Width: 192, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 16}
	OverscanMode1 = Size{Width: 384, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 4}
	OverscanMode2 = Size{Width: 768, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 2}
	SelfMode      = Size{}
)
var (
	ErrorCpcColorNotFound = errors.New("Cpc color not found")
)

// values 50% RGB = 0x7F
// values 100% RGB = 0xFF
var (
	White         = CpcColor{HardwareNumber: 0, FirmwareNumber: 13, HardwareValues: []uint8{0x40}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0x7F, B: 0x7F}}
	SeaGreen      = CpcColor{HardwareNumber: 2, FirmwareNumber: 19, HardwareValues: []uint8{0x42, 0x51}, Color: color.RGBA{A: 0xFF, R: 0, G: 0xFF, B: 0x7F}}
	PastelYellow  = CpcColor{HardwareNumber: 3, FirmwareNumber: 25, HardwareValues: []uint8{0x43, 0x49}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0x7F}}
	Blue          = CpcColor{HardwareNumber: 4, FirmwareNumber: 1, HardwareValues: []uint8{0x44, 0x50}, Color: color.RGBA{A: 0xFF, R: 0, G: 0, B: 0x7F}}
	Purple        = CpcColor{HardwareNumber: 5, FirmwareNumber: 7, HardwareValues: []uint8{0x45, 0x48}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0, B: 0x7F}}
	Cyan          = CpcColor{HardwareNumber: 6, FirmwareNumber: 10, HardwareValues: []uint8{0x46}, Color: color.RGBA{A: 0xFF, R: 0, G: 0x7F, B: 0x7F}}
	Pink          = CpcColor{HardwareNumber: 7, FirmwareNumber: 16, HardwareValues: []uint8{0x40}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0x7F}}
	BrightYellow  = CpcColor{HardwareNumber: 10, FirmwareNumber: 24, HardwareValues: []uint8{0x4A}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0}}
	BrightWhite   = CpcColor{HardwareNumber: 11, FirmwareNumber: 26, HardwareValues: []uint8{0x4B}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0xFF}}
	BrightRed     = CpcColor{HardwareNumber: 12, FirmwareNumber: 6, HardwareValues: []uint8{0x4C}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0, B: 0}}
	BrightMagenta = CpcColor{HardwareNumber: 13, FirmwareNumber: 8, HardwareValues: []uint8{0x4D}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0, B: 0xFF}}
	Orange        = CpcColor{HardwareNumber: 14, FirmwareNumber: 15, HardwareValues: []uint8{0x4E}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0}}
	PastelMagenta = CpcColor{HardwareNumber: 15, FirmwareNumber: 17, HardwareValues: []uint8{0x4F}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0xFF}}
	BrightGreen   = CpcColor{HardwareNumber: 18, FirmwareNumber: 18, HardwareValues: []uint8{0x52}, Color: color.RGBA{A: 0xFF, R: 0, G: 0xFF, B: 0}}
	BrightCyan    = CpcColor{HardwareNumber: 19, FirmwareNumber: 20, HardwareValues: []uint8{0x53}, Color: color.RGBA{A: 0xFF, R: 0, G: 0xFF, B: 0xFF}}
	Black         = CpcColor{HardwareNumber: 20, FirmwareNumber: 0, HardwareValues: []uint8{0x54}, Color: color.RGBA{A: 0xFF, R: 0, G: 0, B: 0}}
	BrightBlue    = CpcColor{HardwareNumber: 21, FirmwareNumber: 2, HardwareValues: []uint8{0x55}, Color: color.RGBA{A: 0xFF, R: 0, G: 0, B: 0xFF}}
	Green         = CpcColor{HardwareNumber: 22, FirmwareNumber: 9, HardwareValues: []uint8{0x56}, Color: color.RGBA{A: 0xFF, R: 0, G: 0x7F, B: 0}}
	SkyBlue       = CpcColor{HardwareNumber: 23, FirmwareNumber: 11, HardwareValues: []uint8{0x57}, Color: color.RGBA{A: 0xFF, R: 0, G: 0x7F, B: 0xFF}}
	Magenta       = CpcColor{HardwareNumber: 24, FirmwareNumber: 4, HardwareValues: []uint8{0x58}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0, B: 0x7F}}
	PastelGreen   = CpcColor{HardwareNumber: 25, FirmwareNumber: 22, HardwareValues: []uint8{0x59}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0xFF, B: 0x7F}}
	Lime          = CpcColor{HardwareNumber: 26, FirmwareNumber: 21, HardwareValues: []uint8{0x5A}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0xFF, B: 0}}
	PastelCyan    = CpcColor{HardwareNumber: 27, FirmwareNumber: 23, HardwareValues: []uint8{0x5B}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0xFF, B: 0xFF}}
	Red           = CpcColor{HardwareNumber: 28, FirmwareNumber: 3, HardwareValues: []uint8{0x5C}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0, B: 0}}
	Mauve         = CpcColor{HardwareNumber: 29, FirmwareNumber: 5, HardwareValues: []uint8{0x5D}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0, B: 0xFF}}
	Yellow        = CpcColor{HardwareNumber: 30, FirmwareNumber: 12, HardwareValues: []uint8{0x5E}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0x7F, B: 0}}
	PastelBlue    = CpcColor{HardwareNumber: 31, FirmwareNumber: 14, HardwareValues: []uint8{0x5F}, Color: color.RGBA{A: 0xFF, R: 0x7F, G: 0x7F, B: 0xFF}}
)

func NewCpcPlusPalette() color.Palette {
	plusPalette := color.Palette{}
	var r, g, b uint8
	for g = 0; g < 0x10; g++ {
		for r = 0; r < 0x10; r++ {
			for b = 0; b < 0x10; b++ {
				//fmt.Fprintf(os.Stderr,"R:%d,G:%d,B:%d\n",r*0x33,g*0x33,b*0x33)
				plusPalette = append(plusPalette, color.RGBA{R: r * 0x33, B: b * 0x33, G: g * 0x33, A: 0xFF})
			}
		}
	}
	return plusPalette
}

var CpcPlusPalette = NewCpcPlusPalette()

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

func HardwareValueAreEquals(hv []uint8, val uint8) bool {
	for _, v := range hv {
		if v == val {
			return true
		}
	}
	return false
}

func ColorFromHardware(c uint8) (color.Color, error) {
	if HardwareValueAreEquals(White.HardwareValues, c) {
		return White.Color, nil
	}
	if HardwareValueAreEquals(SeaGreen.HardwareValues, c) {
		return SeaGreen.Color, nil
	}
	if HardwareValueAreEquals(PastelYellow.HardwareValues, c) {
		return PastelYellow.Color, nil
	}
	if HardwareValueAreEquals(Blue.HardwareValues, c) {
		return Blue.Color, nil
	}
	if HardwareValueAreEquals(Purple.HardwareValues, c) {
		return Purple.Color, nil
	}
	if HardwareValueAreEquals(Cyan.HardwareValues, c) {
		return Cyan.Color, nil
	}
	if HardwareValueAreEquals(Pink.HardwareValues, c) {
		return Pink.Color, nil
	}
	if HardwareValueAreEquals(BrightYellow.HardwareValues, c) {
		return BrightYellow.Color, nil
	}
	if HardwareValueAreEquals(BrightWhite.HardwareValues, c) {
		return BrightWhite.Color, nil
	}
	if HardwareValueAreEquals(BrightRed.HardwareValues, c) {
		return BrightRed.Color, nil
	}
	if HardwareValueAreEquals(BrightMagenta.HardwareValues, c) {
		return BrightMagenta.Color, nil
	}
	if HardwareValueAreEquals(Orange.HardwareValues, c) {
		return Orange.Color, nil
	}
	if HardwareValueAreEquals(PastelMagenta.HardwareValues, c) {
		return PastelMagenta.Color, nil
	}
	if HardwareValueAreEquals(BrightGreen.HardwareValues, c) {
		return BrightGreen.Color, nil
	}
	if HardwareValueAreEquals(BrightCyan.HardwareValues, c) {
		return BrightCyan.Color, nil
	}
	if HardwareValueAreEquals(Black.HardwareValues, c) {
		return Black.Color, nil
	}
	if HardwareValueAreEquals(BrightBlue.HardwareValues, c) {
		return BrightBlue.Color, nil
	}
	if HardwareValueAreEquals(Green.HardwareValues, c) {
		return Green.Color, nil
	}
	if HardwareValueAreEquals(SkyBlue.HardwareValues, c) {
		return SkyBlue.Color, nil
	}
	if HardwareValueAreEquals(Magenta.HardwareValues, c) {
		return Magenta.Color, nil
	}
	if HardwareValueAreEquals(PastelGreen.HardwareValues, c) {
		return PastelGreen.Color, nil
	}
	if HardwareValueAreEquals(Lime.HardwareValues, c) {
		return Lime.Color, nil
	}
	if HardwareValueAreEquals(PastelCyan.HardwareValues, c) {
		return PastelCyan.Color, nil
	}
	if HardwareValueAreEquals(Red.HardwareValues, c) {
		return Red.Color, nil
	}
	if HardwareValueAreEquals(Mauve.HardwareValues, c) {
		return Mauve.Color, nil
	}
	if HardwareValueAreEquals(Yellow.HardwareValues, c) {
		return Yellow.Color, nil
	}
	if HardwareValueAreEquals(PastelBlue.HardwareValues, c) {
		return PastelBlue.Color, nil
	}
	return nil, ErrorCpcColorNotFound
}

func HardwareValues(c color.Color) ([]uint8, error) {
	if ColorsAreEquals(White.Color, c) {
		return White.HardwareValues, nil
	}
	if ColorsAreEquals(SeaGreen.Color, c) {
		return SeaGreen.HardwareValues, nil
	}
	if ColorsAreEquals(PastelYellow.Color, c) {
		return PastelYellow.HardwareValues, nil
	}
	if ColorsAreEquals(Blue.Color, c) {
		return Blue.HardwareValues, nil
	}
	if ColorsAreEquals(Purple.Color, c) {
		return Purple.HardwareValues, nil
	}
	if ColorsAreEquals(Cyan.Color, c) {
		return Cyan.HardwareValues, nil
	}
	if ColorsAreEquals(Pink.Color, c) {
		return Pink.HardwareValues, nil
	}
	if ColorsAreEquals(BrightYellow.Color, c) {
		return BrightYellow.HardwareValues, nil
	}
	if ColorsAreEquals(BrightWhite.Color, c) {
		return BrightWhite.HardwareValues, nil
	}
	if ColorsAreEquals(BrightRed.Color, c) {
		return BrightRed.HardwareValues, nil
	}
	if ColorsAreEquals(BrightMagenta.Color, c) {
		return BrightMagenta.HardwareValues, nil
	}
	if ColorsAreEquals(Orange.Color, c) {
		return Orange.HardwareValues, nil
	}
	if ColorsAreEquals(PastelMagenta.Color, c) {
		return PastelMagenta.HardwareValues, nil
	}
	if ColorsAreEquals(BrightGreen.Color, c) {
		return BrightGreen.HardwareValues, nil
	}
	if ColorsAreEquals(BrightCyan.Color, c) {
		return BrightCyan.HardwareValues, nil
	}
	if ColorsAreEquals(Black.Color, c) {
		return Black.HardwareValues, nil
	}
	if ColorsAreEquals(BrightBlue.Color, c) {
		return BrightBlue.HardwareValues, nil
	}
	if ColorsAreEquals(Green.Color, c) {
		return Green.HardwareValues, nil
	}
	if ColorsAreEquals(SkyBlue.Color, c) {
		return SkyBlue.HardwareValues, nil
	}
	if ColorsAreEquals(Magenta.Color, c) {
		return Magenta.HardwareValues, nil
	}
	if ColorsAreEquals(PastelGreen.Color, c) {
		return PastelGreen.HardwareValues, nil
	}
	if ColorsAreEquals(Lime.Color, c) {
		return Lime.HardwareValues, nil
	}
	if ColorsAreEquals(PastelCyan.Color, c) {
		return PastelCyan.HardwareValues, nil
	}
	if ColorsAreEquals(Red.Color, c) {
		return Red.HardwareValues, nil
	}
	if ColorsAreEquals(Mauve.Color, c) {
		return Mauve.HardwareValues, nil
	}
	if ColorsAreEquals(Yellow.Color, c) {
		return Yellow.HardwareValues, nil
	}
	if ColorsAreEquals(PastelBlue.Color, c) {
		return PastelBlue.HardwareValues, nil
	}
	return nil, ErrorCpcColorNotFound

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
	return -1, ErrorCpcColorNotFound

}
