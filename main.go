package main 

import (
	"fmt"
	"flag"
	"os"
	cpc "github.com/jeromelesaux/m4client/cpc"
)


var (
	filePath = flag.String("filepath","","Filepath of the Amsdos file.")
)

func main() {
	flag.Parse()
	if *filePath == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	fh, err := os.Open(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while opening file %v, error :%v", *filePath,err)
		os.Exit(-2)
	}
	defer fh.Close()
	h, err := cpc.NewCpcHeader(fh)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while reading Amsdos header (file:%v), error :%v",*filePath, err)
		os.Exit(-2)
	}
	
	fmt.Fprintf(os.Stderr, "Header:(%s)", h.ToString())
	h.PrettyPrint()

}