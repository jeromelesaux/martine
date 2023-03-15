package transformation

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/log"
)

func RollLeft(rla, sla, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {

	images := make([]*image.NRGBA, 0)
	// create downgraded palette image with rra pixels rotated
	// and call n iterations spritetransform with this input generated image
	// save the rotated image as png
	if rla != -1 || sla != -1 {
		log.GetLogger().Info("RLA/SLA: Iterations (%d)\n", iterations)
		for i := 0; i < iterations; i++ {
			nbPixels := 0
			if rla != -1 {
				nbPixels = (rla * (1 + i))
			} else {
				if sla != -1 {
					nbPixels = (sla * (1 + i))
				}
			}
			im := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{downgraded.Bounds().Max.X, downgraded.Bounds().Max.Y}})
			y2 := 0
			for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
				x2 := 0
				for x := downgraded.Bounds().Min.X + nbPixels; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
			if rla != -1 {
				y2 = 0
				for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
					x2 := downgraded.Bounds().Max.X - nbPixels
					for x := downgraded.Bounds().Min.X; x < nbPixels; x++ {
						im.Set(x2, y2, downgraded.At(x, y))
						x2++
					}
					y2++
				}
			}
			images = append(images, im)
			/*	newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
				log.GetLogger().Info( "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
				file.Png(filepath.Join(cont.OutputPath, newFilename), im)
				log.GetLogger().Info( "Tranform image in sprite iteration (%d)\n", i)
				common.ToSpriteAndExport(im, newPalette, size, screenMode, newFilename, false, cont)*/
		}
	}
	return images
}
func RollRight(rra, sra, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {
	images := make([]*image.NRGBA, 0)
	if rra != -1 || sra != -1 {
		log.GetLogger().Info("RRA/SRA: Iterations (%d)\n", iterations)

		for i := 0; i < iterations; i++ {
			nbPixels := 0
			if rra != -1 {
				nbPixels = (rra * (1 + i))
			} else {
				if sra != -1 {
					nbPixels = (sra * (1 + i))
				}
			}
			im := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{downgraded.Bounds().Max.X, downgraded.Bounds().Max.Y}})
			y2 := 0
			for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
				x2 := nbPixels
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X-nbPixels; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
			if rra != -1 {
				y2 = 0
				for y := downgraded.Bounds().Min.Y; y < downgraded.Bounds().Max.Y; y++ {
					x2 := 0
					for x := downgraded.Bounds().Max.X - nbPixels; x < downgraded.Bounds().Max.X; x++ {
						im.Set(x2, y2, downgraded.At(x, y))
						x2++
					}
					y2++
				}
			}
			images = append(images, im)
			/*			newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
						log.GetLogger().Info( "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
						file.Png(filepath.Join(cont.OutputPath, newFilename), im)
						log.GetLogger().Info( "Tranform image in sprite iteration (%d)\n", i)
						common.ToSpriteAndExport(im, newPalette, size, screenMode, newFilename, false, cont)*/
		}
	}
	return images
}

func RollUp(keephigh, losthigh, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {
	images := make([]*image.NRGBA, 0)
	if keephigh != -1 || losthigh != -1 {
		log.GetLogger().Info("keephigh/losthigh: Iterations (%d)\n", iterations)
		for i := 0; i < iterations; i++ {
			nbPixels := 0
			if keephigh != -1 {
				nbPixels = (keephigh * (1 + i))
			} else {
				if losthigh != -1 {
					nbPixels = (losthigh * (1 + i))
				}
			}
			im := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{downgraded.Bounds().Max.X, downgraded.Bounds().Max.Y}})
			y2 := 0
			for y := downgraded.Bounds().Min.Y + nbPixels; y < downgraded.Bounds().Max.Y; y++ {
				x2 := 0
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2++
			}
			if keephigh != -1 {
				for y := downgraded.Bounds().Min.Y; y < nbPixels; y++ {
					x2 := 0
					for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
						im.Set(x2, y2, downgraded.At(x, y))
						x2++
					}
					y2++
				}
			}
			images = append(images, im)
			/*	newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
				log.GetLogger().Info( "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
				file.Png(filepath.Join(cont.OutputPath, newFilename), im)
				log.GetLogger().Info( "Tranform image in sprite iteration (%d)\n", i)
				common.ToSpriteAndExport(im, newPalette, size, screenMode, newFilename, false, cont)*/
		}
	}
	return images
}

func RollLow(keeplow, lostlow, iterations int,
	screenMode uint8,
	size constants.Size,
	downgraded *image.NRGBA,
	newPalette color.Palette) []*image.NRGBA {
	images := make([]*image.NRGBA, 0)
	if keeplow != -1 || lostlow != -1 {
		log.GetLogger().Info("keeplow/lostlow: Iterations (%d)\n", iterations)
		for i := 0; i < iterations; i++ {
			nbPixels := 0
			if keeplow != -1 {
				nbPixels = (keeplow * (1 + i))
			} else {
				if lostlow != -1 {
					nbPixels = (lostlow * (1 + i))
				}
			}
			im := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{downgraded.Bounds().Max.X, downgraded.Bounds().Max.Y}})
			y2 := downgraded.Bounds().Max.Y - 1
			for y := downgraded.Bounds().Max.Y - nbPixels; y >= downgraded.Bounds().Min.Y; y-- {
				x2 := 0
				for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
					im.Set(x2, y2, downgraded.At(x, y))
					x2++
				}
				y2--
			}
			if keeplow != -1 {
				for y := downgraded.Bounds().Max.Y - 1; y >= downgraded.Bounds().Max.Y-nbPixels; y-- {
					x2 := 0
					for x := downgraded.Bounds().Min.X; x < downgraded.Bounds().Max.X; x++ {
						im.Set(x2, y2, downgraded.At(x, y))
						x2++
					}
					y2--
				}
			}
			images = append(images, im)
			/*	newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
				log.GetLogger().Info( "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
				file.Png(filepath.Join(cont.OutputPath, newFilename), im)
				log.GetLogger().Info( "Tranform image in sprite iteration (%d)\n", i)
				common.ToSpriteAndExport(im, newPalette, size, screenMode, newFilename, false, cont)*/
		}
	}
	return images
}
