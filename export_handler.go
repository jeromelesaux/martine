package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
)

func ExportHandler() (*config.MartineConfig, constants.Size) {
	var size constants.Size
	cfg := config.NewMartineConfig(*picturePath, *output)

	if !*reverse {
		switch *mode {
		case 0:
			size = constants.Mode0
			if *overscan {
				size = constants.OverscanMode0
			}
		case 1:
			size = constants.Mode1
			if *overscan {
				size = constants.OverscanMode1
			}
		case 2:
			size = constants.Mode2
			if *overscan {
				size = constants.OverscanMode2
			}
		default:
			if *height == -1 && *width == -1 && !*deltaMode {
				fmt.Fprintf(os.Stderr, "mode %d not defined and no custom width or height\n", *mode)
				usage()
			}
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
			fmt.Fprintf(os.Stderr, "Max width allowed is (%d) your choice (%d), Quiting...\n", size.Width, constants.WidthMax)
			os.Exit(-1)
		}
		if size.Height > constants.HeightMax {
			fmt.Fprintf(os.Stderr, "Max height allowed is (%d) your choice (%d), Quiting...\n", size.Height, constants.HeightMax)
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
	cfg.TileMode = *tileMode
	cfg.RollMode = *rollMode
	cfg.RollIteration = *iterations
	cfg.NoAmsdosHeader = *noAmsdosHeader
	cfg.CpcPlus = *plusMode
	cfg.TileIterationX = *tileIterationX
	cfg.TileIterationY = *tileIterationY
	cfg.Compression = *compress
	cfg.RotationMode = *rotateMode
	cfg.Rotation3DMode = *rotate3dMode
	cfg.Rotation3DType = *rotate3dType
	cfg.Rotation3DX0 = *rotate3dX0
	cfg.Rotation3DY0 = *rotate3dY0
	cfg.M4Host = *m4Host
	cfg.M4RemotePath = *m4RemotePath
	cfg.M4Autoexec = *m4Autoexec
	cfg.ResizingAlgo = resizeAlgo
	cfg.DitheringMultiplier = *ditheringMultiplier
	cfg.DitheringWithQuantification = *withQuantization
	cfg.PalettePath = *palettePath
	cfg.InkPath = *inkPath
	cfg.KitPath = *kitPath
	cfg.RotationRlaBit = *rla
	cfg.RotationSraBit = *sra
	cfg.RotationSlaBit = *sla
	cfg.RotationRraBit = *rra
	cfg.RotationKeephighBit = *keephigh
	cfg.RotationKeeplowBit = *keeplow
	cfg.RotationLosthighBit = *losthigh
	cfg.RotationLostlowBit = *lostlow
	cfg.RotationIterations = *iterations
	cfg.Flash = *flash
	cfg.Sna = *sna
	cfg.SpriteHard = *spriteHard
	cfg.SplitRaster = *splitRasters
	cfg.ZigZag = *zigzag
	cfg.Animate = *doAnimation
	cfg.Reducer = *reducer
	cfg.Json = *jsonOutput
	cfg.Ascii = *txtOutput
	cfg.OneLine = *oneLine
	cfg.OneRow = *oneRow

	if *scanlineSequence != "" {
		sequence := strings.Split(*scanlineSequence, ",")
		for _, v := range sequence {
			line, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bad scanline sequence (%s) error:%v\n", *scanlineSequence, err)
				os.Exit(-1)
			}
			cfg.ScanlineSequence = append(cfg.ScanlineSequence, line)
		}
		modulo := size.Height % len(cfg.ScanlineSequence)
		if modulo != 0 {
			fmt.Fprintf(os.Stderr, "height modulo scanlinesequence is not equal to 0 %d lines and the output image lines is %d\n", len(cfg.ScanlineSequence), size.Height)
			os.Exit(-1)
		}
	}

	if err := cfg.ImportInkSwap(*inkSwap); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse inkswap option with error [%s]\n", err)
		os.Exit(-1)
	}
	if *lineWidth != "" {
		if err := cfg.SetLineWith(*lineWidth); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse linewidth option with error [%s]\n", err)
			os.Exit(-1)
		}
	}

	if *maskSprite != "" {

		v, err := common.ParseHexadecimal8(*maskSprite)
		if err == nil {
			cfg.MaskSprite = uint8(v)
		}
		if cfg.MaskSprite != 0 {
			if *maskOrOperation {
				cfg.MaskOrOperation = true
			}
			if *maskAdOperation {
				cfg.MaskAndOperation = true
			}
			if cfg.MaskAndOperation && cfg.MaskOrOperation {
				fmt.Fprintf(os.Stderr, "Or and And operations are setted, will only apply And operation.\n")
				cfg.MaskOrOperation = false
			}
			if !cfg.MaskAndOperation && !cfg.MaskOrOperation {
				fmt.Fprintf(os.Stderr, "Or and And operations are not setted, will only apply And operation.\n")
				cfg.MaskAndOperation = true
			}
			fmt.Fprintf(os.Stdout, "Applying sprite mask value [#%X] [%.8b] AND = %t, OR =%t\n",
				cfg.MaskSprite,
				cfg.MaskSprite,
				cfg.MaskAndOperation,
				cfg.MaskOrOperation)
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
	if cfg.M4Host != "" {
		cfg.M4 = true
	}

	if *egx1 {
		cfg.EgxFormat = config.Egx1Mode
	}
	if *egx2 {
		cfg.EgxFormat = config.Egx2Mode
	}
	if *mode != -1 {
		cfg.EgxMode1 = uint8(*mode)
	}
	if *mode2 != -1 {
		cfg.EgxMode2 = uint8(*mode2)
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
