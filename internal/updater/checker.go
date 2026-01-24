package updater

// Package update handels updating packages

import (
	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	pkg "github.com/w1lam/Raw-Mod-Installer/internal/packages/fetch"
)

func UpdateChecker(m *manifest.Manifest) (manifest.Updates, error) {
	if m == nil {
		panic("manifest is nil")
	}
	if m.InstalledPackages == nil {
		return manifest.Updates{}, nil
	}

	updates := manifest.Updates{
		ModListUpdate: make(map[string]bool),
		ModUpdates:    make(map[string][]modrinth.UpdateEntry),
	}

	allPackages, err := pkg.GetAllAvailablePackages()
	if err != nil {
		return manifest.Updates{}, err
	}

	modPacks := allPackages["modpacks"]

	for name, pack := range modPacks {
		if _, ok := m.InstalledPackages["modpacks"][name]; ok && m.InstalledPackages["modpacks"][name].InstalledVersion != pack.ListVersion {
			updates.ModListUpdate[name] = true
		}
	}

	for name, mp := range m.InstalledPackages["modpacks"] {
		mpHashes := mp.GetHashes()
		updates.ModUpdates[name], err = modrinth.BatchFetchUpdatesFromHash(mpHashes, mp.InstalledVersion, mp.Loader)
		if err != nil {
			return manifest.Updates{}, err
		}
	}

	return updates, nil
}
