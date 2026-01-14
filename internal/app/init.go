// Package app provides application initialization, exit and state functions.
package app

import (
	"fmt"
	"log"

	"github.com/olekukonko/ts"
	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Packages/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/env"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
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

	fmt.Println("* Creating Installer Directories...")
	if err := filesystem.EnsureDirectories(path); err != nil {
		panic("* Failed to create Raw Mod Installer Directories: " + err.Error())
	}

	fmt.Println("* Loading Manifest...")
	m, err := manifest.Load(path)
	if err != nil {
		m, err = manifest.BuildInitialManifest(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	env.GlobalManifest = m

	InitializeMenus(m)

	return m
}
