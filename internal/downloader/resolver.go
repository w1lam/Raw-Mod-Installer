package downloader

import (
	"path/filepath"

	"github.com/w1lam/Packages/modrinth"
)

func ResolveDownloadItem(entries []modrinth.ModrinthListEntry, mcVersion, loader string) (map[string]DownloadItem, error) {
	bestVersions := modrinth.FetchBestVersions(entries, mcVersion, loader)

	out := map[string]DownloadItem{}

	for _, entry := range entries {
		version := bestVersions[entry.Slug]

		downloadURL := version.Files[0].URL
		sha512 := version.Files[0].Hashes.Sha512
		sha1 := version.Files[0].Hashes.Sha1
		for _, f := range version.Files {
			if f.Primary {
				sha512 = f.Hashes.Sha512
				sha1 = f.Hashes.Sha1
				downloadURL = f.URL
				break
			}
		}

		out[entry.Slug] = DownloadItem{
			ID:       entry.Slug,
			FileName: filepath.Base(downloadURL),
			URL:      downloadURL,
			Sha1:     sha1,
			Sha512:   sha512,
			Version:  version.VersionNumber,
		}
	}

	return out, nil
}
