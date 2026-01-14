package manifest

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Load loads the manifest
func Load(path *paths.Paths) (*Manifest, error) {
	data, err := os.ReadFile(path.ManifestPath)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	m.Paths = path

	return &m, nil
}

// Save saves the manifest to the specified path atomically.
func (m *Manifest) Save() error {
	tmp := m.Paths.ManifestPath + ".tmp"

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall manifest: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write manifest temp file: %s", err)
	}

	return os.Rename(tmp, m.Paths.ManifestPath)
}
