package export

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/screen"

	ovs "github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
)

func ToMode0AndExport(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, cfg *config.MartineConfig) error {
	bw := screen.ToMode0(in, p, cfg)
	return Export(filePath, bw, p, 0, cfg)
}

func ToMode1AndExport(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, cfg *config.MartineConfig) error {
	bw := screen.ToMode1(in, p, cfg)
	return Export(filePath, bw, p, 1, cfg)
}

func ToMode2AndExport(in *image.NRGBA, p color.Palette, size constants.Size, filePath string, cfg *config.MartineConfig) error {
	bw := screen.ToMode2(in, p, cfg)
	return Export(filePath, bw, p, 2, cfg)
}

func Export(filePath string, bw []byte, p color.Palette, screenMode uint8, cfg *config.MartineConfig) error {
	if cfg.Overscan {
		if cfg.EgxFormat == 0 {
			if cfg.ExportAsGoFile {
				dataUp, dataDown, err := ovs.ToGo(bw, screenMode, p, cfg.CpcPlus)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while converting into .goN files file %s error :%v", filePath, err)
					return err
				}
				if err := overscan.SaveGo(filePath, dataUp, dataDown, p, screenMode, cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
					return err
				}
			} else {
				if err := overscan.Overscan(filePath, bw, p, screenMode, cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
					return err
				}
			}
		} else {
			if err := overscan.EgxOverscan(filePath, bw, p, cfg.EgxMode1, cfg.EgxMode2, cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
				return err
			}
		}

	} else {
		if err := ocpartstudio.Scr(filePath, bw, p, screenMode, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := ocpartstudio.Loader(filePath, p, screenMode, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}
	if !cfg.CpcPlus {
		if err := ocpartstudio.Pal(filePath, p, screenMode, false, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := cfg.OsFullPath(filePath, "_palettepal.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
		if err := palette.Ink(filePath, p, screenMode, false, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 = cfg.OsFullPath(filePath, "_paletteink.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	} else {
		if err := palette.Kit(filePath, p, screenMode, false, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := cfg.OsFullPath(filePath, "_palettekit.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file %s error :%v", filePath2, err)
			return err
		}
	}
	return ascii.Ascii(filePath, bw, p, false, cfg)
}
