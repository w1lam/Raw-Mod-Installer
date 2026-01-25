// Package packages handles packages
package packages

import "github.com/w1lam/Packages/modrinth"

// AvailablePackages is a nested map with first key being the package type(name of subfolder inside packages folder) and second key is the package name
type AvailablePackages map[PackageType]map[string]ResolvedPackage

type GithubContentResponse struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sha    string `json:"sha"`
	Size   int    `json:"size"`
	RawURL string `json:"download_url"`
	Type   string `json:"type"`
}

// Pkg is a small pacakge struct used for passing around packages
type Pkg struct {
	Name string
	Type PackageType
}

// ResolvedPackage is a resolved package
type ResolvedPackage struct {
	Type PackageType `json:"pkgType"`

	Name        string `json:"name"`
	ListSource  string
	ListVersion string `json:"listVersion"`
	McVersion   string `json:"mcVersion"`
	Loader      string `json:"loader"`
	Env         string `json:"env"`
	Description string `json:"description"`
	Hash        string `json:"hash"` // sha512

	Entries []modrinth.ModrinthListEntry `json:"entries"`
}
