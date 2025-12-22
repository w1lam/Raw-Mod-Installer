// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"strings"

	"github.com/w1lam/Packages/pkg/modrinth"
)

var GlobalManifest *Manifest

// TYPES

type Manifest struct {
	SchemaVersion  int    `json:"schemaVersion"`
	ProgramVersion string `json:"programVersion"`

	Minecraft MinecraftInfo `json:"minecraft"`
	ModList   ModListInfo   `json:"modList"`

	Mods map[string]ManifestMod `json:"mods"`
}

type MinecraftInfo struct {
	Version       string `json:"version"`
	Loader        string `json:"loader"`
	LoaderVersion string `json:"loaderVersion"`
}

type ModListInfo struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

type ManifestMod struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`

	LatestVer string `json:"latestVersion,omitempty"`
	LocalVer  string `json:"localVersion,omitempty"`

	Source string `json:"source,omitempty"`
	Wiki   string `json:"wiki,omitempty"`
}

type ProgramInfo struct {
	ProgramVersion string
	ModListVersion string
}

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

	modEntries []modrinth.ModEntry,
	modInfos modrinth.ModInfoList,
	resolvedMods modrinth.ResolvedModList,
	localMods []modrinth.LocalMod,
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
	resolvedMap := make(map[string]modrinth.ResolvedMod)
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
