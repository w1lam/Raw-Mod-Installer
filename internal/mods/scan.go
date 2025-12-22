package mods

import (
	"os"
	"path/filepath"
)

// LocalMod is a local mod with file, ID and version
type LocalMod struct {
	File    string
	ID      string
	Version string
}

// GetLocalMods scans the specified directory for .jar files and extracts Fabric mod information.
func GetLocalMods(dirPath string) ([]LocalMod, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var mods []LocalMod

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".jar" {
			continue
		}

		jarPath := filepath.Join(dirPath, entry.Name())

		modJSON, err := ReadFabricModJSON(jarPath)
		if err != nil {
			continue
		}

		mods = append(mods, LocalMod{
			File:    entry.Name(),
			ID:      modJSON.ID,
			Version: modJSON.Version,
		})
	}

	return mods, nil
}
