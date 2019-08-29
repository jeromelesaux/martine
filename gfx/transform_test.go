package gfx

import (
	"testing"
)

func TestPixels(t *testing.T) {

	//t.Logf("74:%.b 74&2:%.8b\n",74,74&2)
	p10 := 15
	p20 := 14
	v := pixelMode0(p10, p20)
	t.Logf("v:%d\n", v)
	p1, p2 := rawPixelMode0(v)
	if p1 != p10 {
		t.Fatalf("expected value %d and gets %d\n", p10, p1)
	}
	if p2 != p20 {
		t.Fatalf("expected value %d and gets %d\n", p20, p2)
	}
}
