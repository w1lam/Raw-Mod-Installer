package services

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// DisablePackage disables the currently enabled modpack
func DisablePackage(pkg packages.Pkg) error {
	gState := state.Get()

	var path *paths.Paths
	gState.Read(func(s *state.State) {
		if s.Manifest().EnabledPackages[pkg.Type] == "" {
			return
		}
		path = s.Manifest().Paths
	})

	err := filesystem.DisablePackageFS(pkg, path)
	if err != nil {
		return fmt.Errorf("failed to move pacakge: %w", err)
	}

	return gState.Write(func(s *state.State) error {
		s.Manifest().EnabledPackages[pkg.Type] = ""
		return s.Manifest().Save()
	})
}
