// Package export provides export utilities for martine output files.
package export

import (
	"image"
	"image/color"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/screen"
	co "github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/log"

	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/png"
)

func ExportPalette(filePath string, bw []byte, p color.Palette, screenMode uint8, dontImportDsk bool, cfg *config.MartineConfig) error {
	if !cfg.ScrCfg.IsPlus {
		if err := ocpartstudio.Pal(filePath, p, screenMode, dontImportDsk, cfg); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := cfg.OsFullPath(filePath, "_palettepal.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath2, err)
			return err
		}
		if err := palette.Ink(filePath, p, screenMode, dontImportDsk, cfg); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 = cfg.OsFullPath(filePath, "_paletteink.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath2, err)
			return err
		}
	} else {
		if err := palette.Kit(filePath, p, screenMode, dontImportDsk, cfg); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
			return err
		}
		filePath2 := cfg.OsFullPath(filePath, "_palettekit.png")
		if err := png.PalToPng(filePath2, p); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath2, err)
			return err
		}
	}
	return ascii.ASCII(filePath, bw, p, dontImportDsk, cfg)
}

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

// Export writes the screen data and palette to the appropriate output files.
// nolint: funlen, gocognit, gocyclo
func Export(filePath string, bw []byte, p color.Palette, screenMode uint8, cfg *config.MartineConfig) error {
	switch cfg.ScrCfg.Type {
	case config.FullscreenFormat:
		if err := overscan.Overscan(filePath, bw, p, screenMode, cfg); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
			return err
		}
	case config.Egx1FullscreenFormat, config.Egx2FullscreenFormat:
		if err := overscan.EgxOverscan(filePath, bw, p, cfg.EgxMode1, cfg.EgxMode2, cfg); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
			return err
		}
		if cfg.ScrCfg.IsExport(config.GoImpdrawExport) {
			data, err := co.ToGo(bw, screenMode, p, cfg)
			if err != nil {
				log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
				return err
			}
			if err := overscan.SaveGo(filePath, data, p, screenMode, cfg); err != nil {
				log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
				return err
			}
		}
	default:
		if err := ocpartstudio.Scr(filePath, bw, p, screenMode, cfg); err != nil {
			log.GetLogger().Error("Error while saving file %s error :%v", filePath, err)
			return err
		}
		if err := ocpartstudio.Loader(filePath, p, screenMode, cfg); err != nil {
			log.GetLogger().Error("Error while saving the loader %s with error %v\n", filePath, err)
			return err
		}
	}

	return ExportPalette(filePath, bw, p, screenMode, false, cfg)
}
