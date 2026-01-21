package app

import (
	"log"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

func refreshMetaData(path *paths.Paths, m *manifest.Manifest, meta *resolve.MetaData) {
	slugs := m.AllInstalledModSlugs()

	missing := meta.FilterMissing(slugs)
	stale := meta.FilterStale(24 * time.Hour)

	needFetch := append(missing, stale...)
	if len(needFetch) == 0 {
		return
	}

	fetched, err := resolve.ResolveMetaData(needFetch)
	if err != nil {
		log.Printf("Metadata refresh failed: %v", err)
		return
	}

	meta.Lock()
	meta.Merge(fetched)
	meta.Unlock()

	err = meta.Save(path)
	if err != nil {
		log.Printf("Metadata save failed: %v", err)
	}
}
