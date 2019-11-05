package constants

import (
	"image/color"
	"sort"
	"testing"
)

func TestDiffcolor(t *testing.T) {
	black := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	white := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	distance := ColorsDistance(black, white)
	t.Logf("distance:%f", distance)
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	distance = ColorsDistance(black, red)
	t.Logf("distance:%f", distance)
}

func TestDiffcolor2(t *testing.T) {
	black := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	white := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	distance := ColorDistance2(black, white)
	t.Logf("distance:%d", distance)
	distance = ColorDistance2(white, black)
	t.Logf("distance:%d", distance)
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	distance = ColorDistance2(black, red)
	t.Logf("distance:%d", distance)
}

func TestSortingPalette(t *testing.T) {
	black := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	blue := color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	white := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	var p color.Palette
	p = append(p, blue)
	p = append(p, black)
	p = append(p, white)
	p = append(p, red)
	t.Logf("palette :%v\n", p)
	sort.Sort(sort.Reverse(ByDistance(p)))
	t.Logf("palette sorted :%v\n", p)

	if !ColorsAreEquals(p[0], white) {
		t.Fatalf("error expected color black and gets %v\n", p[0])
	}
}
