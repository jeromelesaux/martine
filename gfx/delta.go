package gfx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jeromelesaux/martine/common"
	x "github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

var ErrorCanNotProceed = errors.New("Can not proceed treatment")

type DeltaItem struct {
	Byte    byte
	Offsets []uint16
}

type DeltaCollection struct {
	Items []DeltaItem
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
	out := fmt.Sprintf("byte value #%.2x :", di.Byte)
	for _, addr := range di.Offsets {
		out += fmt.Sprintf("\n#%.4x", addr)
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

func (dc *DeltaCollection) Add(b byte, address uint16) {
	for i := 0; i < len(dc.Items); i++ {
		if dc.Items[i].Byte == b {
			dc.Items[i].Offsets = append(dc.Items[i].Offsets, address)
			return
		}
	}
	item := NewDeltaItem()
	item.Offsets = append(item.Offsets, address)
	item.Byte = b
	dc.Items = append(dc.Items, item)
}

func DeltaMode0(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, exportType *x.ExportType) (*DeltaCollection, error) {
	data := NewDeltaCollection()
	if current.Bounds().Max.X != next.Bounds().Max.X {
		return data, ErrorSizeMismatch
	}
	if current.Bounds().Max.Y != next.Bounds().Max.Y {
		return data, ErrorSizeMismatch
	}
	for i := 0; i < current.Bounds().Max.X; i += 2 {
		for j := 0; j < current.Bounds().Max.Y; j++ {
			c1 := current.At(i, j)
			c2 := next.At(i, j)
			i++
			c3 := current.At(i, j)
			c4 := next.At(i, j)
			p1, err := PalettePosition(c1, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, i, j)
				p1 = 0
			}
			p3, err := PalettePosition(c3, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, i, j)
				p3 = 0
			}
			pixel1 := pixelMode0(p1, p3)

			p2, err := PalettePosition(c2, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, i, j)
				p2 = 0
			}
			p4, err := PalettePosition(c4, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, i, j)
				p4 = 0
			}
			pixel2 := pixelMode0(p2, p4)
			if pixel1 != pixel2 {
				addr := CpcScreenAddress(0xc000, i, j, 0, exportType.Overscan)
				data.Add(pixel2, uint16(addr))
			}
		}
	}
	return data, nil
}

func DeltaMode1(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, exportType *x.ExportType) (*DeltaCollection, error) {
	data := NewDeltaCollection()
	if current.Bounds().Max.X != next.Bounds().Max.X {
		return data, ErrorSizeMismatch
	}
	if current.Bounds().Max.Y != next.Bounds().Max.Y {
		return data, ErrorSizeMismatch
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
			p1, err := PalettePosition(c1, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, i, j)
				p1 = 0
			}
			p3, err := PalettePosition(c3, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, i, j)
				p3 = 0
			}
			p5, err := PalettePosition(c5, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, i, j)
				p5 = 0
			}
			p7, err := PalettePosition(c7, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, i, j)
				p7 = 0
			}
			pixel1 := pixelMode1(p1, p3, p5, p7)

			p2, err := PalettePosition(c2, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, i, j)
				p2 = 0
			}
			p4, err := PalettePosition(c4, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, i, j)
				p4 = 0
			}
			p6, err := PalettePosition(c6, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, i, j)
				p6 = 0
			}
			p8, err := PalettePosition(c8, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, i, j)
				p8 = 0
			}
			pixel2 := pixelMode1(p2, p4, p6, p8)
			if pixel1 != pixel2 {
				addr := CpcScreenAddress(0xc000, i, j, 1, exportType.Overscan)
				data.Add(pixel2, uint16(addr))
			}
		}
	}
	return data, nil
}

func DeltaMode2(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, exportType *x.ExportType) (*DeltaCollection, error) {
	data := NewDeltaCollection()
	if current.Bounds().Max.X != next.Bounds().Max.X {
		return data, ErrorSizeMismatch
	}
	if current.Bounds().Max.Y != next.Bounds().Max.Y {
		return data, ErrorSizeMismatch
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
			p1, err := PalettePosition(c1, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c1, i, j)
				p1 = 0
			}
			p3, err := PalettePosition(c3, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c3, i, j)
				p3 = 0
			}
			p5, err := PalettePosition(c5, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c5, i, j)
				p5 = 0
			}
			p7, err := PalettePosition(c7, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c7, i, j)
				p7 = 0
			}
			p9, err := PalettePosition(c9, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c9, i, j)
				p9 = 0
			}
			p11, err := PalettePosition(c11, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c11, i, j)
				p11 = 0
			}
			p13, err := PalettePosition(c13, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c13, i, j)
				p13 = 0
			}
			p15, err := PalettePosition(c15, currentPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c15, i, j)
				p15 = 0
			}
			pixel1 := pixelMode2(p1, p3, p5, p7, p9, p11, p13, p15)

			p2, err := PalettePosition(c2, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, i, j)
				p2 = 0
			}
			p4, err := PalettePosition(c4, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c2, i, j)
				p4 = 0
			}
			p6, err := PalettePosition(c6, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c6, i, j)
				p6 = 0
			}
			p8, err := PalettePosition(c8, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c8, i, j)
				p8 = 0
			}
			p10, err := PalettePosition(c10, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c10, i, j)
				p10 = 0
			}
			p12, err := PalettePosition(c12, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c12, i, j)
				p12 = 0
			}
			p14, err := PalettePosition(c14, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c14, i, j)
				p14 = 0
			}
			p16, err := PalettePosition(c16, nextPalette)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v pixel position(%d,%d) not found in palette\n", c16, i, j)
				p16 = 0
			}
			pixel2 := pixelMode2(p2, p4, p6, p8, p10, p12, p14, p16)
			if pixel1 != pixel2 {
				addr := CpcScreenAddress(0xc000, i, j, 2, exportType.Overscan)
				data.Add(pixel2, uint16(addr))
			}
		}
	}
	return data, nil
}

func (dc *DeltaCollection) Marshall() ([]byte, error) {
	var b bytes.Buffer
	for _, item := range dc.Items {
		occ := len(item.Offsets)
		for i := 0; i < occ; i += 255 {
			if err := binary.Write(&b, binary.LittleEndian, item.Byte); err != nil {
				return b.Bytes(), err
			}
			var nbocc uint8 = 255
			if occ-i < 255 {
				nbocc = uint8(occ - i)
			}
			if err := binary.Write(&b, binary.LittleEndian, nbocc); err != nil {
				return b.Bytes(), err
			}
			for j := i; j < 255 && j < occ; j++ {
				if err := binary.Write(&b, binary.LittleEndian, item.Offsets[j]); err != nil {
					return b.Bytes(), err
				}
			}
		}
	}
	return b.Bytes(), nil
}

func (dc *DeltaCollection) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, item := range dc.Items {
		occ := len(item.Offsets)
		for i := 0; i < occ; i += 255 {
			if err := binary.Write(f, binary.LittleEndian, item.Byte); err != nil {
				return err
			}
			var nbocc uint8 = 255
			if occ-i < 255 {
				nbocc = uint8(occ - i)
			}
			if err := binary.Write(f, binary.LittleEndian, nbocc); err != nil {
				return err
			}
			for j := i; j < 255 && j < occ; j++ {
				if err := binary.Write(f, binary.LittleEndian, item.Offsets[j]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Delta(scr1, scr2 []byte, isSprite bool) *DeltaCollection {
	data := NewDeltaCollection()
	//var line int
	for offset := 0; offset < len(scr1); offset++ {
		if scr1[offset] != scr2[offset] {
			data.Add(scr2[offset], uint16(offset))
		}
	}
	return data
}

func ExportDelta(filename string, dc *DeltaCollection, exportType *x.ExportType) error {
	if err := dc.Save(filename + ".bin"); err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving file (%s) error %v \n", filename+".bin", err)
		return err
	}
	data, err := dc.Marshall()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while marshalling delta structure error :%v\n", err)
		return err
	}

	var emptyPalette []color.Color
	outFilepath := exportType.OutputPath + string(filepath.Separator) + filename + ".txt"
	if err = file.Ascii(outFilepath, data, emptyPalette, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while exporting data as ascii mode file (%s) error :%v\n", outFilepath, err)
		return err
	}
	outFilepath = exportType.OutputPath + string(filepath.Separator) + filename + "c.txt"
	if err = file.AsciiByColumn(outFilepath, data, emptyPalette, exportType); err != nil {
		fmt.Fprintf(os.Stderr, "Error while exporting data as ascii by column mode file (%s) error :%v\n", outFilepath, err)
		return err
	}
	return nil
}

func ProceedDelta(filespath []string, exportType *x.ExportType) error {

	if len(filespath) == 1 {
		var err error
		filespath, err = common.WilcardedFiles(filespath)
		if err != nil {
			return err
		}
		if len(filespath) == 1 || len(filespath) == 0 {
			return ErrorCanNotProceed
		}
	}
	var d1, d2 []byte
	var err error
	var isSprite = false
	fmt.Fprintf(os.Stdout, "%v\n", filespath)
	for i := 0; i < len(filespath)-1; i++ {
		switch strings.ToUpper(filepath.Ext(filespath[i])) {
		case ".WIN":
			d1, err = file.RawWin(filespath[i])
			if err != nil {
				return err
			}
		case ".SCR":
			d1, err = file.RawScr(filespath[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "File (%s) is not a simple screen.\n", filespath[i])
			}
			d1, err = file.RawOverscan(filespath[i])
			if err != nil {
				return err
			}
		default:
			return ErrorCanNotProceed
		}

		switch strings.ToUpper(filepath.Ext(filespath[i+1])) {
		case ".WIN":
			d2, err = file.RawWin(filespath[i+1])
			if err != nil {
				return err
			}
		case ".SCR":
			d2, err = file.RawScr(filespath[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "File (%s) is not a simple screen.\n", filespath[i+1])
			}
			d2, err = file.RawOverscan(filespath[i+1])
			if err != nil {
				return err
			}
		default:
			return ErrorCanNotProceed
		}

		dc := Delta(d1, d2, isSprite)
		fmt.Fprintf(os.Stdout, "files (%s) (%s)", filespath[i], filespath[i+1])
		fmt.Fprintf(os.Stdout, "%d bytes differ from the both images\n", len(dc.Items))
		fmt.Fprintf(os.Stdout, "%d screen addresses are involved\n", dc.NbAdresses())
		fmt.Fprintf(os.Stdout, "Report:\n%s\n", dc.ToString())
		out := exportType.OutputPath + string(filepath.Separator) + fmt.Sprintf("%dto%d", i, (i+1))
		if err := ExportDelta(out, dc, exportType); err != nil {
			return err
		}
	}

	switch strings.ToUpper(filepath.Ext(filespath[len(filespath)-1])) {
	case ".WIN":
		d1, err = file.RawWin(filespath[len(filespath)-1])
		if err != nil {
			return err
		}
	case ".SCR":
		d1, err = file.RawScr(filespath[len(filespath)-1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "File (%s) is not a simple screen.\n", filespath[len(filespath)-1])
		}
		d1, err = file.RawOverscan(filespath[len(filespath)-1])
		if err != nil {
			return err
		}
	default:
		return ErrorCanNotProceed
	}

	switch strings.ToUpper(filepath.Ext(filespath[0])) {
	case ".WIN":
		d2, err = file.RawWin(filespath[0])
		if err != nil {
			return err
		}
	case ".SCR":
		d2, err = file.RawScr(filespath[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "File (%s) is not a simple screen.\n", filespath[0])
		}
		d2, err = file.RawOverscan(filespath[0])
		if err != nil {
			return err
		}
	default:
		return ErrorCanNotProceed
	}

	f1, err := os.Open(filespath[len(filespath)-1])
	if err != nil {
		return err
	}
	defer f1.Close()
	f2, err := os.Open(filespath[0])
	if err != nil {
		return err
	}
	defer f2.Close()
	dc := Delta(d1, d2, isSprite)
	fmt.Fprintf(os.Stdout, "files (%s) (%s)", filespath[len(filespath)-1], filespath[0])
	fmt.Fprintf(os.Stdout, "%d bytes differ from the both images\n", len(dc.Items))
	fmt.Fprintf(os.Stdout, "%d screen addresses are involved\n", dc.NbAdresses())
	fmt.Fprintf(os.Stdout, "Report:\n%s\n", dc.ToString())
	out := exportType.OutputPath + string(filepath.Separator) + fmt.Sprintf("%dto0", len(filespath))
	if err := ExportDelta(out, dc, exportType); err != nil {
		return err
	}

	return nil
}
