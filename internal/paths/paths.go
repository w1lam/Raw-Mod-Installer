// Package paths defines file paths and constants used in the Raw Mod Installer application.
package paths

import (
	"os"
	"path/filepath"
)

var (
	UserProfile, _      = os.UserHomeDir()
	ModFolderPath       = filepath.Join(UserProfile, "AppData", "Roaming", ".minecraft", "mods")
	ModBackupPath       = filepath.Join(UserProfile, "AppData", "Roaming", ".minecraft", "mods_old")
	VerFilePath         = filepath.Join(UserProfile, "AppData", "Roaming", ".minecraft", "mods", "ver.txt")
	TempModDownloadPath = filepath.Join(os.TempDir(), "temp-mod-downloads")
)
