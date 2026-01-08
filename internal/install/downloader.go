// Package install handles downloading mods from given URLs and displays progress.
package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// DownloadModpack is the full modpack download function
func DownloadModpack(modPack modpack.ResolvedModPackList, m *manifest.Manifest) (*manifest.Manifest, error) {
	resolved, err := resolve.ResolveModsConcurrent(modPack, SimpleProgress)
	if err != nil {
		return nil, err
	}

	ctx := buildInstallContext(resolved)

	progressCh := make(chan download.Progress)

	go func() {
		download.DownloadMultipleConcurrent(
			resolve.ResolvedModList(resolved).GetURLs(),
			m.Paths.TempDownloadDir,
			progressCh,
		)
	}()

	successFiles, failedFiles := RenderProgress(progressCh)

	if err := os.Rename(m.Paths.TempDownloadDir, filepath.Join(m.Paths.ModPacksDir, modPack.Name)); err != nil {
		return nil, fmt.Errorf("failed to move downloaded mods to modpacks dir: %s", err)
	}

	if len(failedFiles) > 0 {
		return nil, fmt.Errorf("%d mods failed to install", len(failedFiles))
	}

	UpdateManifestInstalledVersions(m, ctx, modPack.Name, successFiles)

	return m, nil
}
