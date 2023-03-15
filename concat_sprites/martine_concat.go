package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/common"
	"github.com/jeromelesaux/martine/log"
)

var (
	files = flag.String("files", "", "sprites files to concat (wildcard accepted such as ? or * file filename.) ")
	out   = flag.String("out", "", "output file to store data.")
)

func main() {
	flag.Parse()
	var err error
	spritesFiles := []string{*files}
	filespath, err := common.WilcardedFiles(spritesFiles)
	if err != nil {
		log.GetLogger().Error( "error while getting the files path. %v\n", err)
		os.Exit(-1)
	}
	f, err := os.Create(*out)
	if err != nil {
		log.GetLogger().Error( "error while opening file %s error : %v\n.", *out, err)
		os.Exit(-1)
	}
	defer f.Close()
	for _, v := range filespath {
		displaySprite(v, f)
	}
	os.Exit(0)
}

func displaySprite(filePath string, out *os.File) {
	f, err := os.Open(filePath)
	if err != nil {
		log.GetLogger().Error( "Cannot open file (%s) error :%v\n", filePath, err)
		return
	}
	defer f.Close()
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	spriteName := strings.Replace(base, ext, "", 1)
	_, err = out.Write([]byte(fmt.Sprintf("%s\n", strings.ToLower(spriteName))))
	if err != nil {
		log.GetLogger().Error( "Error while writing : %v\n", err)
		return
	}
	scanner := bufio.NewScanner(f)

	scanner.Scan() // remove amsdos header
	for scanner.Scan() {
		in := scanner.Text()

		switch in[0] {
		case ';':
			// end of the data
			return
		case '.':
			continue
		default:
			_, err = out.Write([]byte(fmt.Sprintf("%s\n", in)))
			if err != nil {
				log.GetLogger().Error( "Error while writing : %v\n", err)
				return
			}
		}

	}
}
