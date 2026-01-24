package app

import (
	"log"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func refreshMetaData(path *paths.Paths, m *manifest.Manifest, metaD *meta.MetaData) {
	entryIDs := m.AllInstalledEntries()

	missing := metaD.FilterMissing(entryIDs)
	stale := metaD.FilterStale(24 * time.Hour)

	needFetch := append(missing, stale...)
	if len(needFetch) == 0 {
		return
	}

	fetched, err := meta.ResolveMetaData(needFetch)
	if err != nil {
		log.Printf("Metadata refresh failed: %v", err)
		return
	}

	metaD.Lock()
	metaD.Merge(fetched)
	metaD.Unlock()

	err = metaD.Save(path)
	if err != nil {
		log.Printf("Metadata save failed: %v", err)
	}
}
