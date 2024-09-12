package config

import "github.com/jeromelesaux/martine/export/compression"

type ContainerFormat string

var (
	DskContainer         ContainerFormat = "dsk"
	ExtendedDskContainer ContainerFormat = "edsk"
	SnaContainer         ContainerFormat = "sna"
	M4Container          ContainerFormat = "m4"
)

type ContainerConfig struct {
	Type ContainerFormat
	Path string
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
}

var (
	FullscreenFormat ScreenFormat = "fullscreen"
	SpriteFormat     ScreenFormat = "sprite"
	SpriteHardFormat ScreenFormat = "sprite_hard"
	ScreenOldFormat  ScreenFormat = "screen"
	WindowFormat     ScreenFormat = "window"
	EgxFormat        ScreenFormat = "egx"
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
