package updater

import (
	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
)

func UpdateChecker(m *manifest.Manifest) (manifest.Updates, error) {
	if m == nil {
		panic("manifest is nil")
	}
	if m.InstalledModPacks == nil {
		return manifest.Updates{}, nil
	}

	updates := manifest.Updates{
		ModListUpdate: make(map[string]bool),
		ModUpdates:    make(map[string][]modrinth.UpdateEntry),
	}

	mp, err := modpack.GetAvailableModPacks()
	if err != nil {
		return manifest.Updates{}, err
	}

	for name, pack := range mp {
		if _, ok := m.InstalledModPacks[name]; ok && m.InstalledModPacks[name].InstalledVersion != pack.ListVersion {
			updates.ModListUpdate[name] = true
		}
	}

	for name, mp := range m.InstalledModPacks {
		mpHashes := mp.GetHashes()
		updates.ModUpdates[name], err = modrinth.BatchFetchUpdatesFromHash(mpHashes, mp.InstalledVersion, mp.Loader)
		if err != nil {
			return manifest.Updates{}, err
		}
	}

	return updates, nil
}
