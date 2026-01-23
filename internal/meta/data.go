// Package meta holds metadata
package meta

import (
	"sync"
	"time"
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
