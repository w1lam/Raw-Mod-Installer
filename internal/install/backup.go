package install

import (
	"fmt"
	"os"
	"time"

	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

type BackupPolicy int

const (
	BackupNever BackupPolicy = iota
	BackupIfExists
	BackupOnce
)

func BackupIfNeeded(path *paths.Paths) error {
	if !utils.CheckFileExists(path.ModsDir) {
		return nil
	}

	if utils.CheckFileExists(path.BackupDir) {
		ts := time.Now().Format("20060102150405")
		if err := os.Rename(path.BackupDir, path.BackupDir+"_"+ts); err != nil {
			return err
		}
	}
	return os.Rename(path.ModsDir, path.BackupDir)
}

// Backup OLD creates a backup of the mod folder
func Backup(path *paths.Paths) error {
	if utils.CheckFileExists(path.ModsDir) {
		entries, err := os.ReadDir(path.ModsDir)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			return os.RemoveAll(path.ModsDir)
		}

		if utils.CheckFileExists(path.BackupDir) {
			err := os.Rename(path.BackupDir, path.BackupDir+"_"+time.Now().Format("20060102150405"))
			if err != nil {
				return err
			}
		}

		if err := os.Rename(path.ModsDir, path.BackupDir); err != nil {
			return err
		}
	}
	return nil
}

// RestoreBackup OLD restores the mod folder from backup
func RestoreBackup(path *paths.Paths) error {
	if !utils.CheckFileExists(path.BackupDir) {
		return fmt.Errorf("no backup folder found")
	}

	if utils.CheckFileExists(path.ModsDir) {
		return fmt.Errorf("mods folder already exists, refusing to overwrite")
	}
	return os.Rename(path.BackupDir, path.ModsDir)
}

func rollback(path *paths.Paths, plan InstallPlan, cause error) error {
	if plan.BackupPolicy != BackupNever {
		_ = RestoreBackup(path)
	}
	return cause
}
