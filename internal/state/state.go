// Package state holds environtment structure
package state

import (
	"sync"

	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
)

var ProgramVersion string = "0.0.1"

var (
	globalState *State
	once        sync.Once
)

// State is the global state struct
type State struct {
	mu sync.RWMutex

	manifest *manifest.Manifest
	meta     *meta.MetaData

	availablePackages packages.AvailablePackages
	updates           manifest.Updates
}

// SetState sets the state
func SetState(s *State) {
	once.Do(func() {
		globalState = s
	})
}

// Get gets the state only read or edit inside Read or Write funcs
func Get() *State {
	if globalState == nil {
		panic("env.State not initialized")
	}
	return globalState
}

// NewState creates a new state
func NewState(m *manifest.Manifest, meta *meta.MetaData) *State {
	if m == nil || meta == nil {
		panic("NewState: manifest or meta is nil")
	}

	ap := make(packages.AvailablePackages)
	return &State{
		manifest:          m,
		meta:              meta,
		availablePackages: ap,
		updates:           manifest.Updates{},
	}
}

func SetAvailablePackages(pkgs packages.AvailablePackages) {
	globalState.mu.Lock()
	defer globalState.mu.Unlock()

	globalState.availablePackages = pkgs
}
