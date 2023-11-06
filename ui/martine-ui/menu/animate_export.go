package menu

import "github.com/jeromelesaux/martine/export/compression"

type AnimateExport struct {
	ExportDsk              bool
	ExportText             bool
	ExportWithAmsdosHeader bool
	ExportZigzag           bool
	ExportJson             bool
	ExportCompression      compression.CompressionMethod
	ExportFolderPath       string
	ExportImpdraw          bool
}
