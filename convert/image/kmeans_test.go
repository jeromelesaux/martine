package image_test

import (
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
	"testing"

	"github.com/mccutchen/palettor"
	"github.com/stretchr/testify/require"
)

func TestKmeans(t *testing.T) {
	fr, err := os.Open("../../samples/lena-512.png")
	require.NoError(t, err)
	defer fr.Close()
	img, _, err := image.Decode(fr)
	require.NoError(t, err)

	newPalette, err := palettor.Extract(16, 100, img)
	require.NoError(t, err)
	var p color.Palette

	for _, c := range newPalette.Colors() {
		p = append(p, c)
	}

	newImg := image.NewNRGBA(img.Bounds())
	for x := 0; x <= img.Bounds().Max.X; x++ {
		for y := 0; y <= img.Bounds().Max.Y; y++ {
			c := img.At(x, y)
			nc := p.Convert(c)
			newImg.Set(x, y, nc)
		}
	}

	fw, err := os.Create("lena.png")
	require.NoError(t, err)
	err = png.Encode(fw, newImg)
	require.NoError(t, err)
	os.Remove("lena.png")
}
