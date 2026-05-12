package common

import (
	"os"
	"path/filepath"
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

func TestContainsFilepath(t *testing.T) {
	files := []string{"a.txt", "b.txt"}
	if !ContainsFilepath(files, "a.txt") {
		t.Fatal("expected a.txt to be found")
	}
	if ContainsFilepath(files, "c.txt") {
		t.Fatal("expected c.txt to be absent")
	}
}

func TestSortNameOrdering(t *testing.T) {
	first := sortName("file1.txt")
	second := sortName("file2.txt")
	if first >= second {
		t.Fatalf("expected %q < %q", first, second)
	}
	plain := sortName("file.txt")
	if plain >= first {
		t.Fatalf("expected %q < %q", plain, first)
	}
}

func TestWilcardedFiles(t *testing.T) {
	dir := t.TempDir()
	files := []string{"one.txt", "two.txt", "three.png"}
	for _, name := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result, err := WilcardedFiles([]string{filepath.Join(dir, "*.txt")})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 txt files, got %d", len(result))
	}
	if !ContainsFilepath(result, filepath.Join(dir, "one.txt")) || !ContainsFilepath(result, filepath.Join(dir, "two.txt")) {
		t.Fatalf("expected wildcard results to contain one.txt and two.txt, got %v", result)
	}
}
