package lz4

import (
	"bytes"
	"io"
	"testing"

	"github.com/pierrec/lz4"
)

func TestEncodeDecode(t *testing.T) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	compressed, err := Encode(nil, input)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("expected compressed data to be non-empty")
	}
	if bytes.Equal(compressed, input) {
		t.Fatal("expected compressed data to differ from the input")
	}

	r := bytes.NewReader(compressed)
	zr := lz4.NewReader(r)
	decoded, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if !bytes.Equal(decoded, input) {
		t.Fatalf("decoded output mismatch, got %q, want %q", decoded, input)
	}
}

func TestEncodeEmptyInput(t *testing.T) {
	compressed, err := Encode(nil, []byte{})
	if err != nil {
		t.Fatalf("Encode failed for empty input: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("expected encoded output for empty input to be non-empty")
	}
}
