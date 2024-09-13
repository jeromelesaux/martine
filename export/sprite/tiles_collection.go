package sprite

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/martine/constants"
)

const (
	tileHistoricalFilename = "tiles_historical.th"
)

type TilesHistorical struct {
	Width  int              `json:"width"`
	Height int              `json:"height"`
	Tiles  []TileHistorical `json:"sprites"`
}

type TileHistorical struct {
	Label string       `json:"label"`
	Index int          `json:"index"`
	Tile  *image.NRGBA `json:"sprite"`
}

func NewSpritesHistorical(width, height int) *TilesHistorical {
	return &TilesHistorical{
		Width:  width,
		Height: height,
		Tiles:  make([]TileHistorical, 0),
	}
}

func (s *TilesHistorical) Add(spr TileHistorical) {
	s.Tiles = append(s.Tiles, spr)
}

func (s TilesHistorical) LastIndex() int {
	var index int
	for _, spr := range s.Tiles {
		if spr.Index > index {
			index = spr.Index
		}
	}
	return index
}

func (s *TilesHistorical) Append(n *TilesHistorical) {
	s.Tiles = append(s.Tiles, n.Tiles...)
}

func (s *TilesHistorical) Save(folderpath string) error {
	f, err := os.Create(filepath.Join(folderpath, tileHistoricalFilename))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

func (s TilesHistorical) Tile(i *image.NRGBA) (TileHistorical, bool) {
	for _, spr := range s.Tiles {
		if ImageAreEquals(spr.Tile, i) {
			return spr, true
		}
	}
	return TileHistorical{}, false
}

func (s TilesHistorical) IndexOf(i *image.NRGBA) int {
	for _, spr := range s.Tiles {
		if ImageAreEquals(spr.Tile, i) {
			return spr.Index
		}
	}
	return -1
}

func ImageAreEquals(i0, i1 *image.NRGBA) bool {
	if i1 == nil || i0 == nil {
		return false
	}
	if i0.Bounds().Max.X != i1.Bounds().Max.X || i0.Bounds().Max.Y != i1.Bounds().Max.Y {
		return false
	}
	for y := 0; y < i0.Bounds().Max.Y; y++ {
		for x := 0; x < i0.Bounds().Max.X; x++ {
			if !constants.ColorsAreEquals(i0.At(x, y), i1.At(x, y)) {
				return false
			}
		}
	}
	return true
}
