package common

import "testing"

func TestHexParsing(t *testing.T) {
	a := "0xC000"
	v, err := ParseHexadecimal16(a)
	if err != nil {
		t.Fatal()
	}
	t.Log(v)
	a = "#4000"
	_, err = ParseHexadecimal16(a)
	if err != nil {
		t.Fatal()
	}
}
