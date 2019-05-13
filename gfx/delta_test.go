package gfx

import (
	"math/rand"
	"os"
	"testing"
)

func TestSaveDelta(t *testing.T) {
	d := NewDeltaCollection()
	for i := 0; i < 320; i++ {
		d.Add(240, uint16(rand.Int()))
	}
	if err := d.Save("toto.bin"); err != nil {
		t.Fatalf("expected no error and gets %v\n", err)
	}
	os.Remove("toto.bin")
}
