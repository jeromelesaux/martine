package convert

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/screenverter/gfx"
	"image"
	"image/color"
	"os"
)

func Resize(in image.Image, size gfx.Size) *image.NRGBA {
	fmt.Fprintf(os.Stdout, "* Step 1 * Resizing image to width %d pixels heigh %d\n", size.Width, size.Height)
	return imaging.Resize(in, size.Width, size.Height, imaging.Lanczos)
}

func DowngradingPalette(in *image.NRGBA, size gfx.Size)  *image.NRGBA {
	fmt.Fprintf(os.Stdout, "* Step 2 * Downgrading palette image\n")
	p,out := downgrade(in)
	fmt.Fprintf(os.Stdout,"Downgraded palette contains (%d) colors\n",len(p))
	if len(p) > size.ColorsAvailable {
		fmt.Fprintf(os.Stderr,"Downgraded palette size (%d) is greater than the available colors in this mode (%d)\n",len(p),size.ColorsAvailable)
		phasis := 1
		for {
			fmt.Fprintf(os.Stderr,"Phasis (%d) downgrade colors palette\n",phasis)
			p, out = downgrade(out)
			if len(p) <= size.ColorsAvailable {
				break
			}
			phasis++
		}
	}
	return out
}

func downgrade(in *image.NRGBA) (color.Palette,*image.NRGBA) {
	p := color.Palette{}
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			c := in.At(x, y)
			//r0,g0,b0,a0 := c.RGBA()
			cPalette := gfx.CpcOldPalette.Convert(c)
			//r,g,b,a := cPalette.RGBA()
			//fmt.Fprintf(os.Stderr,"(%d,%d) pixel R(%d)G(%d)B(%d)A(%d) => R(%d)G(%d)B(%d)A(%d)\n",x,y,r0,g0,b0,a0,r,g,b,a)
			in.Set(x,y,cPalette)
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