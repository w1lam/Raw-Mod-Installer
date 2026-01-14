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
}

func ResolveMods(slugs []string, mcVersion, loader string) (map[string]ResolvedMod, error) {
	finished := make(chan struct{})
	go tui.RawSpinner(finished, []rune{'▙', '▛', '▜', '▟'}, config.Style.Margin, "Resolving Mods")
	defer close(finished)

	bestVersions := modrinth.FetchBestVersions(slugs, mcVersion, loader)

	out := map[string]ResolvedMod{}

	for _, slug := range slugs {
		version := bestVersions[slug]

		downloadURL := version.Files[0].URL
		for _, f := range version.Files {
			if f.Primary {
				downloadURL = f.URL
				break
			}
		}

		out[slug] = ResolvedMod{
			Slug:        slug,
			Version:     *version,
			DownloadURL: downloadURL,
		}
	}

	return out, nil
}
