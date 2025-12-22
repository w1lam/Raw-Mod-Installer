package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/modinfo"
)

// WriteModInfoListREADME writes a README.txt file summarizing the mod information in the specified output path.
func WriteModInfoListREADME(outPath string, modInfoList modinfo.ModInfoList) error {
	content := func(mods modinfo.ModInfoList) string {
		var b strings.Builder

		b.WriteString("MODPACK README\n")
		b.WriteString("==============\n\n")
		b.WriteString(fmt.Sprintf("Total Mods: %d\n\n", len(mods)))

		for _, mod := range mods {
			b.WriteString(strings.Repeat("-", len(mod.Title)+4) + "\n")
			b.WriteString("* " + mod.Title + " *\n")
			b.WriteString(strings.Repeat("-", len(mod.Title)+4) + "\n")

			if len(mod.Category) > 0 {
				b.WriteString("📂 " + strings.Join(mod.Category, ", 📂 ") + "\n")
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
	}(modInfoList)

	err1 := os.WriteFile(filepath.Join(outPath, "README.md"), []byte(content), 0o644)
	if err1 != nil {
		return err1
	}

	return nil
}
