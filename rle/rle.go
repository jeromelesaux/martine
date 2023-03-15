package rle

import (
	"bytes"
	"encoding/binary"

	"github.com/jeromelesaux/martine/log"
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

			if nb > 255 {
				for j := 0; j < (nb / 255); j++ {
					out = append(out, 255)
					out = append(out, c)
				}
				if nb%255 != 0 {
					out = append(out, uint8(nb%255))
					out = append(out, c)
				}
			} else {
				out = append(out, uint8(nb))
				out = append(out, c)
			}

			nb = 1
			c = d
			if i+1 == binary.Size(in) {

				if nb > 255 {
					for j := 0; j < (nb / 255); j++ {
						out = append(out, 255)
						out = append(out, c)
					}
					if nb%255 != 0 {
						out = append(out, uint8(nb%255))
						out = append(out, c)
					}
				} else {
					out = append(out, uint8(nb))
					out = append(out, c)
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
			if nb%255 != 0 {
				out = append(out, uint8(nb%255))
				out = append(out, c)
			}
		} else {
			out = append(out, uint8(nb))
			out = append(out, c)
		}
	}
	return out
}

func Encode16(in []byte) []byte {
	out := make([]byte, 0)
	var nb int16 = 1
	c := in[0]
	var d byte
	for i := 1; i < binary.Size(in); i++ {
		d = in[i]
		switch {
		case d != c: // valeurs differentes

			buf := new(bytes.Buffer)
			if err := binary.Write(buf, binary.LittleEndian, nb); err != nil {
				log.GetLogger().Error("Error while copying in byte buffer error :%v\n", err)
			}
			out = append(out, buf.Bytes()...)
			out = append(out, c)

			nb = 1
			c = d
			if i+1 == binary.Size(in) {
				buf := new(bytes.Buffer)
				if err := binary.Write(buf, binary.LittleEndian, nb); err != nil {
					log.GetLogger().Error("Error while copying in byte buffer error :%v\n", err)
				}
				out = append(out, buf.Bytes()...)
				out = append(out, c)
			}
		default:
			nb++
			continue
		}
	}
	if nb > 1 {
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, nb); err != nil {
			log.GetLogger().Error("Error while copying in byte buffer error :%v\n", err)
		}
		out = append(out, buf.Bytes()...)
		out = append(out, c)
	}
	return out
}
