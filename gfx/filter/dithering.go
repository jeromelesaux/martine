package filter

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/esimov/colorquant"
	dither "github.com/esimov/dithergo"
	"github.com/jeromelesaux/martine/export/img"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/proc"
)

var (
	// FloydSteinberg is the Floyd Steinberg matrix
	FloydSteinberg = [][]float32{
		{0, 0, 7.0 / 16.0},
		{3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0},
	}
	// JarvisJudiceNinke is the JarvisJudiceNinke matrix
	JarvisJudiceNinke = [][]float32{
		{0, 0, 0, 7.0 / 48.0, 5.0 / 48.0},
		{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
		{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
	}
	// Stucki is the Stucki matrix
	Stucki = [][]float32{
		{0, 0, 0, 8.0 / 42.0, 4.0 / 42.0},
		{2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0},
		{1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0},
	}
	// Atkinson is the Atkinson matrix
	Atkinson = [][]float32{
		{0, 0, 1.0 / 8.0, 1.0 / 8.0},
		{1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0},
		{0, 1.0 / 8.0, 0, 0},
	}
	// Burkes is the Burkes matrix
	Burkes = [][]float32{
		{0, 0, 0, 8.0 / 32.0, 4.0 / 32.0},
		{2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
	}
	// Sierra is the Sierra matrix
	Sierra = [][]float32{
		{0, 0, 0, 5.0 / 32.0, 3.0 / 32.0},
		{2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
		{0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0},
	}
	// TwoRowSierra is a variant of the Sierrra matrix
	TwoRowSierra = [][]float32{
		{0, 0, 0, 4.0 / 16.0, 3.0 / 16.0},
		{1.0 / 32.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 1.0 / 32.0},
	}
	// SierraLite is a variant of the Sierra matrix
	SierraLite = [][]float32{
		{0, 0, 2.0 / 4.0},
		{1.0 / 4.0, 1.0 / 4.0, 0},
	}
	// Sierra3
	Sierra3 = [][]float32{
		{0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0},
		{2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
		{0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0},
	}
	// bayer 4
	Bayer2 = [][]float32{
		{0, 3},
		{2, 1},
	}
	Bayer3 = [][]float32{
		{0, 6, 4},
		{7, 5, 1},
		{3, 2, 8},
	}
	Bayer4 = [][]float32{
		{0, 12, 3, 15},
		{8, 4, 11, 7},
		{2, 14, 1, 13},
		{10, 6, 9, 5},
	}
	Bayer8 = [][]float32{
		{0, 32, 8, 40, 2, 34, 10, 42},    /* 8x8 Bayer ordered dithering  */
		{48, 16, 56, 24, 50, 18, 58, 26}, /* pattern.  Each input pixel   */
		{12, 44, 4, 36, 14, 46, 6, 38},   /* is scaled to the 0..63 range */
		{60, 28, 52, 20, 62, 30, 54, 22}, /* before looking in this table */
		{3, 35, 11, 43, 1, 33, 9, 41},    /* to determine the action.     */
		{51, 19, 59, 27, 49, 17, 57, 25},
		{15, 47, 7, 39, 13, 45, 5, 37},
		{63, 31, 55, 23, 61, 29, 53, 21},
	}
)

func Dithering(input *image.NRGBA, filter [][]float32, errorMultiplier float32) *image.NRGBA {
	d := dither.Dither{Settings: dither.Settings{Filter: filter}}
	dst := d.Color(input, errorMultiplier)
	png.PngImage("test.png", dst)
	return img.Img2NRGBA(dst)
}

func QuantizeNoDither(in *image.NRGBA, numColors int, pal color.Palette) *image.NRGBA {
	bounds := in.Bounds()
	dst := colorquant.NoDither.Quantize(
		in,
		image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), pal),
		numColors,
		false,
		true)
	return img.Img2NRGBA(dst)
}

func QuantizeWithDither(input *image.NRGBA, filter [][]float32, numColors int, pal color.Palette) *image.NRGBA {
	dither := colorquant.Dither{Filter: filter}
	bounds := input.Bounds()
	out := image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), pal)
	dst := dither.Quantize(input, out, numColors, true, true)
	png.PngImage("test.png", dst)
	return img.Img2NRGBA(dst)
}

type MixingPlan struct {
	Colors [4]uint
	Ratio  float32
}

// https://bisqwit.iki.fi/story/howto/dither/jy/#Algorithms
func BayerDiphering(input *image.NRGBA, filter [][]float32, palette color.Palette) *image.NRGBA {
	image2 := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{input.Bounds().Max.X, input.Bounds().Max.Y}})
	height := image2.Bounds().Max.Y
	width := image2.Bounds().Max.X
	filterRowLenght := len(filter[0]) - 1
	filterLenght := len(filter[0]) * len(filter[0])
	log.GetLogger().Info("Palette length used in Bayer dithering %d\n", len(palette))
	pal := InitPalWithPalette(palette)

	proc.Parallel(0, height, func(yc <-chan int) {
		for y := range yc {
			for x := 0; x < width; x++ {
				temp := input.At(x, y)
				color := qColorToint(temp)
				plan := DeviseBestMixingPlan(color, pal, uint(filterLenght))
				if plan.Ratio == 4.0 { // Tri-tone or quad-tone dithering
					color := rgbToQColor(plan.Colors[((y&1)*2 + (x & 1))])
					image2.Set(x, y, color)
				} else {
					//log.GetLogger().Error("(%d,%d):(%d)(%d)\n",x,y,(x & filterRowLenght),((y & filterRowLenght) << 3))
					mapValue := filter[(x & filterRowLenght)][(y&filterRowLenght)] / float32(filterLenght)
					planIndex := 0
					if mapValue < plan.Ratio {
						planIndex = 1
					}
					color := rgbToQColor(plan.Colors[planIndex])
					image2.Set(x, y, color)
				}
			}
			//log.GetLogger().Info("Analyse done for column %d\n",y)
			log.GetLogger().Info(".")
		}
	})
	log.GetLogger().Info("\n")
	return image2
}

func InitPal() [216]uint {
	var pal [216]uint
	index := 0
	for r := 255; r >= 0; r -= 51 {
		for g := 255; g >= 0; g -= 51 {
			for b := 255; b >= 0; b -= 51 {
				rgb := ((r & 0x0ff) << 16) | ((g & 0x0ff) << 8) | (b & 0x0ff)
				pal[index] = uint(rgb)
				index++
			}
		}
	}
	return pal
}

func InitPalWithPalette(p color.Palette) []uint {
	pal := make([]uint, 0)
	for _, c := range p {
		r, g, b, _ := c.RGBA()
		rgb := ((r & 0x0ff) << 16) | ((g & 0x0ff) << 8) | (b & 0x0ff)
		pal = append(pal, uint(rgb))
	}
	return pal
}

func DeviseBestMixingPlan(color uint, pal []uint, matrixLenght uint) MixingPlan {
	r := color >> 16
	g := (color >> 8) & 0xFF
	b := color & 0xFF
	result := MixingPlan{Colors: [4]uint{0, 0, 0, 0}, Ratio: 0.5}
	var leastPenalty float64 = 1e99

	for index1 := 0; index1 < len(pal); index1++ {
		for index2 := index1; index2 < len(pal); index2++ {
			// Determine the two component colors
			color1 := pal[index1]
			color2 := pal[index2]
			r1 := color1 >> 16
			g1 := (color1 >> 8) & 0xFF
			b1 := color1 & 0xFF
			r2 := color2 >> 16
			g2 := (color2 >> 8) & 0xFF
			b2 := color2 & 0xFF
			var ratio uint = 32
			if color1 != color2 {
				// Determine the ratio of mixing for each channel.
				//   solve(r1 + ratio*(r2-r1)/64 = r, ratio)
				// Take a weighed average of these three ratios according to the
				// perceived luminosity of each channel (according to CCIR 601).
				var cr0 uint
				var cr1 uint
				if r2 != r1 {
					cr0 = 299 * uint(len(pal)) * (r - r1) / (r2 - r1)
					cr1 = 299
				}
				var cg0 uint
				var cg1 uint
				if g2 != g1 {
					cg0 = 587 * uint(len(pal)) * (g - g1) / (g2 - g1)
					cg1 = 587
				}
				var cb0 uint
				var cb1 uint
				if b1 != b2 {
					cb0 = 114 * uint(len(pal)) * (b - b1) / (b2 - b1)
					cb1 = 114
				}

				ratio = (cr0 + cg0 + cb0) / (cr1 + cg1 + cb1)

				/*ratio = ((r2 != r1 ? 299 * 64 * int(r - r1) / int(r2 - r1) : 0)
				                            + (g2 != g1 ? 587 * 64 * int(g - g1) / int(g2 - g1) : 0)
				                            + (b1 != b2 ? 114 * 64 * int(b - b1) / int(b2 - b1) : 0))
				                    / ((r2 != r1 ? 299 : 0)
				                            + (g2 != g1 ? 587 : 0)
											+ (b2 != b1 ? 114 : 0));
				*/

				if ratio > (matrixLenght - 1) {
					ratio = matrixLenght - 1
				}

			}
			// Determine what mixing them in this proportion will produce
			r0 := r1 + ratio*(r2-r1)/matrixLenght
			g0 := g1 + ratio*(g2-g1)/matrixLenght
			b0 := b1 + ratio*(b2-b1)/matrixLenght
			penalty := EvaluateMixingError(r, g, b, r0, g0, b0, r1, g1, b1, r2, g2, b2, float64(ratio)/float64(matrixLenght))
			if penalty < leastPenalty {
				leastPenalty = penalty
				result.Colors[0] = pal[index1]
				result.Colors[1] = pal[index2]
				result.Ratio = float32(ratio) / float32(matrixLenght)
			}
			if index1 != index2 {
				for index3 := 0; index3 < len(pal); index3++ {
					if index3 == index2 || index3 == index1 {
						continue
					}
					// 50% index3, 25% index2, 25% index1
					color3 := pal[index3]
					r3 := color3 >> 16
					g3 := (color3 >> 8) & 0xFF
					b3 := color3 & 0xFF
					r0 := (r1 + r2 + r3*2) / 4
					g0 := (g1 + g2 + g3*2) / 4
					b0 := (b1 + b2 + b3*2) / 4
					penalty = ColorCompare(r, g, b, r0, g0, b0) + ColorCompare(r1, g1, b1, r2, g2, b2)*0.025 + ColorCompare((r1+g1)/2, (g1+g2)/2, (b1+b2)/2, r3, g3, b3)*0.025
					if penalty < leastPenalty {
						leastPenalty = penalty
						result.Colors[0] = pal[index3] // (0,0) index3 occurs twice
						result.Colors[1] = pal[index1] // (0,1)
						result.Colors[2] = pal[index2] // (1,0)
						result.Colors[3] = pal[index3] // (1,1)
						result.Ratio = 4.0
					}
				}
			}
		}
	}
	return result
}

func EvaluateMixingError(r, g, b, r0, g0, b0, r1, g1, b1, r2, g2, b2 uint, ratio float64) float64 {
	abs := ratio - 0.5
	if abs < 0 {
		abs = -abs
	}
	return ColorCompare(r, g, b, r0, g0, b0) + ColorCompare(r1, g1, b1, r2, g2, b2)*0.1*(abs+0.5)
}

func ColorCompare(r1, g1, b1, r2, g2, b2 uint) float64 {

	luma1 := float64((r1*299 + g1*587 + b1*114) / (255.0 * 1000))
	luma2 := float64((r2*299 + g2*587 + b2*114) / (255.0 * 1000))
	lumadiff := luma1 - luma2
	diffR := float64(r1-r2) / 255.0
	diffG := float64(g1-g2) / 255.0
	diffB := float64(b1-b2) / 255.0
	return (diffR*diffR*0.299+diffG*diffG*0.587+diffB*diffB*0.114)*0.75 + lumadiff*lumadiff
}
func rgbToQColor(v uint) color.Color {
	r := v >> 16
	g := (v >> 8) & 0xFF
	b := v & 0xFF
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xFF}
}

func qColorToint(c color.Color) uint {
	var rgb uint
	r, g, b, _ := c.RGBA()
	rgb = uint(((r & 0x0ff) << 16) | ((g & 0x0ff) << 8) | (b & 0x0ff))
	return rgb
}

// Dither represent dithering algorithm implementation
type Dither struct {
	// Matrix is the error diffusion matrix
	Matrix    [][]float32
	animation chan draw.Image
	nbFrames  int
}

// NewDither prepares a dithering algorithm
func NewDither(matrix [][]float32) Dither {
	return Dither{matrix, make(chan draw.Image), 1}
}

// NewDitherAnimation prepares a dithering algorithm and animation
//
// you can retrieve every generated frames thanks to RetrieveFrame
// Note: frames are shared using an unbuffered channel
func NewDitherAnimation(matrix [][]float32, nbFrames int) Dither {
	return Dither{matrix, make(chan draw.Image), nbFrames}
}

// abs gives the absolute value of a signed integer
func abs(x int16) uint16 {
	if x < 0 {
		return uint16(-x)
	}
	return uint16(x)
}

// findColor determines the closest color in a palette given the pixel color and the error
//
// It returns the closest color, the updated error and the distance between the error and the color
func findColor(err color.Color, pix color.Color, pal color.Palette) (color.RGBA, PixelError, uint16) {
	var errR, errG, errB,
		pixR, pixG, pixB,
		colR, colG, colB int16
	_errR, _errG, _errB, _ := err.RGBA()
	_pixR, _pixG, _pixB, _ := pix.RGBA()

	// Low-pass filter
	errR = int16(float32(int16(_errR)) * 0.75)
	errG = int16(float32(int16(_errG)) * 0.75)
	errB = int16(float32(int16(_errB)) * 0.75)

	pixR = int16(uint8(_pixR)) + errR
	pixG = int16(uint8(_pixG)) + errG
	pixB = int16(uint8(_pixB)) + errB

	var index int
	var minDiff uint16 = 1<<16 - 1

	for i, col := range pal {
		_colR, _colG, _colB, _ := col.RGBA()

		colR = int16(uint8(_colR))
		colG = int16(uint8(_colG))
		colB = int16(uint8(_colB))
		var distance = abs(pixR-colR) + abs(pixG-colG) + abs(pixB-colB)

		if distance < minDiff {
			index = i
			minDiff = distance
		}
	}

	_colR, _colG, _colB, _ := pal[index].RGBA()

	colR = int16(uint8(_colR))
	colG = int16(uint8(_colG))
	colB = int16(uint8(_colB))

	return color.RGBA{uint8(colR), uint8(colG), uint8(colB), 255},
		PixelError{float32(pixR - colR),
			float32(pixG - colG),
			float32(pixB - colB),
			1<<16 - 1},
		minDiff
}

func findShift(matrix [][]float32) int {
	for _, v1 := range matrix {
		for j, v2 := range v1 {
			if v2 > 0.0 {
				return -j + 1
			}
		}
	}
	return 0
}

// Draw applies an error diffusion algorithm to the src image
func (dit Dither) Draw(dst draw.Image, rect image.Rectangle, src image.Image, sp image.Point) {
	if _, ok := dst.(*image.Paletted); !ok {
		return
	}
	p := dst.(*image.Paletted).Palette

	err := NewErrorImage(rect)
	shift := findShift(dit.Matrix)

	pixPerFrame := (rect.Dx() * rect.Dy()) / dit.nbFrames

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// using the closest color
			r, e, _ := findColor(err.PixelErrorAt(x, y), src.At(x, y), p)
			dst.Set(x, y, r)
			err.SetPixelError(x, y, e)

			if (y != 0 && x != 0) && (((y*rect.Dy())+x)%pixPerFrame == 0) {
				dit.animation <- dst
			}

			// diffusing the error using the diffusion matrix
			for i, v1 := range dit.Matrix {
				for j, v2 := range v1 {
					err.SetPixelError(x+j+shift, y+i,
						err.PixelErrorAt(x+j+shift, y+i).Add(err.PixelErrorAt(x, y).Mul(v2)))
				}
			}
		}
	}
}

// RetrieveFrame returns the next available frame
func (dit Dither) RetrieveFrame() draw.Image {
	return <-dit.animation
}

// PixelError represents the error for each canal in the image
// when dithering an image
// Errors are floats because they are the result of a division
type PixelError struct {
	// TODO(brouxco): the alpha value does not make a lot of sense in a PixelError
	R, G, B, A float32
}

// RGBA returns the errors for each canal in the image
func (c PixelError) RGBA() (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

// Add adds two PixelError
func (c PixelError) Add(c2 PixelError) PixelError {
	r := c.R + c2.R
	g := c.G + c2.G
	b := c.B + c2.B
	return PixelError{r, g, b, 0}
}

// Mul multiplies two PixelError
func (c PixelError) Mul(v float32) PixelError {
	r := c.R * v
	g := c.G * v
	b := c.B * v
	return PixelError{r, g, b, 0}
}

func pixelErrorModel(c color.Color) color.Color {
	if _, ok := c.(PixelError); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return PixelError{float32(r), float32(g), float32(b), float32(a)}
}

// ErrorImage is an in-memory image whose At method returns dithering.PixelError values
type ErrorImage struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []float32
	// Stride is the Pix stride between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Min & Max values in the image
	Min, Max PixelError
}

// ColorModel returns the ErrorImage color model
func (p *ErrorImage) ColorModel() color.Model {
	return color.ModelFunc(pixelErrorModel)
}

// Bounds returns the domain for which At can return non-zero color
func (p *ErrorImage) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y)
func (p *ErrorImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return PixelError{}
	}
	i := p.PixOffset(x, y)

	r := (p.Pix[i+0]) + float32(math.Abs(float64(p.Min.R)))/(p.Max.R-p.Min.R)*255
	g := (p.Pix[i+1]) + float32(math.Abs(float64(p.Min.G)))/(p.Max.G-p.Min.G)*255
	b := (p.Pix[i+2]) + float32(math.Abs(float64(p.Min.B)))/(p.Max.B-p.Min.B)*255

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

// PixelErrorAt returns the pixel error at (x, y)
func (p *ErrorImage) PixelErrorAt(x, y int) PixelError {
	if !(image.Point{x, y}.In(p.Rect)) {
		return PixelError{}
	}
	i := p.PixOffset(x, y)
	r := p.Pix[i+0]
	g := p.Pix[i+1]
	b := p.Pix[i+2]
	a := p.Pix[i+3]

	return PixelError{r, g, b, a}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ErrorImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set sets the error of the pixel at (x, y)
func (p *ErrorImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.ModelFunc(pixelErrorModel).Convert(c).(PixelError)
	// TODO(brouxco): use min and max functions maybe ?
	if c1.R > p.Max.R {
		p.Max.R = c1.R
	}
	if c1.G > p.Max.G {
		p.Max.G = c1.G
	}
	if c1.B > p.Max.B {
		p.Max.B = c1.B
	}
	if c1.R < p.Min.R {
		p.Min.R = c1.R
	}
	if c1.G < p.Min.G {
		p.Min.G = c1.G
	}
	if c1.B < p.Min.B {
		p.Min.B = c1.B
	}
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
	p.Pix[i+3] = c1.A
}

// SetPixelError sets the error of the pixel at (x, y)
func (p *ErrorImage) SetPixelError(x, y int, c PixelError) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	if c.R > p.Max.R {
		p.Max.R = c.R
	}
	if c.G > p.Max.G {
		p.Max.G = c.G
	}
	if c.B > p.Max.B {
		p.Max.B = c.B
	}
	if c.R < p.Min.R {
		p.Min.R = c.R
	}
	if c.G < p.Min.G {
		p.Min.G = c.G
	}
	if c.B < p.Min.B {
		p.Min.B = c.B
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = c.R
	p.Pix[i+1] = c.G
	p.Pix[i+2] = c.B
	p.Pix[i+3] = c.A
}

// NewErrorImage returns a new ErrorImage image with the given width and height
func NewErrorImage(r image.Rectangle) *ErrorImage {
	w, h := r.Dx(), r.Dy()
	buf := make([]float32, 4*w*h)
	return &ErrorImage{buf, 4 * w, r, PixelError{}, PixelError{}}
}
