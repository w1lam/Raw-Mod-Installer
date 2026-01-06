package manifest

import (
	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// NeedsModListUpdate checks if the mod list version in the manifest
func (m *Manifest) NeedsModListUpdate(remote string) bool {
	return m.ModList.InstalledVersion != remote
}

// SetModInstalled sets the local version of a mod in the manifest.
func (m *Manifest) SetModInstalled(slug, version string) {
	mod := m.Mods[slug]
	mod.InstalledVersion = version
	m.Mods[slug] = mod
}

// ModsSlice returns the mods in the manifest as a slice.
func (m *Manifest) ModsSlice() []ManifestMod {
	mods := make([]ManifestMod, 0, len(m.Mods))
	for _, mod := range m.Mods {
		mods = append(mods, mod)
	}
	return mods
}

func Exists() bool {
	path, err := paths.Resolve()
	if err != nil {
		return false
	}
	return utils.CheckFileExists(path.ManifestPath)
}

// InstalledVersion gets installed version of a mod from manifest
func (m *Manifest) InstalledVersion(slug string) string {
	if mod, ok := m.Mods[slug]; ok {
		return mod.InstalledVersion
	}
	return ""
}

func (m *Manifest) GetModEntries() []modlist.ModEntry {
	var mEntries []modlist.ModEntry

	for _, mod := range m.Mods {
		mEntries = append(mEntries, modlist.ModEntry{
			Loader: m.Minecraft.Loader,
			Slug:   mod.Slug,
		})
	}
	return mEntries
}
