package meta

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// LoadMetaData loaeds metadata
func LoadMetaData(path *paths.Paths) *MetaData {
	data, err := os.ReadFile(path.MetaDataPath)
	if err != nil {
		return nil
	}

	var m MetaData
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}

	return &m
}

// Save metadata
func (md *MetaData) Save(path *paths.Paths) error {
	tmp := path.MetaDataPath + ".tmp"

	data, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall metadata: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write metadata temp file: %s", err)
	}

	return os.Rename(tmp, path.MetaDataPath)
}
