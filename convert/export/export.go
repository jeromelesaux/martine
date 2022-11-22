package export

import (
	"fmt"
	"image/color"
	"os"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
)

func Export(filePath string, bw []byte, p color.Palette, screenMode uint8, ex *config.MartineConfig) error {
	if ex.Overscan {
		if ex.EgxFormat == 0 {
			if ex.ExportAsGoFile {
				/*orig, err := OverscanRawToImg(bw, screenMode, p)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while converting into image file %s error :%v", filePath, err)
					return err
				}

				imgUp, imgDown, err := convert.SplitImage(orig)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while splitting image from  file %s error :%v", filePath, err)
					return err
				}
				config := export.NewMartineContext("", "")
				config.Size = constants.Size{Width: imgUp.Bounds().Max.X, Height: imgUp.Bounds().Max.Y}
				config.Overscan = true
				dataUp := ToMode0(imgUp, p, config)
				dataDown := ToMode0(imgDown, p, config)
				*/
				if err := overscan.SaveGo(filePath, bw, p, screenMode, ex); err != nil {
					fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
					return err
				}
			} else {
				if err := overscan.Overscan(filePath, bw, p, screenMode, ex); err != nil {
					fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
					return err
				}
			}
		} else {
			if err := overscan.EgxOverscan(filePath, bw, p, ex.EgxMode1, ex.EgxMode2, ex); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
				return err
			}
		}

	} else {
		if err := ocpartstudio.Scr(filePath, bw, p, screenMode, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := ocpartstudio.Loader(filePath, p, screenMode, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if !ex.CpcPlus {
		if err := ocpartstudio.Pal(filePath, p, screenMode, false, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := ex.OsFullPath(filePath, "_palettepal.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
		if err := palette.Ink(filePath, p, screenMode, false, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 = ex.OsFullPath(filePath, "_paletteink.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	} else {
		if err := palette.Kit(filePath, p, screenMode, false, ex); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := ex.OsFullPath(filePath, "_palettekit.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	}
	return ascii.Ascii(filePath, bw, p, false, ex)
}
