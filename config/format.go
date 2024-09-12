package config

import (
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

func (cc ContainerConfig) Export(c ContainerFormat) bool {
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
	Path string
	Type PaletteType
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
	ScreenOldFormat  ScreenFormat = "screen"
	WindowFormat     ScreenFormat = "window"
	Egx1Format       ScreenFormat = "egx1"
	Egx2Format       ScreenFormat = "egx2"
	ImpdrawTile      ScreenFormat = "tile"
)

type ScreenExport string

var (
	Overscan          ScreenExport = "overscan" // fichier fullscreen .scr
	ScreenOldExport   ScreenExport = "screen"   // fichier classique .scr
	SpriteExport      ScreenExport = "sprite"
	SpriteHardExport  ScreenExport = "sprite_hard"
	GoImpdrawExport   ScreenExport = "go"   // deux fichiers .go1 et .go2
	AssemblyExport    ScreenExport = "asm"  // texte assembleur
	JsonExport        ScreenExport = "json" // json
	ImpdrawTileExport ScreenExport = "tile"
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
	return ScreenOldFormat == f
}

func (f ScreenFormat) IsWindow() bool {
	return WindowFormat == f
}
