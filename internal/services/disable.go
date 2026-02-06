package services

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

func DisablePackage(pkg packages.Pkg) error {
	gs := state.Get()

	var (
		enabled     bool
		storagePath string
		activePath  string
	)

	gs.Read(func(s *state.State) {
		// Is this package currently enabled?
		if s.Manifest().EnabledPackages[pkg.Type] != pkg.Name {
			return
		}

		if p, ok := s.Manifest().InstalledPackages[pkg.Type][pkg.Name]; ok {
			enabled = true
			storagePath = p.StoragePath
			activePath = p.ActivePath
		}
	})

	if !enabled {
		return nil // already disabled → no-op
	}

	backupPath := activePath + ".bak"

	if err := filesystem.SwapDirs(activePath, storagePath, backupPath); err != nil {
		return fmt.Errorf("failed to disable package %s: %w", pkg.Name, err)
	}

	return gs.Write(func(s *state.State) error {
		s.Manifest().EnabledPackages[pkg.Type] = ""
		return s.Manifest().Save()
	})
}
