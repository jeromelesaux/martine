package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
)

func ExportHandler() (*export.ExportType, constants.Size) {
	var size constants.Size
	exp := export.NewExportType(*picturePath, *output)

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
			exp.CustomDimension = true
			exp.Win = true
			size.Height = *height
			if *width != -1 {
				size.Width = *width
			} else {
				size.Width = 0
			}
		}
		if *width != -1 {
			exp.Win = true
			exp.CustomDimension = true
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

	exp.FilloutGif = *filloutGif
	exp.ExtendedDsk = *extendedDsk
	exp.TileMode = *tileMode
	exp.RollMode = *rollMode
	exp.RollIteration = *iterations
	exp.NoAmsdosHeader = *noAmsdosHeader
	exp.CpcPlus = *plusMode
	exp.TileIterationX = *tileIterationX
	exp.TileIterationY = *tileIterationY
	exp.Compression = *compress
	exp.RotationMode = *rotateMode
	exp.Rotation3DMode = *rotate3dMode
	exp.Rotation3DType = *rotate3dType
	exp.Rotation3DX0 = *rotate3dX0
	exp.Rotation3DY0 = *rotate3dY0
	exp.M4Host = *m4Host
	exp.M4RemotePath = *m4RemotePath
	exp.M4Autoexec = *m4Autoexec
	exp.ResizingAlgo = resizeAlgo
	exp.DitheringMultiplier = *ditheringMultiplier
	exp.DitheringWithQuantification = *withQuantization
	exp.PalettePath = *palettePath
	exp.InkPath = *inkPath
	exp.KitPath = *kitPath
	exp.RotationRlaBit = *rla
	exp.RotationSraBit = *sra
	exp.RotationSlaBit = *sla
	exp.RotationRraBit = *rra
	exp.RotationKeephighBit = *keephigh
	exp.RotationKeeplowBit = *keeplow
	exp.RotationLosthighBit = *losthigh
	exp.RotationLostlowBit = *lostlow
	exp.RotationIterations = *iterations
	exp.Flash = *flash
	exp.Sna = *sna
	exp.SpriteHard = *spriteHard
	exp.SplitRaster = *splitRasters
	exp.ZigZag = *zigzag
	exp.Animate = *doAnimation
	exp.Reducer = *reducer
	exp.Json = *jsonOutput
	exp.Ascii = *txtOutput
	exp.OneLine = *oneLine
	exp.OneRow = *oneRow

	if *scanlineSequence != "" {
		sequence := strings.Split(*scanlineSequence, ",")
		for _, v := range sequence {
			line, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bad scanline sequence (%s) error:%v\n", *scanlineSequence, err)
				os.Exit(-1)
			}
			exp.ScanlineSequence = append(exp.ScanlineSequence, line)
		}
		modulo := size.Height % len(exp.ScanlineSequence)
		if modulo != 0 {
			fmt.Fprintf(os.Stderr, "height modulo scanlinesequence is not equal to 0 %d lines and the output image lines is %d\n", len(exp.ScanlineSequence), size.Height)
			os.Exit(-1)
		}
	}

	if err := exp.ImportInkSwap(*inkSwap); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse inkswap option with error [%s]\n", err)
		os.Exit(-1)
	}
	if *lineWidth != "" {
		if err := exp.SetLineWith(*lineWidth); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse linewidth option with error [%s]\n", err)
			os.Exit(-1)
		}
	}

	if *maskSprite != "" {

		v, err := common.ParseHexadecimal8(*maskSprite)
		if err == nil {
			exp.MaskSprite = uint8(v)
		}
		if exp.MaskSprite != 0 {
			if *maskOrOperation {
				exp.MaskOrOperation = true
			}
			if *maskAdOperation {
				exp.MaskAndOperation = true
			}
			if exp.MaskAndOperation && exp.MaskOrOperation {
				fmt.Fprintf(os.Stderr, "Or and And operations are setted, will only apply And operation.\n")
				exp.MaskOrOperation = false
			}
			if !exp.MaskAndOperation && !exp.MaskOrOperation {
				fmt.Fprintf(os.Stderr, "Or and And operations are not setted, will only apply And operation.\n")
				exp.MaskAndOperation = true
			}
			fmt.Fprintf(os.Stdout, "Applying sprite mask value [#%X] [%.8b] AND = %t, OR =%t\n",
				exp.MaskSprite,
				exp.MaskSprite,
				exp.MaskAndOperation,
				exp.MaskOrOperation)
		}
	}

	if exp.CpcPlus {
		exp.Kit = true
		exp.Pal = false
	}
	exp.Overscan = *overscan
	if exp.Overscan {
		exp.Scr = false
		exp.Kit = true
	}
	if exp.M4Host != "" {
		exp.M4 = true
	}

	if *egx1 {
		exp.EgxFormat = export.Egx1Mode
	}
	if *egx2 {
		exp.EgxFormat = export.Egx2Mode
	}
	if *mode != -1 {
		exp.EgxMode1 = uint8(*mode)
	}
	if *mode2 != -1 {
		exp.EgxMode2 = uint8(*mode2)
	}

	exp.DeltaMode = *deltaMode
	exp.Dsk = *dsk
	return exp, size
}
