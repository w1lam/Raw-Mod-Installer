// Package paths defines file paths and constants used in the Raw Mod Installer application.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	MinecraftDir       string
	ModsDir            string
	ManifestPath       string
	BackupDir          string
	TempDownloadDir    string
	RawModInstallerDir string
	UnloadedModsDir    string
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

	installerDir := filepath.Join(mcDir, ".raw-mod-installer")

	return &Paths{
		MinecraftDir:       mcDir,
		ModsDir:            filepath.Join(mcDir, "mods"),
		ManifestPath:       filepath.Join(installerDir, "manifest.json"),
		BackupDir:          filepath.Join(installerDir, "mods.backup"),
		TempDownloadDir:    filepath.Join(installerDir, "downloaded-mods"),
		RawModInstallerDir: installerDir,
		UnloadedModsDir:    filepath.Join(installerDir, "mods"),
	}, nil
}

//var (
//	userProfile, _      = os.UserHomeDir()
//	ModFolderPath       = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods")
//	ModBackupPath       = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods_old")
//	VerFilePath         = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods", "ver.txt")
//	TempModDownloadPath = filepath.Join(os.TempDir(), "temp-mod-downloads")
//)
