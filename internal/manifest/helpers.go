package manifest

import (
	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// SetModInstalled sets the local version of a mod in the manifest.
func (mp *InstalledModPack) SetModInstalled(slug, version, modPackName string) {
	mod := mp.Mods[slug]
	mod.InstalledVersion = version
	mp.Mods[slug] = mod
}

// GetModList returns the mods from specified modpack as as ManifestMod slice
func (m *Manifest) GetModList(modPackName string) []ManifestMod {
	mods := make([]ManifestMod, 0, len(m.InstalledModPacks[modPackName].Mods))
	for _, mod := range m.InstalledModPacks[modPackName].Mods {
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

// GetInstalledVersionOfMod gets installed version of a mod from a mod pack
func (mp *InstalledModPack) GetInstalledVersionOfMod(slug string) string {
	if mod, ok := mp.Mods[slug]; ok {
		return mod.InstalledVersion
	}
	return ""
}

// GetSlugs gets mod slugs from an installed mod pack
func (mp *InstalledModPack) GetSlugs() []string {
	var slugs []string
	for _, mod := range mp.Mods {
		slugs = append(slugs, mod.Slug)
	}
	return slugs
}

func (mp *InstalledModPack) GetHashes() []string {
	var hashes []string
	for _, mod := range mp.Mods {
		hashes = append(hashes, mod.Sha512)
	}
	return hashes
}

// AllInstalledModSlugs gets all installed mods slugs
func (m *Manifest) AllInstalledModSlugs() []string {
	seen := map[string]bool{}
	var slugs []string

	for _, mp := range m.InstalledModPacks {
		for slug := range mp.Mods {
			if !seen[slug] {
				seen[slug] = true
				slugs = append(slugs, slug)
			}
		}
	}
	return slugs
}

// AllInstalledModHashes gets the hash for all installed mods
func (m *Manifest) AllInstalledModHashes() []string {
	seen := map[string]bool{}
	var hashes []string

	for _, mp := range m.InstalledModPacks {
		for slug, mod := range mp.Mods {
			if !seen[slug] {
				seen[slug] = true
				hashes = append(hashes, mod.Sha512)
			}
		}
	}
	return hashes
}
