package ui

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHexConversion(t *testing.T) {
	a := "c000"
	v, err := strconv.ParseUint(a, 16, 64)
	t.Log(err)
	t.Logf("%d\n", v)
}

func TestWidth(t *testing.T) {
	width := 24
	mode := 1
	var colorPerPixel int

	switch mode {
	case 0:
		colorPerPixel = 2
	case 1:
		colorPerPixel = 4
	case 2:
		colorPerPixel = 8
	}
	remain := width % colorPerPixel
	fmt.Printf("%d\n", remain)

}
