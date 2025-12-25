// Package manifest contains types and functions for mod list manifest files.
package manifest

import (
	"fmt"
	"strings"

	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// TYPES

func normalizeID(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}

// NewManifest creates a new Manifest instance from the provided parameters.
func NewManifest(
	programVersion string,
	modListURL string,
	modListInstalledVersion string,
	mcVersion string,
	loader string,
	loaderVersion string,

	modEntries []modlist.ModEntry,
	modInfos []modrinth.ModInfo,
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
			Source:           modListURL,
			InstalledVersion: modListInstalledVersion,
		},

		Mods: make(map[string]ManifestMod),
	}

	// ModInfo list is assumed to be in the same order as modEntries
	// (this is how your existing pipeline works)
	for i, entry := range modEntries {
		if i >= len(modInfos) {
			break
		}

		info := modInfos[i]

		manifest.Mods[entry.Slug] = ManifestMod{
			Slug:        entry.Slug,
			Title:       info.Title,
			Categories:  info.Category,
			Description: info.Description,
			Source:      info.Source,
			Wiki:        info.Wiki,

			InstalledVersion: "", // populated ONLY after install
		}
	}

	return manifest, nil
}

func BuildManifest(programVersion string) (*Manifest, error) {
	path, err := paths.Resolve()
	if err != nil {
		return nil, err
	}

	// 1. Fetch mod entries (slug + loader)
	modEntries, err := modlist.GetModEntryList(netcfg.ModListURL)
	if err != nil {
		return nil, err
	}

	// 2. Fetch Modrinth metadata (NOT versions)
	fmt.Println("Fetching mod metadata from Modrinth...")
	modInfos, err := resolve.ResolveModInfoList(modEntries, 10)
	if err != nil {
		return nil, err
	}

	// 3. Resolve loader version (Fabric example)
	loaderVersion, err := fabric.GetLatestLocalVersion(config.McVersion)
	if err != nil {
		return nil, err
	}

	// 4. Build manifest
	manifest, err := NewManifest(
		programVersion,
		netcfg.ModListURL,
		"", // mod list not installed yet
		config.McVersion,
		"fabric",
		loaderVersion,
		modEntries,
		modInfos,
	)
	if err != nil {
		return nil, err
	}

	// 5. Save immediately
	if err := Save(path.ManifestPath, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}
