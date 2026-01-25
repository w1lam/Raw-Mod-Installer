package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/w1lam/Packages/modrinth"
)

type DownloadItem struct {
	ID       string // slug
	FileName string
	URL      string
	Sha1     string
	Sha512   string
	Version  string
}

func ResolveDownloadItem(entries []modrinth.ModrinthListEntry, filter modrinth.EntryFilter) (map[string]DownloadItem, error) {
	fmt.Println("Fetching Best Versions...")

	bestVersions := modrinth.FetchBestVersions(entries, filter)
	out := map[string]DownloadItem{}

	fmt.Println("Building Download Item...")
	for _, entry := range entries {
		version, ok := bestVersions[entry.Slug]
		if !ok || version == nil {
			return nil, fmt.Errorf("no compatible version found for %s (mc=%s loader=%s)", entry.Slug, filter.McVersion, filter.Loader)
		}

		if len(version.Files) == 0 {
			return nil, fmt.Errorf("no downloadable files for %s", entry.Slug)
		}

		file := version.Files[0]
		for _, f := range version.Files {
			if f.Primary {
				file = f
				break
			}
		}

		out[entry.Slug] = DownloadItem{
			ID:       entry.Slug,
			FileName: filepath.Base(file.URL),
			URL:      file.URL,
			Sha1:     file.Hashes.Sha1,
			Sha512:   file.Hashes.Sha512,
			Version:  version.VersionNumber,
		}
	}

	return out, nil
}
