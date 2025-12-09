// Package fetch provides functions to fetch mod download URLs from Modrinth.
package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

type MRFile struct {
	URL     string `json:"url"`
	Primary bool   `json:"primary"`
}

type MRVersion struct {
	Files []MRFile `json:"files"`
}

func FetchLatestModrinthDownload(slug, mcVersion, loader string) (string, error) {
	fetch := func(params url.Values) ([]MRVersion, error) {
		base := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", slug)
		finalURL := base + "?" + params.Encode()

		req, _ := http.NewRequest("GET", finalURL, nil)
		req.Header.Set("User-Agent", "MyModInstaller/1.0")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var versions []MRVersion
		if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
			return nil, err
		}
		return versions, nil
	}

	params := url.Values{}
	params.Set("game_versions", fmt.Sprintf("[\"%s\"]", mcVersion))
	params.Set("loaders", fmt.Sprintf("[\"%s\"]", loader))

	versions, err := fetch(params)
	if err == nil && len(versions) > 0 && len(versions[0].Files) > 0 {
		return versions[0].Files[0].URL, nil
	}

	paramsFallback := url.Values{}
	paramsFallback.Set("loaders", fmt.Sprintf("[\"%s\"]", loader))

	versions, err = fetch(paramsFallback)
	if err == nil && len(versions) > 0 && len(versions[0].Files) > 0 {
		return versions[0].Files[0].URL, nil
	}

	return "", fmt.Errorf("no compatible or fallback versions found for %s", slug)
}

type ModEntry struct {
	Loader string
	Slug   string
}

func ParseModList(lines []string) ([]ModEntry, error) {
	mods := []ModEntry{}

	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}

		loader := "fabric"
		slug := line

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			loader = parts[0]
			slug = parts[1]
		}

		if strings.Contains(slug, "@") {
			slug = strings.Split(slug, "@")[0]
		}

		mods = append(mods, ModEntry{Loader: loader, Slug: slug})
	}
	return mods, nil
}

func FetchAllDownloadURLs(mods []ModEntry, mcVersion string) ([]string, error) {
	var urls []string

	for _, mod := range mods {
		url, err := FetchLatestModrinthDownload(mod.Slug, mcVersion, mod.Loader)
		if err != nil {
			return nil, fmt.Errorf("mod %s: %v", mod.Slug, err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func FetchAllConcurrent(
	mods []ModEntry,
	mcVersion string,
	progressFunc func(done, total int, currentMod string),
) ([]string, error) {
	total := len(mods)
	results := make([]string, total)
	errChan := make(chan error, total)

	var wg sync.WaitGroup
	var done int32 = 0

	for i, mod := range mods {
		wg.Add(1)

		go func(i int, mod ModEntry) {
			defer wg.Done()

			url, err := FetchLatestModrinthDownload(mod.Slug, mcVersion, mod.Loader)
			if err != nil {
				errChan <- fmt.Errorf("%s: %w", mod.Slug, err)
				return
			}

			results[i] = url

			atomic.AddInt32(&done, 1)
			progressFunc(int(done), total, mod.Slug)
		}(i, mod)
	}

	wg.Wait()
	close(errChan)

	var combined strings.Builder
	for e := range errChan {
		combined.WriteString(e.Error() + "\n")
	}

	if combined.Len() > 0 {
		return nil, errors.New(combined.String())
	}
	return results, nil
}

func SimpleProgress(done, total int, mod string) {
	fmt.Printf("Fetched %d/%d -> %s\n", done, total, mod)
}
