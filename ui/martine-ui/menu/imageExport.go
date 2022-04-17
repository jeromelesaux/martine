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

func (ie *ImageExport) Reset() {
	ie.ExportDsk = false
	ie.ExportText = false
	ie.ExportWithAmsdosHeader = false
	ie.ExportZigzag = false
	ie.ExportJson = false
	ie.ExportCompression = -1
	ie.ExportToM2 = false
}
