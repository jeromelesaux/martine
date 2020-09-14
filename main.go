package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
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
	byteStatement       = flag.String("s", "", "Byte statement to replace in ascii export (default is db), you can replace or instance by defb or byte")
	picturePath         = flag.String("i", "", "Picture path of the input file.")
	width               = flag.Int("w", -1, "Custom output width in pixels. (Will produce a sprite file .win)")
	height              = flag.Int("h", -1, "Custom output height in pixels. (Will produce a sprite file .win)")
	mode                = flag.Int("m", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tand add -f option for overscan export.\n\t")
	output              = flag.String("o", "", "Output directory")
	overscan            = flag.Bool("f", false, "Overscan mode (default no overscan)")
	resizeAlgorithm     = flag.Int("a", 1, "Algorithm to resize the image (available : \n\t1: NearestNeighbor (default)\n\t2: CatmullRom\n\t3: Lanczos\n\t4: Linear\n\t5: Box\n\t6: Hermite\n\t7: BSpline\n\t8: Hamming\n\t9: Hann\n\t10: Gaussian\n\t11: Blackman\n\t12: Bartlett\n\t13: Welch\n\t14: Cosine\n\t15: MitchellNetravali\n\t")
	help                = flag.Bool("help", false, "Display help message")
	noAmsdosHeader      = flag.Bool("n", false, "No amsdos header for all files (default amsdos header added).")
	plusMode            = flag.Bool("p", false, "Plus mode (means generate an image for CPC Plus Screen)")
	rollMode            = flag.Bool("roll", false, "Roll mode allow to walk and walk into the input file, associated with rla,rra,sra,sla, keephigh, keeplow, losthigh or lostlow options.")
	iterations          = flag.Int("iter", -1, "Iterations number to walk in roll mode, or number of images to generate in rotation mode.")
	rra                 = flag.Int("rra", -1, "Bit rotation on the right and keep pixels")
	rla                 = flag.Int("rla", -1, "Bit rotation on the left and keep pixels")
	sra                 = flag.Int("sra", -1, "Bit rotation on the right and lost pixels")
	sla                 = flag.Int("sla", -1, "Bit rotation on the left and lost pixels")
	losthigh            = flag.Int("losthigh", -1, "Bit rotation on the top and lost pixels")
	lostlow             = flag.Int("lostlow", -1, "Bit rotation on the bottom and lost pixels")
	keephigh            = flag.Int("keephigh", -1, "Bit rotation on the top and keep pixels")
	keeplow             = flag.Int("keeplow", -1, "Bit rotation on the bottom and keep pixels")
	palettePath         = flag.String("pal", "", "Apply the input palette to the image")
	info                = flag.Bool("info", false, "Return the information of the file, associated with -pal and -win options")
	winPath             = flag.String("win", "", "Filepath of the ocp win file")
	dsk                 = flag.Bool("dsk", false, "Copy files in a new CPC image Dsk.")
	tileMode            = flag.Bool("tile", false, "Tile mode to create multiples sprites from a same image.")
	tileIterationX      = flag.Int("iterx", 1, "Number of tiles on a row in the input image.")
	tileIterationY      = flag.Int("itery", 1, "Number of tiles on a column in the input image.")
	compress            = flag.Int("z", -1, "Compression algorithm : \n\t1: rle (default)\n\t2: rle 16bits\n\t3: Lz4 Classic\n\t4: Lz4 Raw\n")
	kitPath             = flag.String("kit", "", "Path of the palette Cpc plus Kit file. (Apply the input kit palette on the image)")
	inkPath             = flag.String("ink", "", "Path of the palette Cpc ink file. (Apply the input ink palette on the image)")
	rotateMode          = flag.Bool("rotate", false, "Allow rotation on the input image, the input image must be a square (width equals height)")
	m4Host              = flag.String("host", "", "Set the ip of your M4.")
	m4RemotePath        = flag.String("remotepath", "", "Remote path on your M4 where you want to copy your files.")
	m4Autoexec          = flag.Bool("autoexec", false, "Execute on your remote CPC the screen file or basic file.")
	rotate3dMode        = flag.Bool("rotate3d", false, "Allow 3d rotation on the input image, the input image must be a square (width equals height)")
	rotate3dType        = flag.Int("rotate3dtype", 0, "Rotation type :\n\t1 rotate on X axis\n\t2 rotate on Y axis\n\t3 rotate reverse X axis\n\t4 rotate left to right on Y axis\n\t5 diagonal rotation on X axis\n\t6 diagonal rotation on Y axis\n")
	rotate3dX0          = flag.Int("rotate3dx0", -1, "X0 coordinate to apply in 3d rotation (default width of the image/2)")
	rotate3dY0          = flag.Int("rotate3dy0", -1, "Y0 coordinate to apply in 3d rotation (default height of the image/2)")
	initProcess         = flag.String("initprocess", "", "Create a new empty process file.")
	processFile         = flag.String("processfile", "", "Process file path to apply.")
	deltaMode           = flag.Bool("delta", false, "Delta mode: compute delta between two files (prefixed by the argument -df)\n\t(ex: -delta -df file1.SCR -df file2.SCR -df file3.SCR).\n\t(ex with wildcard: -delta -df file\\?.SCR or -delta file\\*.SCR")
	ditheringAlgo       = flag.Int("dithering", -1, "Dithering algorithm to apply on input image\nAlgorithms available:\n\t0: FloydSteinberg\n\t1: JarvisJudiceNinke\n\t2: Stucki\n\t3: Atkinson\n\t4: Sierra\n\t5: SierraLite\n\t6: Sierra3\n\t7: Bayer2\n\t8: Bayer3\n\t9: Bayer4\n\t10: Bayer8\n")
	ditheringMultiplier = flag.Float64("multiplier", 1.18, "Error dithering multiplier.")
	withQuantization    = flag.Bool("quantization", false, "Use additionnal quantization for dithering.")
	extendedDsk         = flag.Bool("extendeddsk", false, "Export in a Extended DSK 80 tracks, 10 sectors 400 ko per face")
	reverse             = flag.Bool("reverse", false, "Transform .scr (overscan or not) file with palette (pal or kit file) into png file")
	flash               = flag.Bool("flash", false, "generate flash animation with two ocp screens.\n\t(ex: -m 1 -flash -i input.png -o test -dsk)\n\tor\n\t(ex: -m 1 -flash -i input1.scr -pal input1.pal -m2 0 -i2 input2.scr -pal2 input2.pal -o test -dsk )")
	picturePath2        = flag.String("i2", "", "Picture path of the second input file (flash mode)")
	mode2               = flag.Int("m2", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tmode of the second input file (flash mode)")
	palettePath2        = flag.String("pal2", "", "Apply the input palette to the second image (flash mode)")
	egx1                = flag.Bool("egx1", false, "Create egx 1 output cpc image overscan (option -f) or classical (mix mode 0 / 1).\n\t(ex before generate two images one in mode 1 et one in mode 0\n\tfor instance : martine -i myimage.jpg -m 0 and martine -i myimage.jpg -m 1\n\t: -egx1 -i 1.SCR -m 0 -pal 1.PAL -i2 2.SCR -o test -m2 1 -dsk)\n\tor\n\t(ex automatic egx from image file : -egx1 -i input.png -m 0 -o test -dsk)")
	egx2                = flag.Bool("egx2", false, "Create egx 2 output cpc image overscan (option -f) or classical (mix mode 1 / 2).\n\t(ex before generate two images one in mode 1 et one in mode 2\n\tfor instance : martine -i myimage.jpg -m 0 and martine -i myimage.jpg -m 1\n\t: -egx2 -i 1.SCR -m 0 -pal 1.PAL -i2 2.SCR -o test -m2 1 -dsk)\n\tor\n\t(ex automatic egx from image file : -egx2 -i input.png -m 0 -o test -dsk)")
	sna                 = flag.Bool("sna", false, "Copy files in a new CPC image Sna.")
	spriteHard          = flag.Bool("spritehard", false, "Generate sprite hard for cpc plus.")
	splitRasters        = flag.Bool("splitrasters", false, "Create Split rastered image. (Will produce Overscan output file and .SPL with split rasters file)")
	scanlineSequence    = flag.String("scanlinesequence", "", "Scanline sequence to apply on sprite. for instance : \n\tmartine -i myimage.jpg -w 4 -h 4 -scanlinesequence 0,2,1,3 \n\twill generate a sprite stored with lines order 0 2 1 and 3.\n")
	maskSprite          = flag.String("mask", "", "Mask to apply on each bit of the sprite (to apply an and operation on each pixel with the value #AA [in hexdecimal: #AA or 0xAA, in decimal: 170] ex: martine -i myimage.png -w 40 -h 80 -mask #AA -m 0 -maskand)")
	maskOrOperation     = flag.Bool("maskor", false, "Will apply an OR operation on each byte with the mask")
	maskAdOperation     = flag.Bool("maskand", false, "Will apply an AND operation on each byte with the mask")
	zigzag              = flag.Bool("zigzag", false, "generate data in zigzag order (inc first line and dec next line for tiles)")
	tileMap             = flag.Bool("tilemap", false, "Analyse the input image and generate the tiles, the tile map and gloabl schema.")
	initialAddress      = flag.String("address", "0xC000", "Starting address to display sprite in delta packing")
	animate             = flag.Bool("animate", false, "Will produce an full screen with all sprite on the same image (add -i image.gif or -i *.png)")
	reducer             = flag.Int("reducer", -1, "Reducer mask will reduce original image colors. Available : \n\t1 : lower\n\t2 : medium\n\t3 : strong\n")
	jsonOutput          = flag.Bool("json", false, "Generate json format output.")
	txtOutput           = flag.Bool("txt", false, "Generate text format output.")
	oneLine             = flag.Bool("oneline", false, "Display every other line.")
	oneRow              = flag.Bool("onerow", false, "Display  every other row.")
	impCatcher          = flag.Bool("imp", false, "Will generate sprites as IMP-Catcher format (Impdraw V2).")
	inkSwap             = flag.String("inkswap", "", "Swap ink:\n\tfor instance mode 4 (4 inks) : 0=4,1=3,2=1,4=2\n\twill swap in output image index 0 by 4 and 1 by 3 and so on.")
	appVersion          = "0.27.0rc"
	version             = flag.Bool("version", false, "print martine's version")
)

func usage() {
	fmt.Fprintf(os.Stdout, "martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)\n")
	fmt.Fprintf(os.Stdout, "By Impact Sid (Version:%s)\n", appVersion)
	fmt.Fprintf(os.Stdout, "Special thanks to @Ast (for his support), @Siko and @Tronic for ideas\n")
	fmt.Fprintf(os.Stdout, "usage :\n\n")
	flag.PrintDefaults()
	os.Exit(-1)
}

func printVersion() {
	fmt.Fprintf(os.Stdout, "%s\n", appVersion)
	os.Exit(-1)
}

func main() {
	var size constants.Size
	var filename, extension string
	var screenMode uint8
	var in image.Image

	flag.Var(&deltaFiles, "df", "scr file path to add in delta mode comparison. (wildcard accepted such as ? or * file filename.) ")

	flag.Parse()
	if len(flag.Args()) > 0 {
		firstArg := flag.Args()[0]
		if firstArg[0] != '-' {
			flag.Set("i", firstArg)
			for i := 1; i < len(flag.Args()); i += 2 {
				name := strings.Replace(flag.Arg(i), "-", "", 1)
				var value string
				if len(flag.Args()) > i+1 {
					if flag.Arg(i + 1)[0] == '-' {
						value = "true"
						i--
					} else {
						value = flag.Arg(i + 1)
					}
				} else {
					value = "true"
				}
				flag.Set(name, value)
			}
			flag.Parse()
		}
	}
	if *help {
		usage()
	}
	if *version {
		printVersion()
	}

	if *initProcess != "" {
		_, err := InitProcess(*initProcess)
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
		if err := common.CheckOutput(*output); err != nil {
			fmt.Fprintf(os.Stderr, "Error while getting directory informations :%v, Quiting\n", err)
			os.Exit(-2)
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

	if *scanlineSequence != "" {
		sequence := strings.Split(*scanlineSequence, ",")
		for _, v := range sequence {
			line, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bad scanline sequence (%s) error:%v\n", *scanlineSequence, err)
				os.Exit(-1)
			}
			exportType.ScanlineSequence = append(exportType.ScanlineSequence, line)
		}
		modulo := size.Height % len(exportType.ScanlineSequence)
		if modulo != 0 {
			fmt.Fprintf(os.Stderr, "height modulo scanlinesequence is not equal to 0 %d lines and the output image lines is %d\n", len(exportType.ScanlineSequence), size.Height)
			os.Exit(-1)
		}
	}
	exportType.ExtendedDsk = *extendedDsk
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
	exportType.Flash = *flash
	exportType.Sna = *sna
	exportType.SpriteHard = *spriteHard
	exportType.SplitRaster = *splitRasters
	exportType.ZigZag = *zigzag
	exportType.Animate = *animate
	exportType.Reducer = *reducer
	exportType.Json = *jsonOutput
	exportType.Ascii = *txtOutput
	exportType.OneLine = *oneLine
	exportType.OneRow = *oneRow
	if err := exportType.ImportInkSwap(*inkSwap); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse inkswap option with error [%s]\n", err)
		os.Exit(-1)
	}

	if *maskSprite != "" {

		v, err := common.ParseHexadecimal8(*maskSprite)
		if err == nil {
			exportType.MaskSprite = uint8(v)
		}
		if exportType.MaskSprite != 0 {
			if *maskOrOperation {
				exportType.MaskOrOperation = true
			}
			if *maskAdOperation {
				exportType.MaskAndOperation = true
			}
			if exportType.MaskAndOperation && exportType.MaskOrOperation {
				fmt.Fprintf(os.Stderr, "Or and And operations are setted, will only apply And operation.\n")
				exportType.MaskOrOperation = false
			}
			if !exportType.MaskAndOperation && !exportType.MaskOrOperation {
				fmt.Fprintf(os.Stderr, "Or and And operations are not setted, will only apply And operation.\n")
				exportType.MaskAndOperation = true
			}
			fmt.Fprintf(os.Stdout, "Applying sprite mask value [#%X] [%.8b] AND = %t, OR =%t\n",
				exportType.MaskSprite,
				exportType.MaskSprite,
				exportType.MaskAndOperation,
				exportType.MaskOrOperation)
		}
	}

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

	if *egx1 {
		exportType.EgxFormat = x.Egx1Mode
	}
	if *egx2 {
		exportType.EgxFormat = x.Egx2Mode
	}
	if *mode != -1 {
		exportType.EgxMode1 = uint8(*mode)
	}
	if *mode2 != -1 {
		exportType.EgxMode2 = uint8(*mode2)
	}

	exportType.DeltaMode = *deltaMode
	exportType.Dsk = *dsk

	fmt.Fprintf(os.Stdout, "Informations :\n%s", size.ToString())
	if !*impCatcher && !exportType.DeltaMode && !*reverse && !*animate && strings.ToUpper(extension) != ".SCR" {
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

	// gestion de la taille de l'image en sortie
	if !exportType.CustomDimension && *rotateMode && !exportType.SpriteHard {
		size.Width = in.Bounds().Max.X
		size.Height = in.Bounds().Max.Y
	}
	if *spriteHard {
		size.Width = 16
		size.Height = 16
	}
	exportType.Size = size

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
	if *impCatcher {
		if !exportType.CustomDimension {
			fmt.Fprintf(os.Stderr, "You must set custom width and height.")
			os.Exit(-1)
		}
		sprites := make([]byte, 0)
		fmt.Fprintf(os.Stdout, "[%s]\n", *picturePath)
		spritesPaths, err := common.WilcardedFiles([]string{*picturePath})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting wildcard files %s error : %v\n", *picturePath, err)
		}
		for _, v := range spritesPaths {
			f, err := os.Open(v)
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
			gfx.ApplyOneImage(in,
				exportType,
				filepath.Base(v),
				v,
				*mode,
				screenMode)

			spritePath := exportType.AmsdosFullPath(v, ".WIN")
			data, err := file.RawWin(spritePath)
			sprites = append(sprites, data...)
		}
		finalFile := strings.ReplaceAll(filename, "?", "")
		if err = file.Imp(sprites, uint(exportType.Size.Width), uint(exportType.Size.Height), finalFile, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot export to Imp-Catcher the image %s error %v", *picturePath, err)
		}
		os.Exit(0)
	} else if *reverse {

		outpath := filepath.Join(*output, strings.Replace(strings.ToLower(filename), ".scr", ".png", 1))
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
				fmt.Fprintf(os.Stderr, "For screen or window image, pal or kit file palette is mandatory. (kit file must be associated with -p option)\n")
				os.Exit(-1)
			}
		}
		switch strings.ToUpper(filepath.Ext(filename)) {
		case ".WIN":
			if err := gfx.SpriteToPng(*picturePath, outpath, uint8(*mode), p); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not convert to PNG file (%s) error %v\n", *picturePath, err)
				os.Exit(-1)
			}
		case ".SCR":
			if err := gfx.ScrToPng(*picturePath, outpath, uint8(*mode), p); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot not convert to PNG file (%s) error %v\n", *picturePath, err)
				os.Exit(-1)
			}
		}
		os.Exit(1)
	}
	if exportType.Animate {
		if !exportType.CustomDimension {
			fmt.Fprintf(os.Stderr, "You must set sprite dimensions with option -w and -h (mandatory)\n")
			os.Exit(-1)
		}
		fmt.Fprintf(os.Stdout, "animation output.\n")
		files := []string{*picturePath}
		files, err := common.WilcardedFiles(files)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse wildcard in argument (%s) error %v\n", *picturePath, err)
			os.Exit(-1)
		}
		if err := gfx.Animation(files, screenMode, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Error while proceeding to animate export error : %v\n", err)
			os.Exit(-1)
		}
	} else {
		if exportType.DeltaMode {
			fmt.Fprintf(os.Stdout, "delta files to proceed.\n")
			for i, v := range deltaFiles {
				fmt.Fprintf(os.Stdout, "[%d]:%s\n", i, v)
			}
			screenAddress, err := common.ParseHexadecimal16(*initialAddress)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while parsing (%s) use the starting address #C000, err : %v\n", *initialAddress, err)
				screenAddress = 0xC000
			}
			if *mode == -1 {
				fmt.Fprintf(os.Stderr, "You must set the mode for this feature. (option -m)\n")
				os.Exit(-1)
			}
			if err := gfx.ProceedDelta(deltaFiles, screenAddress, exportType, uint8(*mode)); err != nil {
				fmt.Fprintf(os.Stderr, "error while proceeding delta mode %v\n", err)
				os.Exit(-1)
			}
		} else {
			if *tileMap {
				if !exportType.CustomDimension {
					fmt.Fprintf(os.Stderr, "You must set height and width to define the tile dimensions (options -h and -w)\n")
					os.Exit(-1)
				}
				analyze := gfx.AnalyzeTilesBoard(in, exportType.Size)
				if err := analyze.SaveSchema(filepath.Join(exportType.OutputPath, "tilesmap_schema.png")); err != nil {
					fmt.Fprintf(os.Stderr, "Cannot save tilemap schema error :%v\n", err)
					os.Exit(-1)
				}
				if err := analyze.SaveTilemap(filepath.Join(exportType.OutputPath, "tilesmap.map")); err != nil {
					fmt.Fprintf(os.Stderr, "Cannot save tilemap csv file error :%v\n", err)
					os.Exit(-1)
				}
				for i, v := range analyze.BoardTiles {
					tile := v.Tile.Image()
					tileFilepath := filepath.Join(exportType.OutputPath, fmt.Sprintf("%.2d.png", i))
					f, err := os.Create(tileFilepath)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Cannot create tiles %.2d error %v\n", i, err)
						os.Exit(-1)
					}
					defer f.Close()
					if err := png.Encode(f, tile); err != nil {
						fmt.Fprintf(os.Stderr, "Cannot encode in png tile %.2d error %v\n", i, err)
						os.Exit(-1)
					}

					var palette color.Palette
					if exportType.CpcPlus {
						palette = constants.CpcPlusPalette
					} else {
						palette = constants.CpcOldPalette
					}

					out, palette := gfx.DoDithering(tile, palette, exportType)
					palette, out, err = convert.DowngradingPalette(out, exportType.Size, exportType.CpcPlus)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Cannot downgrade colors palette for this image %s\n", tileFilepath)
					}

					palette = constants.SortColorsByDistance(palette)

					fmt.Fprintf(os.Stdout, "Saving downgraded image into (%s)\n", filename+"_down.png")
					if err := file.Png(tileFilepath+"_down.png", out); err != nil {
						os.Exit(-2)
					}
					if err := gfx.SpriteTransform(tile, palette, exportType.Size, screenMode, tileFilepath, false, exportType); err != nil {
						fmt.Fprintf(os.Stderr, "Cannot create tile from image %s, error :%v\n", tileFilepath, err)
					}
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
						var p color.Palette
						var err error
						if exportType.CpcPlus {
							if *kitPath != "" {
								p, _, err = file.OpenKit(*kitPath)
								if err != nil {
									fmt.Fprintf(os.Stderr, "Error while reading kit file (%s) :%v\n", *kitPath, err)
									os.Exit(-1)
								}
							}
						} else {
							if *palettePath != "" {
								p, _, err = file.OpenPal(*palettePath)
								if err != nil {
									fmt.Fprintf(os.Stderr, "Error while reading palette file (%s) :%v\n", *palettePath, err)
									os.Exit(-1)
								}
							}
						}

						if exportType.EgxFormat > 0 {
							if len(p) == 0 {
								fmt.Fprintf(os.Stderr, "Now colors found in palette, give up treatment.\n")
								os.Exit(-1)
							}
							if err := gfx.Egx(*picturePath, *picturePath2,
								p,
								*mode,
								*mode2,
								exportType); err != nil {
								fmt.Fprintf(os.Stderr, "Error while applying on one image :%v\n", err)
								os.Exit(-1)
							}
						} else {
							if exportType.SplitRaster {
								if exportType.Overscan {
									if err := gfx.DoSpliteRaster(in, screenMode, filename, exportType); err != nil {
										fmt.Fprintf(os.Stderr, "Error while applying splitraster on one image :%v\n", err)
										os.Exit(-1)
									}
								} else {
									fmt.Fprintf(os.Stderr, "Only overscan mode implemented for this feature, %v", gfx.ErrorNotYetImplemented)
								}
							} else {
								if strings.ToUpper(extension) != ".SCR" {
									if err := gfx.ApplyOneImage(in,
										exportType,
										filename, *picturePath,
										*mode,
										screenMode); err != nil {
										fmt.Fprintf(os.Stderr, "Error while applying on one image :%v\n", err)
										os.Exit(-1)
									}
								} else {
									fmt.Fprintf(os.Stderr, "Error while applying on one image : SCR format not used for this treatment\n")
									os.Exit(-1)
								}
							}
						}
					}
				}
			}
		}
	}
	// export into bundle DSK or SNA
	if exportType.Dsk {
		if err := file.ImportInDsk(*picturePath, exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create or write into dsk file error :%v\n", err)
		}
	}
	if exportType.Sna {
		if exportType.Overscan {
			var gfxFile string
			for _, v := range exportType.DskFiles {
				if filepath.Ext(v) == ".SCR" {
					gfxFile = v
					break
				}
			}
			exportType.SnaPath = filepath.Join(*output, "test.sna")
			if err := file.ImportInSna(gfxFile, exportType.SnaPath, screenMode); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create or write into sna file error :%v\n", err)
			}
			fmt.Fprintf(os.Stdout, "Sna saved in file %s\n", exportType.SnaPath)
		} else {
			fmt.Fprintf(os.Stderr, "Feature not implemented for this file.")
			os.Exit(-1)
		}
	}
	if exportType.M4 {
		if err := net.ImportInM4(exportType); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot send to M4 error :%v\n", err)
		}
	}
	os.Exit(0)
}
