// Package app provides application initialization, exit and state functions.
package app

import (
	"fmt"
	"log"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Packages/tui"
	minit "github.com/w1lam/Raw-Mod-Installer/internal/app/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/errors"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
	"github.com/w1lam/Raw-Mod-Installer/internal/verify"
)

func Initialize() {
	tui.EnableANSI()
	tui.HideCursor()

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
		log.Fatal(fmt.Errorf("failed to start menu workers: %w", err))
	}

	fmt.Println(" * Resolving Paths...")
	path, err := paths.Resolve()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to resolve paths: %w", err))
	}

	if err := filesystem.EnsureDirectories(path); err != nil {
		log.Fatal(fmt.Errorf("failed to ensure directories: %w", err))
	}

	// Start error handler
	if err := errors.Start(path.LogPath); err != nil {
		log.Fatal(fmt.Errorf("failed to start error handler: %w", err))
	}

	fmt.Println(" * Loading Manifest...")
	m, err := manifest.Load(path)
	if err != nil {
		m, err = manifest.BuildInitialManifest(state.ProgramVersion, path)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to build initial manifest: %w", err))
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

	// Sets global state
	state.SetState(state.NewState(m, metaD))

	// Backup if first run
	if !m.Initialized {
		go func() {
			res, err := services.PerformInitialBackup(path)
			if err != nil {
				errors.Report("startup.backup", err)
				return
			}

			state.Get().Write(func(s *state.State) error {
				if err := services.ApplyInitialBackup(s.Manifest(), res, path); err != nil {
					errors.Report("startup.backup.apply", err)
				}
				return nil
			})
		}()
	}

	// Verify packages
	verify.VerifyAndReconcile(path)

	// refresh meta data of installed package entries
	go refreshMetaData(path, m, metaD)

	// Initialize menus
	minit.InitializeMenus(m)
}
