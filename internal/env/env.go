// Package env holds environtment structure
package env

import (
	"log"
	"sync"

	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/modpack"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

var (
	ProgramVersion    string = "alpha 0.0.1"
	GlobalManifest    *manifest.Manifest
	GlobalMetaData    *resolve.MetaData
	AvailableModPacks map[string]modpack.ResolvedModPack
	ManMu             sync.Mutex
)

type Env struct {
	Manifest *manifest.Manifest
	Config   *config.Config // so far unsused yet to be implemented
	Logger   *log.Logger    // so far unsused yet to be implemented
}
