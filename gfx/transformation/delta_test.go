package transformation

import (
	"fmt"
	"os"
	"testing"
)

func TestSaveDelta(t *testing.T) {
	d := NewDeltaCollection()
	for i := 0; i < 320; i++ {
		d.Add(0xFF, uint16(i))
	}
	if err := d.Save("delta.bin"); err != nil {
		t.Fatalf("expected no error and gets %v\n", err)
	}
	filesize := 4 + (320 * 2)

	fi, err := os.Lstat("delta.bin")
	if err != nil {
		t.Fatalf("expected no error while getting informations gets :%v\n", err)
	}

	if fi.Size() != int64(filesize) {
		t.Fatalf("expected %d length and gets %d\n", filesize, fi.Size())
	}
	//os.Remove("delta.bin")
}

func TestXandY(t *testing.T) {

	for i := 0; i < 0x4000; i++ {
		y := Y(uint16(i), 0)
		x := X(uint16(i), 0)
		addr := DeltaAddress(int(x), int(y), 0)
		if addr != i {
			t.Fatalf("expected #%.4x and gets #%.4x for x:%d y:%d\n", i, addr, x, y)
		}
	}

	add := 0xe010 //- 0xC000
	x, y, err := CpcCoordinates(0xe010, 0xC000, 0)
	if err != nil {
		t.Fatal()
	}
	res := DeltaAddress(x, y, 0) + 0xC000
	a := DeltaAddress(16, 3, 0) + 0xC000
	fmt.Println(a)
	if res != add {
		t.Fatalf("does not match")
	}
}
