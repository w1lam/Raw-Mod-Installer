// Package state hold global state values
package state

import "github.com/w1lam/Raw-Mod-Installer/internal/manifest"

var (
	ProgramVersion string = "ALPHA 0.0.1"
	GlobalManifest *manifest.Manifest
)
