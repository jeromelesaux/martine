package animate

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
)

func init() {
	_, _ = log.InitLoggerWithFile("test.log")
}

func TestOpenGif(t *testing.T) {
	fr, err := os.Open("../../samples/coke.gif")
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
		rect := image.Rect(0, 0, v.Bounds().Max.X, v.Bounds().Max.Y)
		img := image.NewNRGBA(rect)
		draw.Draw(img, rect, v, rect.Min, draw.Over)
		if err := png.Png(fmt.Sprintf(".././../test/%d.png", i), img); err != nil {
			t.Fatal(err)
		}
	}
	os.RemoveAll("tests")
}

func TestInternalDelta(t *testing.T) {
	_ = os.Mkdir("InternalDelta", 0777)
	ex := &config.MartineConfig{
		Size:            constants.Size{Width: 100, Height: 100, ColorsAvailable: 4},
		CustomDimension: true,
		LineWidth:       0x50,
		OutputPath:      "./InternalDelta",
		OneLine:         false,
		OneRow:          false,
		FilloutGif:      false,
	}
	err := DeltaPacking("../../samples/coke.gif", ex, 0xc010, 1, DeltaExportV1)
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll("InternalDelta")
}

func TestHSPSimpleNodes(t *testing.T) {
	// optimisations :
	// - init value ld bc, #.4x
	// - verifier avec plus de 5 valeurs de byte differents
	c := transformation.NewDeltaCollection()
	c.Add(1, 0x4000)
	c.Add(1, 0x4001)
	c.Add(1, 0x4010)
	c.Add(2, 0x4015)
	c.Add(3, 0x4016)
	c.Add(4, 0x4000)
	c.Add(5, 0x4100)
	c.Add(16, 0x4000)
	c.Add(200, 0x4100)
	c.Add(254, 0x4100)
	optim := NewZ80HspNode(0, 0, true, NoneRegister, nil)
	for _, v := range c.ItemsSortByByte() {
		already := false
		reg := optim.NextRegister()
		for _, offset := range v.Offsets {
			node := NewZ80HspNode(v.Byte, offset, already, reg, nil)
			optim.SetLastNode(node)
			already = true
		}
	}
	code := optim.Code()
	t.Log(code)
}
