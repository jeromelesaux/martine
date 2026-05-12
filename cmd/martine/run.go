package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	ci "github.com/jeromelesaux/martine/convert/image"
	"github.com/jeromelesaux/martine/convert/screen"
	covs "github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/diskimage"
	ovs "github.com/jeromelesaux/martine/export/impdraw/overscan"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/impdraw/tile"
	"github.com/jeromelesaux/martine/export/m4"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/snapshot"
	"github.com/jeromelesaux/martine/export/spritehard"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/gfx/effect"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/gfx/filter"
	gfxsprite "github.com/jeromelesaux/martine/gfx/sprite"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
)

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	in, _, err := image.Decode(f)
	return in, err
}

func setDithering(cfg *config.MartineConfig) error {
	switch *ditheringAlgo {
	case 0:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.FloydSteinberg
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 1:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.JarvisJudiceNinke
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 2:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Stucki
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 3:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Atkinson
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 4:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Sierra
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 5:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.SierraLite
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 6:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Sierra3
		cfg.ScrCfg.Process.Dithering.Type = constants.ErrorDiffusionDither
	case 7:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer2
		cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
	case 8:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer3
		cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
	case 9:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer4
		cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
	case 10:
		cfg.ScrCfg.Process.Dithering.Matrix = filter.Bayer8
		cfg.ScrCfg.Process.Dithering.Type = constants.OrderedDither
	default:
		return fmt.Errorf("dithering matrix not available")
	}
	return nil
}

func runDeltaPacking(cfg *config.MartineConfig, size constants.Size) error {
	screenAddress, err := common.ParseHexadecimal16(*initialAddress)
	cfg.ScrCfg.Size = size
	if err != nil {
		log.GetLogger().Error("Error while parsing (%s) use the starting address #C000, err : %v\n", *initialAddress, err)
		screenAddress = 0xC000
	}

	exportVersion := animate.DeltaExportV1
	if *deltaPacking2 {
		exportVersion = animate.DeltaExportV2
	}

	return animate.DeltaPacking(cfg.ScrCfg.InputPath, cfg, screenAddress, uint8(*mode), exportVersion)
}

func runSplitSpriteBoard(cfg *config.MartineConfig, size constants.Size, filename string, in image.Image) error {
	in, err := loadImage(*picturePath)
	if err != nil {
		return fmt.Errorf("Error while opening file %s, error %w", *picturePath, err)
	}

	img := image.NewNRGBA(image.Rect(0, 0, in.Bounds().Max.X, in.Bounds().Max.Y))
	draw.Draw(img, img.Bounds(), in, in.Bounds().Min, draw.Src)
	pal, _, err := ci.DowngradingPalette(img, constants.Size{ColorsAvailable: size.ColorsAvailable, Width: img.Bounds().Max.X, Height: img.Bounds().Max.Y}, cfg.ScrCfg.IsPlus)
	if err != nil {
		return fmt.Errorf("Cannot downgrade palette %s error %w", *picturePath, err)
	}

	outputSize := constants.Size{Width: size.Width, Height: size.Height}
	raw, _, err := gfxsprite.SplitBoardToSprite(img, pal, *spritesPerColumn, *spritesPerRow, uint8(*mode), *spriteHard, outputSize)
	if err != nil {
		return fmt.Errorf("Cannot split the sprite board %s error %w", *picturePath, err)
	}

	if err := impPalette.SaveKit(filepath.Join(cfg.ScrCfg.OutputPath, "SPRITES.KIT"), pal, !cfg.ScrCfg.NoAmsdosHeader); err != nil {
		return fmt.Errorf("Cannot export palette %s error %w", *picturePath, err)
	}

	if *spriteCompiled {
		spr := make([][]byte, 0)
		for _, v := range raw {
			spr = append(spr, v...)
		}
		diffs := animate.AnalyzeSpriteBoard(spr)
		var code string
		for idx, diff := range diffs {
			var routine string
			if cfg.ScrCfg.Type == config.SpriteHardFormat {
				routine = animate.ExportCompiledSpriteHard(diff)
			} else {
				log.GetLogger().Error("not yet implemented")
				os.Exit(-2)
			}
			code += fmt.Sprintf("spr_%.2d:\n", idx)
			code += routine
		}

		if err := amsdos.SaveStringOSFile(filepath.Join(cfg.ScrCfg.OutputPath, "compiled_sprites.asm"), code); err != nil {
			return fmt.Errorf("error while saving sprite file error %w", err)
		}
	}

	if *spriteOcpWin {
		for idxX, v := range raw {
			for idxY, v0 := range v {
				filename := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("L%.2dC%.2d.WIN", idxX, idxY))
				if err := window.Win(filename, v0, uint8(*mode), cfg.ScrCfg.Size.Width, cfg.ScrCfg.Size.Height, cfg.HasContainerExport(config.DskContainer), cfg); err != nil {
					log.GetLogger().Error("error while exporting sprites error %s\n", err.Error())
				}
			}
		}
	}

	if *spriteFlat {
		buf := make([]byte, 0)
		for _, v := range raw {
			for _, v0 := range v {
				buf = append(buf, v0...)
			}
		}
		filename := filepath.Join(cfg.ScrCfg.OutputPath, "SPRITES.BIN")
		buf, _ = compression.Compress(buf, cfg.ScrCfg.Compression)
		var err error
		if !cfg.ScrCfg.NoAmsdosHeader {
			err = amsdos.SaveAmsdosFile(filename, constants.WindowExtension, buf, 2, 0, 0x4000, 0x4000)
		} else {
			err = amsdos.SaveOSFile(filename, buf)
		}
		if err != nil {
			return fmt.Errorf("Error while saving flat sprites file error %w", err)
		}
		if cfg.HasContainerExport(config.DskContainer) {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				return fmt.Errorf("Cannot export to Imp-Catcher the image %s error %w", filename, err)
			}
		}
	}

	if *impCatcher {
		buf := make([]byte, 0)
		for _, v := range raw {
			for _, v0 := range v {
				buf = append(buf, v0...)
			}
		}
		filename := filepath.Join(cfg.ScrCfg.OutputPath, "sprites.imp")
		if err := tile.Imp(buf, uint(*spritesPerRow*(*spritesPerColumn)), uint(cfg.ScrCfg.Size.Width), uint(cfg.ScrCfg.Size.Height), uint(*mode), filename, cfg); err != nil {
			return fmt.Errorf("Cannot export to Imp-Catcher the image %s error %w", filename, err)
		}
		if cfg.HasContainerExport(config.DskContainer) {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				return fmt.Errorf("Cannot export to Imp-Catcher the image %s error %w", filename, err)
			}
		}
	}

	if *spriteHard {
		data := spritehard.SprImpdraw{}
		for _, v := range raw {
			for _, v0 := range v {
				sh := spritehard.SpriteHard{}
				copy(sh.Data[:], v0[:256])
				data.Data = append(data.Data, sh)
			}
		}
		filename := filepath.Join(cfg.ScrCfg.OutputPath, "sprites.spr")
		if err := spritehard.Spr(filename, data, cfg); err != nil {
			return fmt.Errorf("Cannot export to Imp-Catcher the image %s error %w", filename, err)
		}
		if cfg.HasContainerExport(config.DskContainer) {
			if err := diskimage.ImportInDsk(filename, cfg); err != nil {
				return fmt.Errorf("Cannot export to Imp-Catcher the image %s error %w", filename, err)
			}
		}
	}

	data := make([][]byte, 0)
	for _, v := range raw {
		data = append(data, v...)
	}
	header := fmt.Sprintf("' from file %s\n", cfg.ScrCfg.InputPath)
	code := header + ascii.SpritesHardText(data, cfg.ScrCfg.Compression)
	filename = filepath.Join(cfg.ScrCfg.OutputPath, "SPRITES.ASM")
	if err := amsdos.SaveStringOSFile(filename, code); err != nil {
		return fmt.Errorf("cannot save text data file error %w", err)
	}

	os.Exit(0)
	return nil
}

func runImpCatcherWildcard(cfg *config.MartineConfig, filename string) error {
	if !cfg.CustomDimension {
		return fmt.Errorf("You must set custom width and height.")
	}

	sprites := make([]byte, 0)
	log.GetLogger().Info("[%s]\n", *picturePath)
	spritesPaths, err := common.WilcardedFiles([]string{*picturePath})
	if err != nil {
		return fmt.Errorf("error while getting wildcard files %s error : %w", *picturePath, err)
	}

	for _, v := range spritesPaths {
		in, err := loadImage(v)
		if err != nil {
			return fmt.Errorf("Error while opening file %s, error %w", v, err)
		}
		if err := gfx.ApplyOneImageAndExport(in, cfg, filepath.Base(v), v, uint8(*mode)); err != nil {
			return fmt.Errorf("Cannot apply the image %s error %w", *picturePath, err)
		}
		spritePath := cfg.AmsdosFullPath(v, constants.WindowExtension)
		data, err := window.RawWin(spritePath)
		if err != nil {
			log.GetLogger().Error("Error while extracting raw content, err:%s\n", err)
		}
		sprites = append(sprites, data...)
	}

	finalFile := strings.ReplaceAll(filename, "?", "")
	if err := tile.Imp(sprites, uint(len(spritesPaths)), uint(cfg.ScrCfg.Size.Width), uint(cfg.ScrCfg.Size.Height), uint(*mode), finalFile, cfg); err != nil {
		return fmt.Errorf("Cannot export to Imp-Catcher the image %s error %w", *picturePath, err)
	}

	os.Exit(0)
	return nil
}

func runReverse(cfg *config.MartineConfig, filename, extension string) error {
	outpath := filepath.Join(*output, strings.Replace(strings.ToLower(filename), ".scr", ".png", 1))
	if cfg.ScrCfg.Type == config.FullscreenFormat {
		p, mode, err := ovs.OverscanPalette(*picturePath)
		if err != nil {
			return fmt.Errorf("Cannot get the palette from file (%s) error %w", *picturePath, err)
		}

		if err := covs.OverscanToPng(*picturePath, outpath, mode, p); err != nil {
			return fmt.Errorf("Cannot convert to PNG file (%s) error %w", *picturePath, err)
		}
		os.Exit(1)
		return nil
	}
	if *mode == -1 {
		return fmt.Errorf("Mode is mandatory to convert to PNG")
	}

	var p color.Palette
	var err error
	if *palettePath != "" && !*plusMode {
		p, _, err = ocpartstudio.OpenPal(*palettePath)
		if err != nil {
			return fmt.Errorf("Cannot open palette file (%s) error %w", *palettePath, err)
		}
	} else {
		if *kitPath != "" && *plusMode {
			p, _, err = impPalette.OpenKit(*kitPath)
			if err != nil {
				return fmt.Errorf("Cannot open kit file (%s) error %w", *kitPath, err)
			}
		} else {
			return fmt.Errorf("For screen or window image, pal or kit file palette is mandatory. (kit file must be associated with -p option)")
		}
	}

	switch strings.ToUpper(filepath.Ext(filename)) {
	case constants.WindowExtension:
		if err := sprite.SpriteToPng(*picturePath, outpath, uint8(*mode), p); err != nil {
			return fmt.Errorf("Cannot convert to PNG file (%s) error %w", *picturePath, err)
		}
	case constants.ScrExtension:
		if err := screen.ScrToPng(*picturePath, outpath, uint8(*mode), p); err != nil {
			return fmt.Errorf("Cannot convert to PNG file (%s) error %w", *picturePath, err)
		}
	}

	os.Exit(1)
	return nil
}

func runAnimation(cfg *config.MartineConfig) error {
	if !cfg.CustomDimension {
		return fmt.Errorf("You must set sprite dimensions with option -w and -h (mandatory)")
	}

	log.GetLogger().Info("animation output.\n")
	files := []string{*picturePath}
	files, err := common.WilcardedFiles(files)
	if err != nil {
		return fmt.Errorf("Cannot parse wildcard in argument (%s) error %w", *picturePath, err)
	}
	return animate.Animation(files, uint8(*mode), cfg)
}

func runDeltaMode(cfg *config.MartineConfig, filename string) error {
	log.GetLogger().Info("delta files to proceed.\n")
	for i, v := range deltaFiles {
		log.GetLogger().Info("[%d]:%s\n", i, v)
	}
	screenAddress, err := common.ParseHexadecimal16(*initialAddress)
	if err != nil {
		log.GetLogger().Error("Error while parsing (%s) use the starting address #C000, err : %v\n", *initialAddress, err)
		screenAddress = 0xC000
	}
	if *mode == -1 {
		return fmt.Errorf("You must set the mode for this feature. (option -m)")
	}
	return transformation.ProceedDelta(deltaFiles, screenAddress, cfg, uint8(*mode))
}

func runAnalyzeTilemap(cfg *config.MartineConfig, size constants.Size, filename string, in image.Image) error {
	var criteria common.AnalyseTilemapOption
	switch *analyzeTilemap {
	case string(common.SizeTilemapOption):
		criteria = common.SizeTilemapOption
		log.GetLogger().Info("go to analyze by size\n")
	case string(common.NumberTilemapOption):
		criteria = common.NumberTilemapOption
		log.GetLogger().Info("search for the lower number of tiles\n")
	default:
		return fmt.Errorf("Error tilemap analyze option not found : choose between (%s,%s)", string(common.SizeTilemapOption), string(common.NumberTilemapOption))
	}
	return gfx.AnalyzeTilemap(uint8(*mode), *plusMode, filename, *picturePath, in, cfg, criteria)
}

func runTilemap(cfg *config.MartineConfig, size constants.Size, filename string, in image.Image) error {
	return gfx.Tilemap(uint8(*mode), filename, *picturePath, size, in, cfg, nil)
}

func runTileMode(cfg *config.MartineConfig) error {
	if cfg.TileIterationX == -1 || cfg.TileIterationY == -1 {
		return fmt.Errorf("missing arguments iterx and itery to use with tile mode.")
	}
	return transformation.TileMode(cfg, uint8(*mode), cfg.TileIterationX, cfg.TileIterationY)
}

func runFlashOrApply(cfg *config.MartineConfig, size constants.Size, filename, extension string, in image.Image) error {
	if *flash {
		return effect.Flash(*picturePath, *picturePath2, *palettePath, *palettePath2, *mode, *mode2, cfg)
	}

	var p color.Palette
	var err error
	if cfg.ScrCfg.IsPlus {
		p, _, err = impPalette.OpenKit(*kitPath)
		if *kitPath != "" {
			if err != nil {
				return fmt.Errorf("Error while reading kit file (%s) :%w", *kitPath, err)
			}
		}
	} else if *palettePath != "" {
		p, _, err = ocpartstudio.OpenPal(*palettePath)
		if err != nil {
			return fmt.Errorf("Error while reading palette file (%s) :%w", *palettePath, err)
		}
	}

	if cfg.ScrCfg.Type == config.Egx1Format || cfg.ScrCfg.Type == config.Egx2Format {
		if len(p) == 0 {
			return fmt.Errorf("No colors found in palette, give up treatment.")
		}
		return effect.Egx(*picturePath, *picturePath2, p, *mode, *mode2, cfg)
	}

	if cfg.SplitRaster {
		if cfg.ScrCfg.Type.IsFullScreen() {
			return effect.DoSpliteRaster(in, uint8(*mode), filename, cfg)
		}
		log.GetLogger().Error("Only overscan mode implemented for this feature, %v", errors.ErrorNotYetImplemented)
		return nil
	}

	if strings.EqualFold(extension, constants.ScrExtension) {
		return fmt.Errorf("Error while applying on one image : SCR format not used for this treatment")
	}

	return gfx.ApplyOneImageAndExport(in, cfg, filename, *picturePath, uint8(*mode))
}

func runConversion(cfg *config.MartineConfig, size constants.Size, filename, extension string, in image.Image) error {
	if *deltaPacking || *deltaPacking2 {
		return runDeltaPacking(cfg, size)
	}

	if !*reverse {
		log.GetLogger().Info("Informations :\n%s", size.ToString())
	}

	if *ditheringAlgo != -1 {
		if err := setDithering(cfg); err != nil {
			return err
		}
	}

	if *splitSpriteBoard {
		return runSplitSpriteBoard(cfg, size, filename, in)
	}
	if *impCatcher {
		return runImpCatcherWildcard(cfg, filename)
	}
	if *reverse {
		return runReverse(cfg, filename, extension)
	}
	if cfg.Animate {
		return runAnimation(cfg)
	}
	if cfg.DeltaMode {
		return runDeltaMode(cfg, filename)
	}
	if *analyzeTilemap != "" {
		return runAnalyzeTilemap(cfg, size, filename, in)
	}
	if *tileMap {
		return runTilemap(cfg, size, filename, in)
	}
	if cfg.TileMode {
		return runTileMode(cfg)
	}
	return runFlashOrApply(cfg, size, filename, extension, in)
}

func runContainerExports(cfg *config.MartineConfig, picturePath, output string) error {
	if cfg.HasContainerExport(config.DskContainer) {
		if err := diskimage.ImportInDsk(picturePath, cfg); err != nil {
			return fmt.Errorf("Cannot create or write into dsk file error :%w", err)
		}
	}
	if cfg.HasContainerExport(config.SnaContainer) {
		if cfg.ScrCfg.Type.IsFullScreen() {
			var gfxFile string
			for _, v := range cfg.DskFiles {
				if filepath.Ext(v) == constants.ScrExtension {
					gfxFile = v
					break
				}
			}
			cfg.ContainerCfg.Path = filepath.Join(output, "test.sna")
			if err := snapshot.ImportInSna(gfxFile, cfg.ContainerCfg.Path, uint8(*mode)); err != nil {
				return fmt.Errorf("Cannot create or write into sna file error :%w", err)
			}
			log.GetLogger().Info("Sna saved in file %s\n", cfg.ContainerCfg.Path)
		} else {
			return fmt.Errorf("Feature not implemented for this file")
		}
	}
	if cfg.M4cfg.Enabled {
		if err := m4.ImportInM4(cfg); err != nil {
			return fmt.Errorf("Cannot send to M4 error :%w", err)
		}
	}
	return nil
}
