package menu

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/screen"
	ovs "github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/log"
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

// nolint: funlen, gocognit
func (me *ImageMenu) NewImportButton(dialogSize fyne.Size, modeSelection *widget.Select, refreshUI *widget.Button, win fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("Import", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}
			me.originalImagePath = reader.URI()
			if me.IsFullScreen {

				// open palette widget to get palette
				p, mode, err := overscan.OverscanPalette(me.originalImagePath.Path())
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(fmt.Errorf("no palette found in selected file, try to normal option and open the associated palette"), win)
					return
				}
				img, err := ovs.OverscanToImg(me.originalImagePath.Path(), mode, p)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if len(p) == 0 {
					dialog.ShowError(errors.New("palette is empty"), win)
					return
				}
				me.palette = p
				me.Mode = int(mode)
				modeSelection.SetSelectedIndex(me.Mode)

				me.SetPaletteImage(png.PalToImage(p))
				me.SetOriginalImage(img)
			} else if me.IsSprite {
				// loading sprite file
				if len(me.palette) == 0 {
					dialog.ShowError(errors.New("palette is empty, please import palette first"), win)
					return
				}
				img, size, err := sprite.SpriteToImg(me.originalImagePath.Path(), uint8(me.Mode), me.Palette())
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				me.width.SetText(strconv.Itoa(size.Width))
				me.height.SetText(strconv.Itoa(size.Height))
				me.SetOriginalImage(img)
			} else {
				// loading classical screen
				if len(me.Palette()) == 0 {
					dialog.ShowError(errors.New("palette is empty,  please import palette first, or select fullscreen option to open a fullscreen option"), win)
					return
				}
				img, err := screen.ScrToImg(me.originalImagePath.Path(), uint8(me.Mode), me.Palette())
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				me.SetOriginalImage(img)
			}
			refreshUI.OnTapped()
		}, win)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".scr", ".win", ".bin"}))
		d.Resize(dialogSize)
		d.Show()
	})
}
