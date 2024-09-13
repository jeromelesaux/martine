package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/log"
)

// nolint: funlen, gocognit
func ExportHandler() (*config.MartineConfig, constants.Size) {
	var size constants.Size
	cfg := config.NewMartineConfig(*picturePath, *output)
	size = constants.NewSizeMode(uint8(*mode), *overscan)
	if !*reverse {

		emptySize := constants.Size{}
		if size == emptySize && *height == -1 && *width == -1 && !*deltaMode {
			log.GetLogger().Error("mode %d not defined and no custom width or height\n", *mode)
			usage()
		}

		if *height != -1 {
			cfg.CustomDimension = true
			cfg.ScrCfg.Type = config.WindowFormat
			size.Height = *height
			if *width != -1 {
				size.Width = *width
			} else {
				size.Width = 0
			}
		}
		if *width != -1 {
			cfg.ScrCfg.Type = config.WindowFormat
			cfg.CustomDimension = true
			size.Width = *width
			if *height != -1 {
				size.Height = *height
			} else {
				size.Height = 0
			}
		}

		if size.Width > constants.WidthMax {
			log.GetLogger().Error("Max width allowed is (%d) your choice (%d), Quiting...\n", size.Width, constants.WidthMax)
			os.Exit(-1)
		}
		if size.Height > constants.HeightMax {
			log.GetLogger().Error("Max height allowed is (%d) your choice (%d), Quiting...\n", size.Height, constants.HeightMax)
			os.Exit(-1)
		}
	}

	var resizeAlgo imaging.ResampleFilter
	switch *resizeAlgorithm {
	case 1:
		resizeAlgo = imaging.NearestNeighbor
	case 2:
		resizeAlgo = imaging.CatmullRom
	case 3:
		resizeAlgo = imaging.Lanczos
	case 4:
		resizeAlgo = imaging.Linear
	case 5:
		resizeAlgo = imaging.Box
	case 6:
		resizeAlgo = imaging.Hermite
	case 7:
		resizeAlgo = imaging.BSpline
	case 8:
		resizeAlgo = imaging.Hamming
	case 9:
		resizeAlgo = imaging.Hann
	case 10:
		resizeAlgo = imaging.Gaussian
	case 11:
		resizeAlgo = imaging.Blackman
	case 12:
		resizeAlgo = imaging.Bartlett
	case 13:
		resizeAlgo = imaging.Welch
	case 14:
		resizeAlgo = imaging.Cosine
	case 15:
		resizeAlgo = imaging.MitchellNetravali
	default:
		resizeAlgo = imaging.NearestNeighbor
	}

	cfg.FilloutGif = *filloutGif
	if extendedDsk != nil && *extendedDsk {
		cfg.ContainerCfg.AddExport(config.ExtendedDskContainer)
	}
	cfg.TileMode = *tileMode

	cfg.ScrCfg.NoAmsdosHeader = *noAmsdosHeader
	cfg.ScrCfg.IsPlus = *plusMode
	cfg.TileIterationX = *tileIterationX
	cfg.TileIterationY = *tileIterationY
	cfg.ScrCfg.Compression = compression.ToCompressMethod(*compress)
	cfg.ScrCfg.Process.ResizingAlgo = resizeAlgo
	cfg.ScrCfg.Process.DitheringMultiplier = *ditheringMultiplier
	cfg.ScrCfg.Process.DitheringWithQuantification = *withQuantization
	cfg.M4cfg = config.M4Config{
		Host:       *m4Host,
		RemotePath: *m4RemotePath,
		Autoexec:   *m4Autoexec,
		Enabled:    true,
	}

	var ppath string
	var ptype config.PaletteType
	if palettePath != nil {
		ppath = *palettePath
		ptype = config.PalPalette
	}
	if inkPath != nil {
		ppath = *inkPath
		ptype = config.InkPalette
	}
	if kitPath != nil {
		ppath = *kitPath
		ptype = config.KitPalette
	}
	cfg.PalCfg = config.PaletteConfig{
		Path: ppath,
		Type: ptype,
	}
	cfg.RotateCfg = config.RotateConfig{
		RollMode:            *rollMode,
		RollIteration:       *iterations,
		RotationRlaBit:      *rla,
		RotationSlaBit:      *sla,
		RotationRraBit:      *rra,
		RotationKeephighBit: *keephigh,
		RotationKeeplowBit:  *keeplow,
		RotationLosthighBit: *losthigh,
		RotationLostlowBit:  *lostlow,
		RotationIterations:  *iterations,
		RotationMode:        *rotateMode,
		Rotation3DMode:      *rotate3dMode,
		Rotation3DType:      config.Rotation3DType(*rotate3dType),
		Rotation3DX0:        *rotate3dX0,
		Rotation3DY0:        *rotate3dY0,
	}

	cfg.Flash = *flash
	if sna != nil && *sna {
		cfg.ContainerCfg.AddExport(config.SnaContainer)
	}
	if spriteHard != nil && *spriteHard {
		cfg.ScrCfg.Type = config.SpriteHardFormat
	}

	cfg.SplitRaster = *splitRasters
	cfg.ZigZag = *zigzag
	cfg.Animate = *doAnimation
	cfg.ScrCfg.Process.Reducer = *reducer
	if jsonOutput != nil && *jsonOutput {
		cfg.ScrCfg.AddExport(config.JsonExport)
	}
	if txtOutput != nil && *txtOutput {
		cfg.ScrCfg.AddExport(config.AssemblyExport)
	}
	if exportGoFiles != nil && *exportGoFiles {
		cfg.ScrCfg.AddExport(config.GoImpdrawExport)
	}
	cfg.ScrCfg.Process.OneLine = *oneLine
	cfg.ScrCfg.Process.OneRow = *oneRow

	if *scanlineSequence != "" {
		sequence := strings.Split(*scanlineSequence, ",")
		for _, v := range sequence {
			line, err := strconv.Atoi(v)
			if err != nil {
				log.GetLogger().Error("Bad scanline sequence (%s) error:%v\n", *scanlineSequence, err)
				os.Exit(-1)
			}
			cfg.ScanlineSequence = append(cfg.ScanlineSequence, line)
		}
		modulo := size.Height % len(cfg.ScanlineSequence)
		if modulo != 0 {
			log.GetLogger().Error("height modulo scanlinesequence is not equal to 0 %d lines and the output image lines is %d\n", len(cfg.ScanlineSequence), size.Height)
			os.Exit(-1)
		}
	}

	if err := cfg.ImportInkSwap(*inkSwap); err != nil {
		log.GetLogger().Error("Cannot parse inkswap option with error [%s]\n", err)
		os.Exit(-1)
	}
	if *lineWidth != "" {
		if err := cfg.SetLineWith(*lineWidth); err != nil {
			log.GetLogger().Error("Cannot parse linewidth option with error [%s]\n", err)
			os.Exit(-1)
		}
	}

	if *maskSprite != "" {

		v, err := common.ParseHexadecimal8(*maskSprite)
		if err == nil {
			cfg.ScrCfg.Process.MaskSprite = v
		}
		if cfg.ScrCfg.Process.MaskSprite != 0 {
			if *maskOrOperation {
				cfg.ScrCfg.Process.MaskOrOperation = true
			}
			if *maskAdOperation {
				cfg.ScrCfg.Process.MaskAndOperation = true
			}
			if cfg.ScrCfg.Process.MaskAndOperation && cfg.ScrCfg.Process.MaskOrOperation {
				log.GetLogger().Error("Or and And operations are setted, will only apply And operation.\n")
				cfg.ScrCfg.Process.MaskOrOperation = false
			}
			if !cfg.ScrCfg.Process.MaskAndOperation && !cfg.ScrCfg.Process.MaskOrOperation {
				log.GetLogger().Error("Or and And operations are not setted, will only apply And operation.\n")
				cfg.ScrCfg.Process.MaskAndOperation = true
			}
			log.GetLogger().Info("Applying sprite mask value [#%X] [%.8b] AND = %t, OR =%t\n",
				cfg.ScrCfg.Process.MaskSprite,
				cfg.ScrCfg.Process.MaskSprite,
				cfg.ScrCfg.Process.MaskAndOperation,
				cfg.ScrCfg.Process.MaskOrOperation)
		}
	}

	if cfg.ScrCfg.IsPlus {
		cfg.PalCfg.Type = config.KitPalette
	}
	if overscan != nil && *overscan {
		cfg.ScrCfg.Type = config.FullscreenFormat
	}
	if cfg.ScrCfg.Type == config.FullscreenFormat {
		cfg.PalCfg.Type = config.KitPalette
	}

	if *egx1 {
		cfg.ScrCfg.Type = config.Egx1Format
	}
	if *egx2 {
		cfg.ScrCfg.Type = config.Egx2Format
	}
	if *mode != -1 {
		cfg.EgxMode1 = uint8(*mode)
	}
	if *mode2 != -1 {
		cfg.EgxMode2 = uint8(*mode2)
	}

	if *saturationPal > 0 || *brightnessPal > 0 {
		cfg.ScrCfg.IsPlus = true
		cfg.ScrCfg.Process.Saturation = *saturationPal
		cfg.ScrCfg.Process.Brightness = *brightnessPal
	}

	cfg.DeltaMode = *deltaMode
	if dsk != nil && *dsk {
		cfg.ContainerCfg.AddExport(config.DskContainer)
	}
	return cfg, size
}
