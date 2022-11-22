package transformation_test

import (
	"fmt"
	"image/color"
	"image/jpeg"
	"os"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

func TestMotifs(t *testing.T) {
	fr, err := os.Open("../images/ww.jpg")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer fr.Close()
	in, err := jpeg.Decode(fr)
	if err != nil {
		t.Fatalf("%v", err)
	}
	var p color.Palette = constants.CpcOldPalette
	out := convert.Resize(in, constants.Size{Width: 320, Height: 200}, imaging.NearestNeighbor)
	_, out = convert.DowngradingWithPalette(out, p)
	png.Png("../test/motifs/orig.png", out)
	a := transformation.AnalyzeTilesBoard(out, constants.Size{Width: 4, Height: 4})
	threshold := 27
	board := a.ReduceTilesNumber(float64(threshold))
	fmt.Printf("number sprites inital [%d] [%d] with threshold :%d\n", len(a.BoardTiles), len(board), threshold)
	//a.SaveBoardTile("../test/motifs/", board)
	a.Image("../test/motifs/new.png", board, a.ImageSize)
}
