package compression

import (
	rawlz4 "github.com/bkaradzic/go-lz4"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/lz4"
	"github.com/jeromelesaux/martine/rle"
	zx0 "github.com/jeromelesaux/zx0/encode"
)

type CompressionMethod int

const (
	NONE CompressionMethod = iota
	RLE
	RLE16
	LZ4
	RawLZ4
	ZX0
)

func ToCompressMethod(val int) CompressionMethod {
	switch val {
	case -1:
		return NONE
	case 1:
		return RLE
	case 2:
		return RLE16
	case 3:
		return LZ4
	case 4:
		return RawLZ4
	case 5:
		return ZX0
	default:
		return NONE
	}
}

func Compress(data []byte, compression CompressionMethod) ([]byte, error) {
	var err0 error
	if compression != NONE {
		switch compression {
		case RLE:
			log.GetLogger().Info("Using RLE compression\n")
			data = rle.Encode(data)
		case RLE16:
			log.GetLogger().Info("Using RLE 16 bits compression\n")
			data = rle.Encode16(data)
		case LZ4:
			log.GetLogger().Info("Using LZ4 compression\n")
			var dst []byte
			dst, err := lz4.Encode(dst, data)
			if err != nil {
				log.GetLogger().Error("Error while encoding into LZ4 : %v\n", err)
				err0 = err
			}
			data = dst
		case RawLZ4:
			log.GetLogger().Info("Using LZ4-Raw compression\n")
			var dst []byte
			dst, err := rawlz4.Encode(dst, data)
			if err != nil {
				log.GetLogger().Error("Error while encoding into LZ4 : %v\n", err)
				err0 = err
			}
			data = dst[4:]
		case ZX0:
			log.GetLogger().Info("Using Zx0 cruncher")
			data = zx0.Encode(data)
		}
	}
	return data, err0
}
