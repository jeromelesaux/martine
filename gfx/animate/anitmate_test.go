package animate

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/zx0/encode"
)

func TestAnimate(t *testing.T) {
	e := export.NewMartineConfig("/Users/jeromelesaux/Downloads/bomberman.gif", "animation")
	e.Size = constants.Size{Width: 40, Height: 50, ColorsAvailable: 8}
	var screenMode uint8 = 0
	fs := []string{"/Users/jeromelesaux/Downloads/bomberman.gif"}

	Animation(fs, screenMode, e)
}

func TestDeltaMotif(t *testing.T) {
	err := DeltaMotif("/Users/jeromelesaux/Downloads/triangles.gif", &export.MartineConfig{InputPath: "triangles.gif", OutputPath: "."}, 20, 0xc000, 1)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCompressZx0(t *testing.T) {
	f, err := os.Open("/Users/jeromelesaux/Downloads/cat.scr")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("%v", err)
	}
	compressed := encode.Encode(b)
	err = amsdos.SaveOSFile("/Users/jeromelesaux/Downloads/test.zx0", compressed)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDisplayCode(t *testing.T) {
	fmt.Printf("%s", depackRoutine)
}

func TestMergeGifImages(t *testing.T) {
	fr, err := os.Open("/Users/jeromelesaux/Downloads/Files_Gif/sablier-8.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer fr.Close()
	gifs, err := gif.DecodeAll(fr)
	if err != nil {
		t.Fatal(err)
	}
	imgRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: gifs.Config.Width, Y: gifs.Config.Height}}
	origImg := image.NewRGBA(imgRect)
	draw.Draw(origImg, gifs.Image[0].Bounds(), gifs.Image[0], gifs.Image[0].Bounds().Min, 0)
	savePng("origin.png", origImg)
	previousImg := origImg

	for i := 1; i < len(gifs.Image); i++ {
		img := image.NewRGBA(imgRect)
		draw.Draw(img, previousImg.Bounds(), previousImg, previousImg.Bounds().Min, draw.Over)
		currImg := gifs.Image[i]
		draw.Draw(img, currImg.Bounds(), currImg, currImg.Bounds().Min, draw.Over)
		filename := fmt.Sprintf("origin-%.2d.png", i)
		savePng(filename, img)
		previousImg = img
	}
	t.Log(gifs)
}

func savePng(filename string, img image.Image) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := png.Encode(w, img); err != nil {
		return err
	}
	return nil
}
