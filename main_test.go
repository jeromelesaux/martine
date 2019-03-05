package main

import (
	"testing"
)

var (
	mask10000000 = 0xFF
	mask00000010 = 0x02
	mask4        = 0x04
)

func TestMainBit(t *testing.T) {
	a := mask4

	t.Logf("%b", a)
	a = a >> 1
	t.Logf("%b", a)

	t.Logf("%b", 6)
	t.Logf("4th :%b & %b = %b", 6, 0x0E, (6 & 8)) // 4th bit
	t.Logf("3rd :%b & %b = %b", 6, 0x0D, (6 & 4)) // 3rd bit
	t.Logf("2nd :%b & %b = %b", 6, 0x0B, (6 & 2)) // 2nd bit
	t.Logf("1st :%b & %b = %b", 6, 7, (6 & 1))    // 1st bit
}
