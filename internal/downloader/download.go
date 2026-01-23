// Package downloader provides downloading functionality
package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/w1lam/Packages/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Progress represents the download progress of a file
type Progress struct {
	File   string
	Status string // "downloading", "success", "failure"
	Err    error
}

type DownloaderResults struct {
	TempDir         string
	DownloadedItems map[string]DownloadItem
}

type DownloadItem struct {
	ID       string // slug
	FileName string
	URL      string
	Sha1     string
	Sha512   string
	Version  string
}

// DownloadEntries is the main download function
func DownloadEntries(
	entries []DownloadItem,
	path *paths.Paths,
) (DownloaderResults, error) {
	tempDir, err := os.MkdirTemp(path.ProgramFilesDir, "downloads")
	if err != nil {
		return DownloaderResults{}, err
	}

	progressCh := make(chan Progress)
	var wg sync.WaitGroup

	results := DownloaderResults{
		TempDir:         tempDir,
		DownloadedItems: make(map[string]DownloadItem),
	}

	var mu sync.Mutex

	for _, entry := range entries {
		wg.Add(1)
		entry := entry
		go func() {
			defer wg.Done()
			uri := entry.URL

			fileName := filepath.Base(uri)
			filePath := filepath.Join(tempDir, fileName)

			progressCh <- Progress{File: fileName, Status: "downloading"}

			computedSha, err := download.DownloadFile(filePath, uri)
			if err != nil {
				progressCh <- Progress{File: fileName, Status: "failure", Err: err}
				return
			}

			if entry.Sha512 != "" {
				if computedSha != entry.Sha512 {
					progressCh <- Progress{File: fileName, Status: "failure", Err: fmt.Errorf("SHA512 mismatch")}
					return
				}
			} else if entry.Sha1 == "" {
				if computedSha != entry.Sha1 {
					progressCh <- Progress{File: fileName, Status: "failure", Err: fmt.Errorf("SHA1 mismatch")}
					return
				}
			}

			progressCh <- Progress{File: fileName, Status: "success"}

			mu.Lock()
			results.DownloadedItems[entry.ID] = DownloadItem{
				ID:       entry.ID,
				FileName: fileName,
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

	_, failedFiles := RenderProgress(progressCh, len(entries))

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
