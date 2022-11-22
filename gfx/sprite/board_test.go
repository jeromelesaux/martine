package sprite_test

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/gfx/sprite"
)

func TestSpriteBoard(t *testing.T) {
	fr, err := os.Open("/Users/jls/Desktop/sprites_sonic.png")
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	defer fr.Close()
	im, err := png.Decode(fr)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	img := image.NewNRGBA(image.Rect(0, 0, im.Bounds().Max.X, im.Bounds().Max.Y))
	draw.Draw(img, img.Bounds(), im, im.Bounds().Min, draw.Src)
	p, _, err := convert.DowngradingPalette(img, constants.Size{ColorsAvailable: 16, Width: img.Bounds().Max.X, Height: img.Bounds().Max.Y}, true)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	_, s, err := sprite.SplitBoardToSprite(img, p, 3, 9, 0, false, constants.Size{Width: 32, Height: 64})
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	for i, v := range s {
		for j, v1 := range v {
			filename := fmt.Sprintf("%d-%d.png", i, j)
			fw, err := os.Create(filename)
			if err != nil {
				t.Fatalf("%v", err.Error())
			}
			png.Encode(fw, v1)
		}
	}
}
