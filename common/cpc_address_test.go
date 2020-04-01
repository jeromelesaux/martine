package common

import "testing"

func TestHexParsing(t *testing.T) {
	a := "0xC000"
	v, err := ParseHexadecimal(a)
	if err != nil {
		t.Fatal()
	}
	t.Log(v)
	a = "#4000"
	v, err = ParseHexadecimal(a)
	if err != nil {
		t.Fatal()
	}
}
