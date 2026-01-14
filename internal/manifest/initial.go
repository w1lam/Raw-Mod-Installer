// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"fmt"

	"github.com/w1lam/Packages/fabric"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/programinfo"
)

// BuildInitialManifest builds the initial manifest
func BuildInitialManifest(path *paths.Paths) (*Manifest, error) {
	fmt.Printf(" * Building Initial Manifest...\n")

	var loaderVer string
	var loader string

	if loaderVersion, err := fabric.GetLatestLocalVersion(programinfo.ProgramInfo.McVersion); err != nil {
		fmt.Printf(" * No local Mod Loader found\n")
		loaderVer = "NONE"
		loader = "NONE"
	} else {
		loaderVer = loaderVersion
		loader = "fabric"
	}

	m := Manifest{
		SchemaVersion:  1,
		ProgramVersion: programinfo.ProgramInfo.Version,
		InstalledLoader: LoaderInfo{
			Loader:  loader,
			Version: loaderVer,
		},
		Paths: path,

		EnabledModPack: "",

		InstalledModPacks: make(map[string]InstalledModPack),
	}

	if err := m.Save(); err != nil {
		return nil, err
	}

	return &m, nil
}
