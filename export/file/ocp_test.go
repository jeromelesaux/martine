package file

import (
	"testing"

	"github.com/jeromelesaux/martine/constants"
)

func TestCpcPlusColor(t *testing.T) {
	c1 := constants.CpcPlusColor{G: 0xa, R: 0xf, B: 0xf}
	v := c1.Value()
	c2 := constants.NewRawCpcPlusColor(v)
	t.Logf("%b,%b\n", 0xFF0000, v>>8)
	t.Logf("C1:%s\n", c1.ToString())
	t.Logf("C2:%s\n", c2.ToString())
	t.Logf("\nc1:%s\nc2:%s\nv:%b\n", c1.ToString(), c2.ToString(), v)
	if c1.ToString() != c2.ToString() {
		t.Fatalf("expected value %s and gets %s", c1.ToString(), c2.ToString())
	}
}

func TestKitPalette(t *testing.T) {

}
