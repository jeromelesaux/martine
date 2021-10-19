package animate

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/zx0/encode"
)

func TestAnimate(t *testing.T) {
	e := export.NewMartineContext("/Users/jeromelesaux/Downloads/bomberman.gif", "animation")
	e.Size = constants.Size{Width: 40, Height: 50, ColorsAvailable: 8}
	var screenMode uint8 = 0
	fs := []string{"/Users/jeromelesaux/Downloads/bomberman.gif"}

	Animation(fs, screenMode, e)
}

func TestDeltaMotif(t *testing.T) {
	err := DeltaMotif("/Users/jeromelesaux/Downloads/triangles.gif", &export.MartineContext{InputPath: "triangles.gif", OutputPath: "."}, 20, 0xc000, 1)
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
	fw, err := os.Create("/Users/jeromelesaux/Downloads/test.zx0")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer fw.Close()
	fw.Write(compressed)
}

func TestDisplayCode(t *testing.T) {
	fmt.Printf("%s", depackRoutine)
}
