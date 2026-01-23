// Package paths defines file paths and constants used in the Raw Mod Installer application.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	MinecraftDir     string
	ModsDir          string
	ResourcePacksDir string

	ProgramFilesDir string
	DataDir         string
	ManifestPath    string
	MetaDataPath    string

	ModPacksDir            string
	ModsBackupsDir         string
	ResourceBundlesDir     string
	ResourcePackBackupsDir string
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
	dataDir := filepath.Join(installerDir, "data")

	return &Paths{
		MinecraftDir:     mcDir,
		ModsDir:          filepath.Join(mcDir, "mods"),
		ResourcePacksDir: filepath.Join(mcDir, "resourcepacks"),

		ProgramFilesDir: installerDir,
		DataDir:         dataDir,
		ManifestPath:    filepath.Join(dataDir, "manifest.json"),
		MetaDataPath:    filepath.Join(dataDir, "meta.json"),

		ModPacksDir:    filepath.Join(installerDir, "modpacks"),
		ModsBackupsDir: filepath.Join(installerDir, "modpacks", "backups"),

		ResourceBundlesDir:     filepath.Join(installerDir, "resourcebundles"),
		ResourcePackBackupsDir: filepath.Join(installerDir, "resourcebundles", "backups"),
	}, nil
}

//var (
//	userProfile, _      = os.UserHomeDir()
//	ModFolderPath       = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods")
//	ModBackupPath       = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods_old")
//	VerFilePath         = filepath.Join(userProfile, "AppData", "Roaming", ".minecraft", "mods", "ver.txt")
//	TempModDownloadPath = filepath.Join(os.TempDir(), "temp-mod-downloads")
//)
