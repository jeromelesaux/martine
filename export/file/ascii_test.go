package file

import (
	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"os"
	"testing"
)

func TestAsciiByColumn(t *testing.T) {
	data := []byte{
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
	}
	e := x.NewExportType("input.bin", "./")
	e.Size.Height = 5
	e.Size.Width = 5
	err := AsciiByColumn("test.bin", data, constants.CpcOldPalette, e)
	if err != nil {
		t.Fatalf("expected no error and gets :%v", err)
	}
	os.Remove("TESTC.TXT")
}
