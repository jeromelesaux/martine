package gfx

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func RollLeft(rla, sla, iterations int, screenMode uint8, size Size, downgraded *image.NRGBA, newPalette color.Palette, filename string, exportType *ExportType) {

	// create downgraded palette image with rra pixels rotated
	// and call n iterations spritetransform with this input generated image
	// save the rotated image as png
	if rla != -1 || sla != -1 {
		fmt.Fprintf(os.Stdout, "RLA/SLA: Iterations (%d)\n", iterations)
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
			newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
			fmt.Fprintf(os.Stdout, "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
			Png(exportType.OutputPath+string(filepath.Separator)+newFilename, im)
			fmt.Fprintf(os.Stdout, "Tranform image in sprite iteration (%d)\n", i)
			SpriteTransform(im, newPalette, size, screenMode, newFilename, exportType)
		}
	}
}
func RollRight(rra, sra, iterations int, screenMode uint8, size Size, downgraded *image.NRGBA, newPalette color.Palette, filename string, exportType *ExportType) {
	if rra != -1 || sra != -1 {
		fmt.Fprintf(os.Stdout, "RRA/SRA: Iterations (%d)\n", iterations)

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
			newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
			fmt.Fprintf(os.Stdout, "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
			Png(exportType.OutputPath+string(filepath.Separator)+newFilename, im)
			fmt.Fprintf(os.Stdout, "Tranform image in sprite iteration (%d)\n", i)
			SpriteTransform(im, newPalette, size, screenMode, newFilename, exportType)
		}
	}
}

func RollUp(keephigh, losthigh, iterations int, screenMode uint8, size Size, downgraded *image.NRGBA, newPalette color.Palette, filename string, exportType *ExportType) {
	if keephigh != -1 || losthigh != -1 {
		fmt.Fprintf(os.Stdout, "keephigh/losthigh: Iterations (%d)\n", iterations)
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
			newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
			fmt.Fprintf(os.Stdout, "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
			Png(exportType.OutputPath+string(filepath.Separator)+newFilename, im)
			fmt.Fprintf(os.Stdout, "Tranform image in sprite iteration (%d)\n", i)
			SpriteTransform(im, newPalette, size, screenMode, newFilename, exportType)
		}
	}
}

func RollLow(keeplow, lostlow, iterations int, screenMode uint8, size Size, downgraded *image.NRGBA, newPalette color.Palette, filename string, exportType *ExportType) {
	if keeplow != -1 || lostlow != -1 {
		fmt.Fprintf(os.Stdout, "keeplow/lostlow: Iterations (%d)\n", iterations)
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
			newFilename := strconv.Itoa(i) + strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
			fmt.Fprintf(os.Stdout, "Saving downgraded image iteration (%d) into (%s)\n", i, newFilename)
			Png(exportType.OutputPath+string(filepath.Separator)+newFilename, im)
			fmt.Fprintf(os.Stdout, "Tranform image in sprite iteration (%d)\n", i)
			SpriteTransform(im, newPalette, size, screenMode, newFilename, exportType)
		}
	}
}
