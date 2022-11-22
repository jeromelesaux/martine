package image

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"sort"

	"github.com/disintegration/imaging"
	"github.com/oliamb/cutter"

	"github.com/jeromelesaux/martine/constants"
)

var ErrorCannotDowngradePalette = errors.New("cannot Downgrade colors palette")

func Resize(in image.Image, size constants.Size, algo imaging.ResampleFilter) *image.NRGBA {
	//fmt.Fprintf(os.Stdout, "* Step 1 * Resizing image to width %d pixels heigh %d\n", size.Width, size.Height)
	return imaging.Resize(in, size.Width, size.Height, algo)
}

func Reducer(in *image.NRGBA, reducer int) *image.NRGBA {
	var mask uint8

	switch reducer {
	case 1:
		mask = 8
	case 2:
		mask = 16
	case 3:
		mask = 32
	}

	fmt.Fprintf(os.Stdout, "Applying reducer mask :(%.8b)\n", mask)
	for x := 0; x < in.Bounds().Max.X; x++ {
		for y := 0; y < in.Bounds().Max.Y; y++ {
			c := in.At(x, y)
			r, g, b, a := c.RGBA()
			r2 := xorMask(r, mask)
			g2 := xorMask(g, mask)
			b2 := xorMask(b, mask)
			a2 := xorMask(a, mask)
			c2 := color.NRGBA{R: r2, G: g2, B: b2, A: a2}
			in.Set(x, y, c2)
		}
	}
	return in
}

func xorMask(v uint32, m uint8) uint8 {
	v2 := uint8(v)
	if v2 > m {
		v2 ^= m
	}
	return v2
}

func Max(v0, v1 uint32) uint32 {
	if v0 > v1 {
		return v0
	}
	return v1
}

func Min(v0, v1 uint32) uint32 {
	if v0 < v1 {
		return v0
	}
	return v1
}

func LumSaturation(c color.Color, lumi, satur float64) color.Color {
	var r, g, b float64
	r0, g0, b0, _ := c.RGBA()

	r = float64(r0 >> 8)
	g = float64(g0 >> 8)
	b = float64(b0 >> 8)
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	dif := max - min
	hue := 0.0
	if max > min {
		if max == 0 {
			g = (b-r)/dif*60. + 120.
		} else {
			if b == max {
				g = (r-g)/dif*60. + 240.
			} else {
				t := 0.
				if b > g {
					t = 360.
				}
				g = (g-b)/dif*60. + t
			}
		}
		hue = g
		if hue < 0 {
			hue += 360.
		}
		hue *= 255. / 360.
	}
	sat := satur * (dif / max) * 255.
	bri := lumi * max
	r = bri
	g = bri
	b = bri

	if sat != 0 {
		max = bri
		dif = bri * sat / 255.
		min = bri - dif
		h := hue * 360. / 255.
		if h < 60. {
			r = max
			g = h*dif/60. + min
			b = min
		} else {
			if h < 120. {
				r = -(h-120.)*dif/60. + min
				g = max
				b = min
			} else {
				if h < 180. {
					r = min
					g = max
					b = (h-120.)*dif/60. + min
				} else {
					if h < 240. {
						r = min
						g = -(h-240.)*dif/60. + min
						b = max
					} else {
						if h < 300. {
							r = (h-240.)*dif/60. + min
							g = min
							b = max
						} else {
							if h <= 360. {
								r = max
								g = min
								b = -(h-360.)*dif/60 + min
							} else {
								r = 0
								g = 0
								b = 0

							}
						}
					}
				}
			}
		}
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: math.MaxUint8}
}

func EnhanceBrightness(p color.Palette, saturation, brightness float64) color.Palette {
	for i := 0; i < len(p); i++ {
		p[i] = LumSaturation(p[i], brightness, saturation)
	}
	return p
}

func MonochromePalette(p color.Palette) color.Palette {
	var newPalette color.Palette
	for _, c := range p {
		r, g, b, _ := c.RGBA()

		lum := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
		newPalette = append(newPalette, color.Gray{uint8(lum)})
	}
	return newPalette
}

func ColorMonochromePalette(co color.Color, p color.Palette) color.Palette {
	var newPalette color.Palette
	r0, g0, b0, a0 := co.RGBA()
	p = MonochromePalette(p)
	for _, s := range p {
		r, g, b, _ := s.RGBA()
		nc := color.NRGBA{R: uint8((r + r0) >> 8), G: uint8((g + g0) >> 8), B: uint8((b + b0) >> 8), A: uint8(a0 >> 8)}
		newPalette = append(newPalette, nc)
	}
	return newPalette
}

func DowngradingWithPalette(in *image.NRGBA, p color.Palette) (color.Palette, *image.NRGBA) {
	//	fmt.Fprintf(os.Stdout, "Downgrading image with input palette %d\n", len(p))
	return p, downgradeWithPalette(in, p)
}

func DowngradingPalette(in *image.NRGBA, size constants.Size, isCpcPlus bool) (color.Palette, *image.NRGBA, error) {
	//	fmt.Fprintf(os.Stdout, "* Step 2 * Downgrading palette image\n")
	p, out := downgrade(in, isCpcPlus)
	//	fmt.Fprintf(os.Stdout, "Downgraded palette contains (%d) colors\n", len(p))
	if len(p) > size.ColorsAvailable {
		fmt.Fprintf(os.Stderr, "Downgraded palette size (%d) is greater than the available colors in this mode (%d)\n", len(p), size.ColorsAvailable)
		fmt.Fprintf(os.Stderr, "Check color usage in image.\n")
		colorUsage := computePaletteUsage(out, p)
		//fmt.Println(colorUsage)
		// feed sort palette colors structure
		paletteToReduce := constants.NewPaletteReducer()

		for c, v := range colorUsage {
			paletteToReduce.Cs = append(paletteToReduce.Cs, constants.NewColorReducer(c, v))
		}
		// launch analyse
		newPalette := paletteToReduce.Reduce(size.ColorsAvailable)
		fmt.Fprintf(os.Stdout, "Phasis downgrade colors palette palette (%d)\n", len(newPalette))
		return newPalette, downgradeWithPalette(out, newPalette), nil

	}
	return p, out, nil
}

func computePaletteUsage(in *image.NRGBA, p color.Palette) map[color.Color]int {
	usage := make(map[color.Color]int, 0)
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			usage[c]++
		}
	}
	return usage
}

func downgradeWithPalette(in *image.NRGBA, p color.Palette) *image.NRGBA {
	cache := make(map[color.Color]color.Color, 0)
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			if cc := cache[c]; cc != nil {
				in.Set(x, y, cc)
			} else {
				cPalette := p.Convert(c)
				in.Set(x, y, cPalette)
				cache[c] = cPalette
			}
		}
	}
	return in
}

func ExtractPalette(in *image.NRGBA, isCpcPlus bool, nbColors int) color.Palette {
	p := []color.Color{}
	type ks struct {
		Key   color.Color
		Value int
	}
	cache := make(map[color.Color]int, 0)
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			var cPalette color.Color
			if cc := cache[c]; cc != 0 {
				cache[c]++
			} else {
				if isCpcPlus {
					cPalette = constants.CpcPlusPalette.Convert(c)
				} else {
					cPalette = constants.CpcOldPalette.Convert(c)
				}
				cache[cPalette]++
			}
			in.Set(x, y, cPalette)
		}
	}

	var s []ks
	for k, v := range cache {
		s = append(s, ks{Key: k, Value: v})
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].Value > s[j].Value
	})

	for i, v := range s {
		if i >= nbColors {
			break
		}
		p = append(p, v.Key)
	}
	return p
}

func PaletteUsed(in *image.NRGBA, isCpcPlus bool) color.Palette {
	fmt.Fprintf(os.Stdout, "Define the Palette use in image.\n")
	cache := make(map[color.Color]color.Color, 0)
	p := color.Palette{}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			var cPalette color.Color
			if cc := cache[c]; cc != nil {
				cPalette = cc
			} else {
				if isCpcPlus {
					cPalette = constants.CpcPlusPalette.Convert(c)
				} else {
					cPalette = constants.CpcOldPalette.Convert(c)
				}
				cache[c] = cPalette
			}
			in.Set(x, y, cPalette)
			if !paletteContains(p, cPalette) {
				p = append(p, cPalette)
			}
		}
	}
	return p
}

func downgrade(in *image.NRGBA, isCpcPlus bool) (color.Palette, *image.NRGBA) {
	fmt.Fprintf(os.Stdout, "Plus palette :%d\n", len(constants.CpcPlusPalette))
	cache := make(map[color.Color]color.Color, 0)
	p := color.Palette{}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			var cPalette color.Color
			if cc := cache[c]; cc != nil {
				cPalette = cc
			} else {
				if isCpcPlus {
					cPalette = constants.CpcPlusPalette.Convert(c)
				} else {
					cPalette = constants.CpcOldPalette.Convert(c)
				}
				cache[c] = cPalette
			}
			in.Set(x, y, cPalette)
			if !paletteContains(p, cPalette) {
				p = append(p, cPalette)
			}
		}
	}
	return p, in
}

func paletteContains(p color.Palette, c color.Color) bool {
	for _, cp := range p {
		if cp == c {
			return true
		}
	}
	return false
}

func ConvertPalette(p color.Palette, p0 color.Palette) color.Palette {
	var nP []color.Color
	fmt.Fprintf(os.Stdout, "Converting palette length %d\n", len(p))
	for _, v := range p {
		n := p0.Convert(v)
		nP = append(nP, n)
	}
	return nP
}

func ToCPCPalette(p color.Palette, cpcPalette color.Palette) color.Palette {
	var out color.Palette
	for _, v := range p {
		c := cpcPalette.Convert(v)
		out = append(out, c)
	}
	return out
}

func SplitImage(i image.Image) (*image.NRGBA, *image.NRGBA, error) {
	height := i.Bounds().Max.Y / 2
	width := i.Bounds().Max.X
	image1, err := cutter.Crop(i, cutter.Config{
		Width:  width,
		Height: height,
	})
	if err != nil {
		return nil, nil, err
	}
	image2, err := cutter.Crop(i, cutter.Config{
		Width:  width,
		Height: height,
		Anchor: image.Point{X: 0, Y: height},
	})

	b := image1.Bounds()
	raw1 := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(raw1, raw1.Bounds(), image1, b.Min, draw.Src)

	b = image2.Bounds()
	raw2 := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(raw2, raw2.Bounds(), image2, b.Min, draw.Src)

	f, _ := os.Create("imageup.png")
	png.Encode(f, image1)
	f.Close()

	f2, _ := os.Create("imagedown.png")
	png.Encode(f2, image2)
	f2.Close()
	return raw1, raw2, err
}
