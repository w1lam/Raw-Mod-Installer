package meta

import (
	"time"

	"github.com/w1lam/Packages/modrinth"
)

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
