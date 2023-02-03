package directory

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

var savedDirectoryPath fyne.URI

func DefaultDirectoryURI() (fyne.ListableURI, error) {
	if savedDirectoryPath == nil {
		return nil, fmt.Errorf("empty saved directory path")
	}
	return storage.ListerForURI(savedDirectoryPath)
}

func SetDefaultDirectoryURI(path fyne.URI) {
	p, err := storage.Parent(path)
	if err == nil {
		savedDirectoryPath = p
		return
	}
	savedDirectoryPath = path
}
