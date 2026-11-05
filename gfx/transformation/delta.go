// Package transformation provides image transformation utilities.
package transformation

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/convert/address"
	"github.com/jeromelesaux/martine/convert/palette"
	"github.com/jeromelesaux/martine/convert/pixel"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/log"

	cm "github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/impdraw/overscan"
	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/ocpartstudio/window"
)

type DeltaItem struct {
	Byte    byte
	Offsets []uint16
}

type DeltaCollection struct {
	OccurencePerFrame uint8
	Items             []DeltaItem
}

func (dc *DeltaCollection) ItemsSortByByte() []DeltaItem {
	var items []DeltaItem
	items = append(items, dc.Items...)

	for i := 0; i < len(dc.Items)-1; i++ {
		for j := 0; j < len(dc.Items)-1-i; j++ {
			if items[j].Byte > items[j+1].Byte {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}

	return items
}

/*
	details de la structure
	nombre d'occurence par frame occ
	0	byte à poker (uint8) |  nombre d'offsets  uint16 |  offset 0 uint16 ....... offset (nombre d'offset)
	1
	2
	.
	.
	.
	occ
*/

func (dc *DeltaCollection) Occurrences() int {
	occurrence := len(dc.Items)
	return occurrence
}

func NewDeltaCollection() *DeltaCollection {
	return &DeltaCollection{
		Items: make([]DeltaItem, 0),
	}
}

func NewDeltaItem() DeltaItem {
	return DeltaItem{Offsets: make([]uint16, 0)}
}

func (di *DeltaItem) NbAddresses() int {
	return len(di.Offsets)
}
func (dc *DeltaCollection) NbAdresses() int {
	nb := 0
	for _, item := range dc.Items {
		nb += item.NbAddresses()
	}
	return nb
}

func (di *DeltaItem) ToString() string {
	out := fmt.Sprintf("byte value #%.2x : offsets (%d):", di.Byte, len(di.Offsets))

	for i, addr := range di.Offsets {
		if i%8 == 0 {
			out += "\n"
		}
		out += fmt.Sprintf("#%.4x ", addr)
	}
	return out
}

func (dc *DeltaCollection) ToString() string {
	var out string
	for _, item := range dc.Items {
		out += item.ToString() + "\n"
	}
	return out
}

func (dc *DeltaCollection) Add(b byte, addr uint16) {
	for i := 0; i < len(dc.Items); i++ {
		if dc.Items[i].Byte == b {
			dc.Items[i].Offsets = append(dc.Items[i].Offsets, addr)
			return
		}
	}
	item := NewDeltaItem()
	item.Offsets = append(item.Offsets, addr)
	item.Byte = b
	dc.Items = append(dc.Items, item)
}

func DeltaMode0(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, cfg *config.MartineConfig) (*DeltaCollection, error) {
	data := NewDeltaCollection()
	if current.Bounds().Max.X != next.Bounds().Max.X {
		return data, errors.ErrorSizeMismatch
	}
	if current.Bounds().Max.Y != next.Bounds().Max.Y {
		return data, errors.ErrorSizeMismatch
	}
	for i := 0; i < current.Bounds().Max.X; i += 2 {
		for j := 0; j < current.Bounds().Max.Y; j++ {
			c1 := current.At(i, j)
			c2 := next.At(i, j)
			i++
			c3 := current.At(i, j)
			c4 := next.At(i, j)
			p1 := palette.PalettePositionOrDefault(c1, currentPalette)
			p3 := palette.PalettePositionOrDefault(c3, currentPalette)
			pixel1 := pixel.PixelMode0(p1, p3)

			p2 := palette.PalettePositionOrDefault(c2, nextPalette)
			p4 := palette.PalettePositionOrDefault(c4, nextPalette)
			pixel2 := pixel.PixelMode0(p2, p4)
			if pixel1 != pixel2 {
				addr := address.CpcScreenAddress(0xc000, i, j, 0, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
				data.Add(pixel2, uint16(addr))
			}
		}
	}
	return data, nil
}

// DeltaMode1 computes the delta collection between two mode 1 images.
// nolint: funlen, gocognit
func DeltaMode1(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, cfg *config.MartineConfig) (*DeltaCollection, error) {
	data := NewDeltaCollection()
	if current.Bounds().Max.X != next.Bounds().Max.X {
		return data, errors.ErrorSizeMismatch
	}
	if current.Bounds().Max.Y != next.Bounds().Max.Y {
		return data, errors.ErrorSizeMismatch
	}
	for i := 0; i < current.Bounds().Max.X; i += 4 {
		for j := 0; j < current.Bounds().Max.Y; j++ {
			c1 := current.At(i, j)
			c2 := next.At(i, j)
			i++
			c3 := current.At(i, j)
			c4 := next.At(i, j)
			i++
			c5 := current.At(i, j)
			c6 := next.At(i, j)
			i++
			c7 := current.At(i, j)
			c8 := next.At(i, j)
			p1 := palette.PalettePositionOrDefault(c1, currentPalette)
			p3 := palette.PalettePositionOrDefault(c3, currentPalette)
			p5 := palette.PalettePositionOrDefault(c5, currentPalette)
			p7 := palette.PalettePositionOrDefault(c7, currentPalette)
			pixel1 := pixel.PixelMode1(p1, p3, p5, p7)

			p2 := palette.PalettePositionOrDefault(c2, nextPalette)
			p4 := palette.PalettePositionOrDefault(c4, nextPalette)
			p6 := palette.PalettePositionOrDefault(c6, nextPalette)
			p8 := palette.PalettePositionOrDefault(c8, nextPalette)
			pixel2 := pixel.PixelMode1(p2, p4, p6, p8)
			if pixel1 != pixel2 {
				addr := address.CpcScreenAddress(0xc000, i, j, 1, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
				data.Add(pixel2, uint16(addr))
			}
		}
	}
	return data, nil
}

// DeltaMode2 computes the delta collection between two mode 2 images.
// nolint:funlen, gocognit
func DeltaMode2(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, cfg *config.MartineConfig) (*DeltaCollection, error) {
	data := NewDeltaCollection()
	if current.Bounds().Max.X != next.Bounds().Max.X {
		return data, errors.ErrorSizeMismatch
	}
	if current.Bounds().Max.Y != next.Bounds().Max.Y {
		return data, errors.ErrorSizeMismatch
	}
	for i := 0; i < current.Bounds().Max.X; i += 8 {
		for j := 0; j < current.Bounds().Max.Y; j++ {
			c1 := current.At(i, j)
			c2 := next.At(i, j)
			i++
			c3 := current.At(i, j)
			c4 := next.At(i, j)
			i++
			c5 := current.At(i, j)
			c6 := next.At(i, j)
			i++
			c7 := current.At(i, j)
			c8 := next.At(i, j)
			i++
			c9 := current.At(i, j)
			c10 := next.At(i, j)
			i++
			c11 := current.At(i, j)
			c12 := next.At(i, j)
			i++
			c13 := current.At(i, j)
			c14 := next.At(i, j)
			i++
			c15 := current.At(i, j)
			c16 := next.At(i, j)
			p1 := palette.PalettePositionOrDefault(c1, currentPalette)
			p3 := palette.PalettePositionOrDefault(c3, currentPalette)
			p5 := palette.PalettePositionOrDefault(c5, currentPalette)
			p7 := palette.PalettePositionOrDefault(c7, currentPalette)
			p9 := palette.PalettePositionOrDefault(c9, currentPalette)
			p11 := palette.PalettePositionOrDefault(c11, currentPalette)
			p13 := palette.PalettePositionOrDefault(c13, currentPalette)
			p15 := palette.PalettePositionOrDefault(c15, currentPalette)
			pixel1 := pixel.PixelMode2(p1, p3, p5, p7, p9, p11, p13, p15)

			p2 := palette.PalettePositionOrDefault(c2, nextPalette)
			p4 := palette.PalettePositionOrDefault(c4, nextPalette)
			p6 := palette.PalettePositionOrDefault(c6, nextPalette)
			p8 := palette.PalettePositionOrDefault(c8, nextPalette)
			p10 := palette.PalettePositionOrDefault(c10, nextPalette)
			p12 := palette.PalettePositionOrDefault(c12, nextPalette)
			p14 := palette.PalettePositionOrDefault(c14, nextPalette)
			p16 := palette.PalettePositionOrDefault(c16, nextPalette)
			pixel2 := pixel.PixelMode2(p2, p4, p6, p8, p10, p12, p14, p16)
			if pixel1 != pixel2 {
				addr := address.CpcScreenAddress(0xc000, i, j, 2, cfg.ScrCfg.Type.IsFullScreen(), cfg.DoubleScreenAddress)
				data.Add(pixel2, uint16(addr))
			}
		}
	}
	return data, nil
}

// Marshall serializes the delta collection as bytes (byte value, number of occurrences, offsets).
func (dc *DeltaCollection) Marshall() ([]byte, error) {
	var b bytes.Buffer

	if err := binary.Write(&b, binary.LittleEndian, dc.OccurencePerFrame); err != nil {
		return b.Bytes(), err
	}
	if dc.OccurencePerFrame == 0 { // no difference between transitions
		return b.Bytes(), nil
	}
	// occurencesPerframe doit correspondre au nombre offsets modulo 255 et non au nombre d'items
	for _, item := range dc.Items {
		occ := len(item.Offsets)
		if err := binary.Write(&b, binary.LittleEndian, item.Byte); err != nil {
			return b.Bytes(), err
		}
		if err := binary.Write(&b, binary.LittleEndian, uint16(occ)); err != nil {
			return b.Bytes(), err
		}
		for i := 0; i < occ; i++ {
			value := item.Offsets[i]
			if err := binary.Write(&b, binary.LittleEndian, value); err != nil {
				return b.Bytes(), err
			}
		}
	}
	return b.Bytes(), nil
}

func (dc *DeltaCollection) Save(filename string) error {
	b, err := dc.Marshall()
	if err != nil {
		return err
	}
	return amsdos.SaveOSFile(filename, b)
}

func DeltaAddress(x, y, lineOctetWidth int) int {
	// lineOctetWidth == 0x50
	if y > (8 * 0x800 / lineOctetWidth) {
		log.GetLogger().Error("WARNING: y (%d) is superior to  (8 * 0x800 /#%x)\n", y, lineOctetWidth)
	}
	return (0x800 * (y % 8)) + (lineOctetWidth * (y / 8)) + (x)
}

func X(offset uint16, lineOctetWidth int) uint16 {
	line := Y(offset, lineOctetWidth)
	return uint16(int(offset) - DeltaAddress(0, int(line), lineOctetWidth))
}

func Y(offset uint16, lineOctetWidth int) uint16 {
	line := 0
	for i := 0; i < constants.Mode0.Height; i++ {
		lineAddress := DeltaAddress(0, i, lineOctetWidth)
		if lineAddress > int(offset) {
			line = i - 1
			break
		}
	}
	return uint16(line)
}

func CpcCoordinates(addr, startingAddress uint16, lineOctetWidth int) (xVal, yVal int, err error) {
	for y := 0; y < constants.Mode0.Height; y++ {
		for x := 0; x < constants.Mode0.Width; x++ {
			v := uint16(DeltaAddress(x, y, lineOctetWidth))
			v += startingAddress
			if v == addr {
				xVal = x
				yVal = y
				return
			}
		}
	}
	err = errors.ErrorCoordinatesNotFound
	return
}

func Delta(scr1, scr2 []byte, isSprite bool, size constants.Size, mode uint8, x0, y0 uint16, lineOctetWidth int) *DeltaCollection {
	data := NewDeltaCollection()
	for offset := range scr1 { // a revoir car pour un sprite ce n'est le même mode d'adressage
		if scr1[offset] != scr2[offset] {
			if isSprite {

				y := (offset / (size.Width)) + int(y0)
				x := ((offset + int(x0)) - ((y - int(y0)) * (size.Width)))
				newOffset := DeltaAddress(x, y, lineOctetWidth) + 0xC000
				data.Add(scr2[offset], uint16(newOffset))
			} else {
				data.Add(scr2[offset], uint16(offset))
			}
		}
	}
	data.OccurencePerFrame = uint8(data.Occurrences())
	return data
}

func ExportDelta(filename string, dc *DeltaCollection, mode uint8, cfg *config.MartineConfig) error {
	if err := dc.Save(filename + ".bin"); err != nil {
		log.GetLogger().Error("Error while saving file (%s) error %v \n", filename+".bin", err)
		return err
	}
	data, err := dc.Marshall()
	if err != nil {
		log.GetLogger().Error("Error while marshaling delta structure error :%v\n", err)
		return err
	}

	var emptyPalette []color.Color
	outFilepath := filepath.Join(cfg.ScrCfg.OutputPath, filename+".txt")
	if err = ascii.ASCII(outFilepath, data, emptyPalette, false, cfg); err != nil {
		log.GetLogger().Error("Error while exporting data as ascii mode file (%s) error :%v\n", outFilepath, err)
		return err
	}
	outFilepath = filepath.Join(cfg.ScrCfg.OutputPath, filename+"c.txt")
	if err = ascii.ASCIIByColumn(outFilepath, data, emptyPalette, false, mode, cfg); err != nil {
		log.GetLogger().Error("Error while exporting data as ascii by column mode file (%s) error :%v\n", outFilepath, err)
		return err
	}
	return nil
}

func loadRawFile(path string, size *constants.Size, isSprite *bool) ([]byte, error) {
	switch strings.ToUpper(filepath.Ext(path)) {
	case constants.WindowExtension:
		d, err := window.RawWin(path)
		if err != nil {
			return nil, err
		}
		footer, err := window.OpenWin(path)
		if err != nil {
			return nil, err
		}
		size.Width = int(footer.Width)
		size.Height = int(footer.Height)
		*isSprite = true
		return d, nil
	case constants.ScrExtension:
		_, err := ocpartstudio.RawScr(path)
		if err != nil {
			log.GetLogger().Error("File (%s) is not a simple screen.\n", path)
		}
		return overscan.RawOverscan(path)
	default:
		return nil, errors.ErrorCanNotProceed
	}
}

// ProceedDelta computes and exports delta packing for a list of screen files.
// nolint:funlen, gocognit, gocyclo
func ProceedDelta(filespath []string, initialAddress uint16, cfg *config.MartineConfig, mode uint8) error {
	if len(filespath) == 1 {
		var err error
		filespath, err = cm.WilcardedFiles(filespath)
		if err != nil {
			return err
		}
		if len(filespath) == 1 || len(filespath) == 0 {
			return errors.ErrorCanNotProceed
		}
	}
	var err error
	var isSprite = false
	var size constants.Size
	lineOctetsWidth := cfg.LineWidth
	x0, y0, err := CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		log.GetLogger().Error("error while computing cpc coordinates :%v\n", err)
	}
	log.GetLogger().Info("%v\n", filespath)
	log.GetLogger().Info("Cpc coordinates X:%d,Y:%d [#%.4x]\n", x0, y0, initialAddress)
	for i := 0; i < len(filespath)-1; i++ {
		d1, err := loadRawFile(filespath[i], &size, &isSprite)
		if err != nil {
			return err
		}
		d2, err := loadRawFile(filespath[i+1], &size, &isSprite)
		if err != nil {
			return err
		}
		if len(d1) != len(d2) {
			return errors.ErrorSizeDiffers
		}
		dc := Delta(d1, d2, isSprite, size, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		log.GetLogger().Info("files (%s) (%s)", filespath[i], filespath[i+1])
		log.GetLogger().Info("%d bytes differ from the both images\n", len(dc.Items))
		log.GetLogger().Info("%d screen addresses are involved\n", dc.NbAdresses())
		log.GetLogger().Info("Report:\n%s\n", dc.ToString())
		if dc.OccurencePerFrame != 0 {
			out := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("%.2dto%.2d", i, (i+1)))
			if err := ExportDelta(out, dc, mode, cfg); err != nil {
				return err
			}
		}
	}

	d1, err := loadRawFile(filespath[len(filespath)-1], &size, &isSprite)
	if err != nil {
		return err
	}
	d2, err := loadRawFile(filespath[0], &size, &isSprite)
	if err != nil {
		return err
	}
	dc := Delta(d1, d2, isSprite, size, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	log.GetLogger().Info("files (%s) (%s)", filespath[len(filespath)-1], filespath[0])
	log.GetLogger().Info("%d bytes differ from the both images\n", len(dc.Items))
	log.GetLogger().Info("%d screen addresses are involved\n", dc.NbAdresses())
	log.GetLogger().Info("Report:\n%s\n", dc.ToString())
	if dc.OccurencePerFrame != 0 {
		out := filepath.Join(cfg.ScrCfg.OutputPath, fmt.Sprintf("%.2dto00", len(filespath)-1))
		if err := ExportDelta(out, dc, mode, cfg); err != nil {
			return err
		}
	}
	log.GetLogger().Info("files order : %v\n", filespath)
	log.GetLogger().Info("Starting address to display delta : #%.4X\n", initialAddress)
	return nil
}
