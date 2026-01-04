// Package app provides application initialization, exit and state functions.
package app

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/ts"
	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/ui"
)

func Initialize() ui.Context {
	tui.ClearScreenRaw()

	fmt.Printf("Starting Up...\n")

	// Setting Program Exit Function
	menu.SetProgramExitFunc(func() {
		Exit()
	})

	// Setting width to terminal width
	GetSize, _ := ts.GetSize()
	config.Style.Width = GetSize.Col() + 1

	// Setting TUI Config Variables
	config.Style.Set()

	path, err := paths.Resolve()
	if err != nil {
		log.Fatal(err)
	}

	if !utils.CheckFileExists(path.RawModInstallerDir) {
		err := os.MkdirAll(path.RawModInstallerDir, 0o755)
		if err != nil {
			panic("Failed to create Raw Mod Installer Directory: " + err.Error())
		}
	}

	if !utils.CheckFileExists(path.ModsDir) {
		err := os.MkdirAll(path.ModsDir, 0o755)
		if err != nil {
			fmt.Printf("Failed to create Mods Directory: %s\n", err)
		}
	}

	m, err := manifest.Load(path.ManifestPath)
	if err != nil {
		fmt.Printf("No Manifest Found, Building from Scratch...\n")

		m, err := manifest.BuildManifest("0.0.1")
		if err != nil {
			log.Fatal(err)
		}

		if err := manifest.Save(path.ManifestPath, m); err != nil {
			log.Fatal(err)
		}
	}

	GlobalManifest = m

	return ui.Context{
		Manifest: GlobalManifest,
		Paths:    path,
	}
}
