// Package resolve resolves mod versions and metadata
package resolve

import (
	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Packages/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
)

// ResolvedMod is a fullt resolved mod that has all metadata
type ResolvedMod struct {
	Slug        string             `json:"slug"`
	Version     modrinth.MRVersion `json:"version"`
	DownloadURL string             `json:"url"`
	Sha512      string             `json:"sha512"`
	Sha1        string             `json:"sha1"`
}

func ResolveDownloadItem(entries []modrinth.ModrinthListEntry, mcVersion, loader string) (map[string]ResolvedMod, error) {
	finished := make(chan struct{})
	go tui.RawSpinner(finished, []rune{'▙', '▛', '▜', '▟'}, config.Style.Margin, "Resolving Mods")
	defer close(finished)

	bestVersions := modrinth.FetchBestVersions(entries, mcVersion, loader)

	out := map[string]ResolvedMod{}

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

		out[entry.Slug] = ResolvedMod{
			Slug:        entry.Slug,
			Version:     *version,
			DownloadURL: downloadURL,
			Sha512:      sha512,
			Sha1:        sha1,
		}
	}

	return out, nil
}
