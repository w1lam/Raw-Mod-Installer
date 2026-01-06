package resolve

import (
	"sync"

	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
)

// ResolveModInfoList fetches mod information from Modrinth for a list of mod entries concurrently.
func ResolveModInfoList(modEntryList []modlist.ModEntry, maxConcurrentFetches int) ([]modrinth.ModInfo, error) {
	modInfoList := make([]modrinth.ModInfo, len(modEntryList))

	sem := make(chan struct{}, maxConcurrentFetches)

	var wg sync.WaitGroup
	wg.Add(len(modEntryList))

	errCh := make(chan error, 1)

	for i, mod := range modEntryList {
		i := i
		mod := mod

		go func(i int, mod modlist.ModEntry) {
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
		}(i, mod)
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
