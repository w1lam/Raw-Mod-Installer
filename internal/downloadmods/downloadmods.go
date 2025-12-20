// Package downloadmods handles downloading mods from given URLs and displays progress.
package downloadmods

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	dl "github.com/w1lam/Packages/pkg/download"
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

		err := os.WriteFile("ver.txt", []byte(version), 0o755)
		if err != nil {
			return
		}

		fmt.Printf("\n\nAll %d %s installed successfully! ✓", success, fopts)
	} else {
		fmt.Printf("\n\n%d %s failed to install! ✗ \n", failures, fopts)
	}
}

// DownloadMods downloads the mods from the provided list of resolved mods.
func DownloadMods(urls []string) error {
	// Creates temp mod download dir
	err2 := os.MkdirAll(paths.TempModDownloadPath, os.ModePerm)
	if err2 != nil {
		return fmt.Errorf("ERROR: failed to create temp dir: %v", err2)
	}

	progressCh := make(chan dl.Progress)

	var wg sync.WaitGroup

	// Start downloads concurrently
	wg.Go(func() {
		dl.DownloadMultipleConcurrent(urls, paths.TempModDownloadPath, progressCh)
	})

	// Handle UI printing in main goroutine
	func(progressCh <-chan dl.Progress) {
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
			version, _ := features.GetRemoteVersion(paths.ModListURL)

			err := features.WriteVersionFile(filepath.Join(paths.ModFolderPath, "ver.txt"), version)
			if err != nil {
				fmt.Printf("\n\nFailed to write version file: %v", err)
				return
			}

			fmt.Printf("\n\nAll %d Mods installed successfully! ✓", success)
		} else {
			fmt.Printf("\n\n%d Mods failed to install! ✗ \n", failures)
		}
	}(progressCh)

	wg.Wait()

	return nil
}
