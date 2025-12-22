// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/modinfo"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/mods"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

var GlobalManifest *Manifest

// TYPES

func normalizeID(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}

func MigrateToManifest(
	programVersion string,
	modListURL string,
	modListVersion string,
	mcVersion string,
	loader string,
	loaderVersion string,

	modEntries []modlist.ModEntry,
	modInfos modinfo.ModInfoList,
	resolvedMods resolve.ResolvedModList,
	localMods []mods.LocalMod,
) (*Manifest, error) {
	manifest := &Manifest{
		SchemaVersion:  1,
		ProgramVersion: programVersion,
		Minecraft: MinecraftInfo{
			Version:       mcVersion,
			Loader:        loader,
			LoaderVersion: loaderVersion,
		},
		ModList: ModListInfo{
			Source:  modListURL,
			Version: modListVersion,
		},
		Mods: make(map[string]ManifestMod),
	}

	// Build local version lookup (FAST + SAFE)
	localVersionMap := make(map[string]string)
	for _, lm := range localMods {
		localVersionMap[normalizeID(lm.ID)] = lm.Version
	}

	// Build resolved lookup by slug
	resolvedMap := make(map[string]resolve.ResolvedMod)
	for _, r := range resolvedMods {
		resolvedMap[r.Slug] = r
	}

	// Merge EVERYTHING into manifest
	for i, info := range modInfos {
		slug := modEntries[i].Slug
		resolved := resolvedMap[slug]

		localVer := ""
		if resolved.FabricID != "" {
			if v, ok := localVersionMap[normalizeID(resolved.FabricID)]; ok {
				localVer = v
			}
		}

		manifest.Mods[slug] = ManifestMod{
			Slug:        slug,
			Title:       info.Title,
			Categories:  info.Category,
			Description: info.Description,
			Source:      info.Source,
			Wiki:        info.Wiki,
			LatestVer:   resolved.LatestVer,
			LocalVer:    localVer,
		}
	}

	return manifest, nil
}
