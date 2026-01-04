// Package install handles downloading mods from given URLs and displays progress.
package install

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/w1lam/Packages/pkg/download"
)

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
