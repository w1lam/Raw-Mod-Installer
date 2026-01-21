package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/env"
)

// EnableModPack enables the specified mod pack
func EnableModPack(name string) error {
	env.ManMu.Lock()
	defer env.ManMu.Unlock()

	m := env.GlobalManifest

	if m.EnabledModPack == name {
		fmt.Printf("Mod Pack already enabled\n")
		return nil
	}

	if err := DisableModPack(); err != nil {
		return fmt.Errorf("failed to disable current mod pack: %s", err)
	}

	src := filepath.Join(m.Paths.ModPacksDir, name)
	dst := m.Paths.ModsDir

	_ = os.RemoveAll(dst)

	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("failed to enable modpack \"%s\": %w", name, err)
	}

	m.EnabledModPack = name

	return m.Save()
}

// DisableModPack disables the currently enabled modpack
func DisableModPack() error {
	env.ManMu.Lock()
	defer env.ManMu.Unlock()

	m := env.GlobalManifest

	if m.EnabledModPack == "" {
		return nil
	}

	src := m.Paths.ModsDir
	dst := filepath.Join(m.Paths.ModPacksDir, m.EnabledModPack)

	_ = os.RemoveAll(dst)

	if err := os.Rename(src, dst); err != nil {
		return err
	}

	m.EnabledModPack = ""

	return m.Save()
}
