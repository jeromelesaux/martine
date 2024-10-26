package directory

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/jeromelesaux/martine/log"
	"github.com/jeromelesaux/martine/ui/martine-ui/preferences"
)

var savedDirectoryPath fyne.URI
var openDirectoryPath fyne.URI

// nolint: ireturn
func ExportDirectoryURI() (fyne.ListableURI, error) {
	if savedDirectoryPath == nil {
		prefs, err := preferences.OpenPref()
		if err != nil {
			return nil, err
		}
		dir, err := storage.ParseURI(prefs.SaveDirectory)
		if err != nil {
			return nil, err
		}
		savedDirectoryPath = dir
	}
	return storage.ListerForURI(savedDirectoryPath)
}

func SetExportDirectoryURI(path fyne.URI) {
	savedDirectoryPath = path
	if err := preferences.SaveDirectoryPref(savedDirectoryPath.String()); err != nil {
		log.GetLogger().Error("error while saving preferences %s", err.Error())
	}
}

// nolint: ireturn
func ImportDirectoryURI() (fyne.ListableURI, error) {
	if openDirectoryPath == nil {
		prefs, err := preferences.OpenPref()
		if err != nil {
			return nil, err
		}

		dir, err := storage.ParseURI(prefs.OpenDirectory)
		if err != nil {
			return nil, err
		}
		openDirectoryPath = dir
	}
	return storage.ListerForURI(openDirectoryPath)
}

func SetImportDirectoryURI(path fyne.URI) {
	p, err := storage.Parent(path)
	if err == nil {
		openDirectoryPath = p
		if err := preferences.OpenDirectoryPref(openDirectoryPath.String()); err != nil {
			log.GetLogger().Error("error while saving preferences %s", err.Error())
		}
		return
	}
	openDirectoryPath = path
	if err := preferences.OpenDirectoryPref(openDirectoryPath.String()); err != nil {
		log.GetLogger().Error("error while saving preferences %s", err.Error())
	}
}
