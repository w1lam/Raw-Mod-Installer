package manifest

import (
	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Exists checks if manifest exists
func Exists() bool {
	path, err := paths.Resolve()
	if err != nil {
		return false
	}
	return utils.CheckFileExists(path.ManifestPath)
}

// GetHashes gets hashes of entries in an installed package
func (p *InstalledPackage) GetHashes() []string {
	var hashes []string
	for _, mod := range p.Entries {
		hashes = append(hashes, mod.Sha512)
	}
	return hashes
}

// AllInstalledModSlugs gets all installed mods slugs
func (m *Manifest) AllInstalledEntries() []string {
	seen := map[string]bool{}
	var ids []string

	for _, ap := range m.InstalledPackages {
		for _, p := range ap {
			for _, pe := range p.Entries {
				if !seen[pe.ID] {
					seen[pe.ID] = true
					ids = append(ids, pe.ID)
				}
			}
		}
	}
	return ids
}
