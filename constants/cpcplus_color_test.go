package constants_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/jeromelesaux/martine/constants"
)

func TestCpcColorPlusConvertion(t *testing.T) {
	c := constants.CpcPlusColor{
		R: 128,
		G: 112,
		B: 96,
	}
	v := c.Value()
	fmt.Printf("#%.2X\n", v)
	fmt.Printf("%b\n", v)
	b := c.Bytes()
	fmt.Printf("#%.2X\n", b)
}

func TestReadColorPlus(t *testing.T) {
	r := []byte{0x39, 0x06}
	v := binary.LittleEndian.Uint16(r)
	c := constants.NewRawCpcPlusColor(v)
	vl := c.Value()
	fmt.Printf("#%.2X\n", vl)
	b := c.Bytes()
	fmt.Printf("#%.2X\n", b)
}
