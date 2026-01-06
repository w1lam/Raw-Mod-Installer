package mods

import (
	"fmt"
	"os"

	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// EnableModpack enables mods
func EnableModpack(path *paths.Paths) error {
	if entries, err := os.ReadDir(path.ModsDir); err != nil {
		if len(entries) == 0 {
			err := os.RemoveAll(path.ModsDir)
			if err != nil {
				return fmt.Errorf("failed to remove empty mod dir: %s", err)
			}
		}
	}

	err := Backup(path)
	if err != nil {
		return err
	}

	if err := os.Rename(path.UnloadedModsDir, path.ModsDir); err != nil {
		return fmt.Errorf("failed to move mods to final directory: %w", err)
	}

	return nil
}
