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
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

func Initialize() *manifest.Manifest {
	tui.ClearScreenRaw()

	fmt.Println("* Starting up...")

	// Setting Program Exit Function
	menu.SetProgramExitFunc(func() {
		Exit()
	})

	menu.StartWorkers(4)

	// Setting width to terminal width
	GetSize, _ := ts.GetSize()
	config.Style.Width = GetSize.Col() + 1

	// Setting TUI Config Variables
	config.Style.Set()

	fmt.Println(" * Resolving Paths...")
	path, err := paths.Resolve()
	if err != nil {
		log.Fatal(err)
	}

	if err := filesystem.EnsureDirectories(path); err != nil {
		panic(" * Failed to create Raw Mod Installer Directories: " + err.Error())
	}

	fmt.Println(" * Loading Manifest...")
	m, err := manifest.Load(path)
	if err != nil {
		m, err = manifest.BuildInitialManifest(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(" * Loading Meta Data...")
	meta := resolve.LoadMetaData(path)
	if meta == nil {
		emptyMd := &resolve.MetaData{
			SchemaVersion: 1,
			Mods:          make(map[string]resolve.ModMetaData),
		}
		meta = emptyMd
	}

	env.GlobalMetaData = meta
	env.GlobalManifest = m

	go refreshMetaData(path, m, meta)

	InitializeMenus(m)

	return m
}
