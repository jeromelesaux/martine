package animate

import (
	"fmt"

	"github.com/jeromelesaux/martine/gfx/transformation"
)

func ExportHsp(c *transformation.DeltaCollection) string {
	var code string
	var previous uint16 = 0
	var previousHB uint8 = 0

	for _, v := range c.Items {
		if v.Byte == 0 {
			code += "xor a"
		} else {
			code += fmt.Sprintf("ld a,#%.2x\n", v.Byte)
		}
		for _, v1 := range v.Offsets {
			if previous == (v1 - 1) {
				code += "inc l : "
			} else {
				currentHB := uint8(v1 >> 8)
				currentLB := uint8(v1)
				if previousHB == currentHB {
					code += fmt.Sprintf("ld l,#%.2x : ", currentLB)
				} else {
					code += fmt.Sprintf("ld hl,#%.4x : ", v1)
				}
				previousHB = currentHB
			}
			code += "ld (hl),a\n"
			previous = v1
		}
	}
	code += "ret\n"
	return code
}
