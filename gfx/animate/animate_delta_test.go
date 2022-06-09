package animate

import (
	"fmt"
	"image/gif"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
)

func TestOpenGif(t *testing.T) {
	os.Mkdir("tests", os.ModePerm)
	fr, err := os.Open("../../images/triangles.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer fr.Close()

	g, err := gif.DecodeAll(fr)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range g.Image {
		t.Logf("%v", v.Bounds())
		if err := file.Png(fmt.Sprintf("tests/%d.png", i), v); err != nil {
			t.Fatal(err)
		}
	}
	os.RemoveAll("tests")
}

func TestInternalDelta(t *testing.T) {
	ex := &export.MartineContext{
		Size:            constants.Size{Width: 100, Height: 100, ColorsAvailable: 4},
		CustomDimension: true,
		LineWidth:       0x50,
		OutputPath:      "./",
		OneLine:         false,
		OneRow:          false,
		FilloutGif:      false,
	}
	err := DeltaPacking("/Users/jeromelesaux/Downloads/cigarette-femme.gif", ex, 0xc010, 1, DeltaExportV1)
	if err != nil {
		t.Fatal(err)
	}

}
