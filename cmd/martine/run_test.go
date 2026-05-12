package main

import (
	"testing"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
)

func TestSetDitheringValidAlgorithm(t *testing.T) {
	cfg := config.NewMartineConfig("in.scr", "out")
	originalAlgo := *ditheringAlgo
	defer func() { *ditheringAlgo = originalAlgo }()

	*ditheringAlgo = 0
	if err := setDithering(cfg); err != nil {
		t.Fatalf("expected no error for algorithm 0, got %v", err)
	}
	if cfg.ScrCfg.Process.Dithering.Type != constants.ErrorDiffusionDither {
		t.Fatalf("expected ErrorDiffusionDither, got %v", cfg.ScrCfg.Process.Dithering.Type)
	}

	*ditheringAlgo = 7
	if err := setDithering(cfg); err != nil {
		t.Fatalf("expected no error for algorithm 7, got %v", err)
	}
	if cfg.ScrCfg.Process.Dithering.Type != constants.OrderedDither {
		t.Fatalf("expected OrderedDither, got %v", cfg.ScrCfg.Process.Dithering.Type)
	}
}

func TestSetDitheringInvalidAlgorithm(t *testing.T) {
	cfg := config.NewMartineConfig("in.scr", "out")
	originalAlgo := *ditheringAlgo
	defer func() { *ditheringAlgo = originalAlgo }()

	*ditheringAlgo = 99
	if err := setDithering(cfg); err == nil {
		t.Fatal("expected error for invalid dithering algorithm, got nil")
	}
}

func TestRunContainerExportsNoContainer(t *testing.T) {
	cfg := config.NewMartineConfig("in.scr", "out")
	if err := runContainerExports(cfg, "in.scr", "out"); err != nil {
		t.Fatalf("expected no error when no container export is configured, got %v", err)
	}
}
