package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReadPackageIDFile reads a package id file
func ReadPackageIDFile(path string) (InstalledPackage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return InstalledPackage{}, fmt.Errorf("failed to read package id file: %s", path)
	}

	var out InstalledPackage
	if err := json.Unmarshal(data, &out); err != nil {
		return InstalledPackage{}, fmt.Errorf("failed to unmarshal package id file: %s", path)
	}

	return out, nil
}

// WritePackageIDFile writes ResolbedPackage data to json file
func WritePackageIDFile(pkg InstalledPackage, path string) error {
	outFile := filepath.Join(path, pkg.Name+".id.json")

	marshaled, err := json.MarshalIndent(pkg, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal pkg json: %w", err)
	}

	return os.WriteFile(outFile, marshaled, 0o644)
}
