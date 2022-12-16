package ui

import (
	"errors"
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
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/custom_widget"
	"github.com/jeromelesaux/martine/constants"

	ci "github.com/jeromelesaux/martine/convert/image"
	spr "github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/impdraw/tile"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/export/spritehard"
	"github.com/jeromelesaux/martine/gfx/sprite"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) ApplySprite(s *menu.SpriteMenu) {

	if s.SpriteColumns == 0 || s.SpriteRows == 0 {
		dialog.NewError(errors.New("number of sprites per row or column are not set"), m.window).Show()
		return
	}
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
	if b == nil {
		pi.Hide()
		return
	}
	img := image.NewNRGBA(image.Rect(0, 0, b.Bounds().Max.X, b.Bounds().Max.Y))
	draw.Draw(img, img.Bounds(), b, b.Bounds().Min, draw.Src)
	pal, _, err := ci.DowngradingPalette(img, constants.Size{ColorsAvailable: colorsAvailable, Width: img.Bounds().Max.X, Height: img.Bounds().Max.Y}, s.IsCpcPlus)
	if err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	s.Palette = pal
	size := constants.Size{Width: s.SpriteWidth, Height: s.SpriteHeight}
	raw, sprites, err := sprite.SplitBoardToSprite(s.OriginalBoard.Image, s.Palette, s.SpriteColumns, s.SpriteRows, uint8(s.Mode), s.IsHardSprite, size)
	if err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	s.SpritesCollection = sprites
	s.SpritesData = raw

	icache := custom_widget.NewImageTableCache(s.SpriteColumns, s.SpriteRows, fyne.NewSize(50, 50))

	for x := 0; x < s.SpriteColumns; x++ {
		for y := 0; y < s.SpriteRows; y++ {
			icache.Set(x, y, canvas.NewImageFromImage(s.SpritesCollection[x][y]))
		}
	}
	s.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)
	s.PaletteImage.Image = png.PalToImage(s.Palette)
	s.PaletteImage.Refresh()
	pi.Hide()
}

func (m *MartineUI) newSpriteTab(s *menu.SpriteMenu) fyne.CanvasObject {

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

			s.OriginalBoard.Image = img
			s.OriginalBoard.FillMode = canvas.ImageFillContain
			s.OriginalBoard.Refresh()
			//m.window.Canvas().Refresh(&s.OriginalBoard)
			//m.window.Resize(m.window.Content().Size())
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
		s.SpriteRows = r
	}
	spriteNumberPerRowEntry := widget.NewEntry()
	spriteNumberPerRowEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteColumns = r
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
			s.Mode = 0
		}
	})

	IsCpcPlus := widget.NewCheck("is CPC Plus", func(b bool) {
		s.IsCpcPlus = b
	})

	paletteOpen := NewOpenPaletteButton(s, m.window)
	importOpen := ImportSpriteBoard(m)

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			s.OriginalBoard,
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
					importOpen,
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
					s.PaletteImage,
				),
			),
		),
	)
}

func ImportSpriteBoard(m *MartineUI) fyne.Widget {
	return widget.NewButtonWithIcon("Import", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			filePath := reader.URI()
			if m.sprite.IsHardSprite {
				spritesHard, err := spritehard.OpenSpr(filePath.Path())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				m.sprite.Mode = 0
				m.sprite.SpriteHeight = 16
				m.sprite.SpriteWidth = 16
				m.sprite.SpriteColumns = 8
				var row, col int
				nbRow := len(spritesHard.Data) / m.sprite.SpriteColumns
				if len(spritesHard.Data)%m.sprite.SpriteColumns != 0 {
					nbRow++
				}
				m.sprite.SpriteRows = nbRow
				m.sprite.SpritesCollection = make([][]*image.NRGBA, nbRow)
				m.sprite.SpritesData = make([][][]byte, nbRow)

				for i := 0; i < nbRow; i++ {
					m.sprite.SpritesCollection[i] = make([]*image.NRGBA, m.sprite.SpriteColumns)
					m.sprite.SpritesData[i] = make([][]byte, m.sprite.SpriteColumns)
				}

				for i := 0; i < len(spritesHard.Data); i++ {
					m.sprite.SpritesData[row][col] = append(m.sprite.SpritesData[row][col], spritesHard.Data[i].Data[:]...)
					m.sprite.SpritesCollection[row][col] = spritesHard.Data[i].Image(m.sprite.Palette)
					col++
					if col%m.sprite.SpriteColumns == 0 {
						col = 0
						row++
					}
				}

				icache := custom_widget.NewImageTableCache(m.sprite.SpriteRows, m.sprite.SpriteColumns, fyne.NewSize(50, 50))

				for y := 0; y < m.sprite.SpriteColumns; y++ {
					for x := 0; x < m.sprite.SpriteRows; x++ {
						if m.sprite.SpritesCollection[x][y] != nil {
							icache.Set(x, y, canvas.NewImageFromImage(m.sprite.SpritesCollection[x][y]))
						}
					}
				}
				m.sprite.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)
				//len(spritesHard.Data)/m.sprite.SpriteNumberPerColumn
			} else {
				// load and display .imp file content
				mode := m.sprite.Mode
				footer, err := tile.OpenImp(filePath.Path(), mode)
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}

				m.sprite.SpriteWidth = int(footer.Width)
				m.sprite.SpriteHeight = int(footer.Height)
				m.sprite.SpriteColumns = 8
				data, err := tile.RawImp(filePath.Path())
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}

				spriteLength := int(footer.Height) * int(footer.Width)
				nbRow := (len(data) / spriteLength) / m.sprite.SpriteColumns
				m.sprite.SpriteRows = nbRow
				m.sprite.SpritesCollection = make([][]*image.NRGBA, nbRow)
				m.sprite.SpritesData = make([][][]byte, nbRow)

				for i := 0; i < nbRow; i++ {
					m.sprite.SpritesCollection[i] = make([]*image.NRGBA, m.sprite.SpriteColumns)
					m.sprite.SpritesData[i] = make([][]byte, m.sprite.SpriteColumns)
				}
				var row, col int
				for i := 0; i < len(data); i += spriteLength {
					m.sprite.SpritesData[row][col] = append(m.sprite.SpritesData[row][col], data[i:(i+spriteLength)]...)
					m.sprite.SpritesCollection[row][col] = spr.RawSpriteToImg(data[i:(i+spriteLength)], footer.Height, footer.Width, uint8(m.sprite.Mode), m.sprite.Palette)
					col++
					if col%m.sprite.SpriteColumns == 0 {
						col = 0
						row++
					}
				}

				icache := custom_widget.NewImageTableCache(m.sprite.SpriteRows, m.sprite.SpriteColumns, fyne.NewSize(50, 50))

				for y := 0; y < m.sprite.SpriteColumns; y++ {
					for x := 0; x < m.sprite.SpriteRows; x++ {
						icache.Set(x, y, canvas.NewImageFromImage(m.sprite.SpritesCollection[x][y]))
					}
				}
				m.sprite.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)

			}

		}, m.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".spr", ".imp"}))
		d.Resize(dialogSize)
		d.Show()

	})
}
