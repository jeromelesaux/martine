package animate

import (
	"fmt"

	"github.com/jeromelesaux/martine/gfx/transformation"
)

func ExportCompiledSpriteHard(c *transformation.DeltaCollection) string {
	var code string
	var previous uint8 = 0
	for _, v := range c.Items {
		if v.Byte == 0 {
			code += "xor a"
		} else {
			code += fmt.Sprintf("ld a,#%.2x\n", v.Byte)
		}
		for _, v1 := range v.Offsets {
			if previous == uint8(v1-1) {
				code += "inc l : "
			}
			code += "ld (hl),a\n"
			previous = uint8(v1)
		}
	}
	code += "ret\n"
	return code
}

// sprite width, size to change the line
func ExportCompiledSprite(c *transformation.DeltaCollection) string {
	var code string
	var previous uint16 = 0
	var previousHB uint8 = 0

	code += `; HL contains the start screen address
	ld ix,hl`
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
					code += fmt.Sprintf("ld d,#%.2x : ", currentLB)
					code += "ld e,0 : ld hl, ix : add hl,de\n"
				} else {
					code += fmt.Sprintf("ld de,#%.4x : ", v1)
					code += "ld hl, ix : add hl,de\n"
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

func AnalyzeSpriteBoard(spr [][]byte) []*transformation.DeltaCollection {
	dc := make([]*transformation.DeltaCollection, 0)
	for i := 1; i < len(spr); i++ {
		l := spr[i-1]
		r := spr[i]
		c := transformation.NewDeltaCollection()
		for x := 0; x < len(l); x++ {
			if r[x] != l[x] {
				c.Add(r[x], uint16(x))
			}
		}
		dc = append(dc, c)
	}
	return dc
}
