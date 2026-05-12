package config_test

import (
	"reflect"
	"testing"

	"github.com/jeromelesaux/martine/config"
)

func TestModeMaskSprite(t *testing.T) {
	tests := []struct {
		mode    uint8
		want    []uint8
		wantErr bool
	}{
		{mode: 0, want: []uint8{0xAA, 0x55}, wantErr: false},
		{mode: 1, want: []uint8{0x88, 0x44, 0x24, 0x11}, wantErr: false},
		{mode: 2, want: []uint8{}, wantErr: true},
	}

	for _, tt := range tests {
		got, err := config.ModeMaskSprite(tt.mode)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("ModeMaskSprite(%d) expected error, got nil", tt.mode)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ModeMaskSprite(%d) unexpected error: %v", tt.mode, err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("ModeMaskSprite(%d) = %v, want %v", tt.mode, got, tt.want)
		}
	}
}

func TestMaskIsAllowed(t *testing.T) {
	if !config.MaskIsAllowed(0, 0xAA) {
		t.Fatal("expected mask 0xAA to be allowed for mode 0")
	}
	if config.MaskIsAllowed(0, 0xFF) {
		t.Fatal("expected mask 0xFF to be disallowed for mode 0")
	}
}
