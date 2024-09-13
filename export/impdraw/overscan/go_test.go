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
	_, _ = log.InitLoggerWithFile("test.log")
}

func TestSaveGo(t *testing.T) {
	_ = os.Mkdir("overscanTests", 0777)
	file := "../../../samples/lena-512.png"
	f, _ := os.Open(file)
	defer f.Close()
	img, _ := png.Decode(f)

	cfg := config.NewMartineConfig(file, "./overscanTests")
	cfg.ScrCfg.Type = config.FullscreenFormat
	cfg.ScrCfg.AddExport(config.GoImpdrawExport)
	cfg.ContainerCfg.AddExport(config.DskContainer)
	cfg.ScrCfg.IsPlus = true
	cfg.ScrCfg.Size = constants.NewSizeMode(0, true)
	err := gfx.ApplyOneImageAndExport(img, cfg, "lena", "./overscanTests/", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	// gfx.ApplyOneImageAndExport(img, cfg, cfg.InputPath, filepath.Dir(fileInput), 0, 0)
	err = diskimage.ImportInDsk(filepath.Dir(file), cfg)
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll("overscanTests")
}
