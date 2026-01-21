package filesystem

import (
	"fmt"
	"os"

	"github.com/w1lam/Packages/fabric"
	"github.com/w1lam/Packages/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

type SystemState struct {
	FabricStatus FabricStatus
	MCVersion    string
}

type FabricStatus int

const (
	FabricMissing FabricStatus = iota
	FabricOutdated
	FabricUpToDate
)

// EnsureDirectories ensures all program directories exists
func EnsureDirectories(path *paths.Paths) error {
	if !utils.CheckFileExists(path.MinecraftDir) {
		return fmt.Errorf("minecraft directory not found")
	}

	if !utils.CheckFileExists(path.ProgramFilesDir) {
		err := os.MkdirAll(path.ProgramFilesDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.DataDir) {
		err := os.MkdirAll(path.DataDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.ModPacksDir) {
		err := os.MkdirAll(path.ModPacksDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.ModsBackupsDir) {
		err := os.MkdirAll(path.ModsBackupsDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.ResourceBundlesDir) {
		err := os.MkdirAll(path.ResourceBundlesDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.ResourcePackBackupsDir) {
		err := os.MkdirAll(path.ResourcePackBackupsDir, 0o755)
		if err != nil {
			return err
		}
	}

	return nil
}

func DetectSystem(mcVersion string) (SystemState, error) {
	statusStr, err := fabric.CheckVersions(mcVersion)
	if err != nil {
		return SystemState{}, err
	}

	var status FabricStatus
	switch statusStr {
	case "notInstalled":
		status = FabricMissing
	case "updateFound":
		status = FabricOutdated
	default:
		status = FabricUpToDate
	}

	return SystemState{
		FabricStatus: status,
		MCVersion:    mcVersion,
	}, nil
}

func EnsureFabric(mcVersion string) error {
	state, err := DetectSystem(mcVersion)
	if err != nil {
		return err
	}

	if state.FabricStatus == FabricUpToDate {
		return nil
	}

	jar, err := fabric.GetLatestInstallerJar()
	if err != nil {
		return err
	}

	return fabric.RunInstaller(jar, mcVersion)
}
