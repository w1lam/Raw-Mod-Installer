// Package filesystem handles backups and requirments for the file system
package filesystem

import (
	"fmt"
	"os"
	"time"

	"github.com/w1lam/Packages/utils"
)

// BackupDir moves a dir to backup folder
func BackupDir(src, dst string, rotate bool) error {
	if !utils.CheckFileExists(src) {
		return nil
	}

	if rotate && utils.CheckFileExists(dst) {
		ts := time.Now().Format("20060102150405")
		if err := os.Rename(dst, dst+"_"+ts); err != nil {
			return err
		}
	}

	return os.Rename(src, dst)
}

// ResoreBackupDir restores backup
func RestoreBackupDir(src, dst string) error {
	if !utils.CheckFileExists(src) {
		return fmt.Errorf("backup not found")
	}
	if utils.CheckFileExists(dst) {
		return fmt.Errorf("destination exists")
	}
	return os.Rename(src, dst)
}
