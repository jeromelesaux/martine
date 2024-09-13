package config

import (
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/compression"
)

type ContainerFormat string

var (
	DskContainer         ContainerFormat = "dsk"
	ExtendedDskContainer ContainerFormat = "edsk"
	SnaContainer         ContainerFormat = "sna"
	M4Container          ContainerFormat = "m4"
)

type ContainerConfig struct {
	Type []ContainerFormat
	Path string
}

func (c *ContainerConfig) Reset() {
	c.Type = c.Type[:0]
}

func (cc *ContainerConfig) AddExport(c ContainerFormat) {
	cc.Type = append(cc.Type, c)
}

func (cc *ContainerConfig) RemoveExport(c ContainerFormat) {
	for i, v := range cc.Type {
		if v == c {
			cc.Type = append(cc.Type[:i], cc.Type[i+1:]...)
			return
		}
	}
}

func (cc ContainerConfig) HasExport(c ContainerFormat) bool {
	for _, v := range cc.Type {
		if v == c {
			return true
		}
	}
	return false
}

type M4Config struct {
	Host       string
	RemotePath string
	Autoexec   bool
	Enabled    bool
}

type PaletteType string

type PaletteConfig struct {
	Path    string
	Type    PaletteType
	Palette color.Palette
}

var (
	InkPalette PaletteType = "ink"
	PalPalette PaletteType = "pal"
	KitPalette PaletteType = "kit"
)

type ScreenFormat string

type ScreenConfig struct {
	InputPath      string
	OutputPath     string
	IsPlus         bool
	Type           ScreenFormat
	NoAmsdosHeader bool
	Compression    compression.CompressionMethod
	Export         []ScreenExport
	Size           constants.Size
	amsdosFilename [8]byte
	Process        ScreenProcessing
	Mode           uint8
}

func (s *ScreenConfig) ResetExport() {
	s.Export = s.Export[:0]
}

func (s *ScreenConfig) Reset() {
	s.ResetExport()
	s.IsPlus = false
	s.NoAmsdosHeader = false
}

func (s *ScreenConfig) AddExport(c ScreenExport) {
	s.Export = append(s.Export, c)
}

func (s *ScreenConfig) RemoveExport(c ScreenExport) {
	for i, v := range s.Export {
		if v == c {
			s.Export = append(s.Export[:i], s.Export[i+1:]...)
			return
		}
	}
}

func (s ScreenConfig) IsExport(c ScreenExport) bool {
	for _, v := range s.Export {
		if v == c {
			return true
		}
	}
	return false
}

var (
	FullscreenFormat ScreenFormat = "fullscreen"
	SpriteFormat     ScreenFormat = "sprite"
	SpriteHardFormat ScreenFormat = "sprite_hard"
	OcpScreenFormat  ScreenFormat = "screen"
	WindowFormat     ScreenFormat = "window"
	Egx1Format       ScreenFormat = "egx1"
	Egx2Format       ScreenFormat = "egx2"
	ImpdrawTile      ScreenFormat = "tile"
)

type ScreenExport string

var (
	Overscan               ScreenExport = "overscan" // fichier fullscreen .scr
	OcpScreenExport        ScreenExport = "screen"   // fichier classique .scr
	SpriteExport           ScreenExport = "sprite"
	SpriteHardExport       ScreenExport = "sprite_hard"
	SpriteCompiledExport   ScreenExport = "sprite_compiled"
	GoImpdrawExport        ScreenExport = "go_impdraw" // deux fichiers .go1 et .go2
	AssemblyExport         ScreenExport = "asm"        // texte assembleur
	JsonExport             ScreenExport = "json"       // json
	ImpdrawTileExport      ScreenExport = "tile"
	OcpWindowExport        ScreenExport = "window"
	SpriteImpCatcherExport ScreenExport = "impcatcher"
	SpriteFlatExport       ScreenExport = "flat_sprite"
)

func (f ScreenFormat) IsSprite() bool {
	return SpriteFormat == f
}

func (f ScreenFormat) IsFullScreen() bool {
	return FullscreenFormat == f
}

func (f ScreenFormat) IsSpriteHard() bool {
	return SpriteHardFormat == f
}

func (f ScreenFormat) IsScreen() bool {
	return OcpScreenFormat == f
}

func (f ScreenFormat) IsWindow() bool {
	return WindowFormat == f
}

type Rotation3d string

var (
	RotateXAxis            Rotation3d = "rotate_x_axis"
	RotateYAxis            Rotation3d = "rotate_y_axis"
	ReverseRotateXAxis     Rotation3d = "reverse_rotate_x_axis"
	RotateLeftToRightYAxis Rotation3d = "rotate_left_to_right_y_axis"
	RotateDiagonalXAxis    Rotation3d = "rotate_diagonal_x_axis"
	RotateDiagonalYAxis    Rotation3d = "rotate_diagonal_y_Axis"
)

type RotateConfig struct {
	RotationRraBit      int
	RotationRlaBit      int
	RotationSraBit      int
	RotationSlaBit      int
	RotationLosthighBit int
	RotationLostlowBit  int
	RotationKeephighBit int
	RotationKeeplowBit  int
	RotationIterations  int
	RotationMode        bool
	Rotation3DMode      bool
	Rotation3DX0        int
	Rotation3DY0        int
	Rotation3DType      Rotation3d
	RollMode            bool
	RollIteration       int
}

func Rotation3DType(v int) Rotation3d {
	switch v {
	case 1:
		return RotateXAxis
	case 2:
		return RotateYAxis
	case 3:
		return ReverseRotateXAxis
	case 4:
		return RotateLeftToRightYAxis
	case 5:
		return RotateDiagonalXAxis
	case 6:
		return RotateDiagonalYAxis
	}
	return RotateXAxis
}

type ScreenProcessing struct {
	Saturation                  float64
	Brightness                  float64
	KmeansThreshold             float64
	UseKmeans                   bool
	ResizingAlgo                imaging.ResampleFilter
	ApplyDithering              bool
	DitheringAlgo               int
	DitheringMatrix             [][]float32
	DitheringMultiplier         float64
	DitheringWithQuantification bool
	DitheringType               constants.DitheringType
	OneLine                     bool
	OneRow                      bool
	MaskSprite                  uint8
	MaskOrOperation             bool
	MaskAndOperation            bool
	Reducer                     int
}
