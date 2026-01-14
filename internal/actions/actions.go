package actions

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/env"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

func InstallModPackAction(modPackName string, m *manifest.Manifest) error {
	fmt.Printf("\n\nSTARTING INSTALL ACTION\n\n")
	if m == nil {
		panic("InstallModPackAction: Manifest is nil")
	}

	fmt.Printf("\n\nGETTING AVAILABLE MODPACKS\n\n")
	availableModPacks, err := modpack.GetAvailableModPacks(netcfg.ModPacksURL)
	if err != nil {
		return fmt.Errorf("failed to get availalbe modpacks: %s", err)
	}

	if availableModPacks[modPackName].Name != modPackName {
		return fmt.Errorf("no mod pack with specified name: %s", modPackName)
	}
	fmt.Printf("MOD PACK: %s FOUND", availableModPacks[modPackName].Name)

	plan := installer.InstallPlan{
		Intent:           installer.Install,
		RequestedModPack: availableModPacks[modPackName],
		EnsureFabric:     true,
		BackupPolicy:     filesystem.BackupOnce,
		EnableAfter:      true,
	}

	fmt.Printf("\n\nSTARTING INSTALLATION\n\n")
	if nm, err := installer.InstallModPack(m, plan); err != nil {
		return err
	} else {
		m = nm
		env.GlobalManifest = m
		return m.Save()
	}
}
