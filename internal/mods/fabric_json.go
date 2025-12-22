package mods

import (
	"archive/zip"
	"encoding/json"
	"fmt"
)

// FabricModJSON is a mods json
type FabricModJSON struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// ReadFabricModJSON reads and parses the fabric.mod.json file from a given .jar file.
func ReadFabricModJSON(jarPath string) (*FabricModJSON, error) {
	zr, err := zip.OpenReader(jarPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == "fabric.mod.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			var mod FabricModJSON
			if err := json.NewDecoder(rc).Decode(&mod); err != nil {
				return nil, err
			}

			return &mod, nil
		}
	}

	return nil, fmt.Errorf("fabric.mod.json not found in %s", jarPath)
}
