// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// BuildInitialManifest builds the initial manifest
func BuildInitialManifest(programVer string, path *paths.Paths) (*Manifest, error) {
	fmt.Printf(" * Building Initial Manifest...\n")

	m := Manifest{
		SchemaVersion:    1,
		ProgramVersion:   programVer,
		InstalledLoaders: make(map[string]LoaderInfo),
		Paths:            path,

		EnabledModPack: "",

		InstalledModPacks: make(map[string]InstalledModPack),
	}

	if err := m.Save(); err != nil {
		return nil, err
	}

	return &m, nil
}
