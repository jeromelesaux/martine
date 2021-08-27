package constants

import (
	"encoding/binary"
	"fmt"
	"image/color"
)

type CpcPlusColor struct {
	G uint16
	R uint16
	B uint16
}

func NewRawCpcPlusColor(v uint16) *CpcPlusColor {
	c := &CpcPlusColor{}
	c.B = v & 0xf //1111
	c.R = v >> 4 & 0xf
	c.G = v >> 8 & 0xf
	return c
}

func (c *CpcPlusColor) ToString() string {
	return fmt.Sprintf("R:%.2b(%d),G:%.2b(%d),B:%.2b(%d)", c.R, c.R, c.G, c.G, c.B, c.B)
}

func (c *CpcPlusColor) Value() uint16 {
	v := c.B | c.R<<4 | c.G<<8
	//fmt.Fprintf(os.Stderr, "value(%d,%d,%d)(%b,%b,%b) #%.4x (%.b): %d\n", c.R, c.G, c.B, c.R, c.G, c.B,
	//	v, v, c.B+(c.R*16)+c.G*256)
	return v
}
func (c *CpcPlusColor) Bytes() []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, c.Value())
	//	fmt.Fprintf(os.Stderr, "%b\n", buf)
	return buf
}

func NewCpcPlusColor(c color.Color) CpcPlusColor {
	r, g, b, _ := c.RGBA()
	//	fmt.Fprintf(os.Stderr,"original colors r:%d,g:%d,b:%d\n",r,g,b)
	return CpcPlusColor{G: uint16(g / 4096), R: uint16(r / 4096), B: uint16(b / 4096)}
}

func NewColorCpcPlusColor(c CpcPlusColor) color.Color {
	return color.RGBA{G: uint8(c.G << 4), R: uint8(c.R << 4), B: uint8(c.B << 4), A: 0xff}
}
