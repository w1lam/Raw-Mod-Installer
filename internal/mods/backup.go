// Package mods handles backup and restoration of mod folders.
package mods

import (
	"fmt"
	"os"
	"time"

	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Backup creates a backup of the mod folder
func Backup() error {
	path, err := paths.Resolve()
	if err != nil {
		return err
	}

	if utils.CheckFileExists(path.ModsDir) {
		entries, err := os.ReadDir(path.ModsDir)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			return fmt.Errorf("mods folder is empty, refusing to backup")
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

// RestoreBackup restores the mod folder from backup
func RestoreBackup() error {
	path, err := paths.Resolve()
	if err != nil {
		return err
	}

	if !utils.CheckFileExists(path.BackupDir) {
		return fmt.Errorf("no backup folder found")
	}

	if utils.CheckFileExists(path.ModsDir) {
		return fmt.Errorf("mods folder already exists, refusing to overwrite")
	}
	return os.Rename(path.BackupDir, path.ModsDir)
}

// RemoveAll removes the mod folder
func RemoveAll() error {
	path, err0 := paths.Resolve()
	if err0 != nil {
		return err0
	}

	err := os.RemoveAll(path.ModsDir)
	if err != nil {
		return fmt.Errorf("failed to uninstall mods: %v", err)
	}
	return nil
}
