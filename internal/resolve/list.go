package resolve

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
)

// ResolveModListConcurrent fetches the latest download URLs for all mods concurrently, reporting progress via the provided function.
func ResolveModsConcurrent(
	resolvedModPack modpack.ResolvedModPackList,
	progressFunc func(done, total int, currentMod string),
) ([]ResolvedMod, error) {
	//
	total := len(resolvedModPack.Slugs)
	results := make([]ResolvedMod, total)
	errChan := make(chan error, total)

	var wg sync.WaitGroup
	var done int32

	for i, mod := range resolvedModPack.Slugs {
		wg.Add(1)

		go func(i int, mod string) {
			defer wg.Done()
			defer func() {
				atomic.AddInt32(&done, 1)
				progressFunc(int(done), total, mod)
			}()

			resolved, err := ResolveMod(mod, resolvedModPack.McVersion, resolvedModPack.Loader)
			if err != nil {
				errChan <- fmt.Errorf("%s: %w", mod, err)
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
