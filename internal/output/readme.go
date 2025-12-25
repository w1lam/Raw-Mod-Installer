// Package output contains functions for generating output files.
package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

// WriteModInfoListREADME writes a README.txt file summarizing the mod information in the specified output path.
func WriteModInfoListREADME(outPath string, mods []manifest.ManifestMod) error {
	content := func(mods []manifest.ManifestMod) string {
		var b strings.Builder

		b.WriteString("MODPACK README\n")
		b.WriteString("==============\n\n")
		b.WriteString(fmt.Sprintf("Total Mods: %d\n\n", len(mods)))

		for _, mod := range mods {
			b.WriteString(strings.Repeat("-", len(mod.Title)+4) + "\n")
			b.WriteString("* " + mod.Title + " *\n")
			b.WriteString(strings.Repeat("-", len(mod.Title)+4) + "\n")

			if len(mod.Categories) > 0 {
				b.WriteString("📂 " + strings.Join(mod.Categories, ", 📂 ") + "\n")
			}

			if mod.Description != "" {
				b.WriteString("\n" + mod.Description + "\n")
			}

			if mod.Source != "" {
				b.WriteString("\n🔗 " + mod.Source + "\n")
			}

			if mod.Wiki != "" {
				b.WriteString("📖 " + mod.Wiki + "\n\n")
			}

			b.WriteString("\n\n")
		}

		return b.String()
	}(mods)

	err1 := os.WriteFile(filepath.Join(outPath, "README.md"), []byte(content), 0o644)
	if err1 != nil {
		return err1
	}

	return nil
}
