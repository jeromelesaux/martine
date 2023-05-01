package export

type ExportFormat string

var (
	SpriteFlatExport ExportFormat = "Flat"
	OcpWinExport     ExportFormat = "Files"
	SpriteImpCatcher ExportFormat = "Impcatcher"
	SpriteCompiled   ExportFormat = "Compiled Sprite"
	SpriteHard       ExportFormat = "Sprite Hard"
	OcpScreen        ExportFormat = "Screen"
	Overscan         ExportFormat = "Overscan"
)
