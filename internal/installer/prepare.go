package installer

import (
	"fmt"
	"os"

	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

func rollback(modPack manifest.InstalledModPack, m *manifest.Manifest, plan InstallPlan, cause error) error {
	if modPack.Name == "NONE" {
		return cause
	}

	if plan.BackupPolicy != filesystem.BackupNever {
		_ = filesystem.RestoreModsBackup(modPack, m)
	}
	return cause
}

func prepareFS(m *manifest.Manifest, plan InstallPlan) error {
	if m == nil {
		panic("prepareFS: Manifest is nil")
	}
	if m.Paths.ModsDir == "" ||
		m.Paths.ModsBackupsDir == "" ||
		m.Paths.ModPacksDir == "" {
		return fmt.Errorf("manifest paths not initialized")
	}

	switch plan.BackupPolicy {
	case filesystem.BackupIfExists:
		return filesystem.BackupModsIfNeeded(m)

	case filesystem.BackupOnce:
		if !utils.CheckFileExists(m.Paths.ModsBackupsDir) {
			return filesystem.BackupModsIfNeeded(m)
		}
	}
	if plan.Intent == Reinstall {
		return os.RemoveAll(m.Paths.ModsDir)
	}

	return nil
}
