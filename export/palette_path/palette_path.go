package palettepath

import (
	"image/color"

	"github.com/jeromelesaux/martine/constants"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/log"
)

func Open(p constants.PalettePath) (color.Palette, error) {
	var palette color.Palette
	var err error
	if p.OcpPath != "" {
		log.GetLogger().Info("Input palette to apply : (%s)\n", p.OcpPath)
		palette, _, err = ocpartstudio.OpenPal(p.OcpPath)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", p.OcpPath)
			return palette, err
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}
	if p.InkPath != "" {
		log.GetLogger().Info("Input palette to apply : (%s)\n", p.InkPath)
		palette, _, err = impPalette.OpenInk(p.InkPath)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", p.InkPath)
			return palette, err
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}
	if p.KitPath != "" {
		log.GetLogger().Info("Input plus palette to apply : (%s)\n", p.KitPath)
		palette, _, err = impPalette.OpenKit(p.KitPath)
		if err != nil {
			log.GetLogger().Error("Palette in file (%s) can not be read skipped\n", p.KitPath)
			return palette, err
		} else {
			log.GetLogger().Info("Use palette with (%d) colors \n", len(palette))
		}
	}
	return palette, nil
}
