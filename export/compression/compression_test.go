package compression

import (
	"bytes"
	"testing"

	"github.com/jeromelesaux/martine/log"
)

func TestToCompressMethod(t *testing.T) {
	tests := []struct {
		val  int
		want CompressionMethod
	}{
		{-1, NONE},
		{0, NONE},
		{1, RLE},
		{2, RLE16},
		{3, LZ4},
		{4, RawLZ4},
		{5, ZX0},
		{99, NONE},
	}

	for _, tt := range tests {
		got := ToCompressMethod(tt.val)
		if got != tt.want {
			t.Fatalf("ToCompressMethod(%d) = %v, want %v", tt.val, got, tt.want)
		}
	}
}

func TestCompressNone(t *testing.T) {
	log.Default("test")
	data := []byte("hello world")
	got, err := Compress(data, NONE)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("expected unchanged data, got %v", got)
	}
}
