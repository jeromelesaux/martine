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
	"github.com/jeromelesaux/martine/log"
)

var (
	sprite     = flag.String("sprite", "", "sprite json path")
	deltafiles = flag.String("delta", "", "delta wildcarded json file paths")
	outfile    = flag.String("out", "", "output filepath to store data")
	eol        = "\n"
)

// nolint: funlen
func main() {
	log.Default("")
	flag.Parse()
	if *sprite == "" || *deltafiles == "" {
		log.GetLogger().Info("sprite and delta options are mandarories\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	f, err := os.Open(*sprite)
	if err != nil {
		log.GetLogger().Error("Error while opening file (%s) error %v\n", *sprite, err)
		os.Exit(-1)
	}
	defer f.Close()
	s := &export.Json{}
	if err := json.NewDecoder(f).Decode(s); err != nil {
		log.GetLogger().Error("Cannot decode json file (%s) error :%v\n", *sprite, err)
		os.Exit(-1)
	}
	out := "sprite\n"
	for i := 0; i < len(s.Screen); i += 8 {
		out += fmt.Sprintf("%s ", ascii.ByteToken)
		out += convertScreen(out, s.Screen)
		out += eol
	}
	fl := []string{*deltafiles}
	datafiles, err := common.WilcardedFiles(fl)
	if err != nil {
		log.GetLogger().Error("Cannot parse wildcard files (%s) error :%v\n", *deltafiles, err)
		os.Exit(-1)
	}
	for index, v := range datafiles {
		f, err := os.Open(v)
		if err != nil {
			log.GetLogger().Error("Error while opening file (%s) error %v\n", v, err)
			os.Exit(-1)
		}
		defer f.Close()
		d := &export.Json{}
		if err := json.NewDecoder(f).Decode(d); err != nil {
			log.GetLogger().Error("Cannot decode json file (%s) error :%v\n", v, err)
			os.Exit(-1)
		}
		if len(d.Screen) != 1 {
			out += fmt.Sprintf("delta%.2d\n", index)
			for i := 0; i < len(d.Screen); i += 8 {
				out += fmt.Sprintf("%s ", ascii.ByteToken)
				out += convertScreen(out, d.Screen)
				out += eol
			}
		}
	}
	if *outfile != "" {
		f, err := os.Create(*outfile)
		if err != nil {
			log.GetLogger().Error("Error while creating file (%s) error %v\n", *outfile, err)
			os.Exit(-1)
		}
		defer f.Close()
		_, err = f.WriteString(out)
		if err != nil {
			log.GetLogger().Error("Error while writing file (%s) error %v\n", *outfile, err)
			os.Exit(-1)
		}

	} else {
		log.GetLogger().Infoln(out)
	}
}

func toDamsByte(in string) string {
	return strings.ReplaceAll(in, "0x", "#")
}

func convertScreen(out string, screen []string) string {
	for i := 0; i < len(screen); i += 8 {
		out += fmt.Sprintf("%s ", ascii.ByteToken)
		if i < len(screen) {
			out += toDamsByte(screen[i])
		}
		if i+1 < len(screen) {
			out += ", " + toDamsByte(screen[i+1])
		}
		if i+2 < len(screen) {
			out += ", " + toDamsByte(screen[i+2])
		}
		if i+3 < len(screen) {
			out += ", " + toDamsByte(screen[i+3])
		}
		if i+4 < len(screen) {
			out += ", " + toDamsByte(screen[i+4])
		}
		if i+5 < len(screen) {
			out += ", " + toDamsByte(screen[i+5])
		}
		if i+6 < len(screen) {
			out += ", " + toDamsByte(screen[i+6])
		}
		if i+7 < len(screen) {
			out += ", " + toDamsByte(screen[i+7])
		}
		out += eol
	}
	return out
}
