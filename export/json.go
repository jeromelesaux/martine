// Package export provides export utilities for martine output files.
package export

import (
	"encoding/json"
	"os"
)

// JSONSlice holds a collection of JSON sprites.
type JSONSlice struct {
	Sprites []*JSON `json:"sprites"`
}

type JSON struct {
	Label           string   `json:"label"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
	Screen          []string `json:"screen"`
	Palette         []string `json:"palette"`
	HardwarePalette []string `json:"hardwarepalette"`
}

func NewJSON(label string, width, height int, screen, palette, hardwarepalette []string) *JSON {
	return &JSON{
		Label:           label,
		Width:           width,
		Height:          height,
		Screen:          screen,
		Palette:         palette,
		HardwarePalette: hardwarepalette,
	}
}

func NewJSONSlice() *JSONSlice {
	return &JSONSlice{Sprites: make([]*JSON, 0)}
}

func (j *JSON) Save(file string) error {
	fw, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = fw.Close() }()
	return json.NewEncoder(fw).Encode(j)
}

func (j *JSONSlice) Save(file string) error {
	fw, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = fw.Close() }()
	return json.NewEncoder(fw).Encode(j)
}
