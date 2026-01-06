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
	"github.com/w1lam/Raw-Mod-Installer/internal/output"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
	"github.com/w1lam/Raw-Mod-Installer/internal/ui"
)

func Initialize() ui.Context {
	tui.ClearScreenRaw()

	fmt.Println("* Starting up...")

	// Setting Program Exit Function
	menu.SetProgramExitFunc(func() {
		Exit()
	})

	// Setting width to terminal width
	GetSize, _ := ts.GetSize()
	config.Style.Width = GetSize.Col() + 1

	// Setting TUI Config Variables
	config.Style.Set()

	fmt.Println("* Resolving Paths...")
	path, err := paths.Resolve()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("* Creating Installer Directory...")
	if !utils.CheckFileExists(path.RawModInstallerDir) {
		err := os.MkdirAll(path.RawModInstallerDir, 0o755)
		if err != nil {
			panic("* Failed to create Raw Mod Installer Directory: " + err.Error())
		}
	}

	fmt.Println("* Loading Manifest...")
	m, err := manifest.Load(path.ManifestPath)
	if err != nil {
		m, err = manifest.BuildManifest(state.ProgramVersion)
		if err != nil {
			log.Fatal(err)
		}

		if err := manifest.Save(path.ManifestPath, m); err != nil {
			log.Fatal(err)
		}
	}

	state.GlobalManifest = m

	fmt.Println("* Writing README...")
	if err := output.WriteModInfoListREADME(path.RawModInstallerDir, state.GlobalManifest.ModsSlice()); err != nil {
		fmt.Printf("failed to write modlist README: %s", err)
	}

	tui.ClearScreenRaw()

	return ui.Context{
		Manifest: state.GlobalManifest,
		Paths:    path,
	}
}
