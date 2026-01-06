package install

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

type InstallIntent int

const (
	IntentInstall InstallIntent = iota
	IntentUpdate
	IntentReinstall
)

type InstallPlan struct {
	Intent       InstallIntent
	EnsureFabric bool
	BackupPolicy BackupPolicy
	EnableAfter  bool
}

func ExecutePlan(
	m *manifest.Manifest,
	path *paths.Paths,
	plan InstallPlan,
) (*manifest.Manifest, error) {
	if plan.EnsureFabric {
		if err := EnsureFabric(config.McVersion); err != nil {
			return nil, fmt.Errorf("fabric install failed: %w", err)
		}
	}

	if err := prepareFS(path, plan); err != nil {
		return nil, err
	}

	switch plan.Intent {
	case IntentInstall, IntentReinstall:
		var err error
		m, err = DownloadModpack(m, path)
		if err != nil {
			return nil, rollback(path, plan, err)
		}

	case IntentUpdate:
		var err error
		m, err = UpdateModpack(m, path)
		if err != nil {
			return nil, rollback(path, plan, err)
		}
	}

	if plan.EnableAfter {
		if err := EnableMods(path); err != nil {
			return nil, rollback(path, plan, err)
		}
	}

	if err := manifest.Save(path.ManifestPath, m); err != nil {
		return nil, err
	}

	return m, nil
}
