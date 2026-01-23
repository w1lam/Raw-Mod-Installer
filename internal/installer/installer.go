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
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

type InstallIntent int

const (
	Install InstallIntent = iota
	Reinstall
)

type InstallPlan struct {
	Intent           InstallIntent
	RequestedPackage lists.ResolvedModPack
	EnsureFabric     bool
	BackupPolicy     filesystem.BackupPolicy
	EnableAfter      bool
}

// PackageInstaller executes an InstallPlan and installs a pacakge
func PackageInstaller(
	plan InstallPlan,
) error {
	gState := state.Get()

	var path *paths.Paths
	var installed map[string]manifest.InstalledModPack
	var enabled string

	gState.Read(func(s *state.State) {
		path = s.Manifest().Paths
		installed = s.Manifest().InstalledModPacks
		enabled = s.Manifest().EnabledModPack
	})

	if enabled == plan.RequestedPackage.Name {
		return nil
	}

	if plan.EnsureFabric {
		if err := filesystem.EnsureFabric(plan.RequestedPackage.McVersion); err != nil {
			return fmt.Errorf("fabric install failed: %w", err)
		}
	}

	if err := prepareFS(path, plan); err != nil {
		return err
	}

	resolved, err := resolve.ResolveMods(plan.RequestedPackage.Entries, plan.RequestedPackage.McVersion, plan.RequestedPackage.Loader)
	if err != nil {
		return rollback(installed[enabled], path, plan, err)
	}

	downloads, err := downloader.ModsDownloader(resolved, path)
	if err != nil {
		return rollback(installed[enabled], path, plan, err)
	}

	destDir := filepath.Join(path.ModPacksDir, plan.RequestedPackage.Name)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to clear target modpack dir: %w", err)
	}
	if err := os.Rename(downloads.TempDir, destDir); err != nil {
		return fmt.Errorf("failed to move to target modpack dir: %w", err)
	}

	packHash, err := lists.ComputeDirHash(destDir)
	if err != nil {
		return fmt.Errorf("failed to compute pack hash: %w", err)
	}

	installedMp := manifest.InstalledModPack{
		Name:             plan.RequestedPackage.Name,
		ListSource:       plan.RequestedPackage.ListSource,
		InstalledVersion: plan.RequestedPackage.ListVersion,
		McVersion:        plan.RequestedPackage.McVersion,
		Loader:           plan.RequestedPackage.Loader,
		Hash:             packHash,
		Mods:             downloads.DownloadedMods,
	}

	if plan.EnableAfter {
		err := EnableModPack(plan.RequestedPackage.Name)
		if err != nil {
			return rollback(installed[enabled], path, plan, err)
		}
	}

	return gState.Write(func(s *state.State) error {
		s.Manifest().InstalledModPacks[plan.RequestedPackage.Name] = installedMp
		return s.Manifest().Save()
	})
}
