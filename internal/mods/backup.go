package mods

import (
	"fmt"
	"os"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Mod Backup and Restore Functions

func BackupModFolder() error {
	*Timestamp = time.Now().Format("20060102150405")

	err := os.Rename(paths.ModFolderPath, paths.ModBackupPath)
	if err != nil {
		err1 := os.Rename(paths.ModFolderPath, paths.ModBackupPath+"_"+*Timestamp)
		if err1 != nil {
			return err1
		}
	}

	if _, err := os.Stat(paths.ModFolderPath); err == nil {
		if _, err2 := os.Stat(paths.ModBackupPath); err2 == nil {
			timestamp := time.Now().Format("20060102150405")

			err3 := os.Rename(paths.ModBackupPath, paths.ModBackupPath+"_"+timestamp)
			if err3 != nil {
				return fmt.Errorf("failed to backup existing mod backup folder: %v", err3)
			}
		}

		err := os.Rename(paths.ModFolderPath, paths.ModBackupPath)
		if err != nil {
			timestamp := time.Now().Format("20060102150405")

			err2 := os.Rename(paths.ModFolderPath, paths.ModBackupPath+"_"+timestamp+"bruh")
			if err2 != nil {
				return fmt.Errorf("failed to backup existing mod backup folder: %v", err2)
			}
		}
	}

	return nil
}

func RestoreModBackup() error {
	err := os.Rename(paths.ModBackupPath, paths.ModFolderPath)
	if err != nil {

		err1 := os.Rename(paths.ModBackupPath+"_"+*Timestamp, paths.ModFolderPath)
		if err1 != nil {
			return fmt.Errorf("failed to restore backup mods: %v\n %v", err, err1)
		}
	}
	return nil
}

func UninstallMods() error {
	err := os.RemoveAll(paths.ModFolderPath)
	if err != nil {
		return fmt.Errorf("failed to uninstall mods: %v", err)
	}
	return nil
}
