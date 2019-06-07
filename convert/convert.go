package convert

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"image"
	"image/color"
	"os"
	"sort"
)

var ErrorCannotDowngradePalette = errors.New("Cannot Downgrade colors palette.")

func Resize(in image.Image, size constants.Size, algo imaging.ResampleFilter) *image.NRGBA {
	fmt.Fprintf(os.Stdout, "* Step 1 * Resizing image to width %d pixels heigh %d\n", size.Width, size.Height)
	return imaging.Resize(in, size.Width, size.Height, algo)
}

func DowngradingWithPalette(in *image.NRGBA, p color.Palette) (color.Palette, *image.NRGBA) {
	fmt.Fprintf(os.Stdout, "Downgrading image with input palette %d\n", len(p))
	return p, downgradeWithPalette(in, p)
}



func DowngradingPalette(in *image.NRGBA, size constants.Size, isCpcPlus bool) (color.Palette, *image.NRGBA, error) {
	fmt.Fprintf(os.Stdout, "* Step 2 * Downgrading palette image\n")
	p, out := downgrade(in, isCpcPlus)
	fmt.Fprintf(os.Stdout, "Downgraded palette contains (%d) colors\n", len(p))
	if len(p) > size.ColorsAvailable {
		fmt.Fprintf(os.Stderr, "Downgraded palette size (%d) is greater than the available colors in this mode (%d)\n", len(p), size.ColorsAvailable)
		fmt.Fprintf(os.Stderr, "Check color usage in image.\n")
		colorUsage := computePaletteUsage(out, p)
		fmt.Println(colorUsage)
		n := map[int][]color.Color{}
		var a []int
		for k, v := range colorUsage {
			n[v] = append(n[v], k)
		}
		for k := range n {
			a = append(a, k)
		}
		newPalette := []color.Color{}
		sort.Sort(sort.Reverse(sort.IntSlice(a)))
		for _, k := range a {
			for _, s := range n[k] {
				if len(newPalette) >= size.ColorsAvailable {
					break
				}
				newPalette = append(newPalette, s)
			}
		}

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
		Key color.Color
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
	for k,v := range cache {
		s = append(s, ks{Key:k,Value:v})
	}
	sort.Slice(s, func(i, j int) bool {
        return s[i].Value > s[j].Value
	})
	
	for i,v := range s {
		if i >= nbColors {
			break
		}
		p = append(p,v.Key)
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


func ConvertPalette(p color.Palette, p0 color.Palette ) color.Palette {
	var nP  []color.Color
	fmt.Fprintf(os.Stdout,"Converting palette length %d\n",len(p))
	for _,v := range p {
		n := p0.Convert(v)
		nP = append(nP,n)
	}
	return nP
}