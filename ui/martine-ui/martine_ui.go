package ui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

var (
	modeSelection     *widget.Select
	dialogSize        = fyne.NewSize(800, 800)
	savingDialogSize  = fyne.NewSize(800, 800)
	imagesFilesFilter = storage.NewExtensionFileFilter([]string{".jpg", ".gif", ".png", ".jpeg", ".JPG", ".JPEG", ".GIF", ".PNG"})
	appPrefix         = fmt.Sprintf("Martine (%v)", common.AppVersion)
)

type MartineUI struct {
	window  fyne.Window
	main    *menu.ImageMenu
	tilemap *menu.TilemapMenu
	animate *menu.AnimateMenu
	egx     *menu.DoubleImageMenu
	sprite  *menu.SpriteMenu

	imageExport   *menu.ImageExport
	tilemapExport *menu.ImageExport
	animateExport *menu.AnimateExport
	egxExport     *menu.ImageExport
}

func NewMartineUI() *MartineUI {
	m := &MartineUI{
		main:          menu.NewImageMenu(),
		tilemap:       menu.NewTilemapMenu(),
		animate:       menu.NewAnimateMenu(),
		egx:           menu.NewDoubleImageMenu(),
		sprite:        menu.NewSpriteMenu(),
		imageExport:   &menu.ImageExport{},
		tilemapExport: &menu.ImageExport{},
		animateExport: &menu.AnimateExport{},
		egxExport:     &menu.ImageExport{},
	}
	m.animateExport.ExportCompression = -1
	return m
}

func (m *MartineUI) SetPalette(p color.Palette) {
	m.main.SetPalette(p)
	m.main.SetPaletteImage(png.PalToImage(p))
}

func (m *MartineUI) Load(app fyne.App) {
	_, err := log.InitLoggerWithFile(appPrefix)
	if err != nil {
		panic(err)
	}
	m.window = app.NewWindow("Martine @IMPact v" + common.AppVersion)
	m.window.SetContent(m.NewTabs())
	m.window.Resize(fyne.NewSize(1400, 1000))
	m.window.SetTitle("Martine @IMPact v" + common.AppVersion)
	m.window.Show()
}

func (m *MartineUI) NewTabs() *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Image", m.newImageTransfertTab(m.main)),
		container.NewTabItem("Egx", m.newEgxTab(m.egx)),
		container.NewTabItem("Tile", m.newTilemapTab(m.tilemap)),
		container.NewTabItem("Animate", m.newAnimateTab(m.animate)),
		container.NewTabItem("Sprite Board", m.newSpriteTab(m.sprite)),
		container.NewTabItem("Greedings", m.newGreedings()),
	)
}

func (m *MartineUI) NewConfig(me *menu.ImageMenu, checkOriginalImage bool) *config.MartineConfig {
	if checkOriginalImage && me.OriginalImagePath() == "" {
		return nil
	}
	var cfg *config.MartineConfig
	if checkOriginalImage {
		cfg = config.NewMartineConfig(me.OriginalImagePath(), "")
	} else {
		cfg = config.NewMartineConfig("", "")
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
			dialog.NewError(err, m.window).Show()
			return nil
		}
		height, _, err := me.GetHeight()
		if err != nil {
			dialog.NewError(err, m.window).Show()
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
	} else {
		cfg.DitheringAlgo = -1
	}
	cfg.DitheringWithQuantification = me.WithQuantification
	cfg.OutputPath = m.imageExport.ExportFolderPath
	if checkOriginalImage {
		cfg.InputPath = me.OriginalImagePath()
	}
	cfg.Json = m.imageExport.ExportJson
	cfg.Ascii = m.imageExport.ExportText
	cfg.NoAmsdosHeader = !m.imageExport.ExportWithAmsdosHeader
	cfg.ZigZag = m.imageExport.ExportZigzag
	cfg.Compression = m.imageExport.ExportCompression
	cfg.Dsk = m.imageExport.ExportDsk
	cfg.ExportAsGoFile = m.imageExport.ExportAsGoFiles
	cfg.OneLine = me.OneLine
	cfg.OneRow = me.OneRow
	cfg.UseKmeans = me.UseKmeans
	cfg.KmeansInterations = me.KmeansIteration
	if cfg.UseKmeans && me.KmeansIteration == 0 {
		cfg.KmeansInterations = 100
	}

	return cfg
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	return i, err
}
