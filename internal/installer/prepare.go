package installer

import (
	"fmt"
	"os"

	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func rollback(modPack manifest.InstalledModPack, path *paths.Paths, plan InstallPlan, cause error) error {
	if plan.BackupPolicy != filesystem.BackupNever {
		_ = filesystem.RestoreModsBackup(modPack, path)
	}
	return cause
}

func prepareFS(path *paths.Paths, plan InstallPlan) error {
	if path.ModsDir == "" ||
		path.ModsBackupsDir == "" ||
		path.ModPacksDir == "" {
		return fmt.Errorf("manifest paths not initialized")
	}

	switch plan.BackupPolicy {
	case filesystem.BackupIfExists:
		return filesystem.BackupModsIfNeeded()

	case filesystem.BackupOnce:
		if !utils.CheckFileExists(path.ModsBackupsDir) {
			return filesystem.BackupModsIfNeeded()
		}
	}
	if plan.Intent == Reinstall {
		return os.RemoveAll(path.ModsDir)
	}

	return nil
}
