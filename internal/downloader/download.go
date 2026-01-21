// Package downloader provides downloading functionality
package downloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/w1lam/Packages/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// Progress represents the download progress of a file
type Progress struct {
	File   string
	Status string // "downloading", "success", "failure"
	Err    error
}

type DownloaderResults struct {
	TempDir        string
	DownloadedMods map[string]manifest.ManifestMod
}

// ModsDownloader is the download function for mods
func ModsDownloader(
	resolvedMods map[string]resolve.ResolvedMod,
	m *manifest.Manifest,
) (DownloaderResults, error) {
	tempDir, err := os.MkdirTemp(m.Paths.ProgramFilesDir, "downloads")
	if err != nil {
		return DownloaderResults{}, err
	}

	progressCh := make(chan Progress)
	var wg sync.WaitGroup

	results := DownloaderResults{
		TempDir:        tempDir,
		DownloadedMods: make(map[string]manifest.ManifestMod),
	}

	var mu sync.Mutex

	for _, mod := range resolvedMods {
		wg.Add(1)
		mod := mod
		go func() {
			defer wg.Done()
			uri := mod.DownloadURL

			fileName := filepath.Base(uri)
			filePath := path.Join(tempDir, fileName)

			progressCh <- Progress{File: fileName, Status: "downloading"}

			computedSha, err := download.DownloadFile(filePath, uri)
			if err != nil {
				progressCh <- Progress{File: fileName, Status: "failure", Err: err}
				return
			}

			if mod.Sha512 != "" {
				if computedSha != mod.Sha512 {
					progressCh <- Progress{File: fileName, Status: "failure", Err: fmt.Errorf("SHA512 mismatch")}
					return
				}
			} else if mod.Sha1 == "" {
				if computedSha != mod.Sha1 {
					progressCh <- Progress{File: fileName, Status: "failure", Err: fmt.Errorf("SHA1 mismatch")}
					return
				}
			}

			progressCh <- Progress{File: fileName, Status: "success"}

			mu.Lock()
			results.DownloadedMods[mod.Slug] = manifest.ManifestMod{
				Slug:             mod.Slug,
				FileName:         fileName,
				Sha512:           mod.Sha512,
				Sha1:             mod.Sha1,
				InstalledVersion: mod.Version.VersionNumber,
			}
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(progressCh)
	}()

	_, failedFiles := RenderProgress(progressCh, len(resolvedMods))

	if len(failedFiles) > 0 {
		return DownloaderResults{}, fmt.Errorf("%d mods failed to install", len(failedFiles))
	}

	return results, nil
}

// RenderProgress simplified version of progress bar, pure cli
func RenderProgress(ch <-chan Progress, total int) (successFiles []string, failedFiles []string) {
	success := 0
	failed := 0
	fmt.Println("[Downloading Mods: 0/", total, ", 0%]")

	for p := range ch {
		switch p.Status {
		case "success":
			success++
			successFiles = append(successFiles, p.File)
		case "failure":
			failed++
			failedFiles = append(failedFiles, p.File)
		}

		processed := success + failed
		if total == 0 {
			total = 1
		}

		procent := int(float64(processed) / float64(total) * float64(100))
		if procent > 100 {
			procent = 100
		}

		fmt.Println("[Downloading Mods: ", success, "/", total, ", ", procent, "%]")
	}

	return
}
