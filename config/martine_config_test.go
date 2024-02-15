package config_test

import (
	"testing"

	"github.com/jeromelesaux/martine/config"
)

func TestAmsdosFullpath(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		config := config.NewMartineConfig("", "")
		v := "/Users/jls/disk.bas"
		got := config.GetAmsdosFilename(v, ".BAS")
		if got != "DISK.BAS" {
			t.Fatalf("Expected DISK.BAS and gets %s\n", got)
		}
	})
	t.Run("Success", func(t *testing.T) {
		config := config.NewMartineConfig("", "")
		v := "/Users/jls/disk.bas"
		got := config.GetAmsdosFilename(v, ".bas")
		if got != "DISK.BAS" {
			t.Fatalf("Expected DISK.BAS and gets %s\n", got)
		}
	})

	t.Run("LargeSize", func(t *testing.T) {
		config := config.NewMartineConfig("", "")
		v := "/Users/jls/diskletsseeifhewillremove.bas"
		got := config.GetAmsdosFilename(v, ".bas")
		if got != "DISKLETE.BAS" {
			t.Fatalf("Expected DISKLETE.BAS and gets %s\n", got)
		}
	})

	t.Run("RemoveUnsupportedChar", func(t *testing.T) {
		config := config.NewMartineConfig("", "")
		v := "/Users/jls/disk-.-_.bas"
		got := config.GetAmsdosFilename(v, ".bas")
		if got != "DISK.BAS" {
			t.Fatalf("Expected DISK.BAS and gets %s\n", got)
		}
	})
}
