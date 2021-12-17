package ui

import (
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
)

func NewResizeAlgorithmSelect(me *ImageMenu) *widget.Select {
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
			me.resizeAlgoNumber = 0
			me.resizeAlgo = imaging.NearestNeighbor
		case "CatmullRom":
			me.resizeAlgoNumber = 1
			me.resizeAlgo = imaging.CatmullRom
		case "Lanczos":
			me.resizeAlgoNumber = 2
			me.resizeAlgo = imaging.Lanczos
		case "Linear":
			me.resizeAlgoNumber = 3
			me.resizeAlgo = imaging.Linear
		case "Box":
			me.resizeAlgoNumber = 4
			me.resizeAlgo = imaging.Box
		case "Hermite":
			me.resizeAlgoNumber = 5
			me.resizeAlgo = imaging.Hermite
		case "BSpline":
			me.resizeAlgoNumber = 6
			me.resizeAlgo = imaging.BSpline
		case "Hamming":
			me.resizeAlgoNumber = 7
			me.resizeAlgo = imaging.Hamming
		case "Hann":
			me.resizeAlgoNumber = 8
			me.resizeAlgo = imaging.Hann
		case "Gaussian":
			me.resizeAlgoNumber = 9
			me.resizeAlgo = imaging.Gaussian
		case "Blackman":
			me.resizeAlgoNumber = 10
			me.resizeAlgo = imaging.Blackman
		case "Bartlett":
			me.resizeAlgoNumber = 11
			me.resizeAlgo = imaging.Bartlett
		case "Welch":
			me.resizeAlgoNumber = 12
			me.resizeAlgo = imaging.Welch
		case "Cosine":
			me.resizeAlgoNumber = 13
			me.resizeAlgo = imaging.Cosine
		case "MitchellNetravali":
			me.resizeAlgoNumber = 14
			me.resizeAlgo = imaging.MitchellNetravali
		}
	})

	resize.SetSelected("NearestNeighbor")
	return resize
}
