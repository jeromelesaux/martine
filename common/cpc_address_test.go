package common

import (
	"testing"

	"github.com/jeromelesaux/martine/log"
)

func TestHexParsing(t *testing.T) {
	log.Default("test")
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

	t.Run("DecimalAddress", func(t *testing.T) {
		a := "49152"
		v, err := ParseHexadecimal16(a)
		if err != nil {
			t.Fatal(err)
		}
		if v != 49152 {
			t.Fatalf("Expected 49152 and gets %d", v)
		}
	})

	t.Run("InvalidAddress", func(t *testing.T) {
		a := "not-a-number"
		_, err := ParseHexadecimal16(a)
		if err == nil {
			t.Fatal("expected error for invalid address")
		}
	})

	t.Run("ParseHexadecimal8", func(t *testing.T) {
		a := "0x7F"
		v, err := ParseHexadecimal8(a)
		if err != nil {
			t.Fatal(err)
		}
		if v != 0x7F {
			t.Fatalf("Expected 0x7F and gets %x", v)
		}
	})
}
