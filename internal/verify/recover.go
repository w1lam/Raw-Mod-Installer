package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/errors"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func ScanDirForPackageID(dir string) (manifest.InstalledPackage, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return manifest.InstalledPackage{}, nil
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if strings.HasSuffix(e.Name(), ".id.json") {
			idPath := filepath.Join(dir, e.Name())
			return manifest.ReadPackageIDFile(idPath)
		}
	}

	return manifest.InstalledPackage{}, fmt.Errorf("no id file found in: %s", dir)
}

// ScanActitvePackages scans for all active packages
func ScanActitvePackages(path *paths.Paths) ([]manifest.InstalledPackage, error) {
	var found []manifest.InstalledPackage

	for _, dir := range []string{
		path.ModsDir,
		path.ResourcePacksDir,
		path.ShaderPacksDir,
	} {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			pkg, err := ScanDirForPackageID(filepath.Join(dir, e.Name()))
			if err != nil {
				errors.ReportCtx(
					"startup.scan.active",
					err,
					map[string]string{"dir": dir},
				)
				continue
			}

			found = append(found, pkg)
		}
	}

	return found, nil
}
