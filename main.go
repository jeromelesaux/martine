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
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/export/net"
	"github.com/jeromelesaux/martine/gfx"
)

type stringSlice []string

func (f *stringSlice) String() string {
	return ""
}

func (f *stringSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var deltaFiles stringSlice
var (
	byteStatement       = flag.String("s", "", "Byte statement to replace in ascii export (default is BYTE), you can replace or instance by defb")
	picturePath         = flag.String("i", "", "Picture path of the input file.")
	width               = flag.Int("w", -1, "Custom output width in pixels.")
	height              = flag.Int("h", -1, "Custom output height in pixels.")
	mode                = flag.Int("m", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tand add -f option for overscan export.\n\t")
	output              = flag.String("o", "", "Output directory")
	overscan            = flag.Bool("f", false, "Overscan mode (default no overscan)")
	resizeAlgorithm     = flag.Int("a", 1, "Algorithm to resize the image (available : \n\t1: NearestNeighbor (default)\n\t2: CatmullRom\n\t3: Lanczos\n\t4: Linear\n\t5: Box\n\t6: Hermite\n\t7: BSpline\n\t8: Hamming\n\t9: Hann\n\t10: Gaussian\n\t11: Blackman\n\t12: Bartlett\n\t13: Welch\n\t14: Cosine\n\t15: MitchellNetravali\n\t")
	help                = flag.Bool("help", false, "Display help message")
	noAmsdosHeader      = flag.Bool("n", false, "no amsdos header for all files (default amsdos header added).")
	plusMode            = flag.Bool("p", false, "Plus mode (means generate an image for CPC Plus Screen)")
	rollMode            = flag.Bool("roll", false, "Roll mode allow to walk and walk into the input file, associated with rla,rra,sra,sla, keephigh, keeplow, losthigh or lostlow options.")
	iterations          = flag.Int("iter", -1, "Iterations number to walk in roll mode, or number of images to generate in rotation mode.")
	rra                 = flag.Int("rra", -1, "bit rotation on the right and keep pixels")
	rla                 = flag.Int("rla", -1, "bit rotation on the left and keep pixels")
	sra                 = flag.Int("sra", -1, "bit rotation on the right and lost pixels")
	sla                 = flag.Int("sla", -1, "bit rotation on the left and lost pixels")
	losthigh            = flag.Int("losthigh", -1, "bit rotation on the top and lost pixels")
	lostlow             = flag.Int("lostlow", -1, "bit rotation on the bottom and lost pixels")
	keephigh            = flag.Int("keephigh", -1, "bit rotation on the top and keep pixels")
	keeplow             = flag.Int("keeplow", -1, "bit rotation on the bottom and keep pixels")
	palettePath         = flag.String("pal", "", "Apply the input palette to the image")
	info                = flag.Bool("info", false, "Return the information of the file, associated with -pal and -win options")
	winPath             = flag.String("win", "", "Filepath of the ocp win file")
	dsk                 = flag.Bool("dsk", false, "Copy files in a new CPC image Dsk.")
	tileMode            = flag.Bool("tile", false, "Tile mode to create multiples sprites from a same image.")
	tileIterationX      = flag.Int("iterx", -1, "Number of tiles on a row in the input image.")
	tileIterationY      = flag.Int("itery", -1, "Number of tiles on a column in the input image.")
	compress            = flag.Int("z", -1, "Compression algorithm : \n\t1: rle (default)\n\t2: rle 16bits\n\t3: Lz4 Classic\n\t4: Lz4 Raw\n")
	kitPath             = flag.String("kit", "", "Path of the palette Cpc plus Kit file.")
	inkPath             = flag.String("ink", "", "Path of the palette Cpc ink file.")
	rotateMode          = flag.Bool("rotate", false, "Allow rotation on the input image, the input image must be a square (width equals height)")
	m4Host              = flag.String("host", "", "Set the ip of your M4.")
	m4RemotePath        = flag.String("remotepath", "", "remote path on your M4 where you want to copy your files.")
	m4Autoexec          = flag.Bool("autoexec", false, "Execute on your remote CPC the screen file or basic file.")
	rotate3dMode        = flag.Bool("rotate3d", false, "Allow 3d rotation on the input image, the input image must be a square (width equals height)")
	rotate3dType        = flag.Int("rotate3dtype", 0, "Rotation type :\n\t1 rotate on X axis\n\t2 rotate on Y axis\n\t3 rotate reverse X axis\n\t4 rotate left to right on Y axis\n\t5 diagonal rotation on X axis\n\t6 diagonal rotation on Y axis\n")
	rotate3dX0          = flag.Int("rotate3dx0", -1, "X0 coordinate to apply in 3d rotation (default width of the image/2)")
	rotate3dY0          = flag.Int("rotate3dy0", -1, "Y0 coordinate to apply in 3d rotation (default height of the image/2)")
	initProcess         = flag.String("initprocess", "", "create a new empty process file.")
	processFile         = flag.String("processfile", "", "Process file path to apply.")
	deltaMode           = flag.Bool("delta", false, "delta mode: compute delta between two files (prefixed by the argument -df)\n\t(ex: -delta -df file1.SCR -df file2.SCR -df file3.SCR).\n\t(ex with wildcard: -delta -df file\\?.SCR or -delta file\\*.SCR")
	ditheringAlgo       = flag.Int("dithering", -1, "Dithering algorithm to apply on input image\nAlgorithms available:\n\t0: FloydSteinberg\n\t1: JarvisJudiceNinke\n\t2: Stucki\n\t3: Atkinson\n\t4: Sierra\n\t5: SierraLite\n\t6: Sierra3\n\t7: Bayer2\n\t8: Bayer3\n\t9: Bayer4\n\t10: Bayer8\n")
	ditheringMultiplier = flag.Float64("multiplier", 1.18, "error dithering multiplier.")
	withQuantization    = flag.Bool("quantization", false, "use additionnal quantization for dithering.")
	extendedDsk         = flag.Bool("extendeddsk", false, "Export in a Extended DSK 80 tracks, 10 sectors 400 ko per face")
	reverse             = flag.Bool("reverse", false, "Transform .scr (overscan or not) file with palette (pal or kit file) into png file")
	flash               = flag.Bool("flash", false, "generate flash animation with two ocp screens.\n\t(ex: -m 1 -flash -i input.png -o test -dsk)\n\tor\n\t(ex: -m 1 -flash -i input1.scr -pal input1.pal -m2 0 -i2 input2.scr -pal2 input2.pal -o test -dsk )")
	picturePath2        = flag.String("i2", "", "Picture path of the second input file (flash mode)")
	mode2               = flag.Int("m2", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tmode of the second input file (flash mode)")
	palettePath2        = flag.String("pal2", "", "Apply the input palette to the second image (flash mode)")
	version             = "0.18.1.rc"
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
	var screenMode uint8
	var in image.Image

	flag.Var(&deltaFiles, "df", "scr file path to add in delta mode comparison. (wildcard accepted such as ? or * file filename.) ")
	flag.Parse()

	if *help {
		usage()
	}

	if *initProcess != "" {
		err := InitProcess(*initProcess)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating (%s) process file error :%v\n", *initProcess, err)
			os.Exit(-1)
		}
		os.Exit(0)
	}

	if *processFile != "" {
		proc, err := LoadProcessFile(*processFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while loading (%s) process file error :%v\n", *initProcess, err)
			os.Exit(-1)
		}
		proc.Apply()
		if proc.PicturePath == "" && !proc.Delta {
			err = proc.GenerateRawFile()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while loading (%s) process file error :%v\n", *initProcess, err)
				os.Exit(-1)
			}
		}
	}

	if *info {
		if *palettePath != "" {
			file.PalInformation(*palettePath)
		}
		if *winPath != "" {
			file.WinInformation(*winPath)
		}
		if *kitPath != "" {
			file.KitInformation(*kitPath)
		}
		if *inkPath != "" {
			file.InkInformation(*inkPath)
		}
		os.Exit(0)
	}

	// picture path to convert
	if *picturePath == "" && !*deltaMode {
		fmt.Fprintf(os.Stderr, "No picture to compute (option -picturepath or -delta)\n")
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

	exportType := x.NewExportType(*picturePath, *output)

	if *mode == -1 && !*deltaMode && !*reverse {
		fmt.Fprintf(os.Stderr, "No output mode defined can not choose. Quiting\n")
		usage()
	}
	if !*reverse {
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
			if *height == -1 && *width == -1 && !*deltaMode {
				fmt.Fprintf(os.Stderr, "mode %d not defined and no custom width or height\n", *mode)
				usage()
			}
		}
		if *height != -1 {
			exportType.CustomDimension = true
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
			exportType.CustomDimension = true
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

	if *byteStatement != "" {
		file.ByteToken = *byteStatement
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
	exportType.ExtendedDsk = *extendedDsk
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
	exportType.Rotation3DMode = *rotate3dMode
	exportType.Rotation3DType = *rotate3dType
	exportType.Rotation3DX0 = *rotate3dX0
	exportType.Rotation3DY0 = *rotate3dY0
	exportType.M4Host = *m4Host
	exportType.M4RemotePath = *m4RemotePath
	exportType.M4Autoexec = *m4Autoexec
	exportType.ResizingAlgo = resizeAlgo
	exportType.DitheringMultiplier = *ditheringMultiplier
	exportType.DitheringWithQuantification = *withQuantization
	exportType.PalettePath = *palettePath
	exportType.InkPath = *inkPath
	exportType.KitPath = *kitPath
	exportType.RotationRlaBit = *rla
	exportType.RotationSraBit = *sra
	exportType.RotationSlaBit = *sla
	exportType.RotationRraBit = *rra
	exportType.RotationKeephighBit = *keephigh
	exportType.RotationKeeplowBit = *keeplow
	exportType.RotationLosthighBit = *losthigh
	exportType.RotationLostlowBit = *lostlow
	exportType.RotationIterations = *iterations

	if exportType.CpcPlus {
		exportType.Kit = true
		exportType.Pal = false
	}
	exportType.Overscan = *overscan
	if exportType.Overscan {
		exportType.Scr = false
		exportType.Kit = true
	}
	if exportType.M4Host != "" {
		exportType.M4 = true
	}

	exportType.DeltaMode = *deltaMode
	exportType.Dsk = *dsk

	fmt.Fprintf(os.Stdout, "Informations :\n%s", size.ToString())
	if !exportType.DeltaMode && !*reverse && strings.ToUpper(extension) != ".SCR" {
		f, err := os.Open(*picturePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening file %s, error %v\n", *picturePath, err)
			os.Exit(-2)
		}
		defer f.Close()
		in, _, err = image.Decode(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot decode the image %s error %v", *picturePath, err)
			os.Exit(-2)
		}
	}
	if !exportType.CustomDimension && *rotateMode {
		size.Width = in.Bounds().Max.X
		size.Height = in.Bounds().Max.Y
	}

	if !*deltaMode {
		fmt.Fprintf(os.Stdout, "Filename :%s, extension:%s\n", filename, extension)
	}

	if *ditheringAlgo != -1 {
		switch *ditheringAlgo {
		case 0:
			exportType.DitheringMatrix = gfx.FloydSteinberg
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:FloydSteinberg, Type:ErrorDiffusionDither\n")
		case 1:
			exportType.DitheringMatrix = gfx.JarvisJudiceNinke
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:JarvisJudiceNinke, Type:ErrorDiffusionDither\n")
		case 2:
			exportType.DitheringMatrix = gfx.Stucki
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:Stucki, Type:ErrorDiffusionDither\n")
		case 3:
			exportType.DitheringMatrix = gfx.Atkinson
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:Atkinson, Type:ErrorDiffusionDither\n")
		case 4:
			exportType.DitheringMatrix = gfx.Sierra
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:Sierra, Type:ErrorDiffusionDither\n")
		case 5:
			exportType.DitheringMatrix = gfx.SierraLite
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:SierraLite, Type:ErrorDiffusionDither\n")
		case 6:
			exportType.DitheringMatrix = gfx.Sierra3
			exportType.DitheringType = constants.ErrorDiffusionDither
			fmt.Fprintf(os.Stdout, "Dither:Sierra3, Type:ErrorDiffusionDither\n")
		case 7:
			exportType.DitheringMatrix = gfx.Bayer2
			exportType.DitheringType = constants.OrderedDither
			fmt.Fprintf(os.Stdout, "Dither:Bayer2, Type:OrderedDither\n")
		case 8:
			exportType.DitheringMatrix = gfx.Bayer3
			exportType.DitheringType = constants.OrderedDither
			fmt.Fprintf(os.Stdout, "Dither:Bayer3, Type:OrderedDither\n")
		case 9:
			exportType.DitheringMatrix = gfx.Bayer4
			exportType.DitheringType = constants.OrderedDither
			fmt.Fprintf(os.Stdout, "Dither:Bayer4, Type:OrderedDither\n")
		case 10:
			exportType.DitheringMatrix = gfx.Bayer8
			exportType.DitheringType = constants.OrderedDither
			fmt.Fprintf(os.Stdout, "Dither:Bayer8, Type:OrderedDither\n")
		default:
			fmt.Fprintf(os.Stderr, "Dithering matrix not available.")
			os.Exit(-1)
		}
	}

	if *reverse {

		outpath := *output + string(filepath.Separator) + strings.Replace(strings.ToLower(filename), ".scr", ".png", 1)
		if exportType.Overscan {
			p, mode, err := file.OverscanPalette(*picturePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot get the palette from file (%s) error %v\n", *picturePath, err)
				os.Exit(-1)
			}

			if err := gfx.OverscanToPng(*picturePath, outpath, mode, p); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not convert to PNG file (%s) error %v\n", *picturePath, err)
				os.Exit(-1)
			}
			os.Exit(1)
		}
		if *mode == -1 {
			fmt.Fprintf(os.Stderr, "Mode is mandatory to convert to PNG")
			os.Exit(-1)
		}
		var p color.Palette
		var err error
		if *palettePath != "" && *plusMode == false {
			p, _, err = file.OpenPal(*palettePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot open palette file (%s) error %v\n", *palettePath, err)
				os.Exit(-1)
			}
		} else {
			if *kitPath != "" && *plusMode == true {
				p, _, err = file.OpenKit(*kitPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Cannot open kit file (%s) error %v\n", *kitPath, err)
					os.Exit(-1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "For screen or window image, pal or kit file palette is mandatory.\n")
				os.Exit(-1)
			}
		}

		if err := gfx.ScrToPng(*picturePath, outpath, uint8(*mode), p); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot not convert to PNG file (%s) error %v\n", *picturePath, err)
			os.Exit(-1)
		}
		os.Exit(1)
	}

	if exportType.DeltaMode {
		fmt.Fprintf(os.Stdout, "delta files to proceed.\n")
		for i, v := range deltaFiles {
			fmt.Fprintf(os.Stdout, "[%d]:%s\n", i, v)
		}
		if err := gfx.ProceedDelta(deltaFiles, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "error while proceeding delta mode %v\n", err)
			os.Exit(-1)
		}
	} else {
		if exportType.TileMode {
			if exportType.TileIterationX == -1 || exportType.TileIterationY == -1 {
				fmt.Fprintf(os.Stderr, "missing arguments iterx and itery to use with tile mode.\n")
				usage()
				os.Exit(-1)
			}
			err := gfx.TileMode(exportType, uint8(*mode), exportType.TileIterationX, exportType.TileIterationY)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Tile mode on error : error :%v\n", err)
				os.Exit(-1)
			}
		} else {
			if *flash {
				if err := gfx.Flash(*picturePath, *picturePath2,
					*palettePath, *palettePath2,
					*mode,
					*mode2,
					exportType); err != nil {
					fmt.Fprintf(os.Stderr, "Error while applying on one image :%v\n", err)
					os.Exit(-1)
				}
			} else {
				if err := gfx.ApplyOneImage(in,
					exportType,
					filename, *picturePath,
					*mode,
					screenMode); err != nil {
					fmt.Fprintf(os.Stderr, "Error while applying on one image :%v\n", err)
					os.Exit(-1)
				}
			}
		}
	}
	if exportType.Dsk {
		if err := file.ImportInDsk(exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create or write into dsk file error :%v\n", err)
		}
	}
	if exportType.M4 {
		if err := net.ImportInM4(exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot send to M4 error :%v\n", err)
		}
	}
	os.Exit(0)
}
