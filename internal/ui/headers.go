package ui

import (
	"fmt"
	"strings"

	"github.com/w1lam/Packages/tui"
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
	fmt.Print("\n\n")
}
