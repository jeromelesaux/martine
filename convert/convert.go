package convert

import (
	"image"
	"fmt"
	"os"
	"github.com/jeromelesaux/screenverter/gfx"
	"github.com/disintegration/imaging"
)

func Resize(in image.Image,size gfx.Size) *image.NRGBA {
	fmt.Fprintf(os.Stdout,"Resizing image to width %d pixels heigh %d\n",size.Width,size.Height)
	return  imaging.Resize(in,size.Width,size.Height,imaging.Lanczos)
}

