package constants

import (
	"image/color"
	"testing"
)

func TestDiffcolor(t *testing.T) {
	black := color.NRGBA{R:255,G:255,B:255,A:255}
	white := color.NRGBA{R:0,G:0,B:0,A:255}
	distance := ColorsDistance(black, white)
	t.Logf("distance:%d", distance)
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	distance = ColorsDistance(black,red)
	t.Logf("distance:%d", distance)
}
