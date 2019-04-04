package gfx

import (
	"testing"
)

func TestCpcPlusColor(t *testing.T) {
	c1 := CpcPlusColor{G:0xa,R:0xf,B:0xf}
	v := c1.Value()
	c2 := NewRawCpcPlusColor(v)
	t.Logf("%b,%b\n",0xFF0000,v >> 8)
	t.Logf("\nc1:%s\nc2:%s\nv:%b\n",c1.ToString(),c2.ToString(),v)
}