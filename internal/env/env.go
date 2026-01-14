// Package env holds environtment structure
package env

import (
	"log"

	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

var (
	ProgramVersion string = "alpha 0.0.1"
	GlobalManifest *manifest.Manifest
)

type Env struct {
	Manifest *manifest.Manifest
	Config   *config.Config // so far unsused yet to be implemented
	Logger   *log.Logger    // so far unsused yet to be implemented
}
