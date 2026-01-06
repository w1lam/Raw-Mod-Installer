package install

import (
	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
)

type SystemState struct {
	FabricStatus FabricStatus
	MCVersion    string
}

type FabricStatus int

const (
	FabricMissing FabricStatus = iota
	FabricOutdated
	FabricUpToDate
)

func DetectSystem(mcVersion string) (SystemState, error) {
	statusStr, err := fabric.CheckVersions(mcVersion)
	if err != nil {
		return SystemState{}, err
	}

	var status FabricStatus
	switch statusStr {
	case "notInstalled":
		status = FabricMissing
	case "updateFound":
		status = FabricOutdated
	default:
		status = FabricUpToDate
	}

	return SystemState{
		FabricStatus: status,
		MCVersion:    mcVersion,
	}, nil
}

func EnsureFabric(mcVersion string) error {
	state, err := DetectSystem(mcVersion)
	if err != nil {
		return err
	}

	if state.FabricStatus == FabricUpToDate {
		return nil
	}

	jar, err := fabric.GetLatestInstallerJar()
	if err != nil {
		return err
	}

	return fabric.RunInstaller(jar, mcVersion)
}

func DetectRemovedMods(
	m *manifest.Manifest,
	entries []modlist.ModEntry,
) []string {
	active := make(map[string]bool)
	for _, e := range entries {
		active[e.Slug] = true
	}

	var removed []string
	for slug := range m.Mods {
		if !active[slug] {
			removed = append(removed, slug)
		}
	}

	return removed
}
