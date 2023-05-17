package transformation_test

import (
	"fmt"
	"image/color"
	"image/jpeg"
	"os"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
)

func init() {
	log.InitLoggerWithFile("test.log")
}

func TestMotifs(t *testing.T) {
	err := os.MkdirAll("../test/motifs", 0700)
	if err != nil {
		t.Fatal(err)
	}
	fr, err := os.Open("../../samples/Batman-Neal-Adams.jpg")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer fr.Close()
	in, err := jpeg.Decode(fr)
	if err != nil {
		t.Fatalf("%v", err)
	}
	var p color.Palette = constants.CpcOldPalette
	out := image.Resize(in, constants.Size{Width: 320, Height: 200}, imaging.NearestNeighbor)
	_, out = image.DowngradingWithPalette(out, p)
	err = png.Png("../test/motifs/orig.png", out)
	if err != nil {
		t.Fatal(err)
	}
	a := transformation.AnalyzeTilesBoard(out, constants.Size{Width: 4, Height: 4})
	threshold := 27
	board := a.ReduceTilesNumber(float64(threshold))
	fmt.Printf("number sprites inital [%d] [%d] with threshold :%d\n", len(a.BoardTiles), len(board), threshold)
	// a.SaveBoardTile("../test/motifs/", board)
	err = a.Image("../test/motifs/new.png", board, a.ImageSize)
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll("../test")
}
