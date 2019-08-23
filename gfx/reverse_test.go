package gfx

import (
	"testing"
	"github.com/jeromelesaux/martine/export/file"
)

func TestScrToPng(T *testing.T) {
	p, _, _ := file.OpenPal("batman-n.pal")
	ScrToPng("batman-n.scr","batman.png",0,p)
}