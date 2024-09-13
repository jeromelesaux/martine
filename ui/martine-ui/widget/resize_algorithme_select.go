package widget

import (
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

// nolint: funlen
func NewResizeAlgorithmSelect(me *menu.ImageMenu) *widget.Select {
	resize := widget.NewSelect([]string{"NearestNeighbor",
		"CatmullRom",
		"Lanczos",
		"Linear",
		"Box",
		"Hermite",
		"BSpline",
		"Hamming",
		"Hann",
		"Gaussian",
		"Blackman",
		"Bartlett",
		"Welch",
		"Cosine",
		"MitchellNetravali",
	}, func(s string) {
		switch s {
		case "NearestNeighbor":
			me.ResizeAlgoNumber = 0
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.NearestNeighbor
		case "CatmullRom":
			me.ResizeAlgoNumber = 1
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.CatmullRom
		case "Lanczos":
			me.ResizeAlgoNumber = 2
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Lanczos
		case "Linear":
			me.ResizeAlgoNumber = 3
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Linear
		case "Box":
			me.ResizeAlgoNumber = 4
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Box
		case "Hermite":
			me.ResizeAlgoNumber = 5
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Hermite
		case "BSpline":
			me.ResizeAlgoNumber = 6
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.BSpline
		case "Hamming":
			me.ResizeAlgoNumber = 7
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Hamming
		case "Hann":
			me.ResizeAlgoNumber = 8
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Hann
		case "Gaussian":
			me.ResizeAlgoNumber = 9
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Gaussian
		case "Blackman":
			me.ResizeAlgoNumber = 10
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Blackman
		case "Bartlett":
			me.ResizeAlgoNumber = 11
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Bartlett
		case "Welch":
			me.ResizeAlgoNumber = 12
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Welch
		case "Cosine":
			me.ResizeAlgoNumber = 13
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.Cosine
		case "MitchellNetravali":
			me.ResizeAlgoNumber = 14
			me.Cfg.ScrCfg.Treatment.ResizingAlgo = imaging.MitchellNetravali
		}
	})

	resize.SetSelected("NearestNeighbor")
	return resize
}
