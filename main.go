package main 

import (
	"fmt"
	"flag"
	"os"
	"strings"
	"github.com/jeromelesaux/screenverter/gfx"
)


var (
	picturePath = flag.String("p","","Picture path of the Amsdos file.")
	width  = flag.Int("w",-1,"Custom output width in pixels.")
	height = flag.Int("h",-1,"Custom output height in pixels.")
	mode = flag.String("m","","Output mode to use (mode0,mode1,mode2 or overscan available).")
	output = flag.String("o","","Output directory")
)

func main() {
	 var size gfx.Size
	flag.Parse()
	// picture path to convert 
	if *picturePath == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	// output directory to store results
	if *output != "" {
		fi, err := os.Stat(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error while getting directory informations :%v, Quiting\n",err)
			os.Exit(-2)
		}
		
		if !fi.IsDir() {
			fmt.Fprintf(os.Stderr,"%s is not a directory will store in current directory\n",*output)
			*output = "./"
		}
	}
	if *height != -1 && *width != -1 {
		fmt.Fprintf(os.Stderr,"Use the custom informations (width:%d, height:%d)\n",*width,*height)
		size.Height = *height
		size.Width = *width
	} else {
		if *mode == "" {
			fmt.Fprintf(os.Stderr,"No output mode defined can not choose. Quiting\n")
			flag.PrintDefaults()
			os.Exit(-2)
		}
		switch strings.ToLower(*mode) {
		case "mode0": size = gfx.Mode0
		case "mode1": size = gfx.Mode1
		case "mode2": size = gfx.Mode2
		case "overscan": size = gfx.Overscan
		default: fmt.Fprintf(os.Stderr,"mode %s not defined\n", *mode) 
		flag.PrintDefaults()
		os.Exit(-2)
		}
	}

	fmt.Fprintf(os.Stdout,"Informations :\n%s",size.ToString())

}