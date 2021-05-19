package animate

import (
	"fmt"
	"image"
	"image/gif"
	"os"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/gfx/transformation"
)

func DeltaMotif(gitFilepath string, ex *export.ExportType, threshold int, initialAddress uint16, mode uint8) error {
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
	size := constants.Mode1

	screens := make([]*image.NRGBA, 0)
	// resize all images
	for _, v := range images {
		out := convert.Resize(v, size, imaging.NearestNeighbor)
		screens = append(screens, out)
	}

	// downgrading palette
	customPalette, _, err := convert.DowngradingPalette(screens[0], size, ex.CpcPlus)
	if err != nil {
		return err
	}
	// converting all screens
	for index, v := range screens {
		_, out := convert.DowngradingWithPalette(v, customPalette)
		screens[index] = out
	}

	// recuperation des motifs
	a := transformation.AnalyzeTilesBoard(screens[0], constants.Size{Width: 4, Height: 4})
	refBoard := a.ReduceTilesNumber(float64(threshold))
	btc := make([][]transformation.BoardTile, 0)
	btc = append(btc, refBoard)
	refTiles := transformation.GetUniqTiles(refBoard)

	// application des motifs sur toutes les images
	for i := 1; i < len(screens); i++ {
		ab := transformation.AnalyzeTilesBoardWithTiles(screens[i], constants.Size{Width: 4, Height: 4}, refTiles)
		btc = append(btc, ab.ReduceTilesNumber(float64(threshold)))
	}

	return nil
}
