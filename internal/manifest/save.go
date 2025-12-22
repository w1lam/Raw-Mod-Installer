package manifest

import (
	"encoding/json"
	"os"
)

func WriteManifest(path string, manifest *Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// SaveManifest saves the manifest to the specified path atomically.
func SaveManifest(path string, m *Manifest) error {
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
