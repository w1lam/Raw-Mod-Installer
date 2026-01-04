package install

import (
	"fmt"
	"os"

	"github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// ADD RENAME FUNCTION TO MOVE TEMP FOLDER TO MOD FOLDER

func FullInstall() error {
	path, err := paths.Resolve()
	if err != nil {
		return err
	}

	entries, err := modlist.GetModEntryList(netcfg.ModListURL)
	if err != nil {
		return err
	}

	resolved, err := resolve.ResolveModListConcurrent(entries, config.McVersion, SimpleProgress)
	if err != nil {
		return err
	}

	ctx := buildInstallContext(resolved)

	progressCh := make(chan download.Progress)

	go func() {
		_ = DownloadConcurrent(
			resolve.ResolvedModList(resolved).GetURLs(),
			path.TempDownloadDir,
			progressCh,
		)
	}()

	successFiles, failedFiles := RenderProgress(progressCh)

	if len(failedFiles) > 0 {
		return fmt.Errorf("%d mods failed to install", len(failedFiles))
	}

	UpdateManifestInstalledVersions(app.GlobalManifest, ctx, successFiles)

	version, err := modlist.GetRemoteVersion(netcfg.ModListURL)
	if err == nil {
		app.GlobalManifest.ModList.InstalledVersion = version
	}

	if err := os.Rename(path.TempDownloadDir, path.ModsDir); err != nil {
		return fmt.Errorf("failed to move mods to final directory: %w", err)
	}

	return manifest.Save(path.ManifestPath, app.GlobalManifest)
}
