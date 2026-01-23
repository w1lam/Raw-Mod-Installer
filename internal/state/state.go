// Package state holds environtment structure
package state

import (
	"sync"

	"github.com/w1lam/Raw-Mod-Installer/internal/lists"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
)

var ProgramVersion string = "0.0.1"

var (
	globalState *State
	once        sync.Once
)

type State struct {
	mu sync.RWMutex

	manifest          *manifest.Manifest
	meta              *meta.MetaData
	availablePackages lists.AvailablePackages
	updates           manifest.Updates
}

func (s *State) Packages() lists.AvailablePackages {
	return s.availablePackages
}

func (s *State) MetaData() *meta.MetaData {
	return s.meta
}

func (s *State) Manifest() *manifest.Manifest {
	return s.manifest
}

func SetState(s *State) {
	once.Do(func() {
		globalState = s
	})
}

func Get() *State {
	if globalState == nil {
		panic("env.State not initialized")
	}
	return globalState
}

func (s *State) Read(fn func(*State)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s)
}

func (s *State) Write(fn func(*State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fn(s)
}

func NewState(m *manifest.Manifest, meta *meta.MetaData) *State {
	if m == nil || meta == nil {
		panic("NewState: manifest or meta is nil")
	}

	return &State{
		manifest:          m,
		meta:              meta,
		availablePackages: make(lists.AvailablePackages),
		updates:           manifest.Updates{},
	}
}
