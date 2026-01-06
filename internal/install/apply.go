package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

type installContext struct {
	FileToMod map[string]resolve.ResolvedMod
}

func buildInstallContext(mods []resolve.ResolvedMod) installContext {
	m := make(map[string]resolve.ResolvedMod)

	for _, mod := range mods {
		filename := filepath.Base(mod.DownloadURL)
		m[filename] = mod
	}

	return installContext{FileToMod: m}
}

func UpdateManifestInstalledVersions(
	manifest *manifest.Manifest,
	ctx installContext,
	successFiles []string,
) {
	for _, file := range successFiles {
		mod, ok := ctx.FileToMod[file]
		if !ok {
			continue
		}

		entry := manifest.Mods[mod.Slug]
		entry.InstalledVersion = mod.LatestVer
		manifest.Mods[mod.Slug] = entry
	}
}

func prepareFS(path *paths.Paths, plan InstallPlan) error {
	switch plan.BackupPolicy {
	case BackupIfExists:
		return BackupIfNeeded(path)

	case BackupOnce:
		if !utils.CheckFileExists(path.BackupDir) {
			return BackupIfNeeded(path)
		}
	}
	if plan.Intent == IntentReinstall {
		return os.RemoveAll(path.ModsDir)
	}

	return nil
}

func EnableMods(path *paths.Paths) error {
	if utils.CheckFileExists(path.ModsDir) {
		return fmt.Errorf("mods aleady enabled")
	}

	return os.Rename(path.UnloadedModsDir, path.ModsDir)
}

func DisableMods(path *paths.Paths) error {
	if !utils.CheckFileExists(path.ModsDir) {
		return nil
	}
	return os.Rename(path.ModsDir, path.UnloadedModsDir)
}
