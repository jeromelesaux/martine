package main 

import (
	"fmt"
	"flag"
	"os"
	cpc "github.com/jeromelesaux/m4client/cpc"
)


var (
	picturePath = flag.String("p","","Picture path of the Amsdos file.")
)

func main() {
	flag.Parse()
	if *picturePath == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	fh, err := os.Open(*picturePath)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while opening file %v, error :%v", *picturePath,err)
		os.Exit(-2)
	}
	defer fh.Close()
	h, err := cpc.NewCpcHeader(fh)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while reading Amsdos header (file:%v), error :%v",*picturePath, err)
		os.Exit(-2)
	}
	
	fmt.Fprintf(os.Stderr, "Header:(%s)", h.ToString())
	h.PrettyPrint()

}