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
		enabled    bool
		storageDir string
		activeDir  string
	)

	gs.Read(func(s *state.State) {
		// Is this package currently enabled?
		if s.Manifest().EnabledPackages[pkg.Type] != pkg.Name {
			return
		}

		if p, ok := s.Manifest().InstalledPackages[pkg.Type][pkg.Name]; ok {
			enabled = true
			storageDir = p.StorageDir
			activeDir = p.ActiveDir
		}
	})

	if !enabled {
		return nil // already disabled → no-op
	}

	backupDir := activeDir + ".bak"

	if err := filesystem.SwapDirs(activeDir, storageDir, backupDir); err != nil {
		return fmt.Errorf("failed to disable package %s: %w", pkg.Name, err)
	}

	return gs.Write(func(s *state.State) error {
		s.Manifest().EnabledPackages[pkg.Type] = ""
		return s.Manifest().Save()
	})
}
