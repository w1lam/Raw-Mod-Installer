// Package downloader provides downloading functionality
package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/w1lam/Packages/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/ui"
)

// DownloaderResults is the results of the downloader
type DownloaderResults struct {
	TempDir         string
	DownloadedItems map[string]DownloadItem
}

// DownloadEntries is the main download function
func DownloadEntries(
	entries map[string]DownloadItem,
	path *paths.Paths,
) (DownloaderResults, error) {
	tempDir, err := os.MkdirTemp(path.ProgramFilesDir, "downloads")
	if err != nil {
		return DownloaderResults{}, err
	}

	progressCh := make(chan ui.DownloaderProgress)
	var wg sync.WaitGroup

	results := DownloaderResults{
		TempDir:         tempDir,
		DownloadedItems: make(map[string]DownloadItem),
	}

	var mu sync.Mutex

	for id, entry := range entries {
		wg.Add(1)
		entry := entry
		go func() {
			defer wg.Done()
			uri := entry.URL

			filePath := filepath.Join(tempDir, entry.FileName)

			progressCh <- ui.DownloaderProgress{File: entry.FileName, Status: "downloading"}

			computedSha, err := download.DownloadFile(filePath, uri)
			if err != nil {
				progressCh <- ui.DownloaderProgress{File: entry.FileName, Status: "failure", Err: err}
				return
			}

			if entry.Sha512 != "" {
				if computedSha != entry.Sha512 {
					progressCh <- ui.DownloaderProgress{File: entry.FileName, Status: "failure", Err: fmt.Errorf("SHA512 mismatch")}
					return
				}
			} else if entry.Sha1 == "" {
				if computedSha != entry.Sha1 {
					progressCh <- ui.DownloaderProgress{File: entry.FileName, Status: "failure", Err: fmt.Errorf("SHA1 mismatch")}
					return
				}
			}

			progressCh <- ui.DownloaderProgress{File: entry.FileName, Status: "success"}

			mu.Lock()
			results.DownloadedItems[id] = DownloadItem{
				ID:       id,
				FileName: entry.FileName,
				Sha512:   entry.Sha512,
				Sha1:     entry.Sha1,
				Version:  entry.Version,
			}
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(progressCh)
	}()

	_, failedFiles := ui.RenderDownloaderProgress(progressCh, len(entries))

	if len(failedFiles) > 0 {
		return DownloaderResults{}, fmt.Errorf("%d mods failed to install", len(failedFiles))
	}

	return results, nil
}
