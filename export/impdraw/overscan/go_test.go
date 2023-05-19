package overscan_test

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/diskimage"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/log"
)

func init() {
	log.InitLoggerWithFile("test.log")
}

func TestSaveGo(t *testing.T) {
	fileInput := "../../../samples/lena-512.png"
	f, _ := os.Open(fileInput)
	defer f.Close()
	img, _ := png.Decode(f)

	cfg := config.NewMartineConfig(fileInput, "")
	cfg.Overscan = true
	cfg.ExportAsGoFile = true
	cfg.Dsk = true
	cfg.CpcPlus = true
	cfg.Size = constants.NewSizeMode(0, true)
	err := gfx.ApplyOneImageAndExport(img, cfg, cfg.InputPath, filepath.Dir(fileInput), 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ExportAsGoFile = false
	// gfx.ApplyOneImageAndExport(img, cfg, cfg.InputPath, filepath.Dir(fileInput), 0, 0)
	err = diskimage.ImportInDsk(filepath.Dir(fileInput), cfg)
	if err != nil {
		t.Fatal(err)
	}
}
