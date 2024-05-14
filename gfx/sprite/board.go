package sprite

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/gfx"
)

func SplitBoardToSprite(
	im image.Image,
	p color.Palette,
	row, col int,
	mode uint8,
	isSpriteHard bool,
	size constants.Size,
) ([][][]byte, [][]*image.NRGBA, error) {
	var results [][][]byte
	spriteWidth := im.Bounds().Max.X / col
	spriteHeight := im.Bounds().Max.Y / row
	sprites := make([][]*image.NRGBA, row)
	x := 0
	y := 0
	index := 0

	for j := 0; j < row; j++ {
		for i := 0; i < col; i++ {
			img := image.NewNRGBA(image.Rect(0, 0, spriteWidth, spriteHeight))
			for x0 := 0; x0 < spriteWidth; x0++ {
				for y0 := 0; y0 < spriteHeight; y0++ {
					img.Set(x0, y0, im.At(x+x0, y+y0))
				}
			}
			sprites[index] = append(sprites[index], img)
			x += spriteWidth
		}
		index++
		y += spriteHeight
		x = 0
	}
	cfg := config.NewMartineConfig("", "")
	cfg.CustomDimension = true
	rawSprites := make([][]*image.NRGBA, len(sprites))
	results = make([][][]byte, len(sprites))

	cfg.Size = size
	cfg.SpriteHard = isSpriteHard
	for i := 0; i < len(sprites); i++ {
		results[i] = make([][]byte, len(sprites[i]))
		for j := 0; j < len(sprites[i]); j++ {
			v := sprites[i][j]
			r, sp, _, _, err := gfx.ApplyOneImage(v, cfg, int(mode), p, mode)
			if err != nil {
				return results, sprites, err
			}
			rawSprites[i] = append(rawSprites[i], sp)
			results[i][j] = append(results[i][j], r...)
		}
	}
	/*for _, v := range sprites {
		r, sp, _, _, err := gfx.ApplyOneImage(v, cont, int(mode), p, mode)
		if err != nil {
			return results, sprites, err
		}
		results = append(results, r)
		rawSprites = append(rawSprites, sp)
	}*/
	return results, rawSprites, nil
}
