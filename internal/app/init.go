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

// InitCore core functionality initiation
func InitCore() error {
	path, err := paths.Resolve()
	if err != nil {
		return err
	}

	if err := filesystem.EnsureDirectories(path); err != nil {
		return err
	}

	// Start error handler
	if err := errors.Start(path.LogPath); err != nil {
		log.Fatal(fmt.Errorf("failed to start error handler: %w", err))
	}

	m, err := manifest.Load(path)
	if err != nil {
		m, err = manifest.BuildInitialManifest(state.ProgramVersion, path)
		if err != nil {
			return err
		}
	}

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

	// Verify packages
	verify.VerifyAndReconcile(path)

	return nil
}

type InitMode int

const (
	CLI InitMode = iota
	TUI
)

// InitTUI initializes the tui
func InitTUI() {
	tui.EnableANSI()
	tui.HideCursor()
	tui.ClearScreenRaw()

	fmt.Println("* Starting up...")
	m := state.Get().Manifest()

	// Setting Program Exit Function
	menu.SetProgramExitFunc(Exit)

	// Start menu workers
	menu.StartWorkers(4)

	// Start input checker
	if err := menu.StartInput(); err != nil {
		log.Fatal(fmt.Errorf("failed to start menu workers: %w", err))
	}

	// Backup if first run
	if !m.Initialized {
		go func() {
			res, err := services.PerformInitialBackup(m.Paths)
			if err != nil {
				errors.Report("startup.backup", err)
				return
			}

			if err := state.Get().Write(func(s *state.State) error {
				if err := services.ApplyInitialBackup(s.Manifest(), res, s.Manifest().Paths); err != nil {
					errors.Report("startup.backup.apply", err)
				}
				return nil
			}); err != nil {
				log.Fatal(fmt.Errorf("OH NOOOO ): failed to write initial backup to state: %w", err))
			}
		}()
	}

	// refresh meta data of installed package entries
	go refreshMetaData(state.Get().Manifest().Paths, m, state.Get().MetaData())

	// Initialize menus
	minit.InitializeMenus(state.Get().Manifest())
}
