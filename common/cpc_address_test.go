package common

import "testing"

func TestHexParsing(t *testing.T) {
	t.Run("CStandard", func(t *testing.T) {
		a := "0xC000"
		v, err := ParseHexadecimal16(a)
		if err != nil {
			t.Fatal()
		}
		if v != 0xC000 {
			t.Fatalf("Expected 0xc000 and gets %x", v)
		}
	})

	t.Run("RasmStandard", func(t *testing.T) {
		a := "#4000"
		v, err := ParseHexadecimal16(a)
		if err != nil {
			t.Fatal()
		}
		if v != 0x4000 {
			t.Fatalf("Expected 0x4000 and gets %x", v)
		}
	})

}
