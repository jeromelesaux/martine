package ascii_test

import (
	"os"
	"testing"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/log"
)

func init() {
	log.InitLoggerWithFile("test.log")
}

func TestAsciiByColumn(t *testing.T) {
	log.InitLoggerWithFile("test.log")
	data := []byte{
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x1, 0x2, 0x3, 0x4, 0x5,
	}
	e := config.NewMartineConfig("input.bin", "./")
	e.Size.Height = 5
	e.Size.Width = 5
	err := ascii.AsciiByColumn("test.bin", data, constants.CpcOldPalette, true, 1, e)
	if err != nil {
		t.Fatalf("expected no error and gets :%v", err)
	}
	os.Remove("TESTC.TXT")
}
