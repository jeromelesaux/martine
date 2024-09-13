package constants

import (
	"encoding/binary"
	"fmt"
	"image/color"
)

type CpcPlusColor struct {
	G uint8
	R uint8
	B uint8
}

func NewRawCpcPlusColor(v uint16) *CpcPlusColor {
	c := &CpcPlusColor{}
	c.B = uint8(v & 0xf) //1111
	c.R = uint8(v >> 4 & 0xf)
	c.G = uint8(v >> 8 & 0xf)
	return c
}

func (c *CpcPlusColor) ToString() string {
	return fmt.Sprintf("R:%.2b(%d),G:%.2b(%d),B:%.2b(%d)", c.R, c.R, c.G, c.G, c.B, c.B)
}

func (c *CpcPlusColor) Value() uint16 {
	v := uint16(c.G)<<8 | uint16(c.B) | uint16(c.R)<<4
	//fmt.Fprintf(os.Stderr, "value(%d,%d,%d)(%b,%b,%b) #%.4x (%.b): %d\n", c.R, c.G, c.B, c.R, c.G, c.B,

	return v
}
func (c *CpcPlusColor) Bytes() []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, c.Value())
	return buf
}

func NewCpcPlusColor(c color.Color) CpcPlusColor {
	r, g, b, _ := c.RGBA()
	//	fmt.Fprintf(os.Stderr,"original colors r:%d,g:%d,b:%d\n",r,g,b)
	return CpcPlusColor{G: uint8(g>>8) / 16, R: uint8(r>>8) / 16, B: uint8(b>>8) / 16}
}

func NewColorCpcPlusColor(c CpcPlusColor) color.Color {
	return color.RGBA{G: c.G * 16, R: c.R * 16, B: c.B * 16, A: 0xff}
}
