package ui

import (
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/gfx/filter"
)

func NewDitheringSelect(me *ImageMenu) *widget.Select {
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
			me.ditheringAlgoNumber = 0
			me.ditheringMatrix = filter.FloydSteinberg
			me.ditheringType = constants.ErrorDiffusionDither
		case "JarvisJudiceNinke":
			me.ditheringAlgoNumber = 1
			me.ditheringMatrix = filter.JarvisJudiceNinke
			me.ditheringType = constants.ErrorDiffusionDither
		case "Stucki":
			me.ditheringAlgoNumber = 2
			me.ditheringMatrix = filter.Stucki
			me.ditheringType = constants.ErrorDiffusionDither
		case "Atkinson":
			me.ditheringAlgoNumber = 3
			me.ditheringMatrix = filter.Atkinson
			me.ditheringType = constants.ErrorDiffusionDither
		case "Sierra":
			me.ditheringAlgoNumber = 4
			me.ditheringMatrix = filter.Sierra
			me.ditheringType = constants.ErrorDiffusionDither
		case "SierraLite":
			me.ditheringAlgoNumber = 5
			me.ditheringMatrix = filter.SierraLite
			me.ditheringType = constants.ErrorDiffusionDither
		case "Sierra3":
			me.ditheringAlgoNumber = 6
			me.ditheringMatrix = filter.Sierra3
			me.ditheringType = constants.ErrorDiffusionDither
		case "Bayer2":
			me.ditheringAlgoNumber = 7
			me.ditheringMatrix = filter.Bayer2
			me.ditheringType = constants.OrderedDither
		case "Bayer3":
			me.ditheringAlgoNumber = 8
			me.ditheringMatrix = filter.Bayer3
			me.ditheringType = constants.OrderedDither
		case "Bayer4":
			me.ditheringAlgoNumber = 9
			me.ditheringMatrix = filter.Bayer4
			me.ditheringType = constants.OrderedDither
		case "Bayer8":
			me.ditheringAlgoNumber = 10
			me.ditheringMatrix = filter.Bayer8
			me.ditheringType = constants.OrderedDither
		}
	})
	dithering.SetSelected("FloydSteinberg")
	return dithering
}
