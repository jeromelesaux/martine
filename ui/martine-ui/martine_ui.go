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
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/common"
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
	scale   float32
	variant fyne.ThemeVariant
	window  fyne.Window
	main    *menu.ImageMenu
	tilemap *menu.TilemapMenu
	animate *menu.AnimateMenu
	egx     *menu.DoubleImageMenu
	sprite  *menu.SpriteMenu
	editor  *menu.Editor
}

func NewMartineUI() *MartineUI {

	m := &MartineUI{
		scale:   fyne.CurrentApp().Settings().Scale(),
		variant: fyne.CurrentApp().Settings().ThemeVariant(),
		main:    menu.NewImageMenu(),
		tilemap: menu.NewTilemapMenu(),
		animate: menu.NewAnimateMenu(),
		egx:     menu.NewDoubleImageMenu(),
		sprite:  menu.NewSpriteMenu(),
		editor:  menu.NewEditor(),
	}

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
	m.window.Resize(fyne.NewSize(1000, 600))
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
		container.NewTabItem("Editor", m.editor.New(m.window)),
		container.NewTabItem("Greetings", m.newGreedings()),
	)
}

func (m *MartineUI) SetOriginalImagePath(path fyne.URI) {
	m.main.SetOriginalImagePath(path)
	m.egx.LeftImage.SetOriginalImagePath(path)
	m.egx.RightImage.SetOriginalImagePath(path)
	m.tilemap.SetOriginalImagePath(path)
}

func (m *MartineUI) SetImage(img image.Image) {
	m.main.SetOriginalImage(img)
	m.egx.LeftImage.SetOriginalImage(img)
	m.egx.RightImage.SetOriginalImage(img)
	m.tilemap.SetOriginalImage(img)
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
