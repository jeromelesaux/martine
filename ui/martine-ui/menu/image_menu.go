package menu

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/assembly"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/screen"
	ovs "github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/convert/spritehard"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/diskimage"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	impPalette "github.com/jeromelesaux/martine/export/impdraw/palette"
	"github.com/jeromelesaux/martine/export/m4"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/export/snapshot"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/log"
	dl "github.com/jeromelesaux/martine/ui/martine-ui/dialog"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
)

type ImageMenu struct {
	originalImage       *canvas.Image
	cpcImage            *canvas.Image
	originalImagePath   fyne.URI
	IsCpcPlus           bool
	IsFullScreen        bool
	IsSprite            bool
	IsHardSprite        bool
	IsWin               bool
	Mode                int
	width               *widget.Entry
	height              *widget.Entry
	palette             color.Palette
	Data                []byte
	Downgraded          *image.NRGBA
	DitheringMatrix     [][]float32
	DitheringType       constants.DitheringType
	DitheringAlgoNumber int
	ApplyDithering      bool
	ResizeAlgo          imaging.ResampleFilter
	ResizeAlgoNumber    int
	paletteImage        *canvas.Image
	UsePalette          bool
	DitheringMultiplier float64
	WithQuantification  bool
	Brightness          float64
	Saturation          float64
	Reducer             int
	OneLine             bool
	OneRow              bool
	CmdLineGenerate     string
	UseKmeans           bool
	KmeansThreshold     float64
	Edited              bool
	w                   fyne.Window
}

func NewImageMenu() *ImageMenu {
	return &ImageMenu{
		originalImage: &canvas.Image{},
		cpcImage:      &canvas.Image{},
		paletteImage:  &canvas.Image{},
		width:         widget.NewEntry(),
		height:        widget.NewEntry(),
		Downgraded:    &image.NRGBA{},
	}
}

func (i *ImageMenu) SetWindow(w fyne.Window) {
	i.w = w
}

func (i *ImageMenu) SetPalette(p color.Palette) {
	i.palette = p
	i.SetPaletteImage(png.PalToImage(i.Palette()))
	i.paletteImage.Refresh()
}

func (i *ImageMenu) Palette() color.Palette {
	return i.palette
}

func (i *ImageMenu) SetPaletteImage(img image.Image) {
	i.paletteImage.Image = img
	i.paletteImage.Refresh()
}

func (i *ImageMenu) PaletteImage() *canvas.Image {
	return i.paletteImage
}

func (i *ImageMenu) CpcImage() *canvas.Image {
	return i.cpcImage
}

func (i *ImageMenu) SetCpcImage(img image.Image) {
	i.cpcImage.Image = img
	i.cpcImage.FillMode = canvas.ImageFillStretch
	i.cpcImage.Refresh()
}

func (i *ImageMenu) OriginalImagePath() string {
	if i.originalImagePath == nil {
		return ""
	}
	return i.originalImagePath.Path()
}

func (i *ImageMenu) SetOriginalImagePath(path fyne.URI) {
	i.originalImagePath = path
}

func (i *ImageMenu) Width() *widget.Entry {
	return i.width
}

func (i *ImageMenu) Height() *widget.Entry {
	return i.height
}

func (i *ImageMenu) GetWidth() (int, string, error) {
	v, err := strconv.Atoi(i.width.Text)
	return v, i.width.Text, err
}

func (i *ImageMenu) GetHeight() (int, string, error) {
	v, err := strconv.Atoi(i.height.Text)
	return v, i.height.Text, err
}

// nolint: funlen
func (i *ImageMenu) CmdLine() string {
	exec, err := os.Executable()
	if err != nil {
		log.GetLogger().Error("error while getting executable path :%v\n", err)
		return exec
	}
	if i.originalImagePath != nil {
		exec += " -in " + i.originalImagePath.Path()
	}
	if i.IsCpcPlus {
		exec += " -plus"
	}
	if i.IsFullScreen {
		exec += " -fullscreen"
	}
	if i.IsSprite {
		width, err := strconv.Atoi(i.width.Text)
		if err != nil {
			log.GetLogger().Error("cannot convert width value :%s error :%v\n", i.width.Text, err)
		} else {
			exec += " -width " + strconv.Itoa(width)
		}
		height, err := strconv.Atoi(i.height.Text)
		if err != nil {
			log.GetLogger().Error("cannot convert height value :%s error :%v\n", i.height.Text, err)
		} else {
			exec += " -height " + strconv.Itoa(height)
		}
	}
	if i.IsHardSprite {
		exec += " -spritehard"
	}
	if i.ApplyDithering {
		if i.WithQuantification {
			exec += " -quantization"
		} else {
			exec += " -multiplier " + fmt.Sprintf("%.2f", i.DitheringMultiplier)
		}
		exec += " -dithering " + strconv.Itoa(i.DitheringAlgoNumber)
		// stockage du numéro d'algo
	}
	exec += " -mode " + strconv.Itoa(i.Mode)
	if i.Reducer != 0 {
		exec += " -reducer " + strconv.Itoa(i.Reducer)
	}
	// resize algo
	if i.ResizeAlgoNumber != 0 {
		exec += " -algo " + strconv.Itoa(i.ResizeAlgoNumber)
	}
	if i.Brightness != 0 {
		exec += " -brightness " + fmt.Sprintf("%.2f", i.Brightness)
	}
	if i.Saturation != 0 {
		exec += " -saturation " + fmt.Sprintf("%.2f", i.Saturation)
	}
	if i.OneLine {
		exec += " -oneline"
	}
	if i.OneRow {
		exec += " -onerow"
	}
	i.CmdLineGenerate = exec
	return exec
}

func (me *ImageMenu) SetOriginalImage(img image.Image) {
	me.originalImage.Image = img
	me.originalImage.FillMode = canvas.ImageFillContain
	me.originalImage.Refresh()
}

func (me *ImageMenu) OriginalImage() *canvas.Image {
	return me.originalImage
}

func (me *ImageMenu) SetImagePalette(i image.Image, p color.Palette) {
	me.SetCpcImage(i)
	me.SetPalette(p)
	me.Edited = true
}

// nolint:funlen, gocognit
func (me *ImageMenu) ExportImage(e *ImageExport, w fyne.Window, getCfg func(me *ImageExport, checkOriginalImage bool) *config.MartineConfig) {
	pi := wgt.NewProgressInfinite("Saving...., Please wait.", w)
	pi.Show()
	cfg := getCfg(e, true)
	if cfg == nil {
		return
	}
	if e.ExportText {
		if cfg.Overscan {
			cfg.ExportAsGoFile = true
		}

		out, _, palette, _, err := gfx.ApplyOneImage(me.OriginalImage().Image, cfg, me.Mode, me.Palette(), uint8(me.Mode))
		if err != nil {
			pi.Hide()
			dialog.ShowError(err, w)
			return
		}
		if !cfg.ExportAsGoFile {
			code := ascii.FormatAssemblyDatabyte(out, "\n")
			var palCode string
			if cfg.CpcPlus {
				palCode = ascii.FormatAssemblyCPCPlusPalette(palette, "\n")
			} else {
				palCode = ascii.FormatAssemblyCPCPalette(palette, "\n")
			}

			content := fmt.Sprintf("; Generated by martine\n; from file %s\nImage:\n%s\n\n; palette\npalette: \n%s\n ",
				me.OriginalImagePath(),
				code,
				palCode)
			filename := filepath.Base(me.OriginalImagePath())
			fileExport := e.ExportFolderPath + string(filepath.Separator) + filename + ".asm"
			if err = amsdos.SaveStringOSFile(fileExport, content); err != nil {
				pi.Hide()
				dialog.ShowError(err, w)
				return
			}
		} else {
			go1, go2, err := overscan.OverscanToGo(out)
			if err != nil {
				pi.Hide()
				dialog.ShowError(err, w)
				return
			}
			var decompressRoutine string
			if cfg.Compression != compression.NONE && cfg.Compression != -1 {
				go1, _ = compression.Compress(go1, cfg.Compression)
				go2, _ = compression.Compress(go2, cfg.Compression)
				if cfg.Compression == compression.ZX0 {
					decompressRoutine = assembly.DeltapackRoutine
				}
			}
			code1 := ascii.FormatAssemblyDatabyte(go1, "\n")
			code2 := ascii.FormatAssemblyDatabyte(go2, "\n")
			var palCode string

			if cfg.CpcPlus {
				palCode = ascii.FormatAssemblyCPCPlusPalette(palette, "\n")
			} else {
				palCode = ascii.FormatAssemblyCPCPalette(palette, "\n")
			}
			// add the compression in case of compression
			content := fmt.Sprintf("; Generated by martine\n; from file %s (part go1)\nImage_go1:\n%s\n;(part go2)\nImage_go2:\n%s\n; palette\npalette: \n%s\n %s",
				me.OriginalImagePath(),
				code1,
				code2,
				palCode,
				decompressRoutine)
			filename := filepath.Base(me.OriginalImagePath())
			fileExport := e.ExportFolderPath + string(filepath.Separator) + filename + ".asm"
			if err = amsdos.SaveStringOSFile(fileExport, content); err != nil {
				pi.Hide()
				dialog.ShowError(err, w)
				return
			}

		}
	} else {
		// palette export
		userDir, err := os.UserHomeDir()
		if err != nil {
			pi.Hide()
			dialog.ShowError(err, w)
		}
		tmpPalette := userDir + string(filepath.Separator) + "temporary_palette.kit"
		if err := impPalette.SaveKit(tmpPalette, me.Palette(), false); err != nil {
			pi.Hide()
			dialog.ShowError(err, w)
		}
		if me.UsePalette {
			cfg.KitPath = tmpPalette
		}

		filename := filepath.Base(me.OriginalImagePath())
		if me.Edited {
			err = gfx.ExportRawImage(
				me.CpcImage().Image,
				me.Palette(),
				cfg,
				filename,
				e.ExportFolderPath+string(filepath.Separator)+filename,
				uint8(me.Mode),
			)
		} else {
			err = gfx.ApplyOneImageAndExport(
				me.OriginalImage().Image,
				cfg,
				filename,
				e.ExportFolderPath+string(filepath.Separator)+filename,
				me.Mode,
				uint8(me.Mode))
		}
		os.Remove(tmpPalette)

		if err != nil {
			pi.Hide()
			dialog.NewError(err, w).Show()
			return
		}
		if cfg.Dsk {
			if err := diskimage.ImportInDsk(me.OriginalImagePath(), cfg); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		}
		if cfg.Sna {
			if cfg.Overscan {
				var gfxFile string
				for _, v := range cfg.DskFiles {
					if filepath.Ext(v) == ".SCR" {
						gfxFile = v
						break
					}
				}
				cfg.SnaPath = filepath.Join(e.ExportFolderPath, "test.sna")
				if err := snapshot.ImportInSna(gfxFile, cfg.SnaPath, uint8(me.Mode)); err != nil {
					dialog.NewError(err, w).Show()
					return
				}
			}
		}
	}
	if e.ExportToM2 {
		if err := m4.ImportInM4(cfg); err != nil {
			dialog.NewError(err, w).Show()
			log.GetLogger().Error("Cannot send to M4 error :%v\n", err)
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in folder \n"+e.ExportFolderPath, w)
}

// nolint: funlen
func (me *ImageMenu) NewConfig(ex *ImageExport, checkOriginalImage bool) *config.MartineConfig {
	if checkOriginalImage && me.OriginalImagePath() == "" {
		return nil
	}
	cfg := config.NewMartineConfig("", "")
	if checkOriginalImage {
		cfg = config.NewMartineConfig(me.OriginalImagePath(), "")
	}
	cfg.CpcPlus = me.IsCpcPlus
	cfg.Overscan = me.IsFullScreen
	cfg.DitheringMultiplier = me.DitheringMultiplier
	cfg.Brightness = me.Brightness
	cfg.Saturation = me.Saturation

	if me.Brightness > 0 && me.Saturation == 0 {
		cfg.Saturation = me.Brightness
	}
	if me.Brightness == 0 && me.Saturation > 0 {
		cfg.Brightness = me.Saturation
	}
	cfg.Reducer = me.Reducer
	cfg.Size = constants.NewSizeMode(uint8(me.Mode), me.IsFullScreen)
	if me.IsSprite {
		width, _, err := me.GetWidth()
		if err != nil {
			dialog.NewError(err, me.w).Show()
			return nil
		}
		height, _, err := me.GetHeight()
		if err != nil {
			dialog.NewError(err, me.w).Show()
			return nil
		}
		cfg.Size.Height = height
		cfg.Size.Width = width
		cfg.CustomDimension = true
	}
	if me.IsHardSprite {
		cfg.Size.Height = 16
		cfg.Size.Width = 16
	}
	if me.ApplyDithering {
		cfg.DitheringAlgo = 0
		cfg.DitheringMatrix = me.DitheringMatrix
		cfg.DitheringType = me.DitheringType
		if me.DitheringMultiplier == 0 {
			cfg.DitheringMultiplier = .1
		} else {
			cfg.DitheringMultiplier = me.DitheringMultiplier
		}
	} else {
		cfg.DitheringAlgo = -1
	}
	cfg.DitheringWithQuantification = me.WithQuantification
	cfg.OutputPath = ex.ExportFolderPath
	if checkOriginalImage {
		cfg.InputPath = me.OriginalImagePath()
	}
	cfg.Json = ex.ExportJson
	cfg.Ascii = ex.ExportText
	cfg.NoAmsdosHeader = !ex.ExportWithAmsdosHeader
	cfg.ZigZag = ex.ExportZigzag
	cfg.Compression = ex.ExportCompression
	cfg.Dsk = ex.ExportDsk
	cfg.ExportAsGoFile = ex.ExportAsGoFiles
	cfg.OneLine = me.OneLine
	cfg.OneRow = me.OneRow
	cfg.UseKmeans = me.UseKmeans
	cfg.KmeansThreshold = me.KmeansThreshold
	if cfg.UseKmeans && me.KmeansThreshold == 0 {
		cfg.KmeansThreshold = 0.01
	}
	return cfg
}

// nolint: funlen
func (me *ImageMenu) ExportDialog(ie *ImageExport) {
	m2host := widget.NewEntry()
	m2host.SetPlaceHolder("Set your M2 IP here.")
	ie.Reset()

	cont := container.NewVBox(
		container.NewGridWithRows(4,
			container.NewGridWithRows(3,
				widget.NewLabel("Container:"),
				widget.NewCheck("import all file in Dsk", func(b bool) {
					ie.ExportDsk = b
				}),
				widget.NewCheck("add amsdos header", func(b bool) {
					ie.ExportWithAmsdosHeader = b
				}),
			),
			container.NewGridWithRows(3,
				widget.NewLabel("File type:"),
				widget.NewLabel("default screen .scr file"),
				widget.NewCheck("export as Go1 and Go2 files", func(b bool) {
					ie.ExportAsGoFiles = b
				}),
			),
			container.NewGridWithRows(4,
				widget.NewLabel("Treatment to apply :"),
				widget.NewCheck("export assembly text file", func(b bool) {
					ie.ExportText = b
				}),
				widget.NewCheck("export Json file", func(b bool) {
					ie.ExportJson = b
				}),
				widget.NewCheck("apply zigzag", func(b bool) {
					ie.ExportZigzag = b
				}),
			),
			container.NewGridWithRows(2,
				widget.NewLabel("Send results :"),
				widget.NewCheck("export to M2", func(b bool) {
					ie.ExportToM2 = b
					ie.M2IP = m2host.Text
				}),
			),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					ie.ExportCompression = compression.NONE
				case "rle":
					ie.ExportCompression = compression.RLE
				case "rle 16bits":
					ie.ExportCompression = compression.RLE16
				case "Lz4 Classic":
					ie.ExportCompression = compression.LZ4
				case "Lz4 Raw":
					ie.ExportCompression = compression.RawLZ4
				case "zx0 crunch":
					ie.ExportCompression = compression.ZX0
				}
			}),
		m2host,
		widget.NewButtonWithIcon("Export into folder", theme.DocumentSaveIcon(), func() {
			fo := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, me.w)
					return
				}
				if lu == nil {
					// cancel button
					return
				}
				directory.SetExportDirectoryURI(lu)
				ie.ExportFolderPath = lu.Path()
				log.GetLogger().Infoln(ie.ExportFolderPath)
				// m.ExportOneImage(m.main)
				me.ExportImage(ie, me.w, me.NewConfig)
				// apply and export
			}, me.w)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(me.w.Content().Size())

			dl.CheckAmsdosHeaderExport(ie.ExportDsk, ie.ExportWithAmsdosHeader, fo, me.w)
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, me.w)
	d.Resize(me.w.Canvas().Size())
	d.Show()
}

// nolint:funlen, gocognit
func (i *ImageMenu) NewImportButton(modeSelection *widget.Select, callBack func()) *widget.Button {
	return widget.NewButtonWithIcon("Import", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, i.w)
				return
			}
			if reader == nil {
				return
			}
			directory.SetImportDirectoryURI(reader.URI())
			i.SetOriginalImagePath(reader.URI())
			if i.IsFullScreen {

				// open palette widget to get palette
				p, mode, err := overscan.OverscanPalette(i.OriginalImagePath())
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(fmt.Errorf("no palette found in selected file, try to normal option and open the associated palette"), i.w)
					return
				}
				img, err := ovs.OverscanToImg(i.OriginalImagePath(), mode, p)
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(errors.New("palette is empty"), i.w)
					return
				}
				i.SetPalette(p)
				i.Mode = int(mode)
				modeSelection.SetSelectedIndex(i.Mode)
				i.SetPaletteImage(png.PalToImage(p))
				i.SetOriginalImage(img)
			} else if i.IsSprite {
				// loading sprite file
				if len(i.Palette()) == 0 {
					dialog.ShowError(errors.New("palette is empty, please import palette first"), i.w)
					return
				}
				img, size, err := sprite.SpriteToImg(i.OriginalImagePath(), uint8(i.Mode), i.Palette())
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				i.Width().SetText(strconv.Itoa(size.Width))
				i.Height().SetText(strconv.Itoa(size.Height))
				i.SetOriginalImage(img)
			} else {
				// loading classical screen
				if len(i.Palette()) == 0 {
					dialog.ShowError(errors.New("palette is empty,  please import palette first, or select fullscreen option to open a fullscreen option"), i.w)
					return
				}
				if i.IsWin {
					img, err := screen.WinToImg(i.OriginalImagePath(), uint8(i.Mode), i.Palette())
					if err != nil {
						dialog.ShowError(err, i.w)
						return
					}
					i.SetOriginalImage(img)
				} else {
					// loading classical screen
					if len(i.Palette()) == 0 {
						dialog.ShowError(errors.New("palette is empty,  please import palette first, or select fullscreen option to open a fullscreen option"), i.w)
						return
					}
					if i.IsHardSprite {
						// loading hard sprite
						img, err := spritehard.SpriteHardToImg(i.OriginalImagePath(), i.Palette())
						if err != nil {
							dialog.ShowError(err, i.w)
							return
						}
						i.SetOriginalImage(img)

					} else {
						img, err := screen.ScrToImg(i.OriginalImagePath(), uint8(i.Mode), i.Palette())
						if err != nil {
							dialog.ShowError(err, i.w)
							return
						}
						i.SetOriginalImage(img)
					}
				}
			}
			if callBack != nil {
				callBack()
			}
		}, i.w)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin", ".spr"}))
		d.Resize(i.w.Content().Size())
		d.Show()
	})
}

func (me *ImageMenu) NewFormatRadio() *widget.Select {
	winFormat := widget.NewSelect([]string{"Normal", "Fullscreen", "Window", "Sprite", "Sprite Hard"}, func(s string) {
		switch s {
		case "Normal":
			me.IsFullScreen = false
			me.IsSprite = false
			me.IsWin = false
			me.IsHardSprite = false
		case "Fullscreen":
			me.IsFullScreen = true
			me.IsSprite = false
			me.IsWin = false
			me.IsHardSprite = false
		case "Sprite":
			me.IsFullScreen = false
			me.IsSprite = true
			me.IsWin = false
			me.IsHardSprite = false
		case "Sprite Hard":
			me.IsFullScreen = false
			me.IsSprite = false
			me.IsWin = false
			me.IsHardSprite = true
		case "Window":
			me.IsFullScreen = false
			me.IsSprite = false
			me.IsHardSprite = false
			me.IsWin = true
		}
	})
	winFormat.SetSelected("Normal")
	return winFormat
}
