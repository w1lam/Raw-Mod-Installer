package state

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// AvailablePackages safe packages accessor
func (s *State) AvailablePackages() packages.AvailablePackages {
	return s.availablePackages
}

// MetaData safe meta data accessor
func (s *State) MetaData() *meta.MetaData {
	return s.meta
}

// Manifest safe manifest accessor
func (s *State) Manifest() *manifest.Manifest {
	return s.manifest
}

// Paths small paths reader DO NOT USE INSIDE READ OF WRITE
func Paths() *paths.Paths {
	globalState.mu.RLock()
	defer globalState.mu.RUnlock()
	return globalState.manifest.Paths
}
