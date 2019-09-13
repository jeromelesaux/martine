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

func TestPixelsMode1( t *testing.T) {
	p10 := 1
	p20 := 2
	p30 := 0
	p40 := 3
	v := pixelMode1(p10, p20,p30,p40)
	t.Logf("v:%d, %8b\n", v, v)
	p1, p2, p3, p4 := rawPixelMode1(v)
	if p1 != p10 {
		t.Fatalf("expected value %d and gets %d\n", p10, p1)
	}
	if p2 != p20 {
		t.Fatalf("expected value %d and gets %d\n", p20, p2)
	}
	if p3 != p30 {
		t.Fatalf("expected value %d and gets %d\n", p30, p3)
	}
	if p4 != p40 {
		t.Fatalf("expected value %d and gets %d\n", p40, p4)
	}
}
