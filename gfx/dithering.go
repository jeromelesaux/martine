package gfx

import (
	"image"
	"image/color"
	"github.com/esimov/colorquant"
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
	Bayer4 = [][]float32{{0,3,2,1}}
	Bayer9 = [][]float32{{0,6,4},{7,5,1},{3,2,8}}
	Bayer16 = [][]float32{{0,12,3,15},{8,4,11,7},{2,14,1,13},{10,6,9,5}}
	Bayer64 = [][]float32{{ 0, 32,  8, 40,  2, 34, 10, 42},   /* 8x8 Bayer ordered dithering  */
    {48, 16, 56, 24, 50, 18, 58, 26},   /* pattern.  Each input pixel   */
    {12, 44,  4, 36, 14, 46,  6, 38},   /* is scaled to the 0..63 range */
    {60, 28, 52, 20, 62, 30, 54, 22},   /* before looking in this table */
    { 3, 35, 11, 43,  1, 33,  9, 41},   /* to determine the action.     */
    {51, 19, 59, 27, 49, 17, 57, 25},
    {15, 47,  7, 39, 13, 45,  5, 37},
    {63, 31, 55, 23, 61, 29, 53, 21}}
)

func Dithering(input *image.NRGBA, filter [][]float32, errorMultiplier float32) *image.NRGBA {
	bounds := input.Bounds()
	img := image.NewNRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			img.Set(x, y, pixel)
		}
	}
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()

	// Prepopulate multidimensional slices
	redErrors   := make([][]float32, dx)
	greenErrors := make([][]float32, dx)
	blueErrors  := make([][]float32, dx)
	for x := 0; x < dx; x++ {
		redErrors[x]	= make([]float32, dy)
		greenErrors[x]	= make([]float32, dy)
		blueErrors[x]	= make([]float32, dy)
	/*	for y := 0; y < dy; y++ {
			redErrors[x][y]   = 0
			greenErrors[x][y] = 0
			blueErrors[x][y]  = 0
		}*/
	}

	var qrr, qrg, qrb float32
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			r32, g32, b32, a := img.At(x, y).RGBA()
			r, g, b := float32(uint8(r32)), float32(uint8(g32)), float32(uint8(b32))
			r -= redErrors[x][y] * errorMultiplier
			g -= greenErrors[x][y] * errorMultiplier
			b -= blueErrors[x][y] * errorMultiplier

			// Diffuse the error of each calculation to the neighboring pixels
			if r < 128 {
				qrr = -r
				r = 0
			} else {
				qrr = 255 - r
				r = 255
			}
			if g < 128 {
				qrg = -g
				g = 0
			} else {
				qrg = 255 - g
				g = 255
			}
			if b < 128 {
				qrb = -b
				b = 0
			} else {
				qrb = 255 - b
				b = 255
			}
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

			// Diffuse error in two dimension
			ydim := len(filter) - 1
			xdim := len(filter[0]) / 2
			for xx := 0; xx < ydim + 1; xx++ {
				for yy := -xdim; yy <= xdim - 1; yy++ {
					if y + yy < 0 || dy <= y + yy || x + xx < 0 || dx <= x + xx {
						continue
					}
					// Adds the error of the previous pixel to the current pixel
					redErrors[x+xx][y+yy] 	+= qrr * filter[xx][yy + ydim]
					greenErrors[x+xx][y+yy] += qrg * filter[xx][yy + ydim]
					blueErrors[x+xx][y+yy] 	+= qrb * filter[xx][yy + ydim]
				}
			}
		}
	}
   return img
}



func DitheringColorquant(input *image.NRGBA, filter [][]float32, numColors int) *image.NRGBA {
	dither := colorquant.Dither{Filter:filter}
	bounds := input.Bounds()
	img := image.NewNRGBA(bounds)
	dst := dither.Quantize(input,img,numColors,true,true)
	return dst.(*image.NRGBA)
}