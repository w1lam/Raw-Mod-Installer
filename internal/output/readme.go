// Package output contains functions for generating output files.
package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

// WriteModPackInfoList writes a README.txt file summarizing the mod information in the specified output path.
func WriteModPackInfoList(mods []manifest.ManifestMod, m *manifest.Manifest) error {
	fmt.Printf("NEEDS TO BE CHANGED TO MODPACK SPECIFIC AND STORED IN MODPACKS???")
	metadata := resolve.LoadMetaData(m.Paths)
	if metadata == nil {
		var slugs []string
		for _, m := range mods {
			slugs = append(slugs, m.Slug)
		}

		data, err := resolve.ResolveMetaData(slugs)
		if err != nil {
			return err
		}

		if data != nil {
			metadata = data
			return nil
		}
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	content := func(mods []manifest.ManifestMod) string {
		var b strings.Builder

		b.WriteString("MODPACK README\n")
		b.WriteString("==============\n\n")
		b.WriteString(fmt.Sprintf("Total Mods: %d\n\n", len(mods)))

		for _, mod := range mods {
			b.WriteString(strings.Repeat("-", len(metadata.Mods[mod.Slug].Title)+4) + "\n")
			b.WriteString("* " + metadata.Mods[mod.Slug].Title + " *\n")
			b.WriteString(strings.Repeat("-", len(metadata.Mods[mod.Slug].Title)+4) + "\n")

			if len(metadata.Mods[mod.Slug].Categories) > 0 {
				b.WriteString("📂 " + strings.Join(metadata.Mods[mod.Slug].Categories, ", 📂 ") + "\n")
			}

			if metadata.Mods[mod.Slug].Description != "" {
				b.WriteString("\n" + metadata.Mods[mod.Slug].Description + "\n")
			}

			if metadata.Mods[mod.Slug].Source != "" {
				b.WriteString("\n🔗 " + metadata.Mods[mod.Slug].Source + "\n")
			}

			if metadata.Mods[mod.Slug].Wiki != "" {
				b.WriteString("📖 " + metadata.Mods[mod.Slug].Wiki + "\n\n")
			}

			b.WriteString("\n\n")
		}

		return b.String()
	}(mods)

	err1 := os.WriteFile(filepath.Join(m.Paths.ProgramFilesDir, "README.md"), []byte(content), 0o644)
	if err1 != nil {
		return err1
	}

	return nil
}
