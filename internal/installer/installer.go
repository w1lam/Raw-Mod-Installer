// Package installer provides installer functions
package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/downloader"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/lists"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

type InstallIntent int

const (
	Install InstallIntent = iota
	Reinstall
)

type InstallPlan struct {
	Intent           InstallIntent
	RequestedModPack lists.ResolvedModPack
	EnsureFabric     bool
	BackupPolicy     filesystem.BackupPolicy
	EnableAfter      bool
}

// InstallModPack executes an InstallPlan
func InstallModPack(
	m *manifest.Manifest,
	plan InstallPlan,
) (*manifest.Manifest, error) {
	if m == nil {
		panic("InstallModPack: Manifest is nil")
	}

	if m.EnabledModPack == plan.RequestedModPack.Name {
		return nil, nil
	}

	if plan.EnsureFabric {
		if err := filesystem.EnsureFabric(plan.RequestedModPack.McVersion); err != nil {
			return nil, fmt.Errorf("fabric install failed: %w", err)
		}
	}

	if err := prepareFS(m, plan); err != nil {
		return nil, err
	}

	resolved, err := resolve.ResolveMods(plan.RequestedModPack.Entries, plan.RequestedModPack.McVersion, plan.RequestedModPack.Loader)
	if err != nil {
		return nil, rollback(m.InstalledModPacks[m.EnabledModPack], m, plan, err)
	}

	downloads, err := downloader.ModsDownloader(resolved, m)
	if err != nil {
		return nil, rollback(m.InstalledModPacks[m.EnabledModPack], m, plan, err)
	}

	destDir := filepath.Join(m.Paths.ModPacksDir, plan.RequestedModPack.Name)
	if err := os.RemoveAll(destDir); err != nil {
		return nil, fmt.Errorf("failed to clear target modpack dir: %w", err)
	}
	if err := os.Rename(downloads.TempDir, destDir); err != nil {
		return nil, fmt.Errorf("failed to move to target modpack dir: %w", err)
	}

	packHash, err := lists.ComputeDirHash(destDir)
	if err != nil {
		return nil, fmt.Errorf("failed to compute pack hash: %w", err)
	}

	m.InstalledModPacks[plan.RequestedModPack.Name] = manifest.InstalledModPack{
		Name:             plan.RequestedModPack.Name,
		ListSource:       plan.RequestedModPack.ListSource,
		InstalledVersion: plan.RequestedModPack.ListVersion,
		McVersion:        plan.RequestedModPack.McVersion,
		Loader:           plan.RequestedModPack.Loader,
		Hash:             packHash,
		Mods:             downloads.DownloadedMods,
	}

	if plan.EnableAfter {
		err := EnableModPack(plan.RequestedModPack.Name)
		if err != nil {
			return nil, rollback(m.InstalledModPacks[m.EnabledModPack], m, plan, err)
		}
	}

	return m, m.Save()
}
