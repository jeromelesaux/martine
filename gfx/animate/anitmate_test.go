package animate

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/zx0/encode"
)

func init() {
	_ = os.Mkdir("../../test", os.ModePerm)
}
func TestAnimate(t *testing.T) {
	file := "../../samples/sonic_rotate.gif"
	e := config.NewMartineConfig(file, "../../test")
	e.ScrCfg.Size = constants.Size{Width: 40, Height: 50, ColorsAvailable: 8}
	var screenMode uint8 = 0
	fs := []string{file}

	err := Animation(fs, screenMode, e)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeltaMotif(t *testing.T) {
	err := DeltaMotif(
		"../../samples/coke.gif",
		&config.MartineConfig{
			ScrCfg: config.ScreenConfig{
				InputPath:  "triangles.gif",
				OutputPath: "../../test",
			}},
		20,
		0xc000,
		1)
	if err != nil {
		t.Fatalf("%v", err)
	}
	files, err := filepath.Glob("*.png")
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, v := range files {
		os.Remove(v)
	}
}

func TestCompressZx0(t *testing.T) {
	f, err := os.Open("../../samples/coke.gif")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("%v", err)
	}
	compressed := encode.Encode(b)
	err = amsdos.SaveOSFile("../../test/test.zx0", compressed)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDisplayCode(t *testing.T) {
	fmt.Printf("%s", depackRoutine)
}

func TestMergeGifImages(t *testing.T) {
	fr, err := os.Open("../../samples/coke.gif")
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
	err = savePng("../../test/origin.png", origImg)
	if err != nil {
		t.Fatal(err)
	}
	previousImg := origImg

	for i := 1; i < len(gifs.Image); i++ {
		img := image.NewRGBA(imgRect)
		draw.Draw(img, previousImg.Bounds(), previousImg, previousImg.Bounds().Min, draw.Over)
		currImg := gifs.Image[i]
		draw.Draw(img, currImg.Bounds(), currImg, currImg.Bounds().Min, draw.Over)
		filename := fmt.Sprintf("../../test/origin-%.2d.png", i)
		err = savePng(filename, img)
		if err != nil {
			t.Fatal(err)
		}
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
