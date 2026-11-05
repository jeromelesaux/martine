package config_test

import (
	"testing"

	"github.com/jeromelesaux/martine/config"
)

const diskBAS = "DISK.BAS"

func TestAmsdosFullpath(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		cfg := config.NewMartineConfig("", "")
		v := "/Users/jls/disk.bas"
		got := cfg.GetAmsdosFilename(v, ".BAS")
		if got != diskBAS {
			t.Fatalf("Expected DISK.BAS and gets %s\n", got)
		}
	})
	t.Run("Success", func(t *testing.T) {
		cfg := config.NewMartineConfig("", "")
		v := "/Users/jls/disk.bas"
		got := cfg.GetAmsdosFilename(v, ".bas")
		if got != diskBAS {
			t.Fatalf("Expected DISK.BAS and gets %s\n", got)
		}
	})

	t.Run("LargeSize", func(t *testing.T) {
		cfg := config.NewMartineConfig("", "")
		v := "/Users/jls/diskletsseeifhewillremove.bas"
		got := cfg.GetAmsdosFilename(v, ".bas")
		if got != "DISKLETE.BAS" {
			t.Fatalf("Expected DISKLETE.BAS and gets %s\n", got)
		}
	})

	t.Run("RemoveUnsupportedChar", func(t *testing.T) {
		cfg := config.NewMartineConfig("", "")
		v := "/Users/jls/disk-.-_.bas"
		got := cfg.GetAmsdosFilename(v, ".bas")
		if got != diskBAS {
			t.Fatalf("Expected DISK.BAS and gets %s\n", got)
		}
	})
}
