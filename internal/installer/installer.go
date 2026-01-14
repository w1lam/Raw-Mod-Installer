// Package installer provides installer functions
package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/downloader"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

type InstallIntent int

const (
	Install InstallIntent = iota
	Reinstall
)

type InstallPlan struct {
	Intent           InstallIntent
	RequestedModPack modpack.ResolvedModPack
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

	if plan.EnsureFabric {
		fmt.Printf("\n\nENSURING FABRIC EXISTS\n\n")
		if err := filesystem.EnsureFabric(plan.RequestedModPack.McVersion); err != nil {
			return nil, fmt.Errorf("fabric install failed: %w", err)
		}
	}

	if err := prepareFS(m, plan); err != nil {
		return nil, err
	}

	fmt.Printf("\n\nRESOLVING MODS\n\n")
	resolved, err := resolve.ResolveMods(plan.RequestedModPack.Slugs, plan.RequestedModPack.McVersion, plan.RequestedModPack.Loader)
	if err != nil {
		return nil, rollback(m.InstalledModPacks[m.EnabledModPack], m, plan, err)
	}

	fmt.Printf("\n\nSTARTING DOWNLOAD\n\n")
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

	m.InstalledModPacks[plan.RequestedModPack.Name] = manifest.InstalledModPack{
		Name:             plan.RequestedModPack.Name,
		ListSource:       plan.RequestedModPack.ListSource,
		InstalledVersion: plan.RequestedModPack.ListVersion,
		McVersion:        plan.RequestedModPack.McVersion,
		Loader:           plan.RequestedModPack.Loader,
		Mods:             downloads.DownloadedMods,
	}

	if plan.EnableAfter {
		nm, err := EnableModPack(plan.RequestedModPack.Name, m)
		if err != nil {
			return nil, rollback(m.InstalledModPacks[m.EnabledModPack], m, plan, err)
		}

		m = nm
	}

	return m, m.Save()
}
