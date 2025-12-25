package manifest

import (
	"encoding/json"
	"os"
)

// Save saves the manifest to the specified path atomically.
func Save(path string, m *Manifest) error {
	tmp := path + ".tmp"

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
