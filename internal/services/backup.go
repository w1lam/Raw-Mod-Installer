package services

import (
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

type BackupPolicy int

const (
	BackupNever BackupPolicy = iota
	BackupIfExists
	BackupOnce
)

func BackupPackage(pkg packages.Pkg, policy BackupPolicy) error {
	if policy == BackupNever {
		return nil
	}

	gState := state.Get()
	var path *paths.Paths

	gState.Read(func(s *state.State) {
		path = s.Manifest().Paths
	})

	dst := filepath.Join(path.BackupsDir, pkg.Name)
	rotate := policy == BackupOnce

	return filesystem.BackupDir(path.ModsDir, dst, rotate)
}

// RestorePackageBackup restores the package folder from backup
func RestorePackageBackup(pkg packages.Pkg) error {
	gState := state.Get()

	var path *paths.Paths
	gState.Read(func(s *state.State) {
		path = s.Manifest().Paths
	})
	src := filepath.Join(path.BackupsDir, pkg.Name)
	dst := filepath.Join(path.PackagesDir, string(pkg.Type), pkg.Name)

	return filesystem.RestoreBackupDir(src, dst)
}
