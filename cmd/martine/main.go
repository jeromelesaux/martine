package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/app"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/diskimage"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/m4"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
	"github.com/jeromelesaux/martine/export/snapshot"
	"github.com/jeromelesaux/martine/gfx/animate"
	"github.com/jeromelesaux/martine/log"
	ui "github.com/jeromelesaux/martine/ui/martine-ui"
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
	byteStatement       = flag.String("statement", "", "Byte statement to replace in ascii export (default is db), you can replace or instance by defb or byte")
	picturePath         = flag.String("in", "", "Picture path of the input file.")
	width               = flag.Int("width", -1, "Custom output width in pixels. (Will produce a sprite file .win)")
	height              = flag.Int("height", -1, "Custom output height in pixels. (Will produce a sprite file .win)")
	mode                = flag.Int("mode", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tand add -fullscreen option for overscan export.\n\t")
	output              = flag.String("out", "", "Output directory")
	overscan            = flag.Bool("fullscreen", false, "Overscan mode (default no overscan)")
	resizeAlgorithm     = flag.Int("algo", 1, "Algorithm to resize the image (available : \n\t1: NearestNeighbor (default)\n\t2: CatmullRom\n\t3: Lanczos\n\t4: Linear\n\t5: Box\n\t6: Hermite\n\t7: BSpline\n\t8: Hamming\n\t9: Hann\n\t10: Gaussian\n\t11: Blackman\n\t12: Bartlett\n\t13: Welch\n\t14: Cosine\n\t15: MitchellNetravali\n\t")
	help                = flag.Bool("help", false, "Display help message")
	noAmsdosHeader      = flag.Bool("noheader", false, "No amsdos header for all files (default amsdos header added).")
	plusMode            = flag.Bool("plus", false, "Plus mode (means generate an image for CPC Plus Screen)")
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
	compress            = flag.Int("z", -1, "Compression algorithm : \n\t1: rle (default)\n\t2: rle 16bits\n\t3: Lz4 Classic\n\t4: Lz4 Raw\n\t5: zx0 crunch\n")
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
	flash               = flag.Bool("flash", false, "generate flash animation with two ocp screens.\n\t(ex: -mode 1 -flash -in input.png -out test -dsk)\n\tor\n\t(ex: -mode 1 -flash -i input1.scr -pal input1.pal -mode2 0 -iin2 input2.scr -pal2 input2.pal -out test -dsk )")
	picturePath2        = flag.String("in2", "", "Picture path of the second input file (flash mode)")
	mode2               = flag.Int("mode2", -1, "Output mode to use :\n\t0 for mode0\n\t1 for mode1\n\t2 for mode2\n\tmode of the second input file (flash mode)")
	palettePath2        = flag.String("pal2", "", "Apply the input palette to the second image (flash mode)")
	egx1                = flag.Bool("egx1", false, "Create egx 1 output cpc image overscan (option -fullscreen) or classical (mix mode 0 / 1).\n\t(ex before generate two images one in mode 1 et one in mode 0\n\tfor instance : martine -in myimage.jpg -mode 0 and martine -in myimage.jpg -mode 1\n\t: -egx1 -in 1.SCR -mode 0 -pal 1.PAL -in2 2.SCR -out test -mode2 1 -dsk)\n\tor\n\t(ex automatic egx from image file : -egx1 -in input.png -mode 0 -out test -dsk)")
	egx2                = flag.Bool("egx2", false, "Create egx 2 output cpc image overscan (option -fullscreen) or classical (mix mode 1 / 2).\n\t(ex before generate two images one in mode 1 et one in mode 2\n\tfor instance : martine -in myimage.jpg -mode 0 and martine -in myimage.jpg -mode 1\n\t: -egx2 -in 1.SCR -mode 0 -pal 1.PAL -in2 2.SCR -out test -mode2 1 -dsk)\n\tor\n\t(ex automatic egx from image file : -egx2 -in input.png -mode 0 -out test -dsk)")
	sna                 = flag.Bool("sna", false, "Copy files in a new CPC image Sna.")
	spriteHard          = flag.Bool("spritehard", false, "Generate sprite hard for cpc plus.")
	splitRasters        = flag.Bool("splitrasters", false, "Create Split rastered image. (Will produce Overscan output file and .SPL with split rasters file)")
	scanlineSequence    = flag.String("scanlinesequence", "", "Scanline sequence to apply on sprite. for instance : \n\tmartine -in myimage.jpg -width 4 -height 4 -scanlinesequence 0,2,1,3 \n\twill generate a sprite stored with lines order 0 2 1 and 3.\n")
	maskSprite          = flag.String("mask", "", "Mask to apply on each bit of the sprite (to apply an and operation on each pixel with the value #AA [in hexdecimal: #AA or 0xAA, in decimal: 170] ex: martine -in myimage.png -width 40 -height 80 -mask #AA -mode 0 -maskand)")
	maskOrOperation     = flag.Bool("maskor", false, "Will apply an OR operation on each byte with the mask")
	maskAdOperation     = flag.Bool("maskand", false, "Will apply an AND operation on each byte with the mask")
	zigzag              = flag.Bool("zigzag", false, "generate data in zigzag order (inc first line and dec next line for tiles)")
	tileMap             = flag.Bool("tilemap", false, "Analyze the input image and generate the tiles, the tile map and global schema.\n\t for instance: martine -in board.png -mode 0 -width 8 -height 8 -out folder -dsk\n")
	initialAddress      = flag.String("address", "0xC000", "Starting address to display sprite in delta packing")
	doAnimation         = flag.Bool("animate", false, "Will produce an full screen with all sprite on the same image (add -in image.gif or -in *.png)")
	reducer             = flag.Int("reducer", -1, "Reducer mask will reduce original image colors. Available : \n\t1 : lower\n\t2 : medium\n\t3 : strong\n")
	jsonOutput          = flag.Bool("json", false, "Generate json format output.")
	txtOutput           = flag.Bool("txt", false, "Generate text format output.")
	oneLine             = flag.Bool("oneline", false, "Display every other line.")
	oneRow              = flag.Bool("onerow", false, "Display  every other row.")
	impCatcher          = flag.Bool("imp", false, "Will generate sprites as IMP-Catcher format (Impdraw V2).")
	inkSwap             = flag.String("inkswap", "", "Swap ink:\n\tfor instance mode 4 (4 inks) : 0=3,1=0,2=1,3=2\n\twill swap in output image index 0 by 3 and 1 by 0 and so on.")
	lineWidth           = flag.String("linewidth", "#50", "Line width in hexadecimal to compute the screen address in delta mode.")
	deltaPacking        = flag.Bool("deltapacking", false, "Will generate all the animation code from the followed gif file.")
	deltaPacking2       = flag.Bool("deltapacking2", false, "Will generate all the animation code from the followed gif file (and optimize export).")
	filloutGif          = flag.Bool("fillout", false, "Fill out the gif frames needed some case with deltapacking")
	saturationPal       = flag.Float64("contrast", 0., "apply contrast on the color of the palette on amstrad plus screen. (max value 100 and only on CPC PLUS).")
	brightnessPal       = flag.Float64("brightness", 0., "apply brightness on the color of the palette on amstrad plus screen. (max value 100 and only on CPC PLUS).")
	analyzeTilemap      = flag.String("analyzetilemap", "", "analyze the image to get the most accurate tilemap according to the  criteria :\n\tsize : lower export size\n\tnumber : lower number of tiles")
	exportGoFiles       = flag.Bool("go", false, "Export results as .go1 and .go2 files.")
	splitSpriteBoard    = flag.Bool("split", false, "Split sprite board to sprites.")
	spritesPerRow       = flag.Int("spritesrow", 0, "Number of sprites in the board per row")
	spritesPerColumn    = flag.Int("spritescolumn", 0, "Number of sprites in the board per column")
	spriteFlat          = flag.Bool("flat", false, "Export sprite as flat file.")
	spriteCompiled      = flag.Bool("compiled", false, "Export sprite as compiled sprites.")
	spriteOcpWin        = flag.Bool("ocpwin", false, "Export sprite as OCP win file.")
	version             = flag.Bool("version", false, "print martine's version")
	appPrefix           = fmt.Sprintf("Martine (%v)", common.AppVersion)
	noUI                = flag.Bool("noui", false, "Open Martine UI")
)

func usage() {
	log.GetLogger().Info("martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)\n")
	log.GetLogger().Info("By Impact Sid (Version:%s)\n", common.AppVersion)
	log.GetLogger().Info("Special thanks to @Ast (for his support), @Siko and @Tronic for ideas\n")
	log.GetLogger().Info("usage :\n\n")
	flag.PrintDefaults()
	os.Exit(-1)
}

func printVersion() {
	log.GetLogger().Info("%s\n", common.AppVersion)
	os.Exit(0)
}

/*
@Todo : add zigzag on sprite and sprite hard.
*/

// nolint: funlen, gocognit
func main() {

	var filename, extension string
	var in image.Image

	log.Default(appPrefix)

	flag.Var(&deltaFiles, "df", "scr file path to add in delta mode comparison. (wildcard accepted such as ? or * file filename.) ")

	flag.Parse()
	if *help {
		usage()
	}

	if !*noUI {
		os.Setenv("FYNE_SCALE", "0.7")
		/* main application */
		app := app.NewWithID("Martine @IMPact")
		martineUI := ui.NewMartineUI()
		martineUI.Load(app)
		app.Run()

		os.Exit(0)
	}
	if len(flag.Args()) > 0 {
		firstArg := flag.Args()[0]
		if firstArg[0] != '-' {
			err := flag.Set("i", firstArg)
			if err != nil {
				log.GetLogger().Error("Error :%v\n", err)
				os.Exit(-1)
			}
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
				err = flag.Set(name, value)
				if err != nil {
					log.GetLogger().Error("Error :%v\n", err)
					os.Exit(-1)
				}
			}
			flag.Parse()
		}
	}

	if *version {
		printVersion()
	}

	if *initProcess != "" {
		_, err := InitProcess(*initProcess)
		if err != nil {
			log.GetLogger().Error("Error while creating (%s) process file error :%v\n", *initProcess, err)
			os.Exit(-1)
		}
		os.Exit(0)
	}

	if *processFile != "" {
		proc, err := LoadProcessFile(*processFile)
		if err != nil {
			log.GetLogger().Error("Error while loading (%s) process file error :%v\n", *initProcess, err)
			os.Exit(-1)
		}
		proc.Apply()
		if proc.PicturePath == "" && !proc.Delta {
			err = proc.GenerateRawFile()
			if err != nil {
				log.GetLogger().Error("Error while loading (%s) process file error :%v\n", *initProcess, err)
				os.Exit(-1)
			}
		}
	}

	if *info {
		if *palettePath != "" {
			ocpartstudio.PalInformation(*palettePath)
		}
		if *winPath != "" {
			window.WinInformation(*winPath)
		}
		if *kitPath != "" {
			impPalette.KitInformation(*kitPath)
		}
		if *inkPath != "" {
			impPalette.InkInformation(*inkPath)
		}
		os.Exit(0)
	}

	// picture path to convert
	if *picturePath == "" && !*deltaMode {
		log.GetLogger().Error("No picture to compute (option -picturepath or -delta)\n")
		usage()
	}
	filename = filepath.Base(*picturePath)
	extension = filepath.Ext(*picturePath)

	// output directory to store results
	if *output != "" {
		if err := common.CheckOutput(*output); err != nil {
			log.GetLogger().Error("Error while getting directory informations :%v, Quiting\n", err)
			os.Exit(-2)
		}
	} else {
		*output = "./"
	}

	if *mode == -1 && !*deltaMode && !*reverse {
		log.GetLogger().Error("No output mode defined can not choose. Quiting\n")
		usage()
	}

	cfg, size := ExportHandler()
	if *byteStatement != "" {
		ascii.ByteToken = *byteStatement
	}

	if *deltaPacking || *deltaPacking2 {
		screenAddress, err := common.ParseHexadecimal16(*initialAddress)
		cfg.ScrCfg.Size = size
		if err != nil {
			log.GetLogger().Error("Error while parsing (%s) use the starting address #C000, err : %v\n", *initialAddress, err)
			screenAddress = 0xC000
		}
		var exportVersion animate.DeltaExportFormat = animate.DeltaExportV1

		if *deltaPacking2 {
			exportVersion = animate.DeltaExportV2
		}
		if err := animate.DeltaPacking(cfg.ScrCfg.InputPath, cfg, screenAddress, uint8(*mode), exportVersion); err != nil {
			log.GetLogger().Error("Error while deltapacking error: %v\n", err)
		}
		os.Exit(0)
	}

	if !*reverse {
		log.GetLogger().Info("Informations :\n%s", size.ToString())
	}
	if !*impCatcher && !cfg.DeltaMode && !*reverse && !*doAnimation && !strings.EqualFold(extension, constants.ScrExtension) {
		f, err := os.Open(*picturePath)
		if err != nil {
			log.GetLogger().Error("Error while opening file %s, error %v\n", *picturePath, err)
			os.Exit(-2)
		}
		in, _, err = image.Decode(f)
		if err != nil {
			log.GetLogger().Error("Cannot decode the image %s error %v", *picturePath, err)
			os.Exit(-2)
		}
		f.Close()
	}

	// gestion de la taille de l'image en sortie
	if !cfg.CustomDimension && *rotateMode && !cfg.ScrCfg.Type.IsSpriteHard() {
		size.Width = in.Bounds().Max.X
		size.Height = in.Bounds().Max.Y
	}
	if *spriteHard {
		size.Width = 16
		size.Height = 16
	}
	cfg.ScrCfg.Size = size

	if !*deltaMode {
		log.GetLogger().Info("Filename :%s, extension:%s\n", filename, extension)
	}

	if *ditheringAlgo != -1 {
		if err := setDithering(cfg); err != nil {
			log.GetLogger().Error("%v\n", err)
			os.Exit(-1)
		}
	}

	if err := runConversion(cfg, size, filename, extension, in); err != nil {
		log.GetLogger().Error("Error while running conversion: %v\n", err)
		os.Exit(-1)
	}
	// export into bundle DSK or SNA
	if cfg.HasContainerExport(config.DskContainer) {
		if err := diskimage.ImportInDsk(*picturePath, cfg); err != nil {
			log.GetLogger().Error("Cannot create or write into dsk file error :%v\n", err)
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
			cfg.ContainerCfg.Path = filepath.Join(*output, "test.sna")
			if err := snapshot.ImportInSna(gfxFile, cfg.ContainerCfg.Path, uint8(*mode)); err != nil {
				log.GetLogger().Error("Cannot create or write into sna file error :%v\n", err)
			}
			log.GetLogger().Info("Sna saved in file %s\n", cfg.ContainerCfg.Path)
		} else {
			log.GetLogger().Error("Feature not implemented for this file.")
			os.Exit(-1)
		}
	}
	if cfg.M4cfg.Enabled {
		if err := m4.ImportInM4(cfg); err != nil {
			log.GetLogger().Error("Cannot send to M4 error :%v\n", err)
		}
	}
	os.Exit(0)
}
