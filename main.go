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
	"strings"
)

var (
	picturePath     = flag.String("p", "", "Picture path of the Amsdos file.")
	width           = flag.Int("w", -1, "Custom output width in pixels.")
	height          = flag.Int("h", -1, "Custom output height in pixels.")
	mode            = flag.String("m", "", "Output mode to use (mode0,mode1,mode2 or overscan available).")
	output          = flag.String("o", "", "Output directory")
	overscan        = flag.Bool("f", false, "Overscan mode (default no overscan)")
	resizeAlgorithm = flag.Int("a", 1, "Algorithm to resize the image (available 1: NearestNeighbor (default), 2: CatmullRom, 3: Lanczos, 4: Linear)")
	version         = "0.1Beta"
)

func usage() {
	fmt.Fprintf(os.Stdout, "martine to convert image to Amstrad cpc screen (even overscan)\n")
	fmt.Fprintf(os.Stdout, "By Impact Sid (Version:%s)\n", version)
	fmt.Fprintf(os.Stdout, "usage :\n\n")
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	var size gfx.Size
	var filename, extension string
	flag.Parse()
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
	}
	if *height != -1 && *width != -1 {
		fmt.Fprintf(os.Stderr, "Use the custom informations (width:%d, height:%d)\n", *width, *height)
		size.Height = *height
		size.Width = *width
	} else {
		if *mode == "" {
			fmt.Fprintf(os.Stderr, "No output mode defined can not choose. Quiting\n")
			usage()
		}
		switch strings.ToLower(*mode) {
		case "mode0":
			size = gfx.Mode0
			if *overscan {
				size = gfx.OverscanMode0
			}

		case "mode1":
			size = gfx.Mode1
			if *overscan {
				size = gfx.OverscanMode1
			}
		case "mode2":
			size = gfx.Mode2
			if *overscan {
				size = gfx.OverscanMode2
			}
		default:
			fmt.Fprintf(os.Stderr, "mode %s not defined\n", *mode)
			usage()
		}
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
	default:
		resizeAlgo = imaging.NearestNeighbor
	}

	out := convert.Resize(in, size, resizeAlgo)
	fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", *picturePath+"_resized.png")
	fwr, err := os.Create(*picturePath + "_resized.png")
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

	newPalette, downgraded, err := convert.DowngradingPalette(out, size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", *picturePath)
	}

	fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", *picturePath+"_down.png")
	fwd, err := os.Create(*picturePath + "_down.png")
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

	gfx.Transform(downgraded, newPalette, size, *picturePath)
}
