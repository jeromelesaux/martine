package lz4

import (
	"bytes"
	"github.com/pierrec/lz4"
	"io"
)

func Encode(dst, src []byte) ([]byte, error) {
	header := lz4.Header{}
	r := bytes.NewReader(src)
	var zout bytes.Buffer
	zw := lz4.NewWriter(&zout)
	zw.Header = header
	_, err := io.Copy(zw, r)
	if err != nil {
		return nil, err
	}
	err = zw.Close()
	dst = zout.Bytes()
	return dst, err
}
