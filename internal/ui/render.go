package ui

import (
	"fmt"
	"strings"

	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

func MainMenu(programInfo manifest.ProgramInfo, width int) {
	fmt.Print(strings.Repeat("━", width))
	fmt.Print("\n\n")

	tui.PrintCentered("┏━━━━━━━━━━━━━━━┓", width)
	fmt.Print("\n")
	tui.PrintCentered("┃ MOD INSTALLER ┃", width)
	fmt.Print("\n")
	tui.PrintCentered("┗━━━━━━━━━━━━━━━┛", width)
	fmt.Print("\n")
	tui.PrintCentered("Program Version: "+programInfo.ProgramVersion, width)
	fmt.Print("\n")
	tui.PrintCentered("Mod List Version: "+programInfo.ModListVersion, width)
	fmt.Print("\n\n")
}

func PrintModInfoList(modInfoList modrinth.ModInfoList) {
	for _, modInfo := range modInfoList {
		fmt.Printf(" * %s\n  ", modInfo.Title)

		for _, cat := range modInfo.Category {
			fmt.Printf(" 📂 %s ", cat)
		}

		fmt.Printf("\n   ✎  %s\n", modInfo.Description)

		if modInfo.Source != "" {
			fmt.Printf("   🔗 %s\n", modInfo.Source)
		}

		if modInfo.Wiki != "" {
			fmt.Printf("   📖 %s\n\n", modInfo.Wiki)
		} else {
			fmt.Printf("\n")
		}
	}
}
