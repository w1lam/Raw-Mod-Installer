package manifest

import (
	"fmt"

	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/modrinthsvc"
	"github.com/w1lam/Raw-Mod-Installer/internal/mods"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

func BuildManifestFromScratch(programInfo ProgramInfo) (*Manifest, error) {
	// 1. Mod list version
	modListVersion, err := modrinthsvc.GetRemoteVersion(netcfg.ModListURL)
	if err != nil {
		return nil, err
	}

	// 2. Mod entries
	modEntries, err := modrinthsvc.GetModEntryList(netcfg.ModListURL)
	if err != nil {
		return nil, err
	}

	// 3. Mod info
	fmt.Println("Fetching Mod List Info...")
	modInfoList, err := resolve.FetchModInfoList(modEntries, 10)
	if err != nil {
		return nil, err
	}

	// 4. Resolve mods
	fmt.Println("Resolving Mod List...")
	resolvedMods, err := resolve.FetchModListConcurrent(
		modEntries,
		config.McVersion,
		func(done, total int, mod string) {
			fmt.Printf("Fetched %d/%d -> %s\n", done, total, mod)
		},
	)
	if err != nil {
		return nil, err
	}

	// 5. Local mods
	var localMods []mods.LocalMod
	if IsModListPresent() {
		localMods, err = mods.GetLocalMods(paths.ModFolderPath)
		if err != nil {
			return nil, err
		}
	}

	fabricLoaderVersion, err := fabric.GetLatestLocalVersion(config.McVersion)
	if err != nil {
		return nil, err
	}

	// 6. Build manifest
	return MigrateToManifest(
		programInfo.ProgramVersion,
		netcfg.ModListURL,
		modListVersion,
		config.McVersion,
		"fabric",
		fabricLoaderVersion,
		modEntries,
		modInfoList,
		resolvedMods,
		localMods,
	)
}
