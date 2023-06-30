package directory

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

var savedDirectoryPath fyne.URI
var openDirectoryPath fyne.URI

func ExportDirectoryURI() (fyne.ListableURI, error) {
	if savedDirectoryPath == nil {
		return nil, fmt.Errorf("empty saved directory path")
	}
	return storage.ListerForURI(savedDirectoryPath)
}

func SetExportDirectoryURI(path fyne.URI) {
	p, err := storage.Parent(path)
	if err == nil {
		savedDirectoryPath = p
		return
	}
	savedDirectoryPath = path
}

func ImportDirectoryURI() (fyne.ListableURI, error) {
	if openDirectoryPath == nil {
		return nil, fmt.Errorf("empty saved directory path")
	}
	return storage.ListerForURI(openDirectoryPath)
}

func SetImportDirectoryURI(path fyne.URI) {
	p, err := storage.Parent(path)
	if err == nil {
		openDirectoryPath = p
		return
	}
	openDirectoryPath = path
}
