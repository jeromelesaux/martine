package address

import "testing"

func TestCpcScreenAddress(t *testing.T) {
	tests := []struct {
		name         string
		initial      int
		x, y         int
		mode         uint8
		overscan     bool
		doubleScreen bool
		want         int
	}{
		{"Mode0NoOverscan", 0, 1, 0, 0, false, false, 1},
		{"Mode1NoOverscan", 0, 3, 0, 1, false, false, 1},
		{"Mode2NoOverscan", 0, 7, 0, 2, false, false, 1},
		{"Mode0Overscan", 0, 0, 8, 0, true, false, 0x60},
		{"Mode0OverscanWithBase", 0xC000, 0, 8, 0, true, false, 0xC060},
		{"Mode0OverscanDoubleScreen", 0xC000, 0, 168, 0, true, true, 0xFFE0},
	}

	for _, tt := range tests {
		got := CpcScreenAddress(tt.initial, tt.x, tt.y, tt.mode, tt.overscan, tt.doubleScreen)
		if got != tt.want {
			t.Fatalf("%s: got %#x, want %#x", tt.name, got, tt.want)
		}
	}
}

func TestCpcScreenAddressOffset(t *testing.T) {
	tests := []struct {
		line int
		want int
	}{
		{0, 0},
		{1, 2048},
		{8, 80},
		{23, 14496},
	}

	for _, tt := range tests {
		got := CpcScreenAddressOffset(tt.line)
		if got != tt.want {
			t.Fatalf("line %d: got %d, want %d", tt.line, got, tt.want)
		}
	}
}
