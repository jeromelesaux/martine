package gfx

import (
	"image"
	"image/color"
	"image/draw"
	"github.com/esimov/colorquant"
	"github.com/esimov/dithergo"
)

var (
	// FloydSteinberg is the Floyd Steinberg matrix
	FloydSteinberg = [][]float32{{0, 0, 7.0 / 16.0}, {3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0}}
	// JarvisJudiceNinke is the JarvisJudiceNinke matrix
	JarvisJudiceNinke = [][]float32{{0, 0, 0, 7.0 / 48.0, 5.0 / 48.0}, {3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0}, {1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0}}
	// Stucki is the Stucki matrix
	Stucki = [][]float32{{0, 0, 0, 8.0 / 42.0, 4.0 / 42.0}, {2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0}, {1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0}}
	// Atkinson is the Atkinson matrix
	Atkinson = [][]float32{{0, 0, 1.0 / 8.0, 1.0 / 8.0}, {1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0}, {0, 1.0 / 8.0, 0, 0}}
	// Burkes is the Burkes matrix
	Burkes = [][]float32{{0, 0, 0, 8.0 / 32.0, 4.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}}
	// Sierra is the Sierra matrix
	Sierra = [][]float32{{0, 0, 0, 5.0 / 32.0, 3.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}, {0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0}}
	// TwoRowSierra is a variant of the Sierrra matrix
	TwoRowSierra = [][]float32{{0, 0, 0, 4.0 / 16.0, 3.0 / 16.0}, {1.0 / 32.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 1.0 / 32.0}}
	// SierraLite is a variant of the Sierra matrix
	SierraLite = [][]float32{{0, 0, 2.0 / 4.0}, {1.0 / 4.0, 1.0 / 4.0, 0}}
	// Sierra3 
	Sierra3 = [][]float32{{ 0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0 },{ 2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0 },{ 0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0 }}
	// bayer 4
	Bayer2 = [][]float32{{0,3,2,1}}
	Bayer3 = [][]float32{{0,6,4},{7,5,1},{3,2,8}}
	Bayer4 = [][]float32{{0,12,3,15},{8,4,11,7},{2,14,1,13},{10,6,9,5}}
	Bayer8 = [][]float32{{ 0, 32,  8, 40,  2, 34, 10, 42},   /* 8x8 Bayer ordered dithering  */
    {48, 16, 56, 24, 50, 18, 58, 26},   /* pattern.  Each input pixel   */
    {12, 44,  4, 36, 14, 46,  6, 38},   /* is scaled to the 0..63 range */
    {60, 28, 52, 20, 62, 30, 54, 22},   /* before looking in this table */
    { 3, 35, 11, 43,  1, 33,  9, 41},   /* to determine the action.     */
    {51, 19, 59, 27, 49, 17, 57, 25},
    {15, 47,  7, 39, 13, 45,  5, 37},
    {63, 31, 55, 23, 61, 29, 53, 21}}
)

type DitheringType struct {
	string
}

var ( 
	OrderedDither = DitheringType{"Ordered"}
	ErrorDiffusionDither = DitheringType{"ErrorDiffusion"}
)


func Dithering(input *image.NRGBA, filter [][]float32, errorMultiplier float32) *image.NRGBA{
	dither := dither.Dither{Settings:dither.Settings{Filter:filter}}
	dst:= dither.Color(input,errorMultiplier)
	out := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{dst.Bounds().Max.X, dst.Bounds().Max.Y}})
	draw.Draw(out, out.Bounds(), dst, image.ZP, draw.Src)
	return out
}

func QuantizeNoDither(input *image.NRGBA, numColors int, pal color.Palette) *image.NRGBA {
	bounds := input.Bounds()
	img := image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), pal)
	dst := colorquant.NoDither.Quantize(input, img, numColors, false, true)
	out := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{dst.Bounds().Max.X, dst.Bounds().Max.Y}})
	draw.Draw(out, out.Bounds(), dst, image.ZP, draw.Src)
	return out
}

func QuantizeWithDither(input *image.NRGBA, filter [][]float32, numColors int, pal color.Palette) *image.NRGBA {
	dither := colorquant.Dither{Filter:filter}
	bounds := input.Bounds()
	img := image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), pal)
	dst := dither.Quantize(input,img,numColors,true,true)
	out := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{dst.Bounds().Max.X, dst.Bounds().Max.Y}})
	draw.Draw(out, out.Bounds(), dst, image.ZP, draw.Src)
	return out
}