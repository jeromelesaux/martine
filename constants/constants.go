package constants

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"sort"
	"strconv"
)

type DitheringType struct {
	string
}

var (
	OrderedDither        = DitheringType{"Ordered"}
	ErrorDiffusionDither = DitheringType{"ErrorDiffusion"}
)

type Size struct {
	Width           int
	Height          int
	LinesNumber     int
	ColumnsNumber   int
	ColorsAvailable int
	GatearrayValue  uint8
}

func NewSize(mode uint8) Size {
	s := Size{}
	switch mode {
	case 0:
		s = Mode0
	case 1:
		s = Mode1
	case 2:
		s = Mode2
	}
	return s
}

func NewSizeMode(mode uint8, overscan bool) Size {
	s := Size{}
	switch mode {
	case 0:
		s = Mode0
	case 1:
		s = Mode1
	case 2:
		s = Mode2
	}
	if overscan {
		switch mode {
		case 0:
			s = OverscanMode0
		case 1:
			s = OverscanMode1
		case 2:
			s = OverscanMode2
		}
	}
	return s
}

func (s *Size) ModeWidth(mode uint8) int {
	switch mode {
	case 0:
		return int(math.Ceil(float64(s.Width) / 2.))
	case 1:
		return int(math.Ceil(float64(s.Width) / 4.))
	case 2:
		return int(math.Ceil(float64(s.Width) / 8.))
	}
	return -1
}

type CpcColor struct {
	HardwareNumber int
	HardwareValues []uint8
	FirmwareNumber int
	Color          color.RGBA
}

func (c *CpcColor) ToString() string {
	return strconv.Itoa(c.HardwareNumber) + " firmware color :" + strconv.Itoa(c.FirmwareNumber) + " firmware value :" + strconv.Itoa(int(c.HardwareValues[uint8(0)]))
}

const (
	Rle = iota
	Rle16
)

func (s *Size) ToString() string {
	return fmt.Sprintf("Size:\nWidth (%d) pixels\nHigh (%d) pixels\nNumber of lines (%d)\nNumber of columns (%d)\nColors available in this mode (%d)\n",
		s.Width,
		s.Height,
		s.LinesNumber,
		s.ColumnsNumber,
		s.ColorsAvailable)
}

var (
	Mode0         = Size{Width: 160, Height: 200, LinesNumber: 200, ColumnsNumber: 20, ColorsAvailable: 16, GatearrayValue: 0x9c}
	Mode1         = Size{Width: 320, Height: 200, LinesNumber: 200, ColumnsNumber: 40, ColorsAvailable: 4, GatearrayValue: 0x9d}
	Mode2         = Size{Width: 640, Height: 200, LinesNumber: 200, ColumnsNumber: 80, ColorsAvailable: 2, GatearrayValue: 0x9e}
	OverscanMode0 = Size{Width: 192, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 16}
	OverscanMode1 = Size{Width: 384, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 4}
	OverscanMode2 = Size{Width: 768, Height: 272, LinesNumber: 272, ColumnsNumber: 96, ColorsAvailable: 2}
	SelfMode      = Size{}
	WidthMax      = 768
	HeightMax     = 272
)
var (
	ErrorCpcColorNotFound = errors.New("cpc color not found")
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
	Pink          = CpcColor{HardwareNumber: 7, FirmwareNumber: 16, HardwareValues: []uint8{0x47}, Color: color.RGBA{A: 0xFF, R: 0xFF, G: 0x7F, B: 0x7F}}
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

func diffColor(a, b uint32) int64 {
	if a > b {
		return int64(a - b)
	}
	return int64(b - a)
}

func sqrt(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

var DistanceMax int64 = 584970

type PaletteReducer struct {
	Cs []ColorReducer
}

func NewPaletteReducer() *PaletteReducer {
	return &PaletteReducer{Cs: make([]ColorReducer, 0)}
}

type ColorReducer struct {
	C         color.Color
	Occurence int
	Distances map[color.Color]float64
}

func NewColorReducer(c color.Color, occ int) ColorReducer {
	return ColorReducer{C: c, Occurence: occ, Distances: make(map[color.Color]float64)}
}

func (p *PaletteReducer) ComputeDistances() {
	for index, v := range p.Cs {
		for i, v2 := range p.Cs {
			if index == i {
				continue
			}
			p.Cs[index].Distances[v2.C] = ColorsDistance(v.C, v2.C)
		}
	}
}

type ByOccurence []ColorReducer

func (b ByOccurence) Len() int           { return len(b) }
func (b ByOccurence) Less(i, j int) bool { return b[i].Occurence < b[j].Occurence }
func (b ByOccurence) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (p *PaletteReducer) OccurencesSort() {
	sort.Sort(sort.Reverse(ByOccurence(p.Cs)))
}

func (pr *PaletteReducer) Reduce(nbColors int) color.Palette {
	var p color.Palette

	pr.ComputeDistances()
	pr.OccurencesSort()
	p = append(p, pr.Cs[0].C)
	for i := 1; i < len(pr.Cs); i++ {
		if len(p) < nbColors {
			previous := pr.Cs[i-1]
			current := pr.Cs[i]
			if previous.Distances[current.C] > 10. {
				p = append(p, current.C)
			}
		} else {
			break
		}
	}
	return p
}

// from website https://www.compuphase.com/cmetric.htm
func ColorsDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	rmean := int64(r1>>8+r2>>8) / 2
	r := diffColor(r1>>8, r2>>8)
	g := diffColor(g1>>8, g2>>8)
	b := diffColor(b1>>8, b2>>8)
	distance := sqrt((((512 + rmean) * r * r) >> 8) + (4 * g * g) + (((767 - rmean) * b * b) >> 8))
	//log.GetLogger().Info( "distance :%d distanceMax:%d\n", distance, DistanceMax)
	return float64(distance) / float64(DistanceMax) * 100.
}

func ColorDistance2(c1, c2 color.Color) int64 {
	if c1 == nil || c2 == nil {
		return 0
	}
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	rmean := int64(r1>>8+r2>>8) / 2
	r := int64(r1>>8 - r2>>8)
	g := int64(g1>>8 - g2>>8)
	b := int64(b1>>8 - b2>>8)
	distance := (((512+rmean)*r*r)>>8 + (4 * g * g) + (((767 - rmean) * b * b) >> 8))
	return distance
}

func CpcColorFromHardwareNumber(c int) (CpcColor, error) {
	if White.HardwareNumber == c {
		return White, nil
	}
	if SeaGreen.HardwareNumber == c {
		return SeaGreen, nil
	}
	if PastelYellow.HardwareNumber == c {
		return PastelYellow, nil
	}
	if Blue.HardwareNumber == c {
		return Blue, nil
	}
	if Purple.HardwareNumber == c {
		return Purple, nil
	}
	if Cyan.HardwareNumber == c {
		return Cyan, nil
	}
	if Pink.HardwareNumber == c {
		return Pink, nil
	}
	if BrightYellow.HardwareNumber == c {
		return BrightYellow, nil
	}
	if BrightWhite.HardwareNumber == c {
		return BrightWhite, nil
	}
	if BrightRed.HardwareNumber == c {
		return BrightRed, nil
	}
	if BrightMagenta.HardwareNumber == c {
		return BrightMagenta, nil
	}
	if Orange.HardwareNumber == c {
		return Orange, nil
	}
	if PastelMagenta.HardwareNumber == c {
		return PastelMagenta, nil
	}
	if BrightGreen.HardwareNumber == c {
		return BrightGreen, nil
	}
	if BrightCyan.HardwareNumber == c {
		return BrightCyan, nil
	}
	if Black.HardwareNumber == c {
		return Black, nil
	}
	if BrightBlue.HardwareNumber == c {
		return BrightBlue, nil
	}
	if Green.HardwareNumber == c {
		return Green, nil
	}
	if SkyBlue.HardwareNumber == c {
		return SkyBlue, nil
	}
	if Magenta.HardwareNumber == c {
		return Magenta, nil
	}
	if PastelGreen.HardwareNumber == c {
		return PastelGreen, nil
	}
	if Lime.HardwareNumber == c {
		return Lime, nil
	}
	if PastelCyan.HardwareNumber == c {
		return PastelCyan, nil
	}
	if Red.HardwareNumber == c {
		return Red, nil
	}
	if Mauve.HardwareNumber == c {
		return Mauve, nil
	}
	if Yellow.HardwareNumber == c {
		return Yellow, nil
	}
	if PastelBlue.HardwareNumber == c {
		return PastelBlue, nil
	}
	return CpcColor{HardwareNumber: -1}, ErrorCpcColorNotFound
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

func HardwareNumber(c color.Color) (int, error) {
	if ColorsAreEquals(White.Color, c) {
		return White.HardwareNumber, nil
	}
	if ColorsAreEquals(SeaGreen.Color, c) {
		return SeaGreen.HardwareNumber, nil
	}
	if ColorsAreEquals(PastelYellow.Color, c) {
		return PastelYellow.HardwareNumber, nil
	}
	if ColorsAreEquals(Blue.Color, c) {
		return Blue.HardwareNumber, nil
	}
	if ColorsAreEquals(Purple.Color, c) {
		return Purple.HardwareNumber, nil
	}
	if ColorsAreEquals(Cyan.Color, c) {
		return Cyan.HardwareNumber, nil
	}
	if ColorsAreEquals(Pink.Color, c) {
		return Pink.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightYellow.Color, c) {
		return BrightYellow.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightWhite.Color, c) {
		return BrightWhite.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightRed.Color, c) {
		return BrightRed.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightMagenta.Color, c) {
		return BrightMagenta.HardwareNumber, nil
	}
	if ColorsAreEquals(Orange.Color, c) {
		return Orange.HardwareNumber, nil
	}
	if ColorsAreEquals(PastelMagenta.Color, c) {
		return PastelMagenta.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightGreen.Color, c) {
		return BrightGreen.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightCyan.Color, c) {
		return BrightCyan.HardwareNumber, nil
	}
	if ColorsAreEquals(Black.Color, c) {
		return Black.HardwareNumber, nil
	}
	if ColorsAreEquals(BrightBlue.Color, c) {
		return BrightBlue.HardwareNumber, nil
	}
	if ColorsAreEquals(Green.Color, c) {
		return Green.HardwareNumber, nil
	}
	if ColorsAreEquals(SkyBlue.Color, c) {
		return SkyBlue.HardwareNumber, nil
	}
	if ColorsAreEquals(Magenta.Color, c) {
		return Magenta.HardwareNumber, nil
	}
	if ColorsAreEquals(PastelGreen.Color, c) {
		return PastelGreen.HardwareNumber, nil
	}
	if ColorsAreEquals(Lime.Color, c) {
		return Lime.HardwareNumber, nil
	}
	if ColorsAreEquals(PastelCyan.Color, c) {
		return PastelCyan.HardwareNumber, nil
	}
	if ColorsAreEquals(Red.Color, c) {
		return Red.HardwareNumber, nil
	}
	if ColorsAreEquals(Mauve.Color, c) {
		return Mauve.HardwareNumber, nil
	}
	if ColorsAreEquals(Yellow.Color, c) {
		return Yellow.HardwareNumber, nil
	}
	if ColorsAreEquals(PastelBlue.Color, c) {
		return PastelBlue.HardwareNumber, nil
	}
	return -1, ErrorCpcColorNotFound

}

func FlashColorQuotient(c1, c2 CpcColor) float64 {
	r0, g0, b0, _ := c1.Color.RGBA()
	r1, g1, b1, _ := c2.Color.RGBA()

	return ((CToF(r0) * 30) + (CToF(g0) * 59) + (CToF(b0) * 11)) /
		((CToF(r1) * 30) + (CToF(g1) * 59) + (CToF(b1) * 11))
}

func CToF(c uint32) float64 {
	return (float64(c) / 255. * 100)
}

type ByDistance []color.Color

func (p ByDistance) Len() int { return len(p) }

func (p ByDistance) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ByDistance) Less(i, j int) bool {
	d := ColorDistance2(p[i], p[j])
	return d > 0.
}

func SortColorsByDistance(p color.Palette) color.Palette {
	sort.Sort(sort.Reverse(ByDistance(p)))
	return p
}

type SplitRasterScreen struct {
	Values []SplitRaster
}

func NewSplitRasterScreen() *SplitRasterScreen {
	return &SplitRasterScreen{Values: make([]SplitRaster, 0)}
}

func (srs *SplitRasterScreen) Add(s SplitRaster) bool {
	if srs.IsFull() {
		return false
	}
	srs.Values = append(srs.Values, s)
	return true
}

func (srs *SplitRasterScreen) IsFull() bool {
	return len(srs.Values) >= 256
}

type SplitRaster struct {
	Offset        uint16
	Length        int
	Occurence     int
	HardwareColor []int
	PaletteIndex  []int
}

func (s *SplitRaster) Add(paletteIndex, hardwareColor int) bool {
	if len(s.PaletteIndex) >= s.Length || len(s.HardwareColor) >= s.Length {
		return false
	}
	s.PaletteIndex = append(s.PaletteIndex, paletteIndex)
	s.HardwareColor = append(s.HardwareColor, hardwareColor)
	return true
}

func NewSpliteRaster(offset uint16, length, occurence int) SplitRaster {
	return SplitRaster{
		Offset:        offset,
		Length:        length,
		Occurence:     occurence,
		PaletteIndex:  make([]int, 0),
		HardwareColor: make([]int, 0),
	}
}

func (s *SplitRaster) Boundaries() (uint16, uint16) {
	return s.Offset, s.Offset + uint16(s.Length)
}

func CpcColorStringFromHardwareNumber(c uint8) string {
	if HardwareValueAreEquals(White.HardwareValues, c) {
		return "White"
	}
	if HardwareValueAreEquals(SeaGreen.HardwareValues, c) {
		return "SeaGreen"
	}
	if HardwareValueAreEquals(PastelYellow.HardwareValues, c) {
		return "PastelYellow"
	}
	if HardwareValueAreEquals(Blue.HardwareValues, c) {
		return "Blue"
	}
	if HardwareValueAreEquals(Purple.HardwareValues, c) {
		return "Purple"
	}
	if HardwareValueAreEquals(Cyan.HardwareValues, c) {
		return "Cyan"
	}
	if HardwareValueAreEquals(Pink.HardwareValues, c) {
		return "Pink"
	}
	if HardwareValueAreEquals(BrightYellow.HardwareValues, c) {
		return "BrightYellow"
	}
	if HardwareValueAreEquals(BrightWhite.HardwareValues, c) {
		return "BrightWhite"
	}
	if HardwareValueAreEquals(BrightRed.HardwareValues, c) {
		return "BrightRed"
	}
	if HardwareValueAreEquals(BrightMagenta.HardwareValues, c) {
		return "BrightMagenta"
	}
	if HardwareValueAreEquals(Orange.HardwareValues, c) {
		return "Orange"
	}
	if HardwareValueAreEquals(PastelMagenta.HardwareValues, c) {
		return "PastelMagenta"
	}
	if HardwareValueAreEquals(BrightGreen.HardwareValues, c) {
		return "BrightGreen"
	}
	if HardwareValueAreEquals(BrightCyan.HardwareValues, c) {
		return "BrightCyan"
	}
	if HardwareValueAreEquals(Black.HardwareValues, c) {
		return "Black"
	}
	if HardwareValueAreEquals(BrightBlue.HardwareValues, c) {
		return "BrightBlue"
	}
	if HardwareValueAreEquals(Green.HardwareValues, c) {
		return "Green"
	}
	if HardwareValueAreEquals(SkyBlue.HardwareValues, c) {
		return "SkyBlue"
	}
	if HardwareValueAreEquals(Magenta.HardwareValues, c) {
		return "Magenta"
	}
	if HardwareValueAreEquals(PastelGreen.HardwareValues, c) {
		return "PastelGreen"
	}
	if HardwareValueAreEquals(Lime.HardwareValues, c) {
		return "Lime"
	}
	if HardwareValueAreEquals(PastelCyan.HardwareValues, c) {
		return "PastelCyan"
	}
	if HardwareValueAreEquals(Red.HardwareValues, c) {
		return "Red"
	}
	if HardwareValueAreEquals(Mauve.HardwareValues, c) {
		return "Mauve"
	}
	if HardwareValueAreEquals(Yellow.HardwareValues, c) {
		return "Yellow"
	}
	if HardwareValueAreEquals(PastelBlue.HardwareValues, c) {
		return "PastelBlue"
	}
	return "not defined"
}
