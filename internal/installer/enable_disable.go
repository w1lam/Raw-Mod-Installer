package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// EnableModPack enables the specified mod pack
func EnableModPack(name string) error {
	store := state.Get()
	var path *paths.Paths

	// READ MANIFEST
	store.Read(func(s *state.State) {
		path = s.Manifest().Paths
	})

	src := filepath.Join(path.ModPacksDir, name)
	dst := path.ModsDir

	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("failed to clear mods dir: %w", err)
	}

	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("failed to enable modpack \"%s\": %w", name, err)
	}

	// WRITE MANIFEST
	return store.Write(func(s *state.State) error {
		if s.Manifest().EnabledModPack == name {
			return nil
		}

		s.Manifest().EnabledModPack = name
		return s.Manifest().Save()
	})
}

// DisableModPack disables the currently enabled modpack
func DisableModPack() error {
	store := state.Get()

	var (
		enabled string
		path    *paths.Paths
	)

	store.Read(func(s *state.State) {
		enabled = s.Manifest().EnabledModPack
		path = s.Manifest().Paths
	})

	if enabled == "" {
		return nil
	}

	src := path.ModsDir
	dst := filepath.Join(path.ModPacksDir, enabled)

	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("failed to clear target modpack dir: %w", err)
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}

	return state.Get().Write(func(s *state.State) error {
		if s.Manifest().EnabledModPack == enabled {
			s.Manifest().EnabledModPack = ""
			return s.Manifest().Save()
		}
		return nil
	})
}
