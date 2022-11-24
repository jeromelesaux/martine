package overscan

import (
	"fmt"
	"image/color"
	"os"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/convert/screen"
)

func ToGo(data []byte, screenMode uint8, p color.Palette, cfg *config.MartineConfig) ([]byte, error) {
	orig, err := OverscanRawToImg(data, screenMode, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while converting into image  error :%v", err)
		return nil, err
	}
	cfg.DoubleScreenAddress = true
	switch screenMode {
	case 0:
		data = screen.ToMode0(orig, p, cfg)
	case 1:
		data = screen.ToMode1(orig, p, cfg)
	case 2:
		data = screen.ToMode2(orig, p, cfg)
	}
	cfg.DoubleScreenAddress = false
	return data, nil
}
