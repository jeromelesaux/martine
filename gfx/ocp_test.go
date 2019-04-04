package gfx

import (
	"testing"
)

func TestCpcPlusColor(t *testing.T) {
	c1 := CpcPlusColor{G:15,R:200,B:6}
	v := c1.Value()
	c2 := NewRawCpcPlusColor(v)
	t.Logf("\nc1:%s\nc2:%s\nv:%b\n",c1.ToString(),c2.ToString(),v)
}