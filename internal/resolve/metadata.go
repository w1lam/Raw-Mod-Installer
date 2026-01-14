package resolve

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// ModMetaData is the metadata of a mod
type ModMetaData struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
	Wiki        string   `json:"wiki,omitempty"`
	Source      string   `json:"source,omitempty"`
}

// MetaData is the metadata map of mod metadata
type MetaData struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Mods          map[string]ModMetaData `json:"mods"`
}

// LoadMetaData loaeds metadata
func LoadMetaData(path *paths.Paths) *MetaData {
	data, err := os.ReadFile(path.MetaDataPath)
	if err != nil {
		return nil
	}

	var m MetaData
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}

	return &m
}

// Save metadata
func (md *MetaData) Save(path *paths.Paths) error {
	tmp := path.MetaDataPath + ".tmp"

	data, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall metadata: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write metadata temp file: %s", err)
	}

	return os.Rename(tmp, path.MetaDataPath)
}

// ResolveMetaData resolved metadata of given slugs
func ResolveMetaData(slugs []string) (*MetaData, error) {
	mrProj, err := modrinth.BatchFetchModrinthProjects(slugs)
	if err != nil {
		return nil, err
	}

	out := MetaData{
		SchemaVersion: 1,
		Mods:          make(map[string]ModMetaData),
	}

	for _, p := range mrProj {
		out.Mods[p.Slug] = ModMetaData{
			Slug:        p.Slug,
			Title:       p.Title,
			Categories:  p.Categories,
			Description: p.Description,
			Wiki:        p.Wiki,
			Source:      p.Source,
		}
	}

	return &out, nil
}
