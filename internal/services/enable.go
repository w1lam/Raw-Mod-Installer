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
		storagePath    string
		activePath     string
	)

	gs.Read(func(s *state.State) {
		if p, ok := s.Manifest().InstalledPackages[pkg.Type][pkg.Name]; ok {
			installed = true
			storagePath = p.StoragePath
			activePath = p.ActivePath
		}

		alreadyEnabled = s.Manifest().EnabledPackages[pkg.Type] == pkg.Name
	})
	if storagePath == "" || activePath == "" {
		return fmt.Errorf("invalid package paths for %s", pkg.Name)
	}

	if alreadyEnabled {
		return nil
	}
	if !installed {
		return fmt.Errorf("package not installed: %s", pkg.Name)
	}

	backupPath := storagePath + ".bak"
	err := filesystem.SwapDirs(storagePath, activePath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to move package: %w", err)
	}

	return gs.Write(func(s *state.State) error {
		s.Manifest().EnabledPackages[pkg.Type] = pkg.Name
		return s.Manifest().Save()
	})
}
