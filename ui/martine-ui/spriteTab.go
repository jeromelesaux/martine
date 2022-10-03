package ui

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert"
	"github.com/jeromelesaux/martine/gfx/sprite"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) ApplySprite(s *menu.SpriteMenu) {
	pi := dialog.NewProgressInfinite("Computing....", "Please wait.", m.window)
	pi.Show()
	var colorsAvailable int
	switch s.Mode {
	case 0:
		colorsAvailable = constants.Mode0.ColorsAvailable
	case 1:
		colorsAvailable = constants.Mode1.ColorsAvailable
	case 2:
		colorsAvailable = constants.Mode2.ColorsAvailable
	}
	b := s.OriginalBoard.Image
	img := image.NewNRGBA(image.Rect(0, 0, b.Bounds().Dx(), b.Bounds().Dx()))
	draw.Draw(img, img.Bounds(), b, b.Bounds().Min, draw.Src)
	pal, _, err := convert.DowngradingPalette(img, constants.Size{ColorsAvailable: colorsAvailable, Width: img.Bounds().Max.X, Height: img.Bounds().Max.Y}, s.IsCpcPlus)
	if err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	s.Palette = pal
	size := constants.Size{Width: s.SpriteWidth, Height: s.SpriteHeight}
	raw, sprites, err := sprite.SplitBoardToSprite(s.OriginalBoard.Image, s.Palette, s.SpriteNumberPerRow, s.SpriteNumberPerColumn, uint8(s.Mode), s.IsHardSprite, size)
	if err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	s.SpritesCollection = sprites
	s.SpritesData = raw

	imagesCanvas := custom_widget.NewImageTableCache(s.SpriteNumberPerRow, s.SpriteNumberPerColumn, fyne.NewSize(50, 50))

	for y := 0; y < s.SpriteNumberPerRow; y++ {
		for x := 0; x < s.SpriteNumberPerColumn; x++ {
			imagesCanvas.Set(y, x, canvas.NewImageFromImage(s.SpritesCollection[x][y]))
		}
	}
	s.OriginalImages.Update(imagesCanvas, s.SpriteNumberPerRow, s.SpriteNumberPerColumn)
	pi.Hide()
	refreshUI.OnTapped()
}

func (m *MartineUI) newSpriteTab(s *menu.SpriteMenu) fyne.CanvasObject {

	forceUIRefresh := widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
		s := m.window.Content().Size()
		s.Height += 10.
		s.Width += 10.
		m.window.Resize(s)
		m.window.Resize(m.window.Content().Size())
		m.window.Content().Refresh()
	})
	refreshUI = forceUIRefresh

	openFileWidget := widget.NewButton("Image", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}

			img, err := openImage(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			canvasImg := canvas.NewImageFromImage(img)
			s.OriginalBoard = *canvas.NewImageFromImage(canvasImg.Image)
			s.OriginalBoard.FillMode = canvas.ImageFillContain
			m.window.Canvas().Refresh(&s.OriginalBoard)
			m.window.Resize(m.window.Content().Size())
		}, m.window)
		d.SetFilter(imagesFilesFilter)
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportSpriteBoard(m.sprite, m.window)
	})

	applyButton := widget.NewButtonWithIcon("Apply", theme.VisibilityIcon(), func() {
		fmt.Println("apply.")
		m.ApplySprite(s)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(v string) {
		mode, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", v)
		}
		s.Mode = mode
	})

	modeLabel := widget.NewLabel("Mode:")

	spriteNumberPerRowLabel := widget.NewLabel("Number of sprite per row")
	spriteNumberPerColumnLabel := widget.NewLabel("Number of sprite per column")
	spriteNumberPerColumnEntry := widget.NewEntry()
	spriteNumberPerColumnEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteNumberPerColumn = r
	}
	spriteNumberPerRowEntry := widget.NewEntry()
	spriteNumberPerRowEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteNumberPerRow = r
	}

	spriteWidthSizeLabel := widget.NewLabel("sprite width")
	spriteHeightSizeLabel := widget.NewLabel("sprite height")
	spriteWidthSizeEntry := widget.NewEntry()
	spriteWidthSizeEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteWidth = r
	}

	spriteHeightSizeEntry := widget.NewEntry()
	spriteHeightSizeEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteHeight = r
	}

	isSpriteHard := widget.NewCheck("is sprite hard", func(b bool) {
		s.IsHardSprite = b
		if s.IsHardSprite {
			s.SpriteHeight = 16
			s.SpriteWidth = 16
		}
	})

	IsCpcPlus := widget.NewCheck("is CPC Plus", func(b bool) {
		s.IsCpcPlus = b
	})

	paletteOpen := NewOpenPaletteButton(s, m.window)
	s.PaletteImage = canvas.Image{}

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			&s.OriginalBoard,
			s.OriginalImages,
		),
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.New(
				layout.NewVBoxLayout(),
				container.New(
					layout.NewHBoxLayout(),
					openFileWidget,
					applyButton,
					exportButton,
					paletteOpen,
					refreshUI,
				),
				container.New(
					layout.NewHBoxLayout(),
					modeLabel,
					modes,
				),
				container.New(
					layout.NewHBoxLayout(),
					spriteNumberPerRowLabel,
					spriteNumberPerColumnEntry,
				),
				container.New(
					layout.NewHBoxLayout(),
					spriteNumberPerColumnLabel,
					spriteNumberPerRowEntry,
				),
				container.New(
					layout.NewHBoxLayout(),
					spriteWidthSizeLabel,
					spriteWidthSizeEntry,
				),

				container.New(
					layout.NewHBoxLayout(),
					spriteHeightSizeLabel,
					spriteHeightSizeEntry,
				),
				container.New(
					layout.NewHBoxLayout(),
					isSpriteHard,
				),
				container.New(
					layout.NewHBoxLayout(),
					IsCpcPlus,
				),
			),
			container.New(
				layout.NewGridLayoutWithRows(2),
				container.New(
					layout.NewVBoxLayout(),
					widget.NewLabel("Palette"),
				),
				container.New(
					layout.NewGridLayoutWithRows(2),
					&s.PaletteImage,
				),
			),
		),
	)
}
