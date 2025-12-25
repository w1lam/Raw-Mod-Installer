// Package install handles downloading mods from given URLs and displays progress.
package install

import (
	"fmt"
	"os"
	"sync"

	"github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// DownloadMods downloads the mods from the provided list of resolved mods.
func DownloadMods(urls []string) error {
	path, err := paths.Resolve()
	if err != nil {
		return err
	}

	// Creates temp mod download dir
	err1 := os.MkdirAll(path.TempDownloadDir, os.ModePerm)
	if err1 != nil {
		return fmt.Errorf("ERROR: failed to create temp dir: %v", err1)
	}

	progressCh := make(chan download.Progress)

	var wg sync.WaitGroup

	// Start downloads concurrently
	wg.Go(func() {
		download.DownloadMultipleConcurrent(urls, path.TempDownloadDir, progressCh)
	})

	// Handle UI printing in main goroutine
	func(progressCh <-chan download.Progress) {
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
			version, err := modlist.GetRemoteVersion(netcfg.ModListURL)
			if err != nil {
				fmt.Printf("\n\nFfailed to fetch mod list version: %v", err)
				return
			}

			app.GlobalManifest.ModList.InstalledVersion = version

			if err := manifest.Save(path.ManifestPath, app.GlobalManifest); err != nil {
				fmt.Printf("\n\nFailed to save manifest: %v", err)
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
