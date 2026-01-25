// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Normalize normalizes maps in manifest to make sure all maps are initiaded
func (m *Manifest) Normalize() {
	if m.InstalledPackages == nil {
		m.InstalledPackages = make(map[packages.PackageType]map[string]InstalledPackage)
	}

	for t, pkgs := range m.InstalledPackages {
		if t == "" {
			delete(m.InstalledPackages, t)
			continue
		}
		for name := range pkgs {
			if name == "" {
				delete(pkgs, name)
			}
		}
	}
}

// BuildInitialManifest builds the initial manifest
func BuildInitialManifest(programVer string, path *paths.Paths) (*Manifest, error) {
	fmt.Printf(" * Building Initial Manifest...\n")

	m := Manifest{
		SchemaVersion:    1,
		ProgramVersion:   programVer,
		InstalledLoaders: make(map[string]LoaderInfo),

		EnabledPackages: make(map[packages.PackageType]string),

		InstalledPackages: make(map[packages.PackageType]map[string]InstalledPackage),

		Paths: path,
	}

	if err := m.Save(); err != nil {
		return nil, err
	}

	return &m, nil
}
