package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

// EnableModPack enables the specified mod pack
func EnableModPack(modPack manifest.InstalledModPack, m *manifest.Manifest) (*manifest.Manifest, error) {
	if m.EnabledModPack.Name == modPack.Name {
		fmt.Printf("Mod Pack already enabled\n")
		return m, nil
	}

	if me, err := DisableModPack(m); err != nil {
		return m, fmt.Errorf("failed to disable current mod pack: %s", err)
	} else {
		m = me
	}

	err := os.Rename(filepath.Join(m.Paths.ModPacksDir, modPack.Name), m.Paths.ModsDir)
	if err != nil {
		return m, err
	}

	m.EnabledModPack = &modPack
	return m, m.Save()
}

// DisableModPack disables the currently enabled modpack
func DisableModPack(m *manifest.Manifest) (*manifest.Manifest, error) {
	if !utils.CheckFileExists(m.Paths.ModsDir) || m.EnabledModPack == nil {
		return m, nil
	}

	err := os.Rename(m.Paths.ModsDir, filepath.Join(m.Paths.ModPacksDir, m.EnabledModPack.Name))
	if err != nil {
		return m, err
	}

	m.EnabledModPack = nil
	return m, m.Save()
}
