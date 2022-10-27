package common

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/gfx/errors"
)

func RawSpriteToImg(data []byte, height, width, mode uint8, p color.Palette) *image.NRGBA {
	var out *image.NRGBA
	switch mode {
	case 0:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 2), Y: int(height)}})
		index := 0

		for y := 0; y < int(height); y++ {
			indexX := 0
			for x := 0; x < int(width); x++ {
				val := data[index]
				pp1, pp2 := RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				index++
			}
		}
	case 1:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 4), Y: int(height)}})
		index := 0
		for y := 0; y < int(height); y++ {
			indexX := 0
			for x := 0; x < int(width); x++ {
				val := data[index]
				pp1, pp2, pp3, pp4 := RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				index++
			}
		}
	case 2:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(width * 8), Y: int(height)}})
		index := 0
		for y := 0; y < int(width); y++ {
			indexX := 0
			for x := 0; x < int(width/8); x++ {
				val := data[index]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				out.Set(indexX, y, c5)
				indexX++
				out.Set(indexX, y, c6)
				indexX++
				out.Set(indexX, y, c7)
				indexX++
				out.Set(indexX, y, c8)
				indexX++
				index++
			}
		}
	}
	return out
}

// spriteToImg load a OCP win filepath to image.NRGBA
// using the mode and palette as arguments
func SpriteToImg(winPath string, mode uint8, p color.Palette) (*image.NRGBA, constants.Size, error) {
	var s constants.Size
	footer, err := file.OpenWin(winPath)
	if err != nil {
		return nil, s, err
	}
	var out *image.NRGBA

	d, err := file.RawWin(winPath)
	if err != nil {
		return nil, s, err
	}
	switch mode {
	case 0:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(footer.Width * 2), Y: int(footer.Height)}})
		index := 0
		s.Width = int(footer.Width * 2)
		s.Height = int(footer.Height)
		for y := 0; y < int(footer.Height); y++ {
			indexX := 0
			for x := 0; x < int(footer.Width); x++ {
				val := d[index]
				pp1, pp2 := RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				index++
			}
		}
	case 1:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(footer.Width * 4), Y: int(footer.Height)}})
		index := 0
		s.Width = int(footer.Width * 4)
		s.Height = int(footer.Height)
		for y := 0; y < int(footer.Height); y++ {
			indexX := 0
			for x := 0; x < int(footer.Width); x++ {
				val := d[index]
				pp1, pp2, pp3, pp4 := RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				index++
			}
		}
	case 2:
		out = image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(footer.Width * 8), Y: int(footer.Height)}})
		index := 0
		s.Width = int(footer.Width * 8)
		s.Height = int(footer.Height)
		for y := 0; y < int(footer.Width); y++ {
			indexX := 0
			for x := 0; x < int(footer.Width/8); x++ {
				val := d[index]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(indexX, y, c1)
				indexX++
				out.Set(indexX, y, c2)
				indexX++
				out.Set(indexX, y, c3)
				indexX++
				out.Set(indexX, y, c4)
				indexX++
				out.Set(indexX, y, c5)
				indexX++
				out.Set(indexX, y, c6)
				indexX++
				out.Set(indexX, y, c7)
				indexX++
				out.Set(indexX, y, c8)
				indexX++
				index++
			}
		}
	}
	return out, s, err
}

// spriteToPng will save the sprite OCP win filepath in a png file
// output will be the png filpath, using the mode and palette as arguments
func SpriteToPng(winPath string, output string, mode uint8, p color.Palette) error {
	out, _, err := SpriteToImg(winPath, mode, p)
	if err != nil {
		return err
	}
	return file.Png(output+".png", out)
}

// scrRawToImg will convert the classical OCP screen slice of bytes  into image.NRGBA structure
// using the mode and the palette as arguments
func ScrRawToImg(d []byte, mode uint8, p color.Palette) (*image.NRGBA, error) {
	var m constants.Size
	switch mode {
	case 0:
		m = constants.Mode0
	case 1:
		m = constants.Mode1
	case 2:
		m = constants.Mode2
	default:
		return nil, errors.ErrorUndefinedMode
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 2 {
				val := d[cpcLine+cpcRow]
				pp1, pp2 := RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]

				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				cpcRow++
			}
			cpcRow = 0
		}
	case 1:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				cpcRow++
			}
			cpcRow = 0
		}
	case 2:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				out.Set(x+4, y, c5)
				out.Set(x+5, y, c6)
				out.Set(x+6, y, c7)
				out.Set(x+7, y, c8)
				cpcRow++
			}
			cpcRow = 0
		}
	}
	return out, nil
}

// SrcToImg load the amstrad classical 17ko  screen image to image.NRBGA
// using the mode and palette as arguments
func ScrToImg(scrPath string, mode uint8, p color.Palette) (*image.NRGBA, error) {
	var m constants.Size
	switch mode {
	case 0:
		m = constants.Mode0
	case 1:
		m = constants.Mode1
	case 2:
		m = constants.Mode2
	default:
		return nil, errors.ErrorUndefinedMode
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	d, err := file.RawScr(scrPath)
	if err != nil {
		return nil, err
	}
	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 2 {
				val := d[cpcLine+cpcRow]
				pp1, pp2 := RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]

				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				cpcRow++
			}
			cpcRow = 0
		}
	case 1:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				cpcRow++
			}
			cpcRow = 0
		}
	case 2:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x50 + ((y % 0x8) * 0x800))
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				out.Set(x+4, y, c5)
				out.Set(x+5, y, c6)
				out.Set(x+6, y, c7)
				out.Set(x+7, y, c8)
				cpcRow++
			}
			cpcRow = 0
		}
	}
	return out, nil
}

func ScrToPng(scrPath string, output string, mode uint8, p color.Palette) error {

	out, err := ScrToImg(scrPath, mode, p)
	if err != nil {
		return err
	}
	return file.Png(output, out)
}

// overscanRawToImg will convert fullscreen amstrad screen slice of bytes in image.NRGBA
// using the  screen mode  and the palette as arguments
func OverscanRawToImg(d []byte, mode uint8, p color.Palette) (*image.NRGBA, error) {
	var m constants.Size
	switch mode {
	case 0:
		m = constants.OverscanMode0
	case 1:
		m = constants.OverscanMode1
	case 2:
		m = constants.OverscanMode2
	default:
		return nil, errors.ErrorUndefinedMode
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 2 {
				val := d[cpcLine+cpcRow]
				pp1, pp2 := RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]

				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				cpcRow++
			}
			cpcRow = 0
		}
	case 1:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				cpcRow++
			}
			cpcRow = 0
		}
	case 2:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				out.Set(x+4, y, c5)
				out.Set(x+5, y, c6)
				out.Set(x+6, y, c7)
				out.Set(x+7, y, c8)
				cpcRow++
			}
			cpcRow = 0
		}
	}
	return out, nil
}

// overscanToImg will convert fullscreen amstrad screen filepath 'srcPath' in image.NRGBA
// using the  screen mode  and the palette as arguments
func OverscanToImg(scrPath string, mode uint8, p color.Palette) (*image.NRGBA, error) {
	var m constants.Size
	switch mode {
	case 0:
		m = constants.OverscanMode0
	case 1:
		m = constants.OverscanMode1
	case 2:
		m = constants.OverscanMode2
	default:
		return nil, errors.ErrorUndefinedMode
	}
	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(m.Width), Y: int(m.Height)}})

	d, err := file.RawOverscan(scrPath) // RawOverscan data commence en 0x30
	if err != nil {
		return nil, err
	}

	cpcRow := 0
	switch mode {
	case 0:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 2 {
				val := d[cpcLine+cpcRow]
				pp1, pp2 := RawPixelMode0(val)
				c1 := p[pp1]
				c2 := p[pp2]

				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				cpcRow++
			}
			cpcRow = 0
		}
	case 1:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 4 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4 := RawPixelMode1(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				cpcRow++
			}
			cpcRow = 0
		}
	case 2:
		for y := 0; y < m.Height; y++ {
			cpcLine := ((y/0x8)*0x60 + ((y % 0x8) * 0x800))
			if y > 127 {
				cpcLine += (0x3800)
			}
			for x := 0; x < m.Width; x += 8 {
				val := d[cpcLine+cpcRow]
				pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 := RawPixelMode2(val)
				c1 := p[pp1]
				c2 := p[pp2]
				c3 := p[pp3]
				c4 := p[pp4]
				c5 := p[pp5]
				c6 := p[pp6]
				c7 := p[pp7]
				c8 := p[pp8]
				out.Set(x, y, c1)
				out.Set(x+1, y, c2)
				out.Set(x+2, y, c3)
				out.Set(x+3, y, c4)
				out.Set(x+4, y, c5)
				out.Set(x+5, y, c6)
				out.Set(x+6, y, c7)
				out.Set(x+7, y, c8)
				cpcRow++
			}
			cpcRow = 0
		}
	}
	return out, nil
}

// overscanToPng will convert fullscreen amstrad screen filepath in png file 'output'
// using the  screen mode  and the palette as arguments
func OverscanToPng(scrPath string, output string, mode uint8, p color.Palette) error {
	out, err := OverscanToImg(scrPath, mode, p)
	if err != nil {
		return err
	}
	return file.Png(output, out)
}

func ExportSprite(data []byte, lineSize int, p color.Palette, size constants.Size, mode uint8, filename string, dontImportDsk bool, cont *export.MartineContext) error {
	if err := file.Win(filename, data, mode, lineSize, size.Height, dontImportDsk, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
		return err
	}
	if !cont.CpcPlus {
		if err := file.Pal(filename, p, mode, dontImportDsk, cont); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := cont.OsFullPath(filename, "_palettepal.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Ink(filename, p, 2, dontImportDsk, cont); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath = cont.OsFullPath(filename, "_paletteink.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	} else {
		if err := file.Kit(filename, p, mode, dontImportDsk, cont); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filename, err)
			return err
		}
		filePath := cont.OsFullPath(filename, "_palettekit.png")
		if err := file.PalToPng(filePath, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
	}
	if err := file.Ascii(filename, data, p, dontImportDsk, cont); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving ascii file for (%s) error :%v\n", filename, err)
	}
	return file.AsciiByColumn(filename, data, p, dontImportDsk, mode, cont)
}

func Export(filePath string, bw []byte, p color.Palette, screenMode uint8, ex *export.MartineContext) error {
	if ex.Overscan {
		if ex.EgxFormat == 0 {
			if ex.ExportAsGoFile {
				if err := file.SaveGo(filePath, bw, p, screenMode, ex); err != nil {
					fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
					return err
				}
			} else {
				if err := file.Overscan(filePath, bw, p, screenMode, ex); err != nil {
					fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
					return err
				}
			}
		} else {
			if err := file.EgxOverscan(filePath, bw, p, ex.EgxMode1, ex.EgxMode2, ex); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
				return err
			}
		}

	} else {
		if err := file.Scr(filePath, bw, p, screenMode, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := file.Loader(filePath, p, screenMode, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if !ex.CpcPlus {
		if err := file.Pal(filePath, p, screenMode, false, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := ex.OsFullPath(filePath, "_palettepal.png")
		if err := file.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
		if err := file.Ink(filePath, p, screenMode, false, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 = ex.OsFullPath(filePath, "_paletteink.png")
		if err := file.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	} else {
		if err := file.Kit(filePath, p, screenMode, false, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := ex.OsFullPath(filePath, "_palettekit.png")
		if err := file.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	}
	return file.Ascii(filePath, bw, p, false, ex)
}

// PalettePosition returns the position of the color c in the palette
// overwise ErrorColorNotFound error
func PalettePosition(c color.Color, p color.Palette) (int, error) {
	r, g, b, a := c.RGBA()
	for index, cp := range p {
		//fmt.Fprintf(os.Stdout,"index(%d), c:%v,cp:%v\n",index,c,cp)
		rp, gp, bp, ap := cp.RGBA()
		if r == rp && g == gp && b == bp && a == ap {
			//fmt.Fprintf(os.Stdout,"Position found")
			return index, nil
		}
	}
	return -1, errors.ErrorColorNotFound
}

// PixelMode0 converts palette position into byte in  screen mode 0
func PixelMode0(pp1, pp2 int) byte {
	var pixel byte
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	if uint8(pp1)&2 == 2 {
		pixel += 8
	}
	if uint8(pp1)&4 == 4 {
		pixel += 32
	}
	if uint8(pp1)&8 == 8 {
		pixel += 2
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp2)&2 == 2 {
		pixel += 4
	}
	if uint8(pp2)&4 == 4 {
		pixel += 16
	}
	if uint8(pp2)&8 == 8 {
		pixel++
	}
	return pixel
}

// PixelMode1 converts palette position into byte in  screen mode 1
func PixelMode1(pp1, pp2, pp3, pp4 int) byte {
	var pixel byte
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	if uint8(pp1)&2 == 2 {
		pixel += 8
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp2)&2 == 2 {
		pixel += 4
	}
	if uint8(pp3)&1 == 1 {
		pixel += 32
	}
	if uint8(pp3)&2 == 2 {
		pixel += 2
	}
	if uint8(pp4)&1 == 1 {
		pixel += 16
	}
	if uint8(pp4)&2 == 2 {
		pixel++
	}
	return pixel
}

// PixelMode converts palette position into byte in screen mode 2
func PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 int) byte {
	var pixel byte
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp3)&1 == 1 {
		pixel += 32
	}
	if uint8(pp4)&1 == 1 {
		pixel += 16
	}
	if uint8(pp5)&1 == 1 {
		pixel += 8
	}
	if uint8(pp6)&1 == 1 {
		pixel += 4
	}
	if uint8(pp7)&1 == 1 {
		pixel += 2
	}
	if uint8(pp8)&1 == 1 {
		pixel++
	}
	return pixel
}

// RawPixelMode2 converts color  byte in palette position in  screen mode 2
func RawPixelMode2(b byte) (pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 int) {
	val := int(b)
	if val-128 >= 0 {
		pp1 = 1
		val -= 128
	}
	if val-64 >= 0 {
		pp2 = 1
		val -= 64
	}
	if val-32 >= 0 {
		pp3 = 1
		val -= 32
	}
	if val-16 >= 0 {
		pp4 = 1
		val -= 16
	}
	if val-8 >= 0 {
		pp5 = 1
		val -= 8
	}
	if val-4 >= 0 {
		pp6 = 1
		val -= 4
	}
	if val-2 >= 0 {
		pp7 = 1
		val -= 2
	}
	if val-1 >= 0 {
		pp8 = 1
	}
	return
}

// RawPixelMode1 converts color  byte in palette position in screen mode 1
func RawPixelMode1(b byte) (pp1, pp2, pp3, pp4 int) {
	val := int(b)
	if val-128 >= 0 {
		pp1 |= 1
		val -= 128
	}
	if val-64 >= 0 {
		pp2 |= 1
		val -= 64
	}
	if val-32 >= 0 {
		pp3 |= 1
		val -= 32
	}
	if val-16 >= 0 {
		pp4 |= 1
		val -= 16
	}
	if val-8 >= 0 {
		pp1 |= 2
		val -= 8
	}
	if val-4 >= 0 {
		pp2 |= 2
		val -= 4
	}
	if val-2 >= 0 {
		pp3 |= 2
		val -= 2
	}
	if val-1 >= 0 {
		pp4 |= 2
	}

	return
}

// RawPixelMode0 converts color  byte in palette position in screen mode 0
func RawPixelMode0(b byte) (pp1, pp2 int) {
	val := int(b)
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-128 >= 0 {
		pp1 |= 1
		val -= 128
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-64 >= 0 {
		pp2 |= 1
		val -= 64
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-32 >= 0 {
		pp1 |= 4
		val -= 32
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-16 >= 0 {
		pp2 |= 4
		val -= 16
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-8 >= 0 {
		pp1 |= 2
		val -= 8
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-4 >= 0 {
		pp2 |= 2
		val -= 4
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-2 >= 0 {
		pp1 |= 8
		val -= 2
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-1 >= 0 {
		pp2 |= 8
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	return
}

// CpcScreenAddress returns the screen address according the screen mode, the initialAddress (always #C000)
// x the column number and y the line number on the screen
func CpcScreenAddress(intialeAddresse int, x, y int, mode uint8, isOverscan bool) int {
	var addr int
	var adjustMode int
	switch mode {
	case 0:
		adjustMode = 2
	case 1:
		adjustMode = 4
	case 2:
		adjustMode = 8
	}
	if isOverscan {
		if y > 127 {
			addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / adjustMode) + (0x3800)
		} else {
			addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / adjustMode)
		}
	} else {
		addr = (0x800 * (y % 8)) + (0x50 * (y / 8)) + ((x + 1) / adjustMode)
	}
	if intialeAddresse == 0 {
		return addr
	}
	return intialeAddresse + addr
}

func CpcScreenAddressOffset(line int) int {
	return int(math.Floor(float64(line)/8)*80) + (line%8)*2048
}
