package ui

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
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
	"github.com/disintegration/imaging"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/log"

	ci "github.com/jeromelesaux/martine/convert/image"
	spr "github.com/jeromelesaux/martine/convert/sprite"
	"github.com/jeromelesaux/martine/export/impdraw/tile"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/export/spritehard"
	"github.com/jeromelesaux/martine/gfx/sprite"
	"github.com/jeromelesaux/martine/ui/martine-ui/directory"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) ApplySprite(s *menu.SpriteMenu) {
	if s.SpriteColumns == 0 || s.SpriteRows == 0 {
		dialog.NewError(errors.New("number of sprites per row or column are not set"), m.window).Show()
		return
	}
	if (s.SpriteWidth == 0 || s.SpriteHeight == 0) && !s.IsHardSprite {
		dialog.ShowError(errors.New("define dimension before"), m.window)
		return
	}
	pi := wgt.NewProgressInfinite("Computing...., Please wait.", m.window)
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
	b := s.OriginalBoard().Image
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
	s.SetPalette(pal)
	size := constants.Size{Width: s.SpriteWidth, Height: s.SpriteHeight}
	raw, sprites, err := sprite.SplitBoardToSprite(s.OriginalBoard().Image, s.Palette(), s.SpriteColumns, s.SpriteRows, uint8(s.Mode), s.IsHardSprite, size)
	if err != nil {
		pi.Hide()
		dialog.NewError(err, m.window).Show()
		return
	}
	s.SpritesCollection = sprites
	s.SpritesData = raw

	icache := wgt.NewImageTableCache(s.SpriteColumns, s.SpriteRows, fyne.NewSize(50, 50))

	for x := 0; x < s.SpriteColumns; x++ {
		for y := 0; y < s.SpriteRows; y++ {
			icache.Set(x, y, canvas.NewImageFromImage(s.SpritesCollection[x][y]))
		}
	}
	s.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)
	s.SetPaletteImage(png.PalToImage(s.Palette()))
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
			directory.SetImportDirectoryURI(reader.URI())
			s.FilePath = reader.URI().Path()
			img, err := openImage(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			s.SetOriginalBoard(img)
		}, m.window)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(imagesFilesFilter)
		d.Resize(dialogSize)
		d.Show()
	})

	exportButton := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		m.exportSpriteBoard(m.sprite, m.window)
	})

	applyButton := widget.NewButtonWithIcon("Apply", theme.VisibilityIcon(), func() {
		log.GetLogger().Infoln("apply.")
		m.ApplySprite(s)
	})

	openFileWidget.Icon = theme.FileImageIcon()

	modes := widget.NewSelect([]string{"0", "1", "2"}, func(v string) {
		mode, err := strconv.Atoi(v)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", v)
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
			log.GetLogger().Error("Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteRows = r
	}
	spriteNumberPerRowEntry := widget.NewEntry()
	spriteNumberPerRowEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", v)
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
			log.GetLogger().Error("Error %s cannot be cast in int\n", v)
			return
		}
		s.SpriteWidth = r
	}

	spriteHeightSizeEntry := widget.NewEntry()
	spriteHeightSizeEntry.OnChanged = func(v string) {
		r, err := strconv.Atoi(v)
		if err != nil {
			log.GetLogger().Error("Error %s cannot be cast in int\n", v)
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
	gifOpen := applySpriteBoardFromGif(s, m)

	return container.New(
		layout.NewGridLayoutWithColumns(2),
		container.New(
			layout.NewGridLayoutWithRows(2),
			s.OriginalBoard(),
			s.OriginalImages,
		),
		container.New(
			layout.NewGridLayoutWithRows(3),
			container.New(
				layout.NewVBoxLayout(),
				container.New(
					layout.NewHBoxLayout(),
					openFileWidget,
					applyButton,
					exportButton,
					paletteOpen,
					importOpen,
					gifOpen,
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
					s.PaletteImage(),
				),
			),
			container.New(
				layout.NewVBoxLayout(),
				widget.NewButton("show cmd", func() {
					e := widget.NewMultiLineEntry()
					e.SetText(s.CmdLine())

					d := dialog.NewCustom("Command line generated",
						"Ok",
						e,
						m.window)
					log.GetLogger().Info("%s\n", s.CmdLine())
					size := m.window.Content().Size()
					size = fyne.Size{Width: size.Width / 2, Height: size.Height / 2}
					d.Resize(size)
					d.Show()
				},
				),
			),
		),
	)
}

func applySpriteBoardFromGif(s *menu.SpriteMenu, m *MartineUI) fyne.Widget {
	return widget.NewButtonWithIcon("From Gif", theme.FileImageIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			if reader == nil {
				return
			}
			if (s.SpriteWidth == 0 || s.SpriteHeight == 0) && !s.IsHardSprite {
				dialog.ShowError(errors.New("define dimension before"), m.window)
				return
			}
			directory.SetImportDirectoryURI(reader.URI())
			filePath := reader.URI()
			fr, err := os.Open(filePath.Path())
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			gifImage, err := gif.DecodeAll(fr)
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			gifImages := ci.GifToImages(*gifImage)
			resized := make([]*image.NRGBA, 0)
			size := constants.Size{Width: s.SpriteWidth, Height: s.SpriteHeight}
			for _, v := range gifImages {
				r := ci.Resize(v, size, imaging.NearestNeighbor)
				resized = append(resized, r)
			}
			cfg := config.NewMartineConfig("", "")
			cfg.CustomDimension = true
			cfg.Size = size
			cfg.SpriteHard = s.IsHardSprite
			var colorsAvailable int
			switch s.Mode {
			case 0:
				colorsAvailable = constants.Mode0.ColorsAvailable
			case 1:
				colorsAvailable = constants.Mode1.ColorsAvailable
			case 2:
				colorsAvailable = constants.Mode2.ColorsAvailable
			}
			img := resized[0]
			pal, _, err := ci.DowngradingPalette(img, constants.Size{ColorsAvailable: colorsAvailable, Width: img.Bounds().Max.X, Height: img.Bounds().Max.Y}, s.IsCpcPlus)
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}
			s.SetPalette(pal)
			raw, sprites, _, _ := gfx.ApplyImages(resized, cfg, s.Mode, pal, uint8(s.Mode))
			s.SpritesCollection = make([][]*image.NRGBA, 1)
			s.SpritesCollection[0] = sprites
			s.SpritesData = make([][][]byte, 1)
			s.SpritesData[0] = raw
			s.SpriteColumns = 1
			s.SpriteRows = len(resized)
			icache := wgt.NewImageTableCache(s.SpriteColumns, s.SpriteRows, fyne.NewSize(50, 50))

			for x := 0; x < s.SpriteColumns; x++ {
				for y := 0; y < s.SpriteRows; y++ {
					icache.Set(x, y, canvas.NewImageFromImage(s.SpritesCollection[x][y]))
				}
			}
			s.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)
			s.SetPaletteImage(png.PalToImage(s.Palette()))
		}, m.window)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".gif"}))
		d.Resize(dialogSize)
		d.Show()
	})
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
			directory.SetImportDirectoryURI(reader.URI())
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
					m.sprite.SpritesCollection[row][col] = spritesHard.Data[i].Image(m.sprite.Palette())
					col++
					if col%m.sprite.SpriteColumns == 0 {
						col = 0
						row++
					}
				}

				icache := wgt.NewImageTableCache(m.sprite.SpriteRows, m.sprite.SpriteColumns, fyne.NewSize(50, 50))

				for y := 0; y < m.sprite.SpriteColumns; y++ {
					for x := 0; x < m.sprite.SpriteRows; x++ {
						if m.sprite.SpritesCollection[x][y] != nil {
							icache.Set(x, y, canvas.NewImageFromImage(m.sprite.SpritesCollection[x][y]))
						}
					}
				}
				m.sprite.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)
				// len(spritesHard.Data)/m.sprite.SpriteNumberPerColumn
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
					m.sprite.SpritesCollection[row][col] = spr.RawSpriteToImg(data[i:(i+spriteLength)], footer.Height, footer.Width, uint8(m.sprite.Mode), m.sprite.Palette())
					col++
					if col%m.sprite.SpriteColumns == 0 {
						col = 0
						row++
					}
				}

				icache := wgt.NewImageTableCache(m.sprite.SpriteRows, m.sprite.SpriteColumns, fyne.NewSize(50, 50))

				for y := 0; y < m.sprite.SpriteColumns; y++ {
					for x := 0; x < m.sprite.SpriteRows; x++ {
						icache.Set(x, y, canvas.NewImageFromImage(m.sprite.SpritesCollection[x][y]))
					}
				}
				m.sprite.OriginalImages.Update(icache, icache.ImagesPerRow, icache.ImagesPerColumn)

			}
		}, m.window)
		path, err := directory.ImportDirectoryURI()
		if err == nil {
			d.SetLocation(path)
		}
		d.SetFilter(storage.NewExtensionFileFilter([]string{".spr", ".imp"}))
		d.Resize(dialogSize)
		d.Show()
	})
}
