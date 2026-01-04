package ui

import (
	"fmt"
	"strings"

	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

func StartHeader(m *manifest.Manifest) {
	if m == nil {
		panic("StartHeader: Manifest is nil")
	}
	fmt.Print(strings.Repeat("━", config.Style.Width))
	fmt.Print("\n\n")

	tui.PrintCentered("┏━━━━━━━━━━━━━━━┓", config.Style.Width)
	fmt.Print("\n")
	tui.PrintCentered("┃ MOD INSTALLER ┃", config.Style.Width)
	fmt.Print("\n")
	tui.PrintCentered("┗━━━━━━━━━━━━━━━┛", config.Style.Width)
	fmt.Print("\n")
	tui.PrintCentered("Program Version: "+m.ProgramVersion, config.Style.Width)
	fmt.Print("\n")
	tui.PrintCentered("Mod List Version: "+m.ModList.InstalledVersion, config.Style.Width)
	fmt.Print("\n\n")
}

func PrintModInfoList(mods []manifest.ManifestMod) {
	for _, mod := range mods {
		fmt.Printf(" * %s\n  ", mod.Title)

		for _, cat := range mod.Categories {
			fmt.Printf(" 📂 %s ", cat)
		}

		fmt.Printf("\n   ✎  %s\n", mod.Description)

		if mod.Source != "" {
			fmt.Printf("   🔗 %s\n", mod.Source)
		}

		if mod.Wiki != "" {
			fmt.Printf("   📖 %s\n\n", mod.Wiki)
		} else {
			fmt.Printf("\n")
		}
	}
}
