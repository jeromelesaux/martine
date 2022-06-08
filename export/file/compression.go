package file

import (
	"fmt"
	"os"

	rawlz4 "github.com/bkaradzic/go-lz4"
	"github.com/jeromelesaux/martine/lz4"
	"github.com/jeromelesaux/martine/rle"
	zx0 "github.com/jeromelesaux/zx0/encode"
)

func Compress(data []byte, compression int) ([]byte, error) {
	var err0 error
	if compression != -1 {
		switch compression {
		case 1:
			fmt.Fprintf(os.Stdout, "Using RLE compression\n")
			data = rle.Encode(data)
		case 2:
			fmt.Fprintf(os.Stdout, "Using RLE 16 bits compression\n")
			data = rle.Encode16(data)
		case 3:
			fmt.Fprintf(os.Stdout, "Using LZ4 compression\n")
			var dst []byte
			dst, err := lz4.Encode(dst, data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while encoding into LZ4 : %v\n", err)
				err0 = err
			}
			data = dst
		case 4:
			fmt.Fprintf(os.Stdout, "Using LZ4-Raw compression\n")
			var dst []byte
			dst, err := rawlz4.Encode(dst, data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while encoding into LZ4 : %v\n", err)
				err0 = err
			}
			data = dst[4:]
		case 5:
			fmt.Fprintf(os.Stdout, "Using Zx0 cruncher")
			data = zx0.Encode(data)
		}
	}
	return data, err0
}
