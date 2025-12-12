// Package download handles downloading mods from given URLs and displays progress.
package download

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/w1lam/Packages/pkg/dl"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

type Progress struct {
	File   string
	Status string // "downloading", "success", "failure"
	Err    error
}

func DisplayProgress(progressCh <-chan Progress, fopts string, listURL string) {
	success, failures, active := 0, 0, 0

	for p := range progressCh {
		switch p.Status {
		case "downloading":
			active++
			fmt.Print("\n ◌ ", p.File, "...")
		case "success":
			active--
			success++
			fmt.Print("\n ● ", p.File, " ✓")
		case "failure":
			active--
			failures++
			fmt.Print("\n ✗ ", p.File, ": ", p.Err)
		}
		fmt.Print(" [Active: ", active, " | Success: ", success, " | Failures: ", failures, "]")
	}

	if failures == 0 {
		version, _ := features.GetRemoteVersion(listURL)
		os.WriteFile("ver.txt", []byte(version), 0o755)
		fmt.Printf("\n\nAll %d %s installed successfully! ✓", success, fopts)
	} else {
		fmt.Printf("\n\n%d %s failed to install! ✗ \n", failures, fopts)
	}
}

func DownloadFromListConcurrent(list []string, progressCh chan<- Progress) {
	var wg sync.WaitGroup

	for _, uri := range list {
		wg.Add(1)
		go func(uri string) {
			defer wg.Done()
			fileName := filepath.Base(uri)
			progressCh <- Progress{File: fileName, Status: "downloading"}

			err := dl.DownloadFile(fileName, uri)
			if err != nil {
				progressCh <- Progress{File: fileName, Status: "failure", Err: err}
				return
			}
			progressCh <- Progress{File: fileName, Status: "success"}
		}(uri)
	}

	wg.Wait()
	close(progressCh)
}

// Main Mod Download Function

func DownloadMods(listURL string, fopts string) error {
	fmt.Printf("Fetching Mod List...\n")
	slugList, err := dl.GetList(listURL)
	if err != nil {
		return err
	}

	fmt.Printf("Parsing Mod List...\n")
	parsedList, err := modrinth.ParseModList(slugList)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching Download URLs...\n")
	fetchedURLs, err := modrinth.FetchAllConcurrent(parsedList, "1.21.10", modrinth.SimpleProgress)
	if err != nil {
		return err
	}

	// Creates temp mod download dir
	err2 := os.MkdirAll(paths.TempModDownloadPath, os.ModePerm)
	if err2 != nil {
		return fmt.Errorf("ERROR: failed to change to temp dir: %v", err2)
	}
	err3 := os.Chdir(paths.TempModDownloadPath)
	if err3 != nil {
		return fmt.Errorf("ERROR: failed to change to temp mod download dir: %v", err3)
	}

	progressCh := make(chan Progress)

	var wg sync.WaitGroup

	// Start downloads concurrently
	wg.Go(func() {
		DownloadFromListConcurrent(fetchedURLs, progressCh)
	})

	// Handle UI printing in main goroutine
	DisplayProgress(progressCh, fopts, listURL)
	wg.Wait()

	os.Chdir(paths.UserProfile)

	return nil
}
