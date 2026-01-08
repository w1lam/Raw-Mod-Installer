package install

import (
	"fmt"
	"sync"

	"github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
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

// ResolveUpdates resolved updates for specified modpack
func ResolveUpdates(
	m *manifest.Manifest,
	modPack *manifest.InstalledModPack,
) ([]resolve.ResolvedMod, error) {
	var (
		results []resolve.ResolvedMod
		wg      sync.WaitGroup
		mu      sync.Mutex
		errCh   = make(chan error, 1)
		sem     = make(chan struct{}, 8)
	)

	for _, s := range modPack.GetSlugs() {
		s := s
		wg.Add(1)

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			latest, err := modrinth.FetchLatestModrinthVersion(
				s,
				modPack.McVersion,
				modPack.Loader,
			)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			needs := m.EnabledModPack.InstalledVersion == "" || latest.VersionNumber != m.EnabledModPack.InstalledVersion
			if !needs {
				return
			}

			resolved := resolve.ResolvedMod{
				Slug:        s,
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

// UpdateModpack upates the specified mod pack
func UpdateModpack(modPack manifest.InstalledModPack, m *manifest.Manifest) (*manifest.Manifest, error) {
	updates, err := ResolveUpdates(m, m.EnabledModPack)
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return m, nil
	}

	progressCh := make(chan download.Progress)

	go func() {
		download.DownloadMultipleConcurrent(
			resolve.ResolvedModList(updates).GetURLs(),
			m.Paths.TempDownloadDir,
			progressCh,
		)
	}()

	success, failed := RenderProgress(progressCh)
	if len(failed) > 0 {
		return nil, fmt.Errorf("%d updates failed", len(failed))
	}

	UpdateManifestInstalledVersions(m, buildInstallContext(updates), modPack.Name, success)

	if modPacks, err := modpack.GetAvailableModPacks(netcfg.ModPacksURL); err == nil {
		m.EnabledModPack.InstalledVersion = modPacks[modPack.Name].ListVersion
	} else {
		return m, err
	}

	return m, nil
}
