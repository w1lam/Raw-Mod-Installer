// Package resolve provides functions to resolve mod information concurrently.
package resolve

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
)

// FetchModInfoList fetches mod information from Modrinth for a list of mod entries concurrently.
func ResolveModInfoList(modEntryList []modlist.ModEntry, maxConcurrentFetches int) ([]modrinth.ModInfo, error) {
	modInfoList := make([]modrinth.ModInfo, len(modEntryList))

	sem := make(chan struct{}, maxConcurrentFetches)

	var wg sync.WaitGroup
	wg.Add(len(modEntryList))

	errCh := make(chan error, 1)

	for i, mod := range modEntryList {
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			info, err := modrinth.FetchModInfo(mod.Slug)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			modInfoList[i] = info
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
		return modInfoList, nil
	}
}

// ResolveModListConcurrent fetches the latest download URLs for all mods concurrently, reporting progress via the provided function.
func ResolveModListConcurrent(
	mods []modlist.ModEntry,
	mcVersion string,
	progressFunc func(done, total int, currentMod string),
) ([]ResolvedMod, error) {
	//
	total := len(mods)
	results := make([]ResolvedMod, total)
	errChan := make(chan error, total)

	var wg sync.WaitGroup
	var done int32

	for i, mod := range mods {
		wg.Add(1)

		go func(i int, mod modlist.ModEntry) {
			defer wg.Done()
			defer func() {
				atomic.AddInt32(&done, 1)
				progressFunc(int(done), total, mod.Slug)
			}()

			resolved, err := ResolveMod(mod.Slug, mcVersion, mod.Loader)
			if err != nil {
				errChan <- fmt.Errorf("%s: %w", mod.Slug, err)
				return
			}

			results[i] = resolved
		}(i, mod)
	}

	wg.Wait()
	close(errChan)

	var combined strings.Builder
	for err := range errChan {
		combined.WriteString(err.Error())
		combined.WriteByte('\n')
	}

	if combined.Len() > 0 {
		return nil, errors.New(combined.String())
	}
	return results, nil
}
