package file

import (
	"testing"

	"github.com/jeromelesaux/martine/constants"
)

func TestPaletteOutput(t *testing.T) {
	p := constants.CpcOldPalette
	PalToPng("test.png", p[0:16])

}
