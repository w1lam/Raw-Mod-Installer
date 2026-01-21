package resolve

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

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
	UpdatedAt   time.Time
}

// MetaData is the metadata map of mod metadata
type MetaData struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Mods          map[string]ModMetaData `json:"mods"`
	sync.Mutex
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
			UpdatedAt:   time.Now(),
		}
	}

	return &out, nil
}

func (md *MetaData) FilterMissing(slugs []string) []string {
	var missing []string
	for _, slug := range slugs {
		if _, ok := md.Mods[slug]; !ok {
			missing = append(missing, slug)
		}
	}
	return missing
}

func (md *MetaData) FilterStale(threshold time.Duration) []string {
	cutoff := time.Now().Add(-threshold)
	var stale []string

	for slug, meta := range md.Mods {
		if meta.UpdatedAt.Before(cutoff) {
			stale = append(stale, slug)
		}
	}
	return stale
}

func (md *MetaData) Merge(nmd *MetaData) {
	if md.Mods == nil {
		md.Mods = make(map[string]ModMetaData)
	}

	for slug, incoming := range nmd.Mods {
		existing, exists := md.Mods[slug]
		if !exists {
			md.Mods[slug] = incoming
			continue
		}

		merged := existing

		if incoming.Title != "" {
			merged.Title = incoming.Title
		}
		if incoming.Description != "" {
			merged.Description = incoming.Description
		}
		if len(incoming.Categories) > 0 {
			merged.Categories = incoming.Categories
		}
		if incoming.Wiki != "" {
			merged.Wiki = incoming.Wiki
		}
		if incoming.Source != "" {
			merged.Source = incoming.Source
		}

		merged.UpdatedAt = incoming.UpdatedAt

		md.Mods[slug] = merged
	}
}
