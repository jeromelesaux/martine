package rle

import (
	"encoding/binary"
)

func Encode(in []byte) []byte {
	out := make([]byte, 0)
	var nb uint8 = 1
	c := in[0]
	var d byte
	for i := 1; i < binary.Size(in); i++ {
		d = in[i]
		switch {
		case d != c: // valeurs differentes
			out = append(out, nb)
			out = append(out, c)
			nb = 1
			c = d
			if i+1 == binary.Size(in) {
				out = append(out, nb)
				out = append(out, d)
			}
		default:
			nb++
			continue
		}
	}
	if nb > 1 {
		out = append(out, nb)
		out = append(out, c)
	}
	return out
}
