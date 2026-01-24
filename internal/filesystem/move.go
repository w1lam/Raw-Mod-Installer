package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func EnablePackageFS(pkg packages.Pkg, path *paths.Paths) error {
	behavior := packages.PackageBehaviors[pkg.Type]

	src := filepath.Join(behavior.StorageDir(path), pkg.Name)
	dst := behavior.ActiveDir(path)

	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("failed to clear mods dir: %w", err)
	}

	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("failed to enable package \"%s\". pkgType %s : %w", pkg.Name, pkg.Type, err)
	}

	return nil
}

func DisablePackageFS(pkg packages.Pkg, path *paths.Paths) error {
	behavior := packages.PackageBehaviors[pkg.Type]

	src := behavior.ActiveDir(path)
	dst := filepath.Join(behavior.StorageDir(path), pkg.Name)

	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("failed to clear target modpack dir: %w", err)
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}
