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
			me.ResizeAlgo = imaging.NearestNeighbor
		case "CatmullRom":
			me.ResizeAlgoNumber = 1
			me.ResizeAlgo = imaging.CatmullRom
		case "Lanczos":
			me.ResizeAlgoNumber = 2
			me.ResizeAlgo = imaging.Lanczos
		case "Linear":
			me.ResizeAlgoNumber = 3
			me.ResizeAlgo = imaging.Linear
		case "Box":
			me.ResizeAlgoNumber = 4
			me.ResizeAlgo = imaging.Box
		case "Hermite":
			me.ResizeAlgoNumber = 5
			me.ResizeAlgo = imaging.Hermite
		case "BSpline":
			me.ResizeAlgoNumber = 6
			me.ResizeAlgo = imaging.BSpline
		case "Hamming":
			me.ResizeAlgoNumber = 7
			me.ResizeAlgo = imaging.Hamming
		case "Hann":
			me.ResizeAlgoNumber = 8
			me.ResizeAlgo = imaging.Hann
		case "Gaussian":
			me.ResizeAlgoNumber = 9
			me.ResizeAlgo = imaging.Gaussian
		case "Blackman":
			me.ResizeAlgoNumber = 10
			me.ResizeAlgo = imaging.Blackman
		case "Bartlett":
			me.ResizeAlgoNumber = 11
			me.ResizeAlgo = imaging.Bartlett
		case "Welch":
			me.ResizeAlgoNumber = 12
			me.ResizeAlgo = imaging.Welch
		case "Cosine":
			me.ResizeAlgoNumber = 13
			me.ResizeAlgo = imaging.Cosine
		case "MitchellNetravali":
			me.ResizeAlgoNumber = 14
			me.ResizeAlgo = imaging.MitchellNetravali
		}
	})

	resize.SetSelected("NearestNeighbor")
	return resize
}
