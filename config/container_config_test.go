package config_test

import (
	"testing"

	"github.com/jeromelesaux/martine/config"
)

func TestContainerConfigExports(t *testing.T) {
	cfg := config.NewMartineConfig("in.scr", "out")
	if cfg.HasContainerExport(config.DskContainer) {
		t.Fatal("expected no dsk export by default")
	}

	cfg.ContainerCfg.AddExport(config.DskContainer)
	if !cfg.HasContainerExport(config.DskContainer) {
		t.Fatal("expected dsk export after adding it")
	}

	cfg.ContainerCfg.AddExport(config.SnaContainer)
	if !cfg.ContainerCfg.HasExport(config.SnaContainer) {
		t.Fatal("expected sna export after adding it")
	}

	cfg.ContainerCfg.RemoveExport(config.DskContainer)
	if cfg.HasContainerExport(config.DskContainer) {
		t.Fatal("expected dsk export to be removed")
	}
}

func TestNewMartineConfigDefaults(t *testing.T) {
	cfg := config.NewMartineConfig("input.png", "output")
	if cfg.Tiles == nil {
		t.Fatal("expected Tiles slice to be initialized")
	}
	if cfg.InkSwapper == nil {
		t.Fatal("expected InkSwapper map to be initialized")
	}
	if cfg.LineWidth != 0x50 {
		t.Fatalf("expected LineWidth 0x50, got %d", cfg.LineWidth)
	}
}
