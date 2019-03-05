package convert

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/gfx"
	"image"
	"image/color"
	"os"
	"sort"
)

var ErrorCannotDowngradePalette = errors.New("Cannot Downgrade colors palette.")

func Resize(in image.Image, size gfx.Size, algo imaging.ResampleFilter) *image.NRGBA {
	fmt.Fprintf(os.Stdout, "* Step 1 * Resizing image to width %d pixels heigh %d\n", size.Width, size.Height)
	return imaging.Resize(in, size.Width, size.Height, algo)
}

func DowngradingPalette(in *image.NRGBA, size gfx.Size) (color.Palette, *image.NRGBA, error) {
	fmt.Fprintf(os.Stdout, "* Step 2 * Downgrading palette image\n")
	p, out := downgrade(in)
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
			if len(newPalette) >= size.ColorsAvailable {
				break
			}
			for _, s := range n[k] {
				newPalette = append(newPalette, s)
			}
		}

		fmt.Fprintf(os.Stderr, "Phasis downgrade colors palette palette (%d)\n", len(newPalette))
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
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			cPalette := p.Convert(c)
			in.Set(x, y, cPalette)
		}
	}
	return in
}

func downgrade(in *image.NRGBA) (color.Palette, *image.NRGBA) {
	p := color.Palette{}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			//r0,g0,b0,a0 := c.RGBA()
			cPalette := gfx.CpcOldPalette.Convert(c)
			//r,g,b,a := cPalette.RGBA()
			//fmt.Fprintf(os.Stderr,"(%d,%d) pixel R(%d)G(%d)B(%d)A(%d) => R(%d)G(%d)B(%d)A(%d)\n",x,y,r0,g0,b0,a0,r,g,b,a)
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
