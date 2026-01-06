package ui

import (
	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// MainMenuState is the Main Menus state struct
type MainMenuState struct {
	HasManifest      bool
	HasModsInstalled bool

	ModlistUpdateAvailable bool
	FabricUpdateAvailable  bool
	ModUpdates             []resolve.ResolvedMod
}

// ComputeMainMenuState computes the main menus state
func ComputeMainMenuState(m *manifest.Manifest) (*MainMenuState, error) {
	state := &MainMenuState{
		HasManifest: m != nil,
	}

	if m == nil {
		return state, nil
	}

	state.HasModsInstalled = m.ModList.InstalledVersion != ""

	if remote, err := modlist.GetRemoteVersion(netcfg.ModListURL); err == nil {
		state.ModlistUpdateAvailable = remote != m.ModList.InstalledVersion
	}

	if latest, err := fabric.GetLatestLocalVersion(m.Minecraft.Version); err == nil {
		state.FabricUpdateAvailable = latest != m.Minecraft.LoaderVersion
	}

	for _, mod := range m.Mods {
		latest, err := resolve.ResolveMod(mod.Slug, m.Minecraft.Version, m.Minecraft.Loader)
		if err == nil && latest.LatestVer != mod.InstalledVersion {
			state.ModUpdates = append(state.ModUpdates, latest)
		}
	}

	return state, nil
}
