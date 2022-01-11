package menu

type ImageExport struct {
	ExportDsk              bool
	ExportText             bool
	ExportWithAmsdosHeader bool
	ExportZigzag           bool
	ExportJson             bool
	ExportCompression      int
	ExportFolderPath       string
	M2IP                   string
	ExportToM2             bool
}
