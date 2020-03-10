package file

import (
	"testing"

	"github.com/jeromelesaux/martine/constants"
)

func TestPaletteOutput(t *testing.T) {
	p := constants.CpcOldPalette
	if err := PalToPng("test.png", p[0:16]); err != nil {
		t.Fatalf("error while generating the palette with error :%v\n", err)
	}

}
