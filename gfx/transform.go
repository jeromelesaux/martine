package gfx

import (
	"bytes"
	"image"
	"image/color"
)
func Transform(in *image.NRGBA,p color.Palette, size Size, filepath string) error {
	switch size {
	case Mode0:
		return TransformMode0(in,p,filepath)
	}
	return nil
}

func TransformMode0(in *image.NRGBA, p color.Palette, filepath string) error {
	pixel := make([]byte,3)
	bw := bytes.NewBuffer()
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
				r,g,b,_ := in.At(x,y).RGBA()
				pixel[0], pixel[1], pixel[2] = byte(r), byte(g), byte(b)
			}	
	}
	return nil
}