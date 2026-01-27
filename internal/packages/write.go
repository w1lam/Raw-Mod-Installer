package packages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WritePackageIDFile writes ResolbedPackage data to json file
func WritePackageIDFile(pkg ResolvedPackage, path string) error {
	outFile := filepath.Join(path, pkg.Name+".id.json")

	marshaled, err := json.MarshalIndent(pkg, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal pkg json: %w", err)
	}

	return os.WriteFile(outFile, marshaled, 0o644)
}
