package installer

import (
	"fmt"
	"path/filepath"

	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
)

func rollback(pkg manifest.InstalledPackage, path *paths.Paths, plan InstallPlan, cause error) error {
	if plan.BackupPolicy != services.BackupNever {
		err := services.RestorePackageBackup(packages.Pkg{
			Name: pkg.Name,
			Type: pkg.Type,
		})
		if err != nil {
			return fmt.Errorf("failed to restore rollback backup: %w", err)
		}
	}
	return cause
}

func prepareFS(path *paths.Paths, plan InstallPlan) error {
	if path.ModsDir == "" ||
		path.BackupsDir == "" ||
		path.PackagesDir == "" {
		return fmt.Errorf("manifest paths not initialized")
	}

	pkg := packages.Pkg{
		Name: plan.RequestedPackage.Name,
		Type: plan.RequestedPackage.Type,
	}

	switch plan.BackupPolicy {
	case services.BackupIfExists:
		return services.BackupPackage(pkg, plan.BackupPolicy)

	case services.BackupOnce:
		pkgPath := filepath.Join(path.PackagesDir, pkg.Name)
		if !utils.CheckFileExists(pkgPath) {
			return services.BackupPackage(pkg, plan.BackupPolicy)
		}
	}

	return nil
}
