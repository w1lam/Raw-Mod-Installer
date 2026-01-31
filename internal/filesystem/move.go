package filesystem

import (
	"fmt"
	"os"
)

// SwapDirs swaps two dirs
func SwapDirs(src, dst, backup string) error {
	// sanity checks
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("source missing: %w", err)
	}

	_ = os.RemoveAll(backup)

	// Step 1: move current active out of the way
	if _, err := os.Stat(dst); err == nil {
		if err := os.Rename(dst, backup); err != nil {
			return fmt.Errorf("failed to backup active dir: %w", err)
		}
	}

	// Step 2: move new package into place
	if err := os.Rename(src, dst); err != nil {
		// rollback
		_ = os.Rename(backup, dst)
		return fmt.Errorf("failed to activate package: %w", err)
	}

	// Step 3: cleanup backup
	_ = os.RemoveAll(backup)
	return nil
}
