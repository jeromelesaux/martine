package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/ascii"
)

var (
	sprite     = flag.String("sprite", "", "sprite json path")
	deltafiles = flag.String("delta", "", "delta wildcarded json file paths")
	outfile    = flag.String("out", "", "output filepath to store data")
	eol        = "\n"
)

func main() {
	flag.Parse()
	if *sprite == "" || *deltafiles == "" {
		fmt.Fprintf(os.Stdout, "sprite and delta options are mandarories\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	f, err := os.Open(*sprite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file (%s) error %v\n", *sprite, err)
		os.Exit(-1)
	}
	defer f.Close()
	s := &export.Json{}
	if err := json.NewDecoder(f).Decode(s); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot decode json file (%s) error :%v\n", *sprite, err)
		os.Exit(-1)
	}
	out := "sprite\n"
	for i := 0; i < len(s.Screen); i += 8 {
		out += fmt.Sprintf("%s ", ascii.ByteToken)
		if i < len(s.Screen) {
			out += toDamsByte(s.Screen[i])
		}
		if i+1 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+1])
		}
		if i+2 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+2])
		}
		if i+3 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+3])
		}
		if i+4 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+4])
		}
		if i+5 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+5])
		}
		if i+6 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+6])
		}
		if i+7 < len(s.Screen) {
			out += ", " + toDamsByte(s.Screen[i+7])
		}
		out += eol
	}
	fl := []string{*deltafiles}
	datafiles, err := common.WilcardedFiles(fl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse wildcard files (%s) error :%v\n", *deltafiles, err)
		os.Exit(-1)
	}
	for index, v := range datafiles {
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
			out += fmt.Sprintf("delta%.2d\n", index)
			for i := 0; i < len(d.Screen); i += 8 {
				out += fmt.Sprintf("%s ", ascii.ByteToken)
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
