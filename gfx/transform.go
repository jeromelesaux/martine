package gfx

import (
	"path/filepath"
	"encoding/binary"
	"image"
	"image/color"
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/jeromelesaux/m4client/cpc"
)

var (
	ColorNotFound = errors.New("Color not found in palette")
	NotYetImplemented = errors.New("Function is not yet implemented")
)

func Transform(in *image.NRGBA,p color.Palette, size Size, filepath string) error {
	switch size {
	case Mode0:
		return TransformMode0(in,p,filepath)
	default: 
		return NotYetImplemented
	}
	return nil
}

func PalettePosition(c color.Color, p color.Palette) (int,error) {
	r,g,b,a := c.RGBA()
	for index, cp := range p {
		//fmt.Fprintf(os.Stdout,"index(%d), c:%v,cp:%v\n",index,c,cp)
		rp, gp, bp, ap  := cp.RGBA()
		if r == rp && g == gp && b == bp && a == ap {
			//fmt.Fprintf(os.Stdout,"Position found")
			return index,nil
		}
	}
	return -1,ColorNotFound
}

func TransformMode0(in *image.NRGBA, p color.Palette, filePath string) error {
	 bw := make([]byte,0)
	fmt.Fprintf(os.Stdout,"Informations palette (%d) for image (%d,%d)\n",len(p),in.Bounds().Max.X, in.Bounds().Max.Y)
	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x+=2 {
				c1 := in.At(x,y)
				pp1,err := PalettePosition(c1,p)
				if err != nil {
					fmt.Fprintf(os.Stderr,"%v pixel position(%d,%d) not found in palette\n",c1,x,y)
					continue
				} 
				//fmt.Fprintf(os.Stdout,"(%d,%d), %v, position palette %d\n",x,y,c1,pp1)
				c2 := in.At(x+1,y)
				pp2,err := PalettePosition(c2,p)
				if err != nil {
					fmt.Fprintf(os.Stderr,"%v pixel position(%d,%d) not found in palette\n",c2,x+1,y)
					continue
				} 
				//fmt.Fprintf(os.Stdout,"(%d,%d), %v, position palette %d\n",x+1,y,c1,pp2)
				var pixel0, pixel1 byte
				pixel0 = (uint8(pp2) & 128)<<7  + (uint8(pp1) & 32)<<4  + (uint8(pp2) & 8)<<1 + (uint8(pp1) & 2)>>2
				pixel1 = (uint8(pp2) & 64 )<<6 + (uint8(pp1) & 16)<<3  + (uint8(pp2) & 4) + (uint8(pp1) & 1)>>3 
				bw = append(bw , pixel0)
				bw = append(bw, pixel1)
			}	
	}

	header := cpc.CpcHead{}

	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	cpcFilename := strings.Replace(filename,extension, ".SCR",-1)
	copy(header.Filename[:],cpcFilename)
	header.Type = 1
	header.User = 0
	header.Address = 0x4000
	header.Size = 0x4000
	header.Size2 = 0x400
	header.Checksum = int16(header.ComputedChecksum16())
	fw, err := os.Create(cpcFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while creating file (%s) error :%s",cpcFilename,err)
		return err
	}
	binary.Write(fw,binary.LittleEndian,header)
	fw.Write(bw)
	fw.Close()

	return nil
}