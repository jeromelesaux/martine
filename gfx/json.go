package gfx

import (
	"encoding/json"
	"os"
)

type JsonSlice struct {
	Sprites []*Json `json:"sprites"`
}

type Json struct {
	Label   string   `json:"label"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
	Screen  []string `json:"screen"`
	Palette []string `json:"palette"`
}

func NewJson(label string, width int, height int, screen []string, palette []string) *Json {
	return &Json{
		Label:   label,
		Width:   width,
		Height:  height,
		Screen:  screen,
		Palette: palette,
	}
}

func NewJsonSlice() *JsonSlice {
	return &JsonSlice{Sprites:make([]*Json,0)}
}

func (j *Json) Save(file string) error {
	fw, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fw.Close()
	return json.NewEncoder(fw).Encode(j)
}

func (j *JsonSlice)Save(file string) error {
	fw, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fw.Close()
	return json.NewEncoder(fw).Encode(j)
}