package widget

import (
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/gfx/filter"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

// nolint: funlen
func NewDitheringSelect(me *menu.ImageMenu) *widget.Select {
	dithering := widget.NewSelect([]string{
		"FloydSteinberg",
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
			me.Cfg.ScrCfg.Process.DitheringAlgo = 0
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.FloydSteinberg
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "JarvisJudiceNinke":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 1
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.JarvisJudiceNinke
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "Stucki":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 2
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Stucki
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "Atkinson":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 3
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Atkinson
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "Sierra":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 4
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Sierra
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "SierraLite":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 5
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.SierraLite
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "Sierra3":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 6
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Sierra3
			me.Cfg.ScrCfg.Process.DitheringType = constants.ErrorDiffusionDither
		case "Bayer2":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 7
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Bayer2
			me.Cfg.ScrCfg.Process.DitheringType = constants.OrderedDither
		case "Bayer3":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 8
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Bayer3
			me.Cfg.ScrCfg.Process.DitheringType = constants.OrderedDither
		case "Bayer4":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 9
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Bayer4
			me.Cfg.ScrCfg.Process.DitheringType = constants.OrderedDither
		case "Bayer8":
			me.Cfg.ScrCfg.Process.DitheringAlgo = 10
			me.Cfg.ScrCfg.Process.DitheringMatrix = filter.Bayer8
			me.Cfg.ScrCfg.Process.DitheringType = constants.OrderedDither
		}
	})
	dithering.SetSelected("FloydSteinberg")
	return dithering
}
