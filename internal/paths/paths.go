// Package paths defines file paths and constants used in the Raw Mod Installer application.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	MinecraftDir    string
	ModsDir         string
	ManifestPath    string
	BackupDir       string
	TempDownloadDir string
}

func DefaultMinecraftDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", ".minecraft"), nil
	case "linux":
		return filepath.Join(home, ".minecraft"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "minecraft"), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Resolve() (*Paths, error) {
	mcDir, err := DefaultMinecraftDir()
	if err != nil {
		return nil, err
	}

	return &Paths{
		MinecraftDir:    mcDir,
		ModsDir:         filepath.Join(mcDir, "mods"),
		ManifestPath:    filepath.Join(mcDir, "mods", "manifest.json"),
		BackupDir:       filepath.Join(mcDir, "mods.backup"),
		TempDownloadDir: filepath.Join(os.TempDir(), "mod-installer"),
	}, nil
}

//var (
//	userProfile, _      = os.UserHomeDir()
//	ModFolderPath       = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods")
//	ModBackupPath       = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods_old")
//	VerFilePath         = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods", "ver.txt")
//	TempModDownloadPath = filepath.Join(os.TempDir(), "temp-mod-downloads")
//)
