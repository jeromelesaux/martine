package img

import "image"

func Img2NRGBA(in image.Image) *image.NRGBA {
	out := image.NewNRGBA(in.Bounds())
	for x := 0; x < in.Bounds().Max.X; x++ {
		for y := 0; y < in.Bounds().Max.Y; y++ {
			c := in.At(x, y)
			out.Set(x, y, c)
		}
	}
	return out
}
