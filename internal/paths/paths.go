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
	ModListURL          = "https://raw.githubusercontent.com/w1lam/mods/refs/heads/main/mod-list.txt"
	TempModDownloadPath = filepath.Join(os.TempDir(), "temp-mod-downloads")
)
