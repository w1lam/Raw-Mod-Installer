// Package fabric provides functions to fetch the latest Fabric Loader and Installer versions.
package fabric

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CheckFabricVersions compares the local Fabric Loader version with the latest available version for the given Minecraft version.
func CheckFabricVersions(mcVersion string) (string, error) {
	latestVer, err := GetLatestLoaderVersion(mcVersion)
	if err != nil {
		return "", err
	}

	localVer, err1 := GetLocalFabricVersion(mcVersion)
	if err1 != nil {
		return "", err1
	}

	switch {

	case latestVer > localVer:
		return "updateFound", nil

	case latestVer == localVer:
		return "upToDate", nil

	default:
		return "", nil
	}
}

// GetMinecraftFolder returns the path to the Minecraft folder based on the operating system.
func GetMinecraftFolder() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), ".minecraft")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "minecraft")
	default:
		return filepath.Join(os.Getenv("HOME"), ".minecraft")
	}
}

// GetLocalFabricVersion retrieves the installed Fabric Loader version for the specified Minecraft version.
func GetLocalFabricVersion(mcVersion string) (string, error) {
	mcFolder := GetMinecraftFolder()
	versionDir := filepath.Join(mcFolder, "versions")

	entries, err := os.ReadDir(versionDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), "fabric-loader") && strings.HasSuffix(entry.Name(), mcVersion) {
			// folder name format: fabric-loader-0.18.1-1.21.10
			parts := strings.Split(entry.Name(), "-")
			for i, part := range parts {
				if part == "loader" && i+1 < len(parts) {
					return parts[i+1], nil
				}
			}
		}
	}

	return "", fmt.Errorf("no fabric loader installed for MC %s", mcVersion)
}

// LoaderData represents the JSON structure for Fabric Loader version data.
type LoaderData struct {
	Loader struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
		Maven   string `json:"maven"`
	} `json:"loader"`
}

// GetLatestLoaderVersion fetches the latest Fabric Loader version for the specified Minecraft version.
func GetLatestLoaderVersion(mcVersion string) (string, error) {
	url := fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s", mcVersion)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var loaders []LoaderData
	if err := json.NewDecoder(resp.Body).Decode(&loaders); err != nil {
		return "", err
	}

	if len(loaders) == 0 {
		return "", fmt.Errorf("no loader versions found for MC %s", mcVersion)
	}

	// The API returns the latest version first
	return loaders[0].Loader.Version, nil
}

// FabricInstaller represents the JSON structure for Fabric Installer data.
type FabricInstaller struct {
	URL     string `json:"url"`
	Maven   string `json:"maven"`
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// GertLatestFabricInstallerJar fetches the latest Fabric Installer JAR and saves it to a temporary location.
func GetLatestFabricInstallerJar() (string, error) {
	resp, err := http.Get("https://meta.fabricmc.net/v2/versions/installer")
	if err != nil {
		return "", fmt.Errorf("failed to fetch installer list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		return "", fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
	}

	var installers []FabricInstaller
	if err := json.NewDecoder(resp.Body).Decode(&installers); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	if len(installers) == 0 {
		return "", fmt.Errorf("installer list is empty")
	}

	// The API returns newest first already.
	inst := installers[0]

	jarURL := inst.URL
	if jarURL == "" {
		return "", fmt.Errorf("installer entry has no URL")
	}

	// Download to temp dir
	tmpPath := filepath.Join(os.TempDir(), "fabric-installer.jar")

	jarResp, err := http.Get(jarURL)
	if err != nil {
		return "", fmt.Errorf("failed to download jar: %w", err)
	}
	defer jarResp.Body.Close()

	if jarResp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(jarResp.Body, 500))
		return "", fmt.Errorf("bad jar status %d: %s", jarResp.StatusCode, string(body))
	}

	out, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp jar: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, jarResp.Body); err != nil {
		return "", fmt.Errorf("failed writing jar: %w", err)
	}

	// Validate: JAR must be a valid ZIP
	if z, err := zip.OpenReader(tmpPath); err != nil {
		return "", fmt.Errorf("file is not a valid .jar (zip): %w", err)
	} else {
		z.Close()
	}

	return tmpPath, nil
}

// RunFabricInstaller runs the Fabric installer JAR for the specified Minecraft version.
func RunFabricInstaller(jarPath, mcVersion string) error {
	java, err := exec.LookPath("java.exe")
	if err != nil {
		return err
	}
	cmd := exec.Command(java, "-jar", jarPath, "client", "-mcversion", mcVersion)
	err2 := cmd.Start()
	if err2 != nil {
		return err2
	}

	return nil
}
