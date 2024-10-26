package preferences

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/martine/common"
)

const (
	prefFilename = ".martine.cfg"
)

var (
	userPreferenceCache UserPreferences
)

type UserPreferences struct {
	OpenDirectory string `json:"open_directory"`
	SaveDirectory string `json:"save_directory"`
	Version       string `json:"version"`
}

func SaveDirectoryPref(path string) error {
	userPreferenceCache.SaveDirectory = path
	return Save()
}
func OpenDirectoryPref(path string) error {
	userPreferenceCache.OpenDirectory = path
	return Save()
}

func Save() error {
	userPreferenceCache.Version = common.AppVersion
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(home, prefFilename))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&userPreferenceCache); err != nil {
		return err
	}
	return nil
}

func OpenPref() (UserPreferences, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return UserPreferences{}, err
	}
	prefPath := filepath.Join(home, prefFilename)
	if _, err := os.Stat(prefPath); errors.Is(err, os.ErrNotExist) {
		return UserPreferences{
			OpenDirectory: home,
			SaveDirectory: home,
		}, nil
	}

	f, err := os.Open(prefPath)
	if err != nil {
		return UserPreferences{}, err
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&userPreferenceCache); err != nil {
		return UserPreferences{}, err
	}
	return userPreferenceCache, nil
}
