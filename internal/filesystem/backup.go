// Package filesystem handles backups and requirments for the file system
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

type BackupPolicy int

const (
	BackupNever BackupPolicy = iota
	BackupIfExists
	BackupOnce
)

func BackupIfNeeded(m *manifest.Manifest) error {
	if !utils.CheckFileExists(m.Paths.ModsDir) {
		return nil
	}

	backupDir := filepath.Join(m.Paths.BackupsDir, "mods.backup")
	if m.EnabledModPack != nil {
		backupDir = filepath.Join(m.Paths.BackupsDir, m.EnabledModPack.Name+".backup")
	}

	if utils.CheckFileExists(backupDir) {
		ts := time.Now().Format("20060102150405")
		if err := os.Rename(backupDir, backupDir+"_"+ts); err != nil {
			return err
		}
	}
	return os.Rename(m.Paths.ModsDir, backupDir)
}

// RestoreBackup restores the mod folder from backup
func RestoreBackup(modPack manifest.InstalledModPack, m *manifest.Manifest) error {
	if modPack.Name == "DEFAULT" {
		return os.Rename(filepath.Join(m.Paths.BackupsDir, "mods.backup"), m.Paths.ModsDir)
	}
	backupPth := m.Paths.BackupsDir + modPack.Name + ".backup"
	if !utils.CheckFileExists(backupPth) {
		return fmt.Errorf("no backup folder found")
	}

	if utils.CheckFileExists(m.Paths.ModsDir) {
		return fmt.Errorf("mods folder already exists, refusing to overwrite")
	}
	return os.Rename(backupPth, m.Paths.ModsDir)
}
