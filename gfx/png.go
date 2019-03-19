package gfx

import (
	"fmt"
	"image"
	"os"
	"image/png"
)

func Png(filePath string, im *image.NRGBA) error {
    fwd, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create new image (%s) error %v\n", filePath, err)
		return err
	}

	if err := png.Encode(fwd, im); err != nil {
		fwd.Close()
		fmt.Fprintf(os.Stderr, "Cannot create new image (%s) as png error %v\n", filePath, err)
		return err
	}
	fwd.Close()
	return nil
}