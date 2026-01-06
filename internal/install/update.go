package install

import (
	"fmt"
	"sync"

	"github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

type UpdateDecision int

const (
	NoChange UpdateDecision = iota
	NeedsInstall
	NeedsUpdate
)

type UpdateCandidate struct {
	Slug      string
	Installed string
	Latest    string
	Resolved  resolve.ResolvedMod
	Decision  UpdateDecision
}

func ResolveUpdates(
	m *manifest.Manifest,
	entries []modlist.ModEntry,
	mcVersion string,
) ([]resolve.ResolvedMod, error) {
	var (
		results []resolve.ResolvedMod
		wg      sync.WaitGroup
		mu      sync.Mutex
		errCh   = make(chan error, 1)
		sem     = make(chan struct{}, 8)
	)

	for _, entry := range entries {
		entry := entry
		wg.Add(1)

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			latest, err := modrinth.FetchLatestModrinthVersion(
				entry.Slug,
				mcVersion,
				entry.Loader,
			)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			installed := m.InstalledVersion(entry.Slug)

			needs := installed == "" || latest.VersionNumber != installed
			if !needs {
				return
			}

			resolved := resolve.ResolvedMod{
				Slug:        entry.Slug,
				DownloadURL: latest.Files[0].URL,
				LatestVer:   latest.VersionNumber,
			}

			mu.Lock()
			results = append(results, resolved)
			mu.Unlock()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-done:
		return results, nil

	}
}

func UpdateModpack(m *manifest.Manifest, path *paths.Paths) (*manifest.Manifest, error) {
	updates, err := ResolveUpdates(m, m.GetModEntries(), config.McVersion)
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return m, nil
	}

	progressCh := make(chan download.Progress)

	go func() {
		_ = DownloadConcurrent(
			resolve.ResolvedModList(updates).GetURLs(),
			path.TempDownloadDir,
			progressCh,
		)
	}()

	success, failed := RenderProgress(progressCh)
	if len(failed) > 0 {
		return nil, fmt.Errorf("%d updates failed", len(failed))
	}

	UpdateManifestInstalledVersions(m, buildInstallContext(updates), success)

	if version, err := modlist.GetRemoteVersion(netcfg.ModListURL); err == nil {
		m.ModList.InstalledVersion = version
	}

	return m, nil
}
