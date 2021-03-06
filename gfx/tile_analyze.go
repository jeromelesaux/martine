package gfx

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/jeromelesaux/martine/constants"
	"github.com/pbnjay/pixfont"
)

var (
	ErrorSizeOverflow = errors.New("Size overflow the image size capacity.")
)

type TilePosition struct {
	PixelX int
	PixelY int
}

func (s *TilePosition) String() string {
	return fmt.Sprintf("[X:%d,Y:%d]", s.PixelX, s.PixelY)
}

type BoardTile struct {
	Occurence     int
	TilePositions []TilePosition
	Tile          *Tile
}

func (b *BoardTile) String() string {
	out := fmt.Sprintf("Occurence:%d \n", b.Occurence)
	for _, v := range b.TilePositions {
		out += v.String() + "\n"
	}
	return out
}

func (a *AnalyzeBoard) Analyse(sprite *Tile, x, y int) int {
	spriteExists := false
	var spriteIndex int
	for i, v := range a.BoardTiles {
		s := v.Tile
		if TilesAreEquals(s, sprite) {
			spriteExists = true
			a.SetAddTile(x, y, i)
			spriteIndex = i
			break
		}
	}
	if !spriteExists {
		a.NewTile(sprite, x, y)
		spriteIndex = len(a.BoardTiles)
	}
	return spriteIndex
}

func (a *AnalyzeBoard) SetAddTile(x, y, index int) {
	//	a.TileMap[len(a.TileMap)] = append(a.TileMap[len(a.TileMap)], index)
	a.BoardTiles[index].Occurence++
	a.BoardTiles[index].TilePositions = append(a.BoardTiles[index].TilePositions, TilePosition{PixelX: x, PixelY: y})
}

func (a *AnalyzeBoard) AddTile(sprite *Tile, x, y int) {
	for i, v := range a.BoardTiles {
		if TilesAreEquals(v.Tile, sprite) {
			a.BoardTiles[i].Occurence++
			a.BoardTiles[i].TilePositions = append(a.BoardTiles[i].TilePositions, TilePosition{PixelX: x, PixelY: y})
			break
		}
	}
}

func (a *AnalyzeBoard) NewTile(sprite *Tile, x, y int) {
	//	a.TileMap[index] = append(a.TileMap[index], len(a.BoardSprites))
	b := BoardTile{TilePositions: make([]TilePosition, 0), Tile: sprite, Occurence: 1}
	b.TilePositions = append(b.TilePositions, TilePosition{PixelX: x, PixelY: y})
	a.BoardTiles = append(a.BoardTiles, b)
}

type AnalyzeBoard struct {
	BoardTiles []BoardTile
	TileSize   constants.Size
	ImageSize  constants.Size
	TileMap    [][]int
}

func (a *AnalyzeBoard) String() string {
	out := a.TileSize.ToString()
	for _, v := range a.BoardTiles {
		out += v.String()
	}
	return out
}

type Tile struct {
	Size   constants.Size
	Colors [][]color.Color
}

func (t *Tile) Image() *image.NRGBA {
	im := image.NewNRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: t.Size.Width, Y: t.Size.Height},
		})
	for y := 0; y < t.Size.Height; y++ {
		for x := 0; x < t.Size.Width; x++ {
			im.Set(x, y, t.Colors[x][y])
		}
	}
	return im
}

func NewTile(size constants.Size) *Tile {
	colors := make([][]color.Color, size.Width)
	for i := 0; i < size.Width; i++ {
		colors[i] = make([]color.Color, size.Height)
	}
	return &Tile{
		Size:   size,
		Colors: colors,
	}
}

func extractTile(im image.Image, size constants.Size, posX, posY int) (*Tile, error) {
	sprite := NewTile(size)
	var xSpr, ySpr int
	onError := false
	for y := posY; y < (posY + size.Height); y++ {
		if y >= im.Bounds().Max.Y {
			onError = true
			break
		}
		xSpr = 0
		for x := posX; x < (posX + size.Width); x++ {
			if x >= im.Bounds().Max.X {
				onError = true
				break
			}
			c := im.At(x, y)
			sprite.Colors[xSpr][ySpr] = c
			xSpr++
		}
		ySpr++
	}
	if onError {
		return sprite, ErrorSizeOverflow
	}
	return sprite, nil
}

func AnalyzeTilesBoard(im image.Image, size constants.Size) *AnalyzeBoard {
	nbTileW := int(im.Bounds().Max.X / size.Width)
	nbTileH := int(im.Bounds().Max.Y/size.Height) - 1
	board := &AnalyzeBoard{
		TileSize:   size,
		ImageSize:  constants.Size{Width: im.Bounds().Max.X, Height: im.Bounds().Max.Y},
		BoardTiles: make([]BoardTile, 0),
		TileMap:    make([][]int, nbTileH),
	}
	for i := 0; i < nbTileH; i++ {
		board.TileMap[i] = make([]int, nbTileW)
	}
	sprt0, err := extractTile(im, size, 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, 0, 0, err)
	}
	board.NewTile(sprt0, 0, 0)
	board.TileMap[0][0] = 0

	indexX := 1
	for x := size.Width; x < im.Bounds().Max.X; x += size.Width {
		indexY := 0
		for y := size.Height; y < im.Bounds().Max.Y; y += size.Height {
			sprt, err := extractTile(im, size, x, y)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while extracting tile size(%d,%d) at position (%d,%d) error :%v\n", size.Width, size.Height, x, y, err)
				break
			}
			index := board.Analyse(sprt, x, y)
			board.TileMap[indexY][indexX] = index
			indexY++
		}
		indexX++
	}
	return board
}

func TilesAreEquals(s1, s2 *Tile) bool {
	if s1.Size.Width != s2.Size.Width || s1.Size.Height != s2.Size.Height {
		return false
	}
	for y := 0; y < s1.Size.Height; y++ {
		for x := 0; x < s1.Size.Width; x++ {
			if !constants.ColorsAreEquals(s1.Colors[x][y], s2.Colors[x][y]) {
				return false
			}
		}
	}
	return true
}

func (a *AnalyzeBoard) SaveSchema(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	totWidth := (16 * len(a.BoardTiles)) + 20 + (20 * len(a.TileMap[0]))
	totHeight := (20 * len(a.BoardTiles)) + 20 + (40 * len(a.TileMap))
	im := image.NewNRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: totWidth, Y: totHeight},
		})
	draw.Draw(im, im.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	x0 := 5
	y0 := 5
	fontColor := color.Black

	title0 := "Tiles found and index name."
	pixfont.DrawString(im, x0, y0, title0, fontColor)
	y0 += 30

	for index, v := range a.BoardTiles {
		sprt := v.Tile
		// draw the sprite
		for y := 0; y < v.Tile.Size.Height; y++ {
			for x := 0; x < v.Tile.Size.Width; x++ {
				im.Set(x+x0, y+y0, sprt.Colors[x][y])
			}
		}

		// draw sprite label
		label := fmt.Sprintf(" Tile %.2d occurence %d", index, v.Occurence)
		pixfont.DrawString(im, x0+sprt.Size.Width+5, y0+3, label, fontColor)
		y0 += sprt.Size.Height + 5

	}

	x0 = 10
	y0 += 30
	title := " Tiles Map by tile index."
	pixfont.DrawString(im, x0, y0, title, fontColor)
	y0 += 30
	for _, v := range a.TileMap {
		for _, val := range v {
			label := fmt.Sprintf("%.2d", val)
			pixfont.DrawString(im, x0, y0, label, fontColor)
			x0 += 30
		}
		x0 = 10
		y0 += 10
	}

	return png.Encode(f, im)
}

func (a *AnalyzeBoard) SaveTilemap(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range a.TileMap {
		for _, val := range v {
			f.WriteString(fmt.Sprintf("%.2d", val) + ",")
		}
		f.WriteString("\n")
	}
	return nil
}
