package manifest

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/Packages/pkg/fundamental"
	"github.com/w1lam/Raw-Mod-Installer/internal/modrinthsvc"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func IsModListPresent() bool {
	return fundamental.CheckFileExists(paths.VerFilePath)
}

func WriteVersionFile(path, version string) error {
	// Ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	// Write atomically
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(version), 0o644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func CheckForModlistUpdate() (bool, error) {
	if _, err := os.Stat(paths.VerFilePath); err == nil {

		remoteVer, err := modrinthsvc.GetRemoteVersion(netcfg.ModListURL)
		localVer := GetLocalVersion()

		switch {
		case err != nil:
			return false, err

		case remoteVer != localVer:
			return true, nil

		case remoteVer == localVer:
			return false, nil

		default:
			return false, err
		}
	}
	return false, nil
}

func GetLocalVersion() string {
	data, err := os.ReadFile(paths.VerFilePath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
