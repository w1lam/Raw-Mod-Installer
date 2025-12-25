package resolve

import (
	"log"
	"strings"

	"github.com/w1lam/Packages/pkg/modrinth"
)

// ResolvedMod represents a mod with its slug, latest version, and download URL.
type ResolvedMod struct {
	Slug        string
	ProjectID   string
	FabricID    string
	LatestVer   string
	LocalVer    string
	DownloadURL string
}

// ResolvedModList is a list of ResolvedMod objects.
type ResolvedModList []ResolvedMod

func (rm ResolvedModList) GetURLs() []string {
	urls := make([]string, 0, len(rm))
	for _, m := range rm {
		urls = append(urls, m.DownloadURL)
	}
	return urls
}

func (rm ResolvedModList) GetLatestVers() []string {
	vers := make([]string, 0, len(rm))
	for _, m := range rm {
		vers = append(vers, m.LatestVer)
	}
	return vers
}

func (rm ResolvedModList) GetLocalVers() []string {
	vers := make([]string, 0, len(rm))
	for _, m := range rm {
		vers = append(vers, m.LocalVer)
	}
	return vers
}

func ResolveMod(slug, mcVersion, loader string) (ResolvedMod, error) {
	project, err := modrinth.FetchModrinthProject(slug)
	if err != nil {
		return ResolvedMod{}, err
	}

	latestVer, err := modrinth.FetchLatestModrinthVersion(project.ID, mcVersion, loader)
	if err != nil {
		log.Printf("⚠ skipping %s: %v", slug, err)
		return ResolvedMod{}, nil
	}

	return ResolvedMod{
		Slug:        slug,
		ProjectID:   project.ID,
		FabricID:    project.Slug, //
		LatestVer:   latestVer.VersionNumber,
		DownloadURL: latestVer.Files[0].URL,
	}, nil
}

func NormalizeID(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), "-", "")
}
