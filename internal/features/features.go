// Package features contains functions related to mod list updates, backups, and restoration.
package features

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/w1lam/Packages/pkg/fetch"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

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

		remoteVer, err := GetRemoteVersion(paths.ModListURL)
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

func GetModEntryList(modListURL string) ([]modrinth.ModEntry, error) {
	rawList, err := fetch.GetList(modListURL)
	if err != nil {
		return nil, err
	}

	modEntryList, err1 := modrinth.ParseModList(rawList)
	if err1 != nil {
		return nil, err1
	}

	return modEntryList, nil
}

func GetRemoteVersion(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		cutLine, _ := strings.CutPrefix(line, "# version:")
		return strings.TrimSpace(cutLine), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}

func GetLocalVersion() string {
	data, err := os.ReadFile(paths.VerFilePath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// Mod Backup and Restore Functions

var Timestamp = new(string)

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
