// Package app provides application initialization, exit and state functions.
package app

import (
	"fmt"
	"log"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Packages/tui"
	minit "github.com/w1lam/Raw-Mod-Installer/internal/app/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

func Initialize() {
	tui.ClearScreenRaw()

	fmt.Println("* Starting up...")

	// Setting Program Exit Function
	menu.SetProgramExitFunc(func() {
		Exit()
	})

	// Start menu workers
	menu.StartWorkers(4)
	// Start input checker
	if err := menu.StartInput(); err != nil {
		log.Fatal(err)
	}

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
		m, err = manifest.BuildInitialManifest(state.ProgramVersion, path)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(" * Loading Meta Data...")
	metaD := meta.LoadMetaData(path)
	if metaD == nil {
		emptyMd := &meta.MetaData{
			SchemaVersion: 1,
			Mods:          make(map[string]meta.ModMetaData),
		}
		metaD = emptyMd
	}

	state.SetState(state.NewState(m, metaD))

	go refreshMetaData(path, m, metaD)

	minit.InitializeMenus(m)
}
