// Package env holds environtment structure
package env

import (
	"sync"

	"github.com/w1lam/Raw-Mod-Installer/internal/lists"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

var (
	ProgramVersion    string = "0.0.1"
	GlobalManifest    *manifest.Manifest
	GlobalMetaData    *resolve.MetaData
	AvailableModPacks map[string]lists.ResolvedModPack
	Updates           manifest.Updates
	ManMu             sync.Mutex
)
