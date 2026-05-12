package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeromelesaux/martine/log"
)

func TestCheckOutputCreatesDirectory(t *testing.T) {
	log.Default("test")
	dir := t.TempDir()
	path := filepath.Join(dir, "newdir")
	if err := CheckOutput(path); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(path); err != nil {
		t.Fatal(err)
	} else if !info.IsDir() {
		t.Fatal("expected path to be a directory")
	}
}

func TestCheckOutputFailsOnFile(t *testing.T) {
	log.Default("test")
	tempFile := filepath.Join(t.TempDir(), "notadir")
	if err := os.WriteFile(tempFile, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CheckOutput(tempFile); err != ErrorIsNotDirectory {
		t.Fatalf("expected ErrorIsNotDirectory, got %v", err)
	}
}
