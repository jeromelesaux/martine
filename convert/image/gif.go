package image

import (
	"image"
	"image/draw"
	"image/gif"
)

func GifToImages(g gif.GIF) []image.Image {
	c := make([]image.Image, 0)
	width := g.Image[0].Bounds().Max.X
	height := g.Image[0].Bounds().Max.Y
	reference := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(reference, reference.Bounds(), g.Image[0], image.Point{0, 0}, draw.Src)
	for i := 1; i < len(g.Image)-1; i++ {
		in := g.Image[i]
		draw.Draw(reference, reference.Bounds(), in, image.Point{0, 0}, draw.Over)
		img := image.NewNRGBA(image.Rect(0, 0, width, height))
		draw.Draw(img, img.Bounds(), reference, image.Point{0, 0}, draw.Over)
		c = append(c, img)
	}
	return c
}
