package gfx

import (
	"math/rand"
	"os"
	"testing"
)

func TestSaveDelta(t *testing.T) {
	d := NewDeltaCollection()
	for i := 0; i < 320; i++ {
		d.Add(byte(rand.Intn(255-1)), uint16(rand.Int()))
	}
	if err := d.Save("delta.bin"); err != nil {
		t.Fatalf("expected no error and gets %v\n", err)
	}
	os.Remove("delta.bin")
}
