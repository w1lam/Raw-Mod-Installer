package core

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
)

func ModListNeedsUpdate(m *manifest.Manifest, modListURL string) (bool, error) {
	remoteVersion, err := modlist.GetRemoteVersion(modListURL)
	if err != nil {
		return false, err
	}

	if m.ModList.InstalledVersion == "" {
		return true, nil
	}

	return m.ModList.InstalledVersion != remoteVersion, nil
}
