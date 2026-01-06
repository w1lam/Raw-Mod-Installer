// Package install handles downloading mods from given URLs and displays progress.
package install

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// DownloadModpack is the full modpack download function
func DownloadModpack(m *manifest.Manifest, path *paths.Paths) (*manifest.Manifest, error) {
	entries, err := modlist.GetModEntryList(netcfg.ModListURL)
	if err != nil {
		return nil, err
	}

	resolved, err := resolve.ResolveModListConcurrent(entries, config.McVersion, SimpleProgress)
	if err != nil {
		return nil, err
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

	if err := os.Rename(path.TempDownloadDir, path.UnloadedModsDir); err != nil {
		return nil, fmt.Errorf("failed to move downloaded mods to installer mod dir: %s", err)
	}

	if len(failedFiles) > 0 {
		return nil, fmt.Errorf("%d mods failed to install", len(failedFiles))
	}

	UpdateManifestInstalledVersions(m, ctx, successFiles)

	if version, err := modlist.GetRemoteVersion(netcfg.ModListURL); err != nil {
		m.ModList.InstalledVersion = version
	}

	return m, nil
}

// DownloadConcurrent downloads mods from urls concurrently
func DownloadConcurrent(
	urls []string,
	destPath string,
	progressCh chan<- download.Progress,
) error {
	if err := os.MkdirAll(destPath, 0o755); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, uri := range urls {
		uri := uri
		wg.Add(1)

		go func() {
			defer wg.Done()

			file := filepath.Base(uri)
			target := filepath.Join(destPath, file)

			progressCh <- download.Progress{
				File:   file,
				Status: "downloading",
			}

			if err := download.DownloadFile(target, uri); err != nil {
				progressCh <- download.Progress{
					File:   file,
					Status: "failure",
					Err:    err,
				}
				return
			}

			progressCh <- download.Progress{
				File:   file,
				Status: "success",
			}
		}()
	}

	go func() {
		wg.Wait()
		close(progressCh)
	}()

	return nil
}
