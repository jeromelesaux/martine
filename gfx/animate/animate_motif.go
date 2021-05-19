package animate

import (
	"fmt"
	"image/gif"
	"os"

	"github.com/jeromelesaux/martine/export"
)

func DeltaMotif(gitFilepath string, ex *export.ExportType, initialAddress uint16, mode uint8) error {
	var isSprite = true
	var maxImages = 22
	if !ex.CustomDimension && !ex.SpriteHard {
		isSprite = false
	}
	fr, err := os.Open(gitFilepath)
	if err != nil {
		return err
	}
	defer fr.Close()
	gifImages, err := gif.DecodeAll(fr)
	if err != nil {
		return err
	}
	images := ConvertToImage(*gifImages)
	var pad int = 1
	if len(images) > maxImages {
		fmt.Fprintf(os.Stderr, "Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(images) / maxImages
	}

	if isSprite {
		pad++
	}

	// convertion de la palette en mode 1
	//	p := constants.CpcOldPalette
	// transformation des images en mode 1
	//	size := constants.Mode1
	//	DowngradingWithPalette()
	return nil
}
