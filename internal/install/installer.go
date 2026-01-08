package install

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
)

type InstallIntent int

const (
	IntentInstall InstallIntent = iota
	IntentUpdate
	IntentReinstall
)

type InstallPlan struct {
	Intent           InstallIntent
	RequestedModPack modpack.ResolvedModPackList
	EnsureFabric     bool
	BackupPolicy     filesystem.BackupPolicy
	EnableAfter      bool
}

// ExecuteInstallerPlan executes an InstallPlan
func ExecuteInstallerPlan(
	m *manifest.Manifest,
	plan InstallPlan,
) (*manifest.Manifest, error) {
	if plan.EnsureFabric {
		if err := filesystem.EnsureFabric(config.McVersion); err != nil {
			return nil, fmt.Errorf("fabric install failed: %w", err)
		}
	}

	if err := prepareFS(m, plan); err != nil {
		return nil, err
	}

	switch plan.Intent {
	case IntentInstall, IntentReinstall:
		var err error
		m, err = DownloadModpack(plan.RequestedModPack, m)
		if err != nil {
			return nil, rollback(*m.EnabledModPack, m, plan, err)
		}

	case IntentUpdate:
		var err error
		m, err = UpdateModpack(*m.EnabledModPack, m)
		if err != nil {
			return nil, rollback(*m.EnabledModPack, m, plan, err)
		}
	}

	if plan.EnableAfter {
		if me, err := EnableModPack(m.InstalledModPacks[plan.RequestedModPack.Name], m); err != nil {
			return nil, rollback(*m.EnabledModPack, m, plan, err)
		} else {
			m = me
		}
	}

	if err := m.Save(); err != nil {
		return nil, err
	}

	return m, nil
}
