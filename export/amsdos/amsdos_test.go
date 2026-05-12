package amsdos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeromelesaux/martine/log"
)

func TestAmsdosFilename(t *testing.T) {
	got := AmsdosFilename("/tmp/my-file_name.txt", ".BAS")
	if got != "MYFILENA.BAS" {
		t.Fatalf("expected MYFILENA.BAS, got %s", got)
	}
}

func TestSaveOSFileAndSaveStringOSFile(t *testing.T) {
	log.Default("test")
	dir := t.TempDir()
	path1 := filepath.Join(dir, "output.bin")
	path2 := filepath.Join(dir, "output.txt")

	data := []byte{1, 2, 3, 4}
	if err := SaveOSFile(path1, data); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(path1); err != nil {
		t.Fatal(err)
	} else if info.Size() == 0 {
		t.Fatal("expected non-empty binary file")
	}

	text := "hello"
	if err := SaveStringOSFile(path2, text); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path2)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != text {
		t.Fatalf("expected %s, got %s", text, string(b))
	}
}
