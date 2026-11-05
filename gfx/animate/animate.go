// Package animate provides animation and delta packing utilities.
package animate

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/sprite"
	p "github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/log"
)

func Animation(filepaths []string, screenMode uint8, export *config.MartineConfig) error {
	sizeScreen := constants.NewSizeMode(screenMode, true)

	export.ScrCfg.Type = config.FullscreenFormat
	board, palette, err := concatSprites(filepaths, sizeScreen, export.ScrCfg.Size, screenMode, export)
	if err != nil {
		log.GetLogger().Error("Cannot concat content of files %v error :%v\n", filepaths, err)
		return err
	}
	if err := gfx.Transform(board, palette, sizeScreen, filepath.Join(export.ScrCfg.OutputPath, "board.png"), export); err != nil {
		log.GetLogger().Error("Can not transform to image error : %v\n", err)
		return err
	}
	return nil
}

func processSpriteFrame(in image.Image, filename string, palette color.Palette, screenMode uint8, export *config.MartineConfig) (*image.NRGBA, color.Palette, error) {
	out := ci.Resize(in, export.ScrCfg.Size, export.ScrCfg.Process.ResizingAlgo)
	log.GetLogger().Info("Saving resized image into (%s)\n", filename+"_resized.png")
	if err := p.Png(filepath.Join(export.ScrCfg.OutputPath, filename+"_resized.png"), out); err != nil {
		return nil, nil, err
	}
	var downgraded *image.NRGBA
	var newPalette color.Palette
	var err error
	if len(palette) > 0 {
		newPalette, downgraded = ci.DowngradingWithPalette(out, palette)
	} else {
		newPalette, downgraded, err = ci.DowngradingPalette(out, export.ScrCfg.Size, export.ScrCfg.IsPlus)
		if err != nil {
			log.GetLogger().Error("Cannot downgrade colors palette for this image\n")
		}
	}
	newPalette = constants.SortColorsByDistance(newPalette)
	log.GetLogger().Info("Saving downgraded image into (%s)\n", filename+"_down.png")
	if err := p.Png(filepath.Join(export.ScrCfg.OutputPath, filename+"_down.png"), downgraded); err != nil {
		return nil, nil, err
	}
	if err := sprite.ToSpriteAndExport(downgraded, newPalette, export.ScrCfg.Size, screenMode, filename, true, export); err != nil {
		log.GetLogger().Error("error while transform in sprite error : %v\n", err)
	}
	return downgraded, newPalette, nil
}

// nolint: funlen, gocognit, gocyclo
func concatSprites(filepaths []string, sizeScreen, spriteSize constants.Size, screenMode uint8, export *config.MartineConfig) (*image.NRGBA, color.Palette, error) {
	nbImgWidth := sizeScreen.Width / spriteSize.Width
	largeMarge := (sizeScreen.Width - (spriteSize.Width * nbImgWidth)) / nbImgWidth

	board := image.NewNRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: sizeScreen.Width, Y: sizeScreen.Height}})
	var palette, newPalette color.Palette
	palette = gfx.LoadPalette(export)

	var startX, startY int
	nbLarge := 0
	for index0, v := range filepaths {

		if strings.ToUpper(filepath.Ext(v)) == ".GIF" {
			f, err := os.Open(v)
			if err != nil {
				return board, newPalette, err
			}
			g, err := gif.DecodeAll(f)
			if err != nil {
				return board, newPalette, err
			}
			_ = f.Close()
			for index, in := range g.Image {
				if in.Rect.Max.X != g.Config.Width || in.Rect.Max.Y != g.Config.Height {
					newIm := image.NewPaletted(
						image.Rectangle{image.Point{X: 0, Y: 0},
							image.Point{X: g.Config.Width, Y: g.Config.Height}},
						in.Palette)
					draw.Draw(newIm, in.Rect, in, in.Rect.Min, draw.Src)
					in = newIm
				}
				filename := fmt.Sprintf("%.2d", index)
				downgraded, newPalette, err := processSpriteFrame(in, filename, palette, screenMode, export)
				if err != nil {
					return board, newPalette, err
				}
				contour := image.Rectangle{Min: image.Point{X: startX, Y: startY}, Max: image.Point{X: startX + spriteSize.Width, Y: startY + spriteSize.Height}}
				draw.Draw(board, contour, downgraded, image.Point{0, 0}, draw.Src)
				nbLarge++
				if nbLarge >= nbImgWidth {
					nbLarge = 0
					startX = 0
					startY += spriteSize.Height
				} else {
					startX += spriteSize.Width + largeMarge
				}
			}
		} else if strings.ToUpper(filepath.Ext(v)) == ".PNG" {
			f, err := os.Open(v)
			if err != nil {
				return board, newPalette, err
			}
			in, err := png.Decode(f)
			_ = f.Close()
			if err != nil {
				log.GetLogger().Error("Error while reading png file (%s) error %v, skipping\n", v, err)
				continue
			}
			filename := fmt.Sprintf("%.2d", index0)
			downgraded, newPalette, err := processSpriteFrame(in, filename, palette, screenMode, export)
			if err != nil {
				return board, newPalette, err
			}
			contour := image.Rectangle{Min: image.Point{X: startX, Y: startY}, Max: image.Point{X: startX + spriteSize.Width, Y: startY + spriteSize.Height}}
			draw.Draw(board, contour, downgraded, image.Point{0, 0}, draw.Src)
			nbLarge++
			if nbLarge >= nbImgWidth {
				nbLarge = 0
				startX = 0
				startY += spriteSize.Height
			} else {
				startX += spriteSize.Width + largeMarge
			}
		} else {
			log.GetLogger().Error("File is not a image file compatible (%s) skipping.\n", v)
		}

	}
	if err := p.Png(filepath.Join(export.ScrCfg.OutputPath, "board.png"), board); err != nil {
		return nil, nil, err
	}
	return board, newPalette, nil
}
