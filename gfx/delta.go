package gfx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"
)

type DeltaItem struct {
	Byte      byte
	Addresses []uint16
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
	return DeltaItem{Addresses: make([]uint16, 0)}
}

func (di *DeltaItem) NbAddresses() int {
	return len(di.Addresses)
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
	for _, addr := range di.Addresses {
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
			dc.Items[i].Addresses = append(dc.Items[i].Addresses, address)
			return
		}
	}
	item := NewDeltaItem()
	item.Addresses = append(item.Addresses, address)
	item.Byte = b
	dc.Items = append(dc.Items, item)
}

func DeltaMode0(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, exportType *ExportType) (*DeltaCollection, error) {
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

func DeltaMode1(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, exportType *ExportType) (*DeltaCollection, error) {
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

func DeltaMode2(current *image.NRGBA, currentPalette color.Palette, next *image.NRGBA, nextPalette color.Palette, exportType *ExportType) (*DeltaCollection, error) {
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
		occ := len(item.Addresses)
		iter := occ / 255
		if occ%255 != 0 {
			iter++
		}
		for i := 0; i < iter; i++ {
			if err := binary.Write(&b, binary.LittleEndian, item.Byte); err != nil {
				return b.Bytes(), err
			}
			left := occ - (i * 255)
			if left < 255 {
				if err := binary.Write(&b, binary.LittleEndian, uint8(left)); err != nil {
					return b.Bytes(), err
				}
				for j := 0 + (i * 255); j < len(item.Addresses); j++ {
					if err := binary.Write(&b, binary.LittleEndian, item.Addresses[j]); err != nil {
						return b.Bytes(), err
					}
				}
			} else {
				if err := binary.Write(&b, binary.LittleEndian, uint8(254)); err != nil {
					return b.Bytes(), err
				}
				for j := 0 + (i * 255); j < 255; j++ {
					if err := binary.Write(&b, binary.LittleEndian, item.Addresses[j]); err != nil {
						return b.Bytes(), err
					}
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
		occ := len(item.Addresses)
		iter := occ / 255
		if occ%255 != 0 {
			iter++
		}
		for i := 0; i < iter; i++ {
			if err = binary.Write(f, binary.LittleEndian, item.Byte); err != nil {
				return err
			}
			left := occ - (i * 255)
			if left < 255 {
				if err = binary.Write(f, binary.LittleEndian, uint8(left)); err != nil {
					return err
				}
				for j := 0 + (i * 255); j < len(item.Addresses); j++ {
					if err = binary.Write(f, binary.LittleEndian, item.Addresses[j]); err != nil {
						return err
					}
				}
			} else {
				if err = binary.Write(f, binary.LittleEndian, uint8(254)); err != nil {
					return err
				}
				for j := 0 + (i * 255); j < 255; j++ {
					if err = binary.Write(f, binary.LittleEndian, item.Addresses[j]); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func Delta(scr1, scr2 []byte) *DeltaCollection {
	data := NewDeltaCollection()
	//var line int
	for offset := 0; offset < 0x4000; offset++ {
		//offsetLine := (offset % 80)
		if scr1[offset] != scr2[offset] {
			//offsetScreen := CpcScreenAddressOffset(line)
			addr := 0xC000 + offset //  + offsetLine + offsetScreen
			data.Add(scr2[offset], uint16(addr))
		}
		/*if offsetLine == 0 {
			line++
		}*/
	}
	return data
}
