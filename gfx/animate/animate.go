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

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/common"
)

func Animation(filepaths []string, screenMode uint8, export *export.ExportType) error {
	var sizeScreen constants.Size
	switch screenMode {
	case 0:
		sizeScreen = constants.OverscanMode0
	case 1:
		sizeScreen = constants.OverscanMode1
	case 2:
		sizeScreen = constants.OverscanMode2
	default:
		sizeScreen = constants.OverscanMode0
	}
	export.Overscan = true
	board, palette, err := concatSprites(filepaths, sizeScreen, export.Size, screenMode, export)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot concat content of files %v error :%v\n", filepaths, err)
		return err
	}
	if err := gfx.Transform(board, palette, sizeScreen, filepath.Join(export.OutputPath, "board.png"), export); err != nil {
		fmt.Fprintf(os.Stderr, "Can not transform to image error : %v\n", err)
		return err
	}
	return nil
}

func concatSprites(filepaths []string, sizeScreen, spriteSize constants.Size, screenMode uint8, export *export.ExportType) (*image.NRGBA, color.Palette, error) {
	nbImgWidth := int(sizeScreen.Width / spriteSize.Width)
	//nbImgHeight := int(sizeScreen.Height / size.Height)
	largeMarge := (sizeScreen.Width - (spriteSize.Width * nbImgWidth)) / nbImgWidth

	board := image.NewNRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: sizeScreen.Width, Y: sizeScreen.Height}})
	var palette, newPalette color.Palette
	if export.PalettePath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", export.PalettePath)
		palette, _, err := file.OpenPal(export.PalettePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", export.PalettePath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if export.InkPath != "" {
		fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", export.InkPath)
		palette, _, err := file.OpenInk(export.InkPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", export.InkPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	if export.KitPath != "" {
		fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", export.KitPath)
		palette, _, err := file.OpenKit(export.KitPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", export.KitPath)
		} else {
			fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
		}
	}
	var startX, startY int
	nbLarge := 0
	for index0, v := range filepaths {

		if strings.ToUpper(filepath.Ext(v)) == ".GIF" {
			f, err := os.Open(v)
			if err != nil {
				return board, newPalette, err
			}
			defer f.Close()
			g, err := gif.DecodeAll(f)
			if err != nil {
				return board, newPalette, err
			}

			for index, in := range g.Image {
				// gif change size between frame.
				// create a new blank image with size from config.width config.height
				// add position into blank image the paletted image rect
				if in.Rect.Max.X != g.Config.Width || in.Rect.Max.Y != g.Config.Height {
					newIm := image.NewPaletted(
						image.Rectangle{image.Point{X: 0, Y: 0},
							image.Point{X: g.Config.Width, Y: g.Config.Height}},
						in.Palette)
					draw.Draw(newIm, in.Rect, in, in.Rect.Min, draw.Src)
					in = newIm
				}
				var downgraded *image.NRGBA
				filename := fmt.Sprintf("%.2d", index)
				out := convert.Resize(in, export.Size, export.ResizingAlgo)
				fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
				if err := file.Png(filepath.Join(export.OutputPath, filename+"_resized.png"), out); err != nil {
					os.Exit(-2)
				}

				if len(palette) > 0 {
					newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
				} else {
					newPalette, downgraded, err = convert.DowngradingPalette(out, export.Size, export.CpcPlus)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", v)
					}
				}

				newPalette = constants.SortColorsByDistance(newPalette)

				fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
				if err := file.Png(filepath.Join(export.OutputPath, filename+"_down.png"), downgraded); err != nil {
					os.Exit(-2)
				}

				if err := common.ToSpriteAndExport(downgraded, newPalette, export.Size, screenMode, filename, true, export); err != nil {
					fmt.Fprintf(os.Stderr, "error while transform in sprite error : %v\n", err)
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
		} else {
			if strings.ToUpper(filepath.Ext(v)) == ".PNG" {
				f, err := os.Open(v)
				if err != nil {
					return board, newPalette, err
				}
				defer f.Close()
				in, err := png.Decode(f)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while reading png file (%s) error %v, skipping\n", v, err)
					continue
				}
				var downgraded *image.NRGBA
				filename := fmt.Sprintf("%.2d", index0)
				out := convert.Resize(in, export.Size, export.ResizingAlgo)
				fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
				if err := file.Png(filepath.Join(export.OutputPath, filename+"_resized.png"), out); err != nil {
					os.Exit(-2)
				}

				if len(palette) > 0 {
					newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
				} else {
					newPalette, downgraded, err = convert.DowngradingPalette(out, export.Size, export.CpcPlus)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", v)
					}
				}

				newPalette = constants.SortColorsByDistance(newPalette)

				fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
				if err := file.Png(filepath.Join(export.OutputPath, filename+"_down.png"), downgraded); err != nil {
					os.Exit(-2)
				}

				if err := common.ToSpriteAndExport(downgraded, newPalette, export.Size, screenMode, filename, true, export); err != nil {
					fmt.Fprintf(os.Stderr, "error while transform in sprite error : %v\n", err)
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
				fmt.Fprintf(os.Stderr, "File is not a image file compatible (%s) skipping.\n", v)
			}
		}

	}
	if err := file.Png(filepath.Join(export.OutputPath, "board.png"), board); err != nil {
		os.Exit(-2)
	}
	return board, newPalette, nil
}
