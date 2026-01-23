package lists

import "github.com/w1lam/Packages/modrinth"

// AvailablePackages is a nested map with first key being the package type(name of subfolder inside packages folder) and second key is the package name
type AvailablePackages map[string]map[string]ResolvedPackage

type ResolvedPackage struct {
	Type string

	Name        string
	ListSource  string
	ListVersion string
	McVersion   string
	Loader      string
	Env         string
	Description string
	Hash        string // Sha512

	Entries []modrinth.ModrinthListEntry
}

type GithubContentResponse struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sha    string `json:"sha"`
	Size   int    `json:"size"`
	RawURL string `json:"download_url"`
	Type   string `json:"type"`
}
