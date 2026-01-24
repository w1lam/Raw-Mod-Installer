package services

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// EnablePackage enables the specified package
func EnablePackage(pkg packages.Pkg) error {
	gState := state.Get()

	var alreadyEnabled bool
	var path *paths.Paths
	gState.Read(func(s *state.State) {
		alreadyEnabled = s.Manifest().EnabledPackages[pkg.Type] == pkg.Name
		path = s.Manifest().Paths
	})

	if alreadyEnabled {
		return nil
	}

	err := filesystem.EnablePackageFS(pkg, path)
	if err != nil {
		return fmt.Errorf("failed to move package: %w", err)
	}

	return gState.Write(func(s *state.State) error {
		s.Manifest().EnabledPackages[pkg.Type] = pkg.Name
		return s.Manifest().Save()
	})
}
