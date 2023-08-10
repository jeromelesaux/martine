package transformation_test

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

func TestBoardSprite(t *testing.T) {
	f, err := os.Open("../../samples/mario-level1.png")
	if err != nil {
		t.Fatalf("Cannot open file error %v\n", err)
	}
	defer f.Close()
	im, err := png.Decode(f)
	if err != nil {
		t.Fatalf("Cannot decode png file error :%v\n", err)
	}
	a := transformation.AnalyzeTilesBoard(im, constants.Size{Width: 16, Height: 16})
	t.Log(a.String())
	fmt.Println(a.String())
	err = a.SaveSchema("../../test/alexkidd_board.png")
	if err != nil {
		t.Fatal(err)
	}
	err = a.SaveTilemap("../../test/alexkidd.map")
	if err != nil {
		t.Fatal(err)
	}
}
