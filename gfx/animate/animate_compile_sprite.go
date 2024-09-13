package animate

import (
	"fmt"

	"github.com/jeromelesaux/martine/gfx/transformation"
)

type Z80Register string

var (
	A            Z80Register = "a"
	B            Z80Register = "b"
	C            Z80Register = "c"
	D            Z80Register = "d"
	E            Z80Register = "e"
	H            Z80Register = "h"
	L            Z80Register = "l"
	NoneRegister Z80Register = ""
)

type Z80HspNode struct {
	register          Z80Register
	byte              uint8
	offset            uint16
	next              *Z80HspNode
	previous          *Z80HspNode
	samePreviousValue bool
}

func NewZ80HspNode(byte uint8, offset uint16, samevalue bool, register Z80Register, next *Z80HspNode) *Z80HspNode {
	return &Z80HspNode{
		register:          register,
		byte:              byte,
		next:              next,
		offset:            offset,
		samePreviousValue: samevalue,
	}
}

func (z *Z80HspNode) SetLastNode(node *Z80HspNode) {
	if z.next == nil {
		z.next = node
		node.previous = z
		return
	} else {
		next := z.next
		next.SetLastNode(node)
	}
}

func (z *Z80HspNode) NextRegister() Z80Register {
	registers := make([]Z80Register, 0)
	return z.internalNextRegister(registers)
}

func (z *Z80HspNode) internalNextRegister(registers []Z80Register) Z80Register {

	registers = append(registers, z.register)
	if z.next != nil {
		return z.next.internalNextRegister(registers)
	}
	if registers[len(registers)-1] == NoneRegister {
		return A
	}
	if registers[len(registers)-1] == A {
		return B
	}
	if registers[len(registers)-1] == B {
		return C
	}
	if registers[len(registers)-1] == C {
		return D
	}
	if registers[len(registers)-1] == D {
		return E
	}
	if registers[len(registers)-1] == E {
		return A
	}
	return A
}

// nolint: gocognit
func (z *Z80HspNode) InitOpcode() string {

	switch z.register {
	case A:
		if z.samePreviousValue {
			return ""
		}
		return fmt.Sprintf("ld a,%d\n", z.byte)
	case B:
		if z.next != nil && z.next.register == C {
			return fmt.Sprintf("ld bc,#%.4x\n", (uint16(z.byte)<<8)+uint16(z.next.byte))
		}
		if z.samePreviousValue {
			return ""
		}
		return fmt.Sprintf("ld b,%d\n", z.byte)
	case C:
		if z.previous != nil && z.previous.register == B {
			return ""
		}
		return fmt.Sprintf("ld c,%d\n", z.byte)
	case D:
		if z.next != nil && z.next.register == E {
			return fmt.Sprintf("ld de,#%.4x\n", (uint16(z.byte)<<8)+uint16(z.next.byte))
		}
		if z.samePreviousValue {
			return ""
		}
		return fmt.Sprintf("ld d,%d\n", z.byte)
	case E:
		if z.previous != nil && z.previous.register == D {
			return ""
		}
		if z.samePreviousValue {
			return ""
		}
		return fmt.Sprintf("ld e,%d\n", z.byte)
	case H:
		if z.next != nil && z.next.register == L {
			return fmt.Sprintf("ld hl,#%.4x\n", (uint16(z.byte)<<8)+uint16(z.next.byte))
		}
		if z.samePreviousValue {
			return ""
		}
		return fmt.Sprintf("ld h,%d\n", z.byte)
	case L:
		if z.previous != nil && z.previous.register == H {
			return ""
		}
		if z.samePreviousValue {
			return ""
		}
		return fmt.Sprintf("ld l,%d\n", z.byte)
	}
	return ""
}

func (z *Z80HspNode) OffsetInit() string {
	if z.previous != nil {
		if z.offset-1 != z.previous.offset {
			return fmt.Sprintf("ld l,#%.2x\n", uint8(z.offset))
		} else {
			return "inc l\n"
		}
	}
	return fmt.Sprintf("ld l,#%.2x\n", uint8(z.offset))
}

func (z *Z80HspNode) ValueOpcode() string {
	switch z.register {
	case A:
		return "ld (hl),a\n"
	case B:
		return "ld (hl),b\n"
	case C:
		return "ld (hl),c\n"
	case D:
		return "ld (hl),d\n"
	case E:
		return "ld (hl),e\n"
	default:
		return ""
	}
}

func (z *Z80HspNode) Code() string {
	if z.next != nil {
		return z.next.internalCode("")
	}
	return ""
}

func (z *Z80HspNode) internalCode(code string) string {
	code += z.InitOpcode()
	code += z.OffsetInit()
	code += z.ValueOpcode()
	if z.next == nil {
		return code
	}
	return z.next.internalCode(code)
}

func ExportCompiledSpriteHard(c *transformation.DeltaCollection) string {
	items := c.ItemsSortByByte()
	optim := NewZ80HspNode(0, 0, true, NoneRegister, nil)
	for _, v := range items {
		var already = false
		reg := optim.NextRegister()
		for _, offset := range v.Offsets {
			node := NewZ80HspNode(v.Byte, offset, already, reg, nil)
			optim.SetLastNode(node)
			already = true
		}
	}
	code := optim.Code()
	code += "ret\n"
	return code
}

// sprite width, size to change the line
func ExportCompiledSprite(c *transformation.DeltaCollection) string {
	items := c.ItemsSortByByte()
	optim := NewZ80HspNode(0, 0, true, NoneRegister, nil)
	for _, v := range items {
		var already = false
		reg := optim.NextRegister()
		for _, offset := range v.Offsets {
			node := NewZ80HspNode(v.Byte, offset, already, reg, nil)
			optim.SetLastNode(node)
			already = true
		}
	}
	code := optim.Code()
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
