package widget

import (
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/gfx/filter"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func NewDitheringSelect(me *menu.ImageMenu) *widget.Select {
	dithering := widget.NewSelect([]string{"FloydSteinberg",
		"JarvisJudiceNinke",
		"Stucki",
		"Atkinson",
		"Sierra",
		"SierraLite",
		"Sierra3",
		"Bayer2",
		"Bayer3",
		"Bayer4",
		"Bayer8",
	}, func(s string) {
		switch s {
		case "FloydSteinberg":
			me.DitheringAlgoNumber = 0
			me.DitheringMatrix = filter.FloydSteinberg
			me.DitheringType = constants.ErrorDiffusionDither
		case "JarvisJudiceNinke":
			me.DitheringAlgoNumber = 1
			me.DitheringMatrix = filter.JarvisJudiceNinke
			me.DitheringType = constants.ErrorDiffusionDither
		case "Stucki":
			me.DitheringAlgoNumber = 2
			me.DitheringMatrix = filter.Stucki
			me.DitheringType = constants.ErrorDiffusionDither
		case "Atkinson":
			me.DitheringAlgoNumber = 3
			me.DitheringMatrix = filter.Atkinson
			me.DitheringType = constants.ErrorDiffusionDither
		case "Sierra":
			me.DitheringAlgoNumber = 4
			me.DitheringMatrix = filter.Sierra
			me.DitheringType = constants.ErrorDiffusionDither
		case "SierraLite":
			me.DitheringAlgoNumber = 5
			me.DitheringMatrix = filter.SierraLite
			me.DitheringType = constants.ErrorDiffusionDither
		case "Sierra3":
			me.DitheringAlgoNumber = 6
			me.DitheringMatrix = filter.Sierra3
			me.DitheringType = constants.ErrorDiffusionDither
		case "Bayer2":
			me.DitheringAlgoNumber = 7
			me.DitheringMatrix = filter.Bayer2
			me.DitheringType = constants.OrderedDither
		case "Bayer3":
			me.DitheringAlgoNumber = 8
			me.DitheringMatrix = filter.Bayer3
			me.DitheringType = constants.OrderedDither
		case "Bayer4":
			me.DitheringAlgoNumber = 9
			me.DitheringMatrix = filter.Bayer4
			me.DitheringType = constants.OrderedDither
		case "Bayer8":
			me.DitheringAlgoNumber = 10
			me.DitheringMatrix = filter.Bayer8
			me.DitheringType = constants.OrderedDither
		}
	})
	dithering.SetSelected("FloydSteinberg")
	return dithering
}
