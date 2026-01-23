// Package filesystem handles backups and requirments for the file system
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

type BackupPolicy int

const (
	BackupNever BackupPolicy = iota
	BackupIfExists
	BackupOnce
)

func BackupModsIfNeeded() error {
	gState := state.Get()
	var path *paths.Paths
	var enabled string

	gState.Read(func(s *state.State) {
		path = s.Manifest().Paths
		enabled = s.Manifest().EnabledModPack
	})

	if !utils.CheckFileExists(path.ModsDir) {
		return nil
	}

	backupDir := filepath.Join(path.ModsBackupsDir, "mods.backup")
	if enabled != "" {
		backupDir = filepath.Join(path.ModsBackupsDir, enabled+".backup")
	}

	if utils.CheckFileExists(backupDir) {
		ts := time.Now().Format("20060102150405")
		if err := os.Rename(backupDir, backupDir+"_"+ts); err != nil {
			return err
		}
	}
	return os.Rename(path.ModsDir, backupDir)
}

// RestoreModsBackup restores the mod folder from backup
func RestoreModsBackup(modPack manifest.InstalledModPack, path *paths.Paths) error {
	if modPack.Name == "DEFAULT" {
		return os.Rename(filepath.Join(path.ModsBackupsDir, "mods.backup"), path.ModsDir)
	}
	backupPth := path.ModsBackupsDir + modPack.Name + ".backup"
	if !utils.CheckFileExists(backupPth) {
		return fmt.Errorf("no backup folder found")
	}

	if utils.CheckFileExists(path.ModsDir) {
		return fmt.Errorf("mods folder already exists, refusing to overwrite")
	}
	return os.Rename(backupPth, path.ModsDir)
}
