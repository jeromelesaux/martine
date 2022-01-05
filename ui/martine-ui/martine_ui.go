package ui

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

var (
	refreshUI        *widget.Button
	modeSelection    *widget.Select
	paletteSelection *widget.Select
	dialogSize       = fyne.NewSize(800, 800)
)

type MartineUI struct {
	window  fyne.Window
	main    *menu.ImageMenu
	tilemap *menu.TilemapMenu

	exportDsk              bool
	exportText             bool
	exportWithAmsdosHeader bool
	exportZigzag           bool
	exportJson             bool
	exportCompression      int
	exportFolderPath       string
	m2IP                   string
	exportToM2             bool
}

func NewMartineUI() *MartineUI {
	return &MartineUI{
		main:    &menu.ImageMenu{},
		tilemap: &menu.TilemapMenu{},
	}
}

func (m *MartineUI) SetPalette(p color.Palette) {

	m.main.Palette = p
	m.main.PaletteImage = *canvas.NewImageFromImage(file.PalToImage(p))
	refreshUI.OnTapped()
}

func (m *MartineUI) Load(app fyne.App) {
	m.window = app.NewWindow("Martine @IMPact v" + common.AppVersion)
	m.window.SetContent(m.NewTabs())
	m.window.Resize(fyne.NewSize(1400, 1000))
	m.window.SetTitle("Martine @IMPact v" + common.AppVersion)
	m.window.Show()
}

func (m *MartineUI) NewTabs() *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Image", m.newImageTransfertTab(m.main)),
		container.NewTabItem("Tile", m.newTilemapTab(m.tilemap)),
	)
}

func (m *MartineUI) NewContext(me *menu.ImageMenu) *export.MartineContext {
	if me.OriginalImagePath == nil {
		return nil
	}
	context := export.NewMartineContext(me.OriginalImagePath.Path(), "")
	context.CpcPlus = me.IsCpcPlus
	context.Overscan = me.IsFullScreen
	context.DitheringMultiplier = me.DitheringMultiplier
	context.Brightness = me.Brightness
	context.Saturation = me.Saturation
	if me.Brightness > 0 && me.Saturation == 0 {
		context.Saturation = me.Brightness
	}
	if me.Brightness == 0 && me.Saturation > 0 {
		context.Brightness = me.Saturation
	}
	context.Reducer = me.Reducer
	var size constants.Size
	switch me.Mode {
	case 0:
		size = constants.Mode0
		if me.IsFullScreen {
			size = constants.OverscanMode0
		}
	case 1:
		size = constants.Mode1
		if me.IsFullScreen {
			size = constants.OverscanMode1
		}
	case 2:
		size = constants.Mode2
		if me.IsFullScreen {
			size = constants.OverscanMode2
		}
	}
	context.Size = size
	if me.IsSprite {
		width, err := strconv.Atoi(me.Width.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return nil
		}
		height, err := strconv.Atoi(me.Height.Text)
		if err != nil {
			dialog.NewError(err, m.window).Show()
			return nil
		}
		context.Size.Height = height
		context.Size.Width = width
	}
	if me.IsHardSprite {
		context.Size.Height = 16
		context.Size.Width = 16
	}

	if me.ApplyDithering {
		context.DitheringAlgo = 0
		context.DitheringMatrix = me.DitheringMatrix
		context.DitheringType = me.DitheringType
	} else {
		context.DitheringAlgo = -1
	}
	context.DitheringWithQuantification = me.WithQuantification
	context.OutputPath = m.exportFolderPath
	context.InputPath = me.OriginalImagePath.Path()
	context.Json = m.exportJson
	context.Ascii = m.exportText
	context.NoAmsdosHeader = !m.exportWithAmsdosHeader
	context.ZigZag = m.exportZigzag
	context.Compression = m.exportCompression
	context.Dsk = m.exportDsk
	return context
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
