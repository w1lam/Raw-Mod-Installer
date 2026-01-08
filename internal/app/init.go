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
)

func Initialize() *manifest.Manifest {
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
	if !utils.CheckFileExists(path.ProgramFilesDir) {
		err := os.MkdirAll(path.ProgramFilesDir, 0o755)
		if err != nil {
			panic("* Failed to create Raw Mod Installer Directory: " + err.Error())
		}
	}

	fmt.Println("* Loading Manifest...")
	m, err := manifest.Load(path.ManifestPath)
	if err != nil {
		m, err = manifest.BuildInitialManifest(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	state.GlobalManifest = m

	fmt.Println("* Writing README...")
	if err := output.WriteModInfoListREADME(path.ProgramFilesDir, state.GlobalManifest.ModsSlice()); err != nil {
		fmt.Printf("failed to write modlist README: %s", err)
	}

	tui.ClearScreenRaw()

	return m
}
