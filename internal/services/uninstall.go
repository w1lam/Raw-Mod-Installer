package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// UninstallPackage uninstalls a pacakge
func UninstallPackage(pkg packages.Pkg) error {
	st := state.Get()

	var path *paths.Paths
	var enabled string
	var ok bool

	st.Read(func(s *state.State) {
		if s.Manifest().InstalledPackages == nil {
			return
		}

		_, ok = s.Manifest().InstalledPackages[pkg.Type][pkg.Name]
		path = s.Manifest().Paths
		enabled = s.Manifest().EnabledPackages[pkg.Type]
	})

	if !ok {
		return fmt.Errorf("package not instaleld: %s", pkg.Name)
	}

	if pkg.Name == enabled {
		if err := DisablePackage(pkg); err != nil {
			return err
		}
	}

	pkgTypeDir := packages.PackageBehaviors[pkg.Type]
	target := filepath.Join(path.PackagesDir, pkgTypeDir.StorageDir(path), pkg.Name)
	err := os.RemoveAll(target)
	if err != nil {
		return err
	}

	return st.Write(func(s *state.State) error {
		delete(s.Manifest().InstalledPackages[pkg.Type], pkg.Name)
		return s.Manifest().Save()
	})
}
