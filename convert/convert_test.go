package convert

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeromelesaux/martine/constants"
)

func TestMask(t *testing.T) {
	bl := constants.Green
	_, g, _, _ := bl.Color.RGBA()
	r2 := uint8(g)
	r := r2
	t.Logf("%b\n", g)
	//r2 |= r2 - 128
	//r = r ^ 128
	if r > 128 {
		r ^= 128
	}
	fmt.Fprintf(os.Stdout, "%b - %b = %b\n", r, 128, r^128)
	fmt.Fprintf(os.Stdout, "r:%b,%d\n", r2, r2)
}
