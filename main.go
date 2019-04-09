package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/gfx"
)

var (
	byteStatement   = flag.String("s", "", "Byte statement to replace in ascii export (default is BYTE), you can replace or instance by defb")
	picturePath     = flag.String("i", "", "Picture path of the input file.")
	width           = flag.Int("w", -1, "Custom output width in pixels.")
	height          = flag.Int("h", -1, "Custom output height in pixels.")
	mode            = flag.Int("m", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tand add -f option for overscan export.\n\t")
	output          = flag.String("o", "", "Output directory")
	overscan        = flag.Bool("f", false, "Overscan mode (default no overscan)")
	resizeAlgorithm = flag.Int("a", 1, "Algorithm to resize the image (available : \n\t1: NearestNeighbor (default)\n\t2: CatmullRom\n\t3: Lanczos\n\t4: Linear\n\t5: Box\n\t6: Hermite\n\t7: BSpline\n\t8: Hamming\n\t9: Hann\n\t10: Gaussian\n\t11: Blackman\n\t12: Bartlett\n\t13: Welch\n\t14: Cosine\n\t")
	help            = flag.Bool("help", false, "Display help message")
	noAmsdosHeader  = flag.Bool("n", false, "no amsdos header for all files (default amsdos header added).")
	plusMode        = flag.Bool("p", false, "Plus mode (means generate an image for CPC Plus Screen)")
	rollMode        = flag.Bool("roll", false, "Roll mode allow to walk and walk into the input file, associated with rla,rra,sra,sla, keephigh, keeplow, losthigh or lostlow options.")
	iterations      = flag.Int("iter", -1, "Iterations number to walk in roll mode, or number of images to generate in rotation mode.")
	rra             = flag.Int("rra", -1, "bit rotation on the right and keep pixels")
	rla             = flag.Int("rla", -1, "bit rotation on the left and keep pixels")
	sra             = flag.Int("sra", -1, "bit rotation on the right and lost pixels")
	sla             = flag.Int("sla", -1, "bit rotation on the left and lost pixels")
	losthigh        = flag.Int("losthigh", -1, "bit rotation on the top and lost pixels")
	lostlow         = flag.Int("lostlow", -1, "bit rotation on the bottom and lost pixels")
	keephigh        = flag.Int("keephigh", -1, "bit rotation on the top and keep pixels")
	keeplow         = flag.Int("keeplow", -1, "bit rotation on the bottom and keep pixels")
	palettePath     = flag.String("pal", "", "Apply the input palette to the image")
	info            = flag.Bool("info", false, "Return the information of the file, associated with -pal and -win options")
	winPath         = flag.String("win", "", "Filepath of the ocp win file")
	dsk             = flag.Bool("dsk", false, "Copy files in a new CPC image Dsk.")
	tileMode        = flag.Bool("tile", false, "Tile mode to create multiples sprites from a same image.")
	tileIterationX  = flag.Int("iterx", -1, "Number of tiles on a row in the input image.")
	tileIterationY  = flag.Int("itery", -1, "Number of tiles on a column in the input image.")
	compress        = flag.Int("z", -1, "Compression algorithm : \n\t1: rle (default)\n\t2: rle 16bits\n\t3: Lz4 Classic\n\t4: Lz4 Raw\n")
	kitPath         = flag.String("kit", "", "Path of the palette Cpc plus Kit file.")
	inkPath         = flag.String("ink", "", "Path of the palette Cpc ink file.")
	rotateMode      = flag.Bool("rotate", false, "Allow rotation on the input image, the input image must be a square (width equals height)")
	version         = "0.13"
)

func usage() {
	fmt.Fprintf(os.Stdout, "martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)\n")
	fmt.Fprintf(os.Stdout, "By Impact Sid (Version:%s)\n", version)
	fmt.Fprintf(os.Stdout, "Special thanks to @Ast (for his support), @Siko and @Tronic for ideas\n")
	fmt.Fprintf(os.Stdout, "usage :\n\n")
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	var size constants.Size
	var filename, extension string
	var customDimension bool
	var screenMode uint8
	var palette color.Palette
	flag.Parse()

	if *help {
		usage()
	}

	if *info {
		if *palettePath != "" {
			gfx.PalInformation(*palettePath)
		}
		if *winPath != "" {
			gfx.WinInformation(*winPath)
		}
		if *kitPath != "" {
			gfx.KitInformation(*kitPath)
		}
		if *inkPath != "" {
			gfx.InkInformation(*inkPath)
		}
		os.Exit(0)
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

	exportType := gfx.NewExportType(*picturePath, *output)

	if *mode == -1 {
		fmt.Fprintf(os.Stderr, "No output mode defined can not choose. Quiting\n")
		usage()
	}
	switch *mode {
	case 0:
		size = constants.Mode0
		screenMode = 0
		if *overscan {
			size = constants.OverscanMode0
		}
	case 1:
		size = constants.Mode1
		screenMode = 1
		if *overscan {
			size = constants.OverscanMode1
		}
	case 2:
		screenMode = 2
		size = constants.Mode2
		if *overscan {
			size = constants.OverscanMode2
		}
	default:
		if *height == -1 && *width == -1 {
			fmt.Fprintf(os.Stderr, "mode %d not defined and no custom width or height\n", *mode)
			usage()
		}
	}
	if *height != -1 {
		customDimension = true
		exportType.Win = true
		size.Height = *height
		if *width != -1 {
			size.Width = *width
		} else {
			size.Width = 0
		}
	}
	if *width != -1 {
		exportType.Win = true
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
	exportType.Size = size
	exportType.TileMode = *tileMode
	exportType.RollMode = *rollMode
	exportType.RollIteration = *iterations
	exportType.NoAmsdosHeader = *noAmsdosHeader
	exportType.CpcPlus = *plusMode
	exportType.TileIterationX = *tileIterationX
	exportType.TileIterationY = *tileIterationY
	exportType.Compression = *compress
	exportType.RotationMode = *rotateMode
	if exportType.CpcPlus {
		exportType.Kit = true
		exportType.Pal = false
	}
	exportType.Overscan = *overscan
	if exportType.Overscan {
		exportType.Scr = false
		exportType.Kit = true
	}
	exportType.Dsk = *dsk

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

	if ! customDimension && *rotateMode {
		size.Width = in.Bounds().Max.X
		size.Height = in.Bounds().Max.Y
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

	if exportType.TileMode {
		if exportType.TileIterationX == -1 || exportType.TileIterationY == -1 {
			fmt.Fprintf(os.Stderr, "missing arguments iterx and itery to use with tile mode.\n")
			usage()
			os.Exit(-1)
		}
		err := gfx.TileMode(exportType, uint8(*mode), exportType.TileIterationX, exportType.TileIterationY, resizeAlgo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Tile mode on error : error :%v\n", err)
			os.Exit(-1)
		}
	} else {

		if *palettePath != "" {
			fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", *palettePath)
			palette, _, err = gfx.OpenPal(*palettePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", *palettePath)
			} else {
				fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
			}
		}
		if *inkPath != "" {
			fmt.Fprintf(os.Stdout, "Input palette to apply : (%s)\n", *inkPath)
			palette, _, err = gfx.OpenInk(*inkPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", *inkPath)
			} else {
				fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
			}
		}
		if *kitPath != "" {
			fmt.Fprintf(os.Stdout, "Input plus palette to apply : (%s)\n", *kitPath)
			palette, _, err = gfx.OpenKit(*kitPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Palette in file (%s) can not be read skipped\n", *palettePath)
			} else {
				fmt.Fprintf(os.Stdout, "Use palette with (%d) colors \n", len(palette))
			}
		}

		out := convert.Resize(in, size, resizeAlgo)
		fmt.Fprintf(os.Stdout, "Saving resized image into (%s)\n", filename+"_resized.png")
		if err := gfx.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_resized.png", out); err != nil {
			os.Exit(-2)
		}

		var newPalette color.Palette
		var downgraded *image.NRGBA
		if len(palette) > 0 {
			newPalette, downgraded = convert.DowngradingWithPalette(out, palette)
		} else {
			newPalette, downgraded, err = convert.DowngradingPalette(out, size, exportType.CpcPlus)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", *picturePath)
			}

		}
		fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
		if err := gfx.Png(exportType.OutputPath+string(filepath.Separator)+filename+"_down.png", downgraded); err != nil {
			os.Exit(-2)
		}

		if exportType.RollMode {
			if *rla != -1 || *sla != -1 {
				gfx.RollLeft(*rla, *sla, *iterations, screenMode, size, downgraded, newPalette, filename, exportType)
			} else {
				if *rra != -1 || *sra != -1 {
					gfx.RollRight(*rra, *sra, *iterations, screenMode, size, downgraded, newPalette, filename, exportType)
				}
			}
			if *keephigh != -1 || *losthigh != -1 {
				gfx.RollUp(*keephigh, *losthigh, *iterations, screenMode, size, downgraded, newPalette, filename, exportType)
			} else {
				if *keeplow != -1 || *lostlow != -1 {
					gfx.RollLow(*keeplow, *lostlow, *iterations, screenMode, size, downgraded, newPalette, filename, exportType)
				}
			}
		} else {
			if !customDimension {
				gfx.Transform(downgraded, newPalette, size, *picturePath, exportType)
			} else {
				fmt.Fprintf(os.Stdout, "Transform image in sprite.\n")
				gfx.SpriteTransform(downgraded, newPalette, size, screenMode, filename, exportType)
			}
		}
		if exportType.RotationMode {
			if err := gfx.Rotate(downgraded, newPalette, size, uint8(*mode), *picturePath, resizeAlgo, exportType); err != nil {
				fmt.Fprintf(os.Stderr, "Error while perform rotation on image (%s) error :%v\n", *picturePath, err)
			}
		}
	}
	if exportType.Dsk {
		if err := gfx.ImportInDsk(exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create or write into dsk file error :%v\n", err)
		}
	}
	os.Exit(0)
}
