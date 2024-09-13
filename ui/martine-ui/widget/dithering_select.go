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
			me.Cfg.ScrCfg.Process.Dithering.Algo = 0
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.FloydSteinberg
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "JarvisJudiceNinke":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 1
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.JarvisJudiceNinke
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "Stucki":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 2
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Stucki
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "Atkinson":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 3
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Atkinson
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "Sierra":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 4
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Sierra
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "SierraLite":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 5
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.SierraLite
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "Sierra3":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 6
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Sierra3
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
		case "Bayer2":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 7
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer2
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
		case "Bayer3":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 8
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer3
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
		case "Bayer4":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 9
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer4
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
		case "Bayer8":
			me.Cfg.ScrCfg.Process.Dithering.Algo = 10
			me.Cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer8
			me.Cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
		}
	})
	dithering.SetSelected("FloydSteinberg")
	return dithering
}
