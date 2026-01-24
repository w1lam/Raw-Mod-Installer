// Package installer provides installer functions
package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/downloader"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

type InstallPlan struct {
	RequestedPackage packages.ResolvedPackage
	BackupPolicy     services.BackupPolicy
}

// PackageInstaller executes an InstallPlan and installs a pacakge
func PackageInstaller(
	plan InstallPlan,
) error {
	gState := state.Get()
	behavior := packages.PackageBehaviors[plan.RequestedPackage.Type]

	var path *paths.Paths
	var installed map[string]manifest.InstalledPackage
	var enabled string

	gState.Read(func(s *state.State) {
		path = s.Manifest().Paths
		installed = s.Manifest().InstalledPackages[plan.RequestedPackage.Type]
		enabled = s.Manifest().EnabledPackages[plan.RequestedPackage.Type]
	})

	if _, ok := installed[plan.RequestedPackage.Name]; ok {
		return fmt.Errorf("package already installed")
	}

	if behavior.EnsureLoader {
		if err := filesystem.EnsureFabric(plan.RequestedPackage.McVersion); err != nil {
			return fmt.Errorf("fabric install failed: %w", err)
		}
	}

	if err := prepareFS(path, plan); err != nil {
		return err
	}

	filter := modrinth.EntryFilter{
		McVersion:   plan.RequestedPackage.McVersion,
		ProjectType: string(plan.RequestedPackage.Type),
		Loader:      plan.RequestedPackage.Loader,
	}

	resolved, err := downloader.ResolveDownloadItem(plan.RequestedPackage.Entries, filter)
	if err != nil {
		return rollback(installed[enabled], path, plan, err)
	}

	downloads, err := downloader.DownloadEntries(resolved, path)
	if err != nil {
		return rollback(installed[enabled], path, plan, err)
	}

	destDir := filepath.Join(path.PackagesDir, string(plan.RequestedPackage.Type), plan.RequestedPackage.Name)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to clear target modpack dir: %w", err)
	}
	if err := os.Rename(downloads.TempDir, destDir); err != nil {
		return fmt.Errorf("failed to move to target modpack dir: %w", err)
	}

	packHash, err := filesystem.ComputeDirHash(destDir)
	if err != nil {
		return fmt.Errorf("failed to compute pack hash: %w", err)
	}

	downloadedEntries := make(map[string]manifest.PackageEntry)
	for n, item := range downloads.DownloadedItems {
		downloadedEntries[n] = manifest.PackageEntry{
			ID:               item.ID,
			FileName:         item.FileName,
			Sha512:           item.Sha512,
			Sha1:             item.Sha1,
			InstalledVersion: item.Version,
		}
	}

	installedMp := manifest.InstalledPackage{
		Name:             plan.RequestedPackage.Name,
		ListSource:       plan.RequestedPackage.ListSource,
		InstalledVersion: plan.RequestedPackage.ListVersion,
		McVersion:        plan.RequestedPackage.McVersion,
		Loader:           plan.RequestedPackage.Loader,
		Hash:             packHash,
		Entries:          downloadedEntries,
	}

	if packages.PackageBehaviors[plan.RequestedPackage.Type].EnableAfter {
		err := services.EnablePackage(packages.Pkg{Name: plan.RequestedPackage.Name, Type: plan.RequestedPackage.Type})
		if err != nil {
			return rollback(installed[enabled], path, plan, err)
		}
	}

	return gState.Write(func(s *state.State) error {
		s.Manifest().InstalledPackages[plan.RequestedPackage.Type][plan.RequestedPackage.Name] = installedMp
		return s.Manifest().Save()
	})
}
