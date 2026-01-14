package manifest

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Manifest is the manifest for all global information required by the program
type Manifest struct {
	SchemaVersion   int        `json:"schemaVersion"`
	ProgramVersion  string     `json:"programVersion"`
	InstalledLoader LoaderInfo `json:"installedLoader"`

	EnabledModPack string `json:"enabledModpack"`

	InstalledModPacks map[string]InstalledModPack `json:"installedModPacks"`
	Paths             *paths.Paths                `json:"-"`
}

// LoaderInfo is the information about the mod loader
type LoaderInfo struct {
	Loader  string `json:"loader"`
	Version string `json:"version"`
}

// InstalledModPack is an installed mod pack which holds all information about the mod pack, including all mods in form of map of ManifestMod with mods slug as key
type InstalledModPack struct {
	Name             string                 `json:"name"`
	ListSource       string                 `json:"listSource"`
	InstalledVersion string                 `json:"version"`
	McVersion        string                 `json:"mcVersion"`
	Loader           string                 `json:"loader"`
	Mods             map[string]ManifestMod `json:"installedMods"`
}

// ManifestMod is a mod entry in the manifest that holds all information about a mod
type ManifestMod struct {
	Slug             string `json:"slug"`
	FileName         string `json:"fileName"`
	InstalledVersion string `json:"InstalledVersion"`
}
