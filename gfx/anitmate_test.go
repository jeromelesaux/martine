package gfx

import (
	"testing"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
)

func TestAnimate(t *testing.T) {
	e := export.NewExportType("/Users/jeromelesaux/Downloads/bomberman.gif", "animation")
	e.Size = constants.Size{Width: 40, Height: 50, ColorsAvailable: 8}
	var screenMode uint8 = 0
	fs := []string{"/Users/jeromelesaux/Downloads/bomberman.gif"}

	Animation(fs, screenMode, e)
}
