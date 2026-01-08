package install

import (
	"os"

	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

func prepareFS(m *manifest.Manifest, plan InstallPlan) error {
	switch plan.BackupPolicy {
	case filesystem.BackupIfExists:
		return filesystem.BackupIfNeeded(m)

	case filesystem.BackupOnce:
		if !utils.CheckFileExists(m.Paths.BackupsDir) {
			return filesystem.BackupIfNeeded(m)
		}
	}
	if plan.Intent == IntentReinstall {
		return os.RemoveAll(m.Paths.ModsDir)
	}

	return nil
}
