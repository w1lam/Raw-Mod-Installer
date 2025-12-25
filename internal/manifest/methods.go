package manifest

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
