package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/w1lam/Packages/modrinth"
)

// DownloadItem represents a single item to be downloaded
type DownloadItem struct {
	ID       string // ID/slug
	FileName string
	URL      string
	Sha1     string
	Sha512   string
	Version  string
	Type     modrinth.EntryType
}

// ResolveDownloadItem resolves the download items that the downloader consumes
func ResolveDownloadItem(entries []modrinth.Entry, filter modrinth.Filter) (map[string]DownloadItem, error) {
	bestVersions := modrinth.FetchBestVersions(entries, filter)
	out := map[string]DownloadItem{}

	for _, entry := range entries {
		version, ok := bestVersions[entry.ID]
		if !ok || version == nil {
			return nil, fmt.Errorf("no compatible version found for %s (mc=%s loader=%s)", entry.ID, filter.McVersion, filter.Loader)
		}

		if len(version.Files) == 0 {
			return nil, fmt.Errorf("no downloadable files for %s", entry.ID)
		}

		file := version.Files[0]
		for _, f := range version.Files {
			if f.Primary {
				file = f
				break
			}
		}

		out[entry.ID] = DownloadItem{
			ID:       entry.ID,
			FileName: filepath.Base(file.URL),
			URL:      file.URL,
			Sha1:     file.Hashes.Sha1,
			Sha512:   file.Hashes.Sha512,
			Version:  version.VersionNumber,
		}
	}

	return out, nil
}
