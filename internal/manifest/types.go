package manifest

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Manifest is the manifest for all global information required by the program
type Manifest struct {
	SchemaVersion  int           `json:"schemaVersion"`
	ProgramVersion string        `json:"programVersion"`
	Minecraft      MinecraftInfo `json:"minecraft"`
	Paths          *paths.Paths

	EnabledModPack *InstalledModPack `json:"enabledModpack"`

	AvailableModPacks map[string]modpack.ResolvedModPackList
	InstalledModPacks map[string]InstalledModPack `json:"installedModPacks"`
}

// MinecraftInfo is the information about minecraft
type MinecraftInfo struct {
	Version       string `json:"version"`
	Loader        string `json:"loader"`
	LoaderVersion string `json:"loaderVersion"`
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
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Categories       []string `json:"categories"`
	Description      string   `json:"description"`
	InstalledVersion string   `json:"InstalledVersion"`

	Source string `json:"source,omitempty"`
	Wiki   string `json:"wiki,omitempty"`
}
