// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"fmt"

	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/programinfo"
)

// BuildInitialManifest builds the initial manifest
func BuildInitialManifest(path *paths.Paths) (*Manifest, error) {
	fmt.Printf(" * Building Initial Manifest...\n")

	availableModPacks, err := modpack.GetAvailableModPacks(netcfg.ModPacksURL)
	if err != nil {
		fmt.Printf(" * Failed to Fetch available Mod Packs: %s\n", err)
	}

	var loaderVer string
	var loader string

	if loaderVersion, err := fabric.GetLatestLocalVersion(programinfo.ProgramInfo.McVersion); err != nil {
		fmt.Printf(" * No local Mod Loader found\n")
		loaderVer = ""
		loader = "NONE"
	} else {
		loaderVer = loaderVersion
		loader = "fabric"
	}

	m := Manifest{
		SchemaVersion:  1,
		ProgramVersion: programinfo.ProgramInfo.Version,
		Minecraft: MinecraftInfo{
			Version:       "",
			Loader:        loader,
			LoaderVersion: loaderVer,
		},
		Paths: path,

		EnabledModPack: nil,

		AvailableModPacks: availableModPacks,
		InstalledModPacks: nil,
	}

	if err := m.Save(); err != nil {
		return nil, err
	}

	return &m, nil
}
