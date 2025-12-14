// Package downloadmods handles downloading mods from given URLs and displays progress.
package downloadmods

import (
	"fmt"
	"os"
	"sync"

	dl "github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Packages/pkg/fetch"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func DisplayDownloadProgress(progressCh <-chan dl.Progress, fopts string, listURL string) {
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

// Main Mod Download Function

func DownloadMods(listURL string, fopts string) error {
	// Fetching Mod List
	fmt.Printf("Fetching Mod List...\n")
	slugList, err := fetch.GetList(listURL)
	if err != nil {
		return err
	}

	// Parsing Mod List
	fmt.Printf("Parsing Mod List...\n")
	parsedList, err := modrinth.ParseModList(slugList)
	if err != nil {
		return err
	}

	// Fetching Download URLs
	fmt.Printf("Fetching Download URLs...\n")
	fetchedURLs, err := modrinth.FetchAllConcurrent(parsedList, "1.21.10", modrinth.SimpleProgress)
	if err != nil {
		return err
	}

	// Creates temp mod download dir
	err2 := os.MkdirAll(paths.TempModDownloadPath, os.ModePerm)
	if err2 != nil {
		return fmt.Errorf("ERROR: failed to create temp dir: %v", err2)
	}

	progressCh := make(chan dl.Progress)

	var wg sync.WaitGroup

	// Start downloads concurrently
	wg.Go(func() {
		dl.DownloadFromListConcurrent(fetchedURLs, paths.TempModDownloadPath, progressCh)
	})

	// Handle UI printing in main goroutine
	func(progressCh <-chan dl.Progress, fopts string, listURL string) {
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
	}(progressCh, fopts, listURL)

	wg.Wait()

	return nil
}
