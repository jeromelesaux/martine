package animate

import (
	"fmt"
	"testing"

	"github.com/jeromelesaux/martine/gfx/transformation"
)

func TestExportHsp(t *testing.T) {
	c := &transformation.DeltaCollection{Items: []transformation.DeltaItem{}}
	item := transformation.DeltaItem{}
	item.Byte = 254
	item.Offsets = []uint16{0xC000, 0xC010, 0xC011, 0xD122}
	c.Items = append(c.Items, item)
	item.Byte = 15
	item.Offsets = []uint16{0xC001, 0xC002, 0xC050, 0xD000}
	c.Items = append(c.Items, item)
	code := ExportCompiledSprite(c)
	fmt.Println(code)

}
