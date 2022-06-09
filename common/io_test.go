package common

import (
	"testing"
)

func TestStructToBytes(t *testing.T) {
	type s struct {
		A uint8
		B [12]uint8
		C uint16
	}
	p := s{A: 12, B: [12]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, C: 0x2000}
	v, err := StructToBytes(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 15 {
		t.Fatalf("expects length 15 and gets  %d\n", len(v))
	}
	t.Logf("%v\n", v)
}
