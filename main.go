package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/gfx"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

var (
	byteStatement   = flag.String("s", "", "Byte statement to replace in ascii export (default is BYTE), you can replace or instance by defb")
	picturePath     = flag.String("i", "", "Picture path of the input file.")
	width           = flag.Int("w", -1, "Custom output width in pixels.")
	height          = flag.Int("h", -1, "Custom output height in pixels.")
	mode            = flag.Int("m", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tand add -f option for overscan export.")
	output          = flag.String("o", "", "Output directory")
	overscan        = flag.Bool("f", false, "Overscan mode (default no overscan)")
	resizeAlgorithm = flag.Int("a", 1, "Algorithm to resize the image (available : \n\t1: NearestNeighbor (default)\n\t2: CatmullRom\n\t3: Lanczos\n\t4: Linear\n\t5: Box\n\t6: Hermite\n\t7: BSpline\n\t8: Hamming\n\t9: Hann\n\t10: Gaussian\n\t11: Blackman\n\t12: Bartlett\n\t13: Welch\n\t14: Cosine")
	help            = flag.Bool("help", false, "Display help message")
	noAmsdosHeader  = flag.Bool("n", false, "no amsdos header for all files (default amsdos header added).")
	plusMode        = flag.Bool("p", false, "Plus mode (means generate an image for CPC Plus Screen)")
	version         = "0.1.Alpha"
)

func usage() {
	fmt.Fprintf(os.Stdout, "martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)\n")
	fmt.Fprintf(os.Stdout, "By Impact Sid (Version:%s)\n", version)
	fmt.Fprintf(os.Stdout, "usage :\n\n")
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	var size gfx.Size
	var filename, extension string
	var customDimension bool
	var screenMode uint8
	flag.Parse()

	if *help {
		usage()
	}

	// picture path to convert
	if *picturePath == "" {
		usage()
	}
	filename = filepath.Base(*picturePath)
	extension = filepath.Ext(*picturePath)

	// output directory to store results
	if *output != "" {
		fi, err := os.Stat(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while getting directory informations :%v, Quiting\n", err)
			os.Exit(-2)
		}

		if !fi.IsDir() {
			fmt.Fprintf(os.Stderr, "%s is not a directory will store in current directory\n", *output)
			*output = "./"
		}
	} else {
		*output = "./"
	}

	if *mode == -1 {
		fmt.Fprintf(os.Stderr, "No output mode defined can not choose. Quiting\n")
		usage()
	}
	switch *mode {
	case 0:
		size = gfx.Mode0
		screenMode = 0
		if *overscan {
			size = gfx.OverscanMode0
		}

	case 1:
		size = gfx.Mode1
		screenMode = 1
		if *overscan {
			size = gfx.OverscanMode1
		}
	case 2:
		screenMode = 2
		size = gfx.Mode2
		if *overscan {
			size = gfx.OverscanMode2
		}
	default:
		if *height == -1 && *width == -1 {
			fmt.Fprintf(os.Stderr, "mode %d not defined and no custom width or height\n", *mode)
			usage()
		}
	}
	if *height != -1 {
		customDimension = true
		size.Height = *height
		if *width != -1 {
			size.Width = *width
		} else {
			size.Width = 0
		}
	}
	if *width != -1 {
		customDimension = true
		size.Width = *width
		if *height != -1 {
			size.Height = *height
		} else {
			size.Height = 0
		}
	}

	if *byteStatement != "" {
		gfx.ByteToken = *byteStatement
	}

	fmt.Fprintf(os.Stdout, "Informations :\n%s", size.ToString())

	f, err := os.Open(*picturePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file %s, error %v\n", *picturePath, err)
		os.Exit(-2)
	}
	defer f.Close()
	in, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot decode the image %s error %v", *picturePath, err)
		os.Exit(-2)
	}

	fmt.Fprintf(os.Stderr, "Filename :%s, extension:%s\n", filename, extension)

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
	default:
		resizeAlgo = imaging.NearestNeighbor
	}

	out := convert.Resize(in, size, resizeAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", *picturePath+"_resized.png")
	fwr, err := os.Create(*output + string(filepath.Separator) + *picturePath + "_resized.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create new image (%s) error %v\n", *picturePath+"_resized.png", err)
		os.Exit(-2)
	}
	if err := png.Encode(fwr, out); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create new image (%s) as png error %v\n", *picturePath+"_resized.png", err)
		fwr.Close()
		os.Exit(-2)
	}
	fwr.Close()

	newPalette, downgraded, err := convert.DowngradingPalette(out, size, *plusMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", *picturePath)
	}

	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", *picturePath+"_down.png")
	fwd, err := os.Create(*output + string(filepath.Separator) + *picturePath + "_down.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create new image (%s) error %v\n", *picturePath+"_down.png", err)
		os.Exit(-2)
	}

	if err := png.Encode(fwd, downgraded); err != nil {
		fwd.Close()
		fmt.Fprintf(os.Stderr, "Cannot create new image (%s) as png error %v\n", *picturePath+"_down.png", err)
		os.Exit(-2)
	}
	fwd.Close()
	if !customDimension {
		gfx.Transform(downgraded, newPalette, size, *picturePath, *output, *noAmsdosHeader, *plusMode)
	} else {
		fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
		gfx.SpriteTransform(downgraded, newPalette, size, screenMode, *picturePath, *output, *noAmsdosHeader, *plusMode)
	}
}
