package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
)

type Process struct {
	ByteStatement       string   `json:"byteStatement"`
	PicturePath         string   `json:"picturePath"`
	Width               int      `json:"width"`
	Height              int      `json:"height"`
	Mode                int      `json:"mode"`
	Output              string   `json:"outputPath"`
	Overscan            bool     `json:"isOverscan"`
	ResizeAlgorithm     int      `json:"resizeAlgorithm"`
	NoAmsdosHeader      bool     `json:"noAmsdosHeader"`
	PlusMode            bool     `json:"isPlusMode"`
	RollMode            bool     `json:"isRollMode"`
	Iterations          int      `json:"iterations"`
	Rra                 int      `json:"rra"`
	Rla                 int      `json:"rla"`
	Sra                 int      `json:"sra"`
	Sla                 int      `json:"sla"`
	Losthigh            int      `json:"lostHigh"`
	Lostlow             int      `json:"lostLow"`
	Keephigh            int      `json:"keepHigh"`
	Keeplow             int      `json:"keepLow"`
	PalettePath         string   `json:"palettePath"`
	Info                bool     `json:"info"`
	WinPath             string   `json:"winPath"`
	Dsk                 bool     `json:"generateDsk"`
	TileMode            bool     `json:"tileMode"`
	TileIterationX      int      `json:"tileIterationX"`
	TileIterationY      int      `json:"tileIterationY"`
	Compress            int      `json:"compress"`
	KitPath             string   `json:"kitPath"`
	InkPath             string   `json:"inkPath"`
	RotateMode          bool     `json:"isRotateMode"`
	M4Host              string   `json:"m4Host"`
	M4RemotePath        string   `json:"m4RemotePath"`
	M4Autoexec          bool     `json:"m4Autoexec"`
	Rotate3dMode        bool     `json:"isRotate3dMode"`
	Rotate3dType        int      `json:"rotate3dType"`
	Rotate3dX0          int      `json:"rotate3dX0"`
	Rotate3dY0          int      `json:"rotate3dY0"`
	Data                []int    `json:"data"`
	Palette             []int    `json:"palette"`
	Delta               bool     `json:"delta"`
	DitheringAlgo       int      `json:"ditheringAlgo"`
	DitheringMultiplier float64  `json:"ditheringMultiplier"`
	WithQuantization    bool     `json:"withQuantization"`
	ExtendedDsk         bool     `json:"extendedDsk"`
	Reverse             bool     `json:"reverse"`
	Flash               bool     `json:"flash"`
	PicturePath2        string   `json:"picturePath2"`
	Mode2               int      `json:"mode2"`
	PalettePath2        string   `json:"palettePath2"`
	Egx1                bool     `json:"egx1"`
	Egx2                bool     `json:"egx2"`
	DeltaFile           []string `json:"df"`
}

func NewProcess() *Process {
	return &Process{
		Width:               -1,
		Height:              -1,
		Mode:                -1,
		ResizeAlgorithm:     1,
		Iterations:          -1,
		Rra:                 -1,
		Rla:                 -1,
		Sra:                 -1,
		Sla:                 -1,
		Losthigh:            -1,
		Lostlow:             -1,
		Keephigh:            -1,
		Keeplow:             -1,
		TileIterationX:      -1,
		TileIterationY:      -1,
		Compress:            -1,
		Rotate3dType:        -1,
		Rotate3dX0:          -1,
		Rotate3dY0:          -1,
		DitheringAlgo:       -1,
		DitheringMultiplier: 1.18,
		WithQuantization:    false,
		ExtendedDsk:         false,
		Reverse:             false,
		Flash:               false,
		PicturePath2:        "",
		Mode2:               -1,
		PalettePath2:        "",
		Egx1:                false,
		Egx2:                false,
		Data:                make([]int, 0),
		Palette:             make([]int, 0),
		DeltaFile:           make([]string, 0),
	}
}

func InitProcess(filePath string) (*Process, error) {
	p := NewProcess()
	f, err := os.Create(filePath)
	if err != nil {
		return p, err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(p)
	return p, err
}

func LoadProcessFile(filePath string) (*Process, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p := NewProcess()
	err = json.NewDecoder(f).Decode(p)
	return p, err
}

// nolint: funlen
func (p *Process) Apply() {
	*byteStatement = p.ByteStatement
	*picturePath = p.PicturePath
	*width = p.Width
	*height = p.Height
	*mode = p.Mode
	*output = p.Output
	*overscan = p.Overscan
	*resizeAlgorithm = p.ResizeAlgorithm
	*noAmsdosHeader = p.NoAmsdosHeader
	*plusMode = p.PlusMode
	*rollMode = p.RollMode
	*iterations = p.Iterations
	*rra = p.Rra
	*rla = p.Rla
	*sra = p.Sra
	*sla = p.Sla
	*losthigh = p.Losthigh
	*lostlow = p.Lostlow
	*keephigh = p.Keephigh
	*keeplow = p.Keeplow
	*palettePath = p.PalettePath
	*info = p.Info
	*winPath = p.WinPath
	*dsk = p.Dsk
	*tileMode = p.TileMode
	*tileIterationX = p.TileIterationX
	*tileIterationY = p.TileIterationY
	*compress = p.Compress
	*kitPath = p.KitPath
	*inkPath = p.InkPath
	*rotateMode = p.RotateMode
	*m4Host = p.M4Host
	*m4RemotePath = p.M4RemotePath
	*m4Autoexec = p.M4Autoexec
	*rotate3dMode = p.Rotate3dMode
	*rotate3dType = p.Rotate3dType
	*rotate3dX0 = p.Rotate3dX0
	*rotate3dY0 = p.Rotate3dY0
	*deltaMode = p.Delta
	*ditheringAlgo = p.DitheringAlgo
	*ditheringMultiplier = p.DitheringMultiplier
	*withQuantization = p.WithQuantization
	*extendedDsk = p.ExtendedDsk
	*reverse = p.Reverse
	*flash = p.Flash
	*picturePath2 = p.PicturePath2
	*mode2 = p.Mode2
	*palettePath2 = p.PalettePath2
	*egx1 = p.Egx1
	*egx2 = p.Egx2
	for i := 0; i < len(p.DeltaFile); i++ {
		err := deltaFiles.Set(p.DeltaFile[i])
		if err != nil {
			continue
		}
	}
}

func (p *Process) GenerateRawFile() error {
	p.PicturePath = fmt.Sprintf("raw_%.4d.png", os.Getppid())
	in, err := gfx.TransformRawCpcData(p.Data, p.Palette, p.Width, p.Height, p.Mode, p.PlusMode)
	if err != nil {
		return err
	}
	*picturePath = p.PicturePath
	return png.Png(p.PicturePath, in)
}
