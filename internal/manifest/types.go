package manifest

import (
	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Manifest is the manifest for all global information required by the program
type Manifest struct {
	SchemaVersion    int                   `json:"schemaVersion"`
	ProgramVersion   string                `json:"programVersion"`
	InstalledLoaders map[string]LoaderInfo `json:"installedLoader"`

	EnabledPackages map[packages.PackageType]string `json:"enabledPackages"`

	InstalledPackages map[packages.PackageType]map[string]InstalledPackage `json:"installedPackages"`

	Paths *paths.Paths `json:"-"`
}

// LoaderInfo is the information about a mod loader
type LoaderInfo struct {
	Loader        string `json:"loader"`
	McVersion     string `json:"mcVersion"`
	LoaderVersion string `json:"version"`
}

// InstalledPackage is an installed package which holds all information
type InstalledPackage struct {
	Name             string                  `json:"name"`
	Type             packages.PackageType    `json:"type"`
	ListSource       string                  `json:"listSource"`
	InstalledVersion string                  `json:"version"`
	McVersion        string                  `json:"mcVersion"`
	Loader           string                  `json:"loader"`
	Path             string                  `json:"path"`
	Hash             string                  `json:"hash"`
	Entries          map[string]PackageEntry `json:"installedEntries"`
}

// PackageEntry is a mod entry in the manifest that holds all information about an entry
type PackageEntry struct {
	ID               string `json:"id"` // id or slug
	FileName         string `json:"fileName"`
	Sha512           string `json:"sha512"`
	Sha1             string `json:"sha1,omitempty"`
	InstalledVersion string `json:"InstalledVersion"`
}

// Updates is all info on available updates
type Updates struct {
	ModListUpdate map[string]bool                   `json:"-"`
	ModUpdates    map[string][]modrinth.UpdateEntry `json:"-"`
}
