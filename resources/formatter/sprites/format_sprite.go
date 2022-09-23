package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
)

type stringSlice []string

func (f *stringSlice) String() string {
	return ""
}

func (f *stringSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

//var deltaFiles stringSlice
var (
	//	sprites = flag.String("sprites", "", "sprites json path")
	outfile = flag.String("out", "", "output filepath to store data")
	eol     = "\n"
)

func main() {
	var spritesFiles stringSlice
	flag.Var(&spritesFiles, "in", "sprites json path")

	flag.Parse()
	if len(spritesFiles) == 0 {
		fmt.Fprintf(os.Stdout, "sprites is mandarory\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	var out string
	for _, sprite := range spritesFiles {
		fl := []string{sprite}
		datafiles, err := common.WilcardedFiles(fl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse wildcard files (%s) error :%v\n", sprite, err)
			os.Exit(-1)
		}
		for _, v := range datafiles {
			f, err := os.Open(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", v, err)
				os.Exit(-1)
			}
			defer f.Close()
			d := &export.Json{}
			if err := json.NewDecoder(f).Decode(d); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot decode json file (%s) error :%v\n", v, err)
				os.Exit(-1)
			}
			if len(d.Screen) != 1 {
				filename := strings.Replace(filepath.Base(sprite), filepath.Ext(sprite), "", -1)
				out += fmt.Sprintf("%s\n", filename)
				for i := 0; i < len(d.Screen); i += 8 {
					out += fmt.Sprintf("%s ", file.ByteToken)
					if i < len(d.Screen) {
						out += toDamsByte(d.Screen[i])
					}
					if i+1 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+1])
					}
					if i+2 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+2])
					}
					if i+3 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+3])
					}
					if i+4 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+4])
					}
					if i+5 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+5])
					}
					if i+6 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+6])
					}
					if i+7 < len(d.Screen) {
						out += ", " + toDamsByte(d.Screen[i+7])
					}
					out += eol
				}
			}
		}
	}
	if *outfile != "" {
		f, err := os.Create(*outfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating file (%s) error %v\n", *outfile, err)
			os.Exit(-1)
		}
		defer f.Close()
		f.WriteString(out)
	} else {
		fmt.Println(out)
	}
}

func toDamsByte(in string) string {
	return strings.ReplaceAll(in, "0x", "#")
}
