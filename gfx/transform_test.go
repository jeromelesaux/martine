package gfx

import (
	"testing"

	"github.com/jeromelesaux/martine/gfx/common"
)

func TestPixels(t *testing.T) {

	//t.Logf("74:%.b 74&2:%.8b\n",74,74&2)
	p10 := 15
	p20 := 14
	v := common.PixelMode0(p10, p20)
	t.Logf("v:%d\n", v)
	p1, p2 := common.RawPixelMode0(v)
	if p1 != p10 {
		t.Fatalf("expected value %d and gets %d\n", p10, p1)
	}
	if p2 != p20 {
		t.Fatalf("expected value %d and gets %d\n", p20, p2)
	}
}

func TestPixelsMode1(t *testing.T) {
	p10 := 1
	p20 := 2
	p30 := 0
	p40 := 3
	v := common.PixelMode1(p10, p20, p30, p40)
	t.Logf("v:%d, %8b\n", v, v)
	p1, p2, p3, p4 := common.RawPixelMode1(v)
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
func TestPixelsMode2(t *testing.T) {
	p10 := 1
	p20 := 0
	p30 := 1
	p40 := 0
	p50 := 1
	p60 := 0
	p70 := 1
	p80 := 0
	v := common.PixelMode2(p10, p20, p30, p40, p50, p60, p70, p80)
	t.Logf("v:%d, %8b\n", v, v)
	p1, p2, p3, p4, p5, p6, p7, p8 := common.RawPixelMode2(v)
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
	if p5 != p50 {
		t.Fatalf("expected value %d and gets %d\n", p50, p5)
	}
	if p6 != p60 {
		t.Fatalf("expected value %d and gets %d\n", p60, p6)
	}
	if p7 != p70 {
		t.Fatalf("expected value %d and gets %d\n", p70, p7)
	}
	if p8 != p80 {
		t.Fatalf("expected value %d and gets %d\n", p80, p8)
	}
}

func TestMaskChoice(t *testing.T) {
	a := 0xAA
	t.Logf("mask:%b\n", a)
	b := 0x55
	t.Logf("%b\n", b)
	t.Logf("%b\n", (b & a))
	c := 1
	t.Logf("%b\n", (c & a))
	d := 2
	t.Logf("%b\n", (d & a))
	if (a & d) != 2 {
		t.Fatalf("expected value %b and gets %b, comparison with %b\n", 2, (a & d), a)
	}
}
