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
	originalImage     *canvas.Image
	cpcImage          *canvas.Image
	originalImagePath fyne.URI
	Cfg               *config.MartineConfig

	width            *widget.Entry
	height           *widget.Entry
	Data             []byte
	Downgraded       *image.NRGBA
	ResizeAlgoNumber int
	paletteImage     *canvas.Image
	UsePalette       bool
	CmdLineGenerate  string

	Edited bool
	w      fyne.Window
}

func NewImageMenu() *ImageMenu {
	return &ImageMenu{
		originalImage: &canvas.Image{},
		cpcImage:      &canvas.Image{},
		paletteImage:  &canvas.Image{},
		width:         widget.NewEntry(),
		height:        widget.NewEntry(),
		Downgraded:    &image.NRGBA{},
		Cfg:           config.NewMartineConfig("", ""),
	}
}

func (i *ImageMenu) SetWindow(w fyne.Window) {
	i.w = w
}

func (i *ImageMenu) SetPalette(p color.Palette) {
	i.Cfg.PalCfg.Palette = p
	i.SetPaletteImage(png.PalToImage(i.Palette()))
	i.paletteImage.Refresh()
}

func (i *ImageMenu) Palette() color.Palette {
	return i.Cfg.PalCfg.Palette
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
	if i.Cfg.ScrCfg.IsPlus {
		exec += " -plus"
	}

	if i.Cfg.ScrCfg.Type.IsFullScreen() {
		exec += " -fullscreen"
	}
	if i.Cfg.ScrCfg.Type.IsSprite() {
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
	if i.Cfg.ScrCfg.Type.IsSpriteHard() {
		exec += " -spritehard"
	}
	if i.Cfg.ScrCfg.Process.ApplyDithering {
		if i.Cfg.ScrCfg.Process.Dithering.WithQuantification {
			exec += " -quantization"
		} else {
			exec += " -multiplier " + fmt.Sprintf("%.2f", i.Cfg.ScrCfg.Process.Dithering.Multiplier)
		}
		exec += " -dithering " + strconv.Itoa(i.Cfg.ScrCfg.Process.Dithering.Algo)
		// stockage du numéro d'algo
	}
	exec += " -mode " + strconv.Itoa(int(i.Cfg.ScrCfg.Mode))
	if i.Cfg.ScrCfg.Process.Reducer != 0 {
		exec += " -reducer " + strconv.Itoa(i.Cfg.ScrCfg.Process.Reducer)
	}
	// resize algo
	if i.ResizeAlgoNumber != 0 {
		exec += " -algo " + strconv.Itoa(i.ResizeAlgoNumber)
	}
	if i.Cfg.ScrCfg.Process.Brightness != 0 {
		exec += " -brightness " + fmt.Sprintf("%.2f", i.Cfg.ScrCfg.Process.Brightness)
	}
	if i.Cfg.ScrCfg.Process.Saturation != 0 {
		exec += " -saturation " + fmt.Sprintf("%.2f", i.Cfg.ScrCfg.Process.Saturation)
	}
	if i.Cfg.ScrCfg.Process.OneLine {
		exec += " -oneline"
	}
	if i.Cfg.ScrCfg.Process.OneRow {
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

func (me *ImageMenu) GetConfig(checkOriginalImage bool) *config.MartineConfig {
	return me.Cfg
}

// nolint:funlen, gocognit
func (me *ImageMenu) ExportImage(w fyne.Window, getCfg func(checkOriginalImage bool) *config.MartineConfig) {
	pi := wgt.NewProgressInfinite("Saving...., Please wait.", w)
	pi.Show()
	cfg := getCfg(true)
	if cfg == nil {
		return
	}
	cfg.ResetDskFiles()
	if cfg.ScrCfg.IsExport(config.AssemblyExport) {
		if cfg.ScrCfg.Type == config.FullscreenFormat {
			cfg.ScrCfg.AddExport(config.GoImpdrawExport)
		}

		out, _, palette, _, err := gfx.ApplyOneImage(me.OriginalImage().Image, cfg, int(me.Cfg.ScrCfg.Mode), me.Palette(), me.Cfg.ScrCfg.Mode)
		if err != nil {
			pi.Hide()
			dialog.ShowError(err, w)
			return
		}
		if !cfg.ScrCfg.IsExport(config.GoImpdrawExport) {
			code := ascii.FormatAssemblyDatabyte(out, "\n")
			var palCode string
			if cfg.ScrCfg.IsPlus {
				palCode = ascii.FormatAssemblyCPCPlusPalette(palette, "\n")
			} else {
				palCode = ascii.FormatAssemblyCPCPalette(palette, "\n")
			}

			content := fmt.Sprintf("; Generated by martine\n; from file %s\nImage:\n%s\n\n; palette\npalette: \n%s\n ",
				me.OriginalImagePath(),
				code,
				palCode)
			filename := filepath.Base(me.OriginalImagePath())
			fileExport := filepath.Join(cfg.ScrCfg.OutputPath, filename+".asm")
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
			if cfg.ScrCfg.Compression != compression.NONE && cfg.ScrCfg.Compression != -1 {
				go1, _ = compression.Compress(go1, cfg.ScrCfg.Compression)
				go2, _ = compression.Compress(go2, cfg.ScrCfg.Compression)
				if cfg.ScrCfg.Compression == compression.ZX0 {
					decompressRoutine = assembly.DeltapackRoutine
				}
			}
			code1 := ascii.FormatAssemblyDatabyte(go1, "\n")
			code2 := ascii.FormatAssemblyDatabyte(go2, "\n")
			var palCode string

			if cfg.ScrCfg.IsPlus {
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
			fileExport := filepath.Join(cfg.ScrCfg.OutputPath, filename+".asm")
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
		tmpPalette := filepath.Join(userDir, "temporary_palette.kit")
		if err := impPalette.SaveKit(tmpPalette, me.Palette(), false); err != nil {
			pi.Hide()
			dialog.ShowError(err, w)
		}
		if me.UsePalette {
			cfg.PalCfg.Path = tmpPalette
		}

		filename := filepath.Base(me.OriginalImagePath())
		if me.Edited {
			err = gfx.ExportRawImage(
				me.CpcImage().Image,
				me.Palette(),
				cfg,
				filename,
				filepath.Join(cfg.ScrCfg.OutputPath, filename),
				me.Cfg.ScrCfg.Mode,
			)
		} else {
			err = gfx.ApplyOneImageAndExport(
				me.OriginalImage().Image,
				cfg,
				filename,
				filepath.Join(cfg.ScrCfg.OutputPath, filename),
				me.Cfg.ScrCfg.Mode)
		}
		os.Remove(tmpPalette)

		if err != nil {
			pi.Hide()
			dialog.NewError(err, w).Show()
			return
		}
		if cfg.HasContainerExport(config.DskContainer) {
			if err := diskimage.ImportInDsk(filepath.Join(cfg.ScrCfg.OutputPath, "IMG"), cfg); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		}
		if cfg.HasContainerExport(config.SnaContainer) {
			if cfg.ScrCfg.Type == config.FullscreenFormat {
				var gfxFile string
				for _, v := range cfg.DskFiles {
					if filepath.Ext(v) == ".SCR" {
						gfxFile = v
						break
					}
				}
				cfg.ContainerCfg.Path = filepath.Join(cfg.ContainerCfg.Path, "test.sna")
				if err := snapshot.ImportInSna(gfxFile, cfg.ContainerCfg.Path, me.Cfg.ScrCfg.Mode); err != nil {
					dialog.NewError(err, w).Show()
					return
				}
			}
		}
	}
	if cfg.M4cfg.Enabled {
		if err := m4.ImportInM4(cfg); err != nil {
			dialog.NewError(err, w).Show()
			log.GetLogger().Error("Cannot send to M4 error :%v\n", err)
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in folder \n"+cfg.ScrCfg.OutputPath, w)
}

// nolint: funlen
func (me *ImageMenu) NewConfig(checkOriginalImage bool) *config.MartineConfig {
	if checkOriginalImage && me.OriginalImagePath() == "" {
		return me.Cfg
	}

	if checkOriginalImage {
		me.Cfg.ScrCfg.InputPath = me.OriginalImagePath()
	}
	if me.Cfg.ScrCfg.Process.Brightness > 0 && me.Cfg.ScrCfg.Process.Saturation == 0 {
		me.Cfg.ScrCfg.Process.Saturation = me.Cfg.ScrCfg.Process.Brightness
	}
	if me.Cfg.ScrCfg.Process.Brightness == 0 && me.Cfg.ScrCfg.Process.Saturation > 0 {
		me.Cfg.ScrCfg.Process.Brightness = me.Cfg.ScrCfg.Process.Saturation
	}
	me.Cfg.ScrCfg.Size = constants.NewSizeMode(me.Cfg.ScrCfg.Mode, me.Cfg.ScrCfg.Type.IsFullScreen())
	if me.Cfg.ScrCfg.Type.IsSprite() {
		width, _, err := me.GetWidth()
		if err != nil {
			dialog.NewError(err, me.w).Show()
			return me.Cfg
		}
		height, _, err := me.GetHeight()
		if err != nil {
			dialog.NewError(err, me.w).Show()
			return me.Cfg
		}
		me.Cfg.ScrCfg.Size.Height = height
		me.Cfg.ScrCfg.Size.Width = width
		me.Cfg.CustomDimension = true
	}
	if me.Cfg.ScrCfg.Type.IsSpriteHard() {
		me.Cfg.ScrCfg.Size.Height = 16
		me.Cfg.ScrCfg.Size.Width = 16
	}
	if me.Cfg.ScrCfg.Process.ApplyDithering {
		me.Cfg.ScrCfg.Process.Dithering.Algo = 0
		if me.Cfg.ScrCfg.Process.Dithering.Multiplier == 0 {
			me.Cfg.ScrCfg.Process.Dithering.Multiplier = .1
		}
	} else {
		me.Cfg.ScrCfg.Process.Dithering.Algo = -1
	}
	if checkOriginalImage {
		me.Cfg.ScrCfg.InputPath = me.OriginalImagePath()
	}
	if me.Cfg.ScrCfg.Process.Kmeans.Used && me.Cfg.ScrCfg.Process.Kmeans.Threshold == 0 {
		me.Cfg.ScrCfg.Process.Kmeans.Threshold = 0.01
	}
	return me.Cfg
}

// nolint: funlen
func (me *ImageMenu) ExportDialog(cfg *config.MartineConfig, getCfg func(checkOriginalImage bool) *config.MartineConfig) {
	m2host := widget.NewEntry()
	m2host.SetPlaceHolder("Set your M2 IP here.")

	cont := container.NewVBox(
		container.NewGridWithRows(4,
			container.NewGridWithRows(3,
				widget.NewLabel("Container:"),
				widget.NewCheck("import all file in Dsk", func(b bool) {
					if b {
						cfg.ContainerCfg.AddExport(config.DskContainer)
					} else {
						cfg.ContainerCfg.RemoveExport(config.DskContainer)
					}
				}),
				widget.NewCheck("add amsdos header", func(b bool) {
					cfg.ScrCfg.NoAmsdosHeader = !b
				}),
			),
			container.NewGridWithRows(3,
				widget.NewLabel("File type:"),
				widget.NewLabel("default screen .scr file"),
				widget.NewCheck("export as Go1 and Go2 files", func(b bool) {
					if b {
						cfg.ScrCfg.AddExport(config.GoImpdrawExport)
					} else {
						cfg.ScrCfg.RemoveExport(config.GoImpdrawExport)
					}
				}),
			),
			container.NewGridWithRows(4,
				widget.NewLabel("Treatment to apply :"),
				widget.NewCheck("export assembly text file", func(b bool) {
					if b {
						cfg.ScrCfg.AddExport(config.AssemblyExport)
					} else {
						cfg.ScrCfg.RemoveExport(config.AssemblyExport)
					}
				}),
				widget.NewCheck("export Json file", func(b bool) {
					if b {
						cfg.ScrCfg.AddExport(config.JsonExport)
					} else {
						cfg.ScrCfg.RemoveExport(config.JsonExport)
					}
				}),
				widget.NewCheck("apply zigzag", func(b bool) {
					cfg.ZigZag = b
				}),
			),
			container.NewGridWithRows(2,
				widget.NewLabel("Send results :"),
				widget.NewCheck("export to M2", func(b bool) {
					cfg.M4cfg.Enabled = b
					cfg.M4cfg.Host = m2host.Text
				}),
			),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					cfg.ScrCfg.Compression = compression.NONE
				case "rle":
					cfg.ScrCfg.Compression = compression.RLE
				case "rle 16bits":
					cfg.ScrCfg.Compression = compression.RLE16
				case "Lz4 Classic":
					cfg.ScrCfg.Compression = compression.LZ4
				case "Lz4 Raw":
					cfg.ScrCfg.Compression = compression.RawLZ4
				case "zx0 crunch":
					cfg.ScrCfg.Compression = compression.ZX0
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
				cfg.ScrCfg.OutputPath = lu.Path()
				log.GetLogger().Infoln(cfg.ScrCfg.OutputPath)
				// m.ExportOneImage(m.main)
				me.ExportImage(me.w, getCfg)
				// apply and export
			}, me.w)
			d, err := directory.ExportDirectoryURI()
			if err == nil {
				fo.SetLocation(d)
			}
			fo.Resize(me.w.Content().Size())

			dl.CheckAmsdosHeaderExport(cfg.ContainerCfg.HasExport(config.DskContainer), !cfg.ScrCfg.NoAmsdosHeader, fo, me.w)
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
			if !i.Cfg.ScrCfg.Type.IsFullScreen() {
				// load palette from file
				if len(i.Palette()) == 0 {
					dialog.ShowError(errors.New("palette is empty, please import palette first"), i.w)
					return
				}
			}
			switch i.Cfg.ScrCfg.Type {
			case config.FullscreenFormat:
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
				i.Cfg.ScrCfg.Mode = mode
				modeSelection.SetSelectedIndex(int(i.Cfg.ScrCfg.Mode))
				i.SetPaletteImage(png.PalToImage(p))
				i.SetOriginalImage(img)
			case config.SpriteFormat:
				// loading sprite file
				img, size, err := sprite.SpriteToImg(i.OriginalImagePath(), i.Cfg.ScrCfg.Mode, i.Palette())
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				i.Width().SetText(strconv.Itoa(size.Width))
				i.Height().SetText(strconv.Itoa(size.Height))
				i.SetOriginalImage(img)
			case config.WindowFormat:
				img, err := screen.WinToImg(i.OriginalImagePath(), i.Cfg.ScrCfg.Mode, i.Palette())
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				i.SetOriginalImage(img)
			case config.SpriteHardFormat:
				img, err := spritehard.SpriteHardToImg(i.OriginalImagePath(), i.Palette())
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				i.SetOriginalImage(img)
			case config.OcpScreenFormat:
				img, err := screen.ScrToImg(i.OriginalImagePath(), i.Cfg.ScrCfg.Mode, i.Palette())
				if err != nil {
					dialog.ShowError(err, i.w)
					return
				}
				i.SetOriginalImage(img)
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
			me.Cfg.ScrCfg.Type = config.OcpScreenFormat
		case "Fullscreen":
			me.Cfg.ScrCfg.Type = config.FullscreenFormat
		case "Sprite":
			me.Cfg.ScrCfg.Type = config.SpriteFormat
		case "Sprite Hard":
			me.Cfg.ScrCfg.Type = config.SpriteHardFormat
		case "Window":
			me.Cfg.ScrCfg.Type = config.WindowFormat
		}
	})
	winFormat.SetSelected("Normal")
	return winFormat
}
