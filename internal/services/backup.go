package services

import (
	"path/filepath"
	"time"

	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
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

type BackupResult struct {
	Entries []manifest.BackupEntry
}

// BackupIfFirstRun backs up the content of all active package dirs
func PerformInitialBackup(path *paths.Paths) (*BackupResult, error) {
	var result BackupResult

	for _, bt := range []struct {
		Type packages.PackageType
		Dir  string
	}{
		{packages.PackageModPack, path.ModsDir},
		{packages.PackageResourceBundle, path.ResourcePacksDir},
		{packages.PackageShaderBundle, path.ShaderPacksDir},
	} {
		if utils.DirEmpty(bt.Dir) {
			continue
		}

		id := utils.RandomID()
		dst := filepath.Join(path.BackupsDir, id)

		if err := utils.CopyDir(bt.Dir, dst); err != nil {
			return nil, err
		}

		result.Entries = append(result.Entries, manifest.BackupEntry{
			Time: time.Now(),
			Type: bt.Type,
			Path: dst,
			ID:   id,
		})
	}

	return &result, nil
}

func ApplyInitialBackup(m *manifest.Manifest, res *BackupResult, path *paths.Paths) error {
	if m.Initialized {
		return nil
	}

	m.Backups = append(m.Backups, res.Entries...)
	m.Initialized = true

	return m.Save()
}
