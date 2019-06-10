package gfx

import (
	"fmt"
	"github.com/esimov/colorquant"
	"github.com/esimov/dithergo"
	"github.com/jeromelesaux/martine/proc"
	"image"
	"image/color"
	"image/draw"
	"os"
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
	Sierra3 = [][]float32{{0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}, {0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0}}
	// bayer 4
	Bayer2 = [][]float32{{0, 3}, {2, 1}}
	Bayer3 = [][]float32{{0, 6, 4}, {7, 5, 1}, {3, 2, 8}}
	Bayer4 = [][]float32{{0, 12, 3, 15}, {8, 4, 11, 7}, {2, 14, 1, 13}, {10, 6, 9, 5}}
	Bayer8 = [][]float32{{0, 32, 8, 40, 2, 34, 10, 42}, /* 8x8 Bayer ordered dithering  */
		{48, 16, 56, 24, 50, 18, 58, 26}, /* pattern.  Each input pixel   */
		{12, 44, 4, 36, 14, 46, 6, 38},   /* is scaled to the 0..63 range */
		{60, 28, 52, 20, 62, 30, 54, 22}, /* before looking in this table */
		{3, 35, 11, 43, 1, 33, 9, 41},    /* to determine the action.     */
		{51, 19, 59, 27, 49, 17, 57, 25},
		{15, 47, 7, 39, 13, 45, 5, 37},
		{63, 31, 55, 23, 61, 29, 53, 21}}
)

type DitheringType struct {
	string
}

var (
	OrderedDither        = DitheringType{"Ordered"}
	ErrorDiffusionDither = DitheringType{"ErrorDiffusion"}
)

func Dithering(input *image.NRGBA, filter [][]float32, errorMultiplier float32) *image.NRGBA {
	dither := dither.Dither{Settings: dither.Settings{Filter: filter}}
	dst := dither.Color(input, errorMultiplier)
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
	dither := colorquant.Dither{Filter: filter}
	bounds := input.Bounds()
	img := image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), pal)
	dst := dither.Quantize(input, img, numColors, true, true)
	out := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{dst.Bounds().Max.X, dst.Bounds().Max.Y}})
	draw.Draw(out, out.Bounds(), dst, image.ZP, draw.Src)
	return out
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
	fmt.Fprintf(os.Stdout, "Palette lenght used in Bayer dithering %d\n", len(palette))
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
					//fmt.Fprintf(os.Stderr,"(%d,%d):(%d)(%d)\n",x,y,(x & filterRowLenght),((y & filterRowLenght) << 3))
					mapValue := filter[(x & filterRowLenght)][(y&filterRowLenght)] / float32(filterLenght)
					planIndex := 0
					if mapValue < plan.Ratio {
						planIndex = 1
					}
					color := rgbToQColor(plan.Colors[planIndex])
					image2.Set(x, y, color)
				}
			}
			//fmt.Fprintf(os.Stdout,"Analyse done for column %d\n",y)
			fmt.Fprintf(os.Stdout, ".")
		}
	})
	fmt.Fprintf(os.Stdout, "\n")
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
				if ratio < 0 {
					ratio = 0
				} else {
					if ratio > (matrixLenght - 1) {
						ratio = matrixLenght - 1
					}
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
