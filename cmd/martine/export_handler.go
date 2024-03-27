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
			cfg.Win = true
			size.Height = *height
			if *width != -1 {
				size.Width = *width
			} else {
				size.Width = 0
			}
		}
		if *width != -1 {
			cfg.Win = true
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
	cfg.ExtendedDsk = *extendedDsk
	cfg.Transformation.TileMode = *tileMode
	cfg.Transformation.RollMode = *rollMode
	cfg.Transformation.RollIteration = *iterations
	cfg.NoAmsdosHeader = *noAmsdosHeader
	cfg.CpcPlus = *plusMode
	cfg.Transformation.TileIterationX = *tileIterationX
	cfg.Transformation.TileIterationY = *tileIterationY
	cfg.Compression = compression.ToCompressMethod(*compress)
	cfg.Transformation.RotationMode = *rotateMode
	cfg.Transformation.Rotation3DMode = *rotate3dMode
	cfg.Transformation.Rotation3DType = *rotate3dType
	cfg.Transformation.Rotation3DX0 = *rotate3dX0
	cfg.Transformation.Rotation3DY0 = *rotate3dY0
	cfg.M4.Host = *m4Host
	cfg.M4.RemotePath = *m4RemotePath
	cfg.M4.Autoexec = *m4Autoexec
	cfg.ResizingAlgo = resizeAlgo
	cfg.Dithering.DitheringMultiplier = *ditheringMultiplier
	cfg.Dithering.DitheringWithQuantification = *withQuantization
	cfg.PalettePath.OcpPath = *palettePath
	cfg.PalettePath.InkPath = *inkPath
	cfg.PalettePath.KitPath = *kitPath
	cfg.Transformation.RotationRlaBit = *rla
	cfg.Transformation.RotationSraBit = *sra
	cfg.Transformation.RotationSlaBit = *sla
	cfg.Transformation.RotationRraBit = *rra
	cfg.Transformation.RotationKeephighBit = *keephigh
	cfg.Transformation.RotationKeeplowBit = *keeplow
	cfg.Transformation.RotationLosthighBit = *losthigh
	cfg.Transformation.RotationLostlowBit = *lostlow
	cfg.Transformation.RotationIterations = *iterations
	cfg.Flash.Enabled = *flash
	cfg.Sna.Enabled = *sna
	cfg.SpriteHard = *spriteHard
	cfg.SplitRaster = *splitRasters
	cfg.ZigZag = *zigzag
	cfg.Animate = *doAnimation
	cfg.Reducer = *reducer
	cfg.Json = *jsonOutput
	cfg.Ascii = *txtOutput
	cfg.OneLine = *oneLine
	cfg.OneRow = *oneRow
	cfg.ExportAsGoFile = *exportGoFiles

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
			cfg.Mask.MaskSprite = v
		}
		if cfg.Mask.MaskSprite != 0 {
			if *maskOrOperation {
				cfg.Mask.OrOperation = true
			}
			if *maskAdOperation {
				cfg.Mask.AndOperation = true
			}
			if cfg.Mask.AndOperation && cfg.Mask.OrOperation {
				log.GetLogger().Error("Or and And operations are setted, will only apply And operation.\n")
				cfg.Mask.OrOperation = false
			}
			if !cfg.Mask.AndOperation && !cfg.Mask.OrOperation {
				log.GetLogger().Error("Or and And operations are not setted, will only apply And operation.\n")
				cfg.Mask.AndOperation = true
			}
			log.GetLogger().Info("Applying sprite mask value [#%X] [%.8b] AND = %t, OR =%t\n",
				cfg.Mask.MaskSprite,
				cfg.Mask.MaskSprite,
				cfg.Mask.AndOperation,
				cfg.Mask.OrOperation)
		}
	}

	if cfg.CpcPlus {
		cfg.Kit = true
		cfg.Pal = false
	}
	cfg.Overscan = *overscan
	if cfg.Overscan {
		cfg.Scr = false
		cfg.Kit = true
	}
	if cfg.M4.Host != "" {
		cfg.M4.Enabled = true
	}

	if *egx1 {
		cfg.Egx.EgxFormat = config.Egx1Mode
	}
	if *egx2 {
		cfg.Egx.EgxFormat = config.Egx2Mode
	}
	if *mode != -1 {
		cfg.Egx.EgxMode1 = uint8(*mode)
	}
	if *mode2 != -1 {
		cfg.Egx.EgxMode2 = uint8(*mode2)
	}

	if *saturationPal > 0 || *brightnessPal > 0 {
		cfg.CpcPlus = true
		cfg.Saturation = *saturationPal
		cfg.Brightness = *brightnessPal
	}

	cfg.DeltaMode = *deltaMode
	cfg.Dsk = *dsk
	return cfg, size
}
