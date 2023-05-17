package png_test

import (
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/log"
)

func TestPaletteOutput(t *testing.T) {
	log.InitLoggerWithFile("test.log")
	p := constants.CpcOldPalette
	if err := png.PalToPng("test.png", p[0:16]); err != nil {
		t.Fatalf("error while generating the palette with error :%v\n", err)
	}

}
