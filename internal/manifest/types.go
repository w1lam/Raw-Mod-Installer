package manifest

type Manifest struct {
	SchemaVersion  int    `json:"schemaVersion"`
	ProgramVersion string `json:"programVersion"`

	Minecraft MinecraftInfo `json:"minecraft"`
	ModList   ModListInfo   `json:"modList"`

	Mods map[string]ManifestMod `json:"mods"`
}

type MinecraftInfo struct {
	Version       string `json:"version"`
	Loader        string `json:"loader"`
	LoaderVersion string `json:"loaderVersion"`
}

type ModListInfo struct {
	Source           string `json:"source"`
	InstalledVersion string `json:"version"`
}

type ManifestMod struct {
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Categories       []string `json:"categories"`
	Description      string   `json:"description"`
	InstalledVersion string   `json:"InstalledVersion"`

	Source string `json:"source,omitempty"`
	Wiki   string `json:"wiki,omitempty"`
}

type ModList []ManifestMod
