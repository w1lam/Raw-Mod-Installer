package services

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// EnablePackage enables the specified package
func EnablePackage(pkg packages.Pkg) error {
	gs := state.Get()

	var (
		installed      bool
		alreadyEnabled bool
		storageDir     string
		activeDir      string
	)

	gs.Read(func(s *state.State) {
		if p, ok := s.Manifest().InstalledPackages[pkg.Type][pkg.Name]; ok {
			installed = true
			storageDir = p.StorageDir
			activeDir = p.ActiveDir
		}

		alreadyEnabled = s.Manifest().EnabledPackages[pkg.Type] == pkg.Name
	})
	if storageDir == "" || activeDir == "" {
		return fmt.Errorf("invalid package paths for %s", pkg.Name)
	}

	if alreadyEnabled {
		return nil
	}
	if !installed {
		return fmt.Errorf("package not installed: %s", pkg.Name)
	}

	backupDir := storageDir + ".bak"
	err := filesystem.SwapDirs(storageDir, activeDir, backupDir)
	if err != nil {
		return fmt.Errorf("failed to move package: %w", err)
	}

	return gs.Write(func(s *state.State) error {
		s.Manifest().EnabledPackages[pkg.Type] = pkg.Name
		return s.Manifest().Save()
	})
}
