package rle

import (
	"encoding/binary"
)

func Encode(in []byte) []byte {
	out := make([]byte, 0)
	nb := 1
	c := in[0]
	var d byte
	for i := 1; i < binary.Size(in); i++ {
		d = in[i]
		switch {
		case d != c: // valeurs differentes
			if nb > 1 {
				if nb > 255 {
					for j := 0; j < (nb / 255); j++ {
						out = append(out, 255)
						out = append(out, c)
					}
					out = append(out, uint8(nb%255))
					out = append(out, c)
				} else {
					out = append(out, uint8(nb))
					out = append(out, c)
				}
			}
			nb = 1
			c = d
			if i+1 == binary.Size(in) {
				if nb > 1 {
					if nb > 255 {
						for j := 0; j < (nb / 255); j++ {
							out = append(out, 255)
							out = append(out, c)
						}
						out = append(out, uint8(nb%255))
						out = append(out, c)
					} else {
						out = append(out, uint8(nb))
						out = append(out, c)
					}
				}
			}
		default:
			nb++
			continue
		}
	}
	if nb > 1 {
		if nb > 255 {
			for j := 0; j < (nb / 255); j++ {
				out = append(out, 255)
				out = append(out, c)
			}
			out = append(out, uint8(nb%255))
			out = append(out, c)
		} else {
			out = append(out, uint8(nb))
			out = append(out, c)
		}
	}
	return out
}
