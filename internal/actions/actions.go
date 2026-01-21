package actions

import (
	"fmt"
	"time"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Packages/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/env"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
)

func InstallModPackAction(modPackName string) menu.Action {
	if env.GlobalManifest == nil {
		panic("InstallModPackAction: GlobalManifest is nil")
	}

	if env.AvailableModPacks[modPackName].Name != modPackName {
		panic("no mod pack with specified name: " + modPackName)
	}

	plan := installer.InstallPlan{
		Intent:           installer.Install,
		RequestedModPack: env.AvailableModPacks[modPackName],
		EnsureFabric:     true,
		BackupPolicy:     filesystem.BackupOnce,
		EnableAfter:      true,
	}

	return menu.Action{
		Function: func() error {
			m, err := installer.InstallModPack(env.GlobalManifest, plan)
			if err != nil {
				return err
			}
			env.GlobalManifest = m
			return env.GlobalManifest.Save()
		},
		WrapUp: func() {
			fmt.Printf("\n* %s Installation Complete!\n", modPackName)
			time.Sleep(time.Second * 3)
			tui.ClearScreenRaw()
			menu.RenderCurrentMenu()
		},
	}
}

func EnableModPackAction(modPackName string) menu.Action {
	return menu.Action{
		Function: func() error {
			err := installer.EnableModPack(modPackName)
			if err != nil {
				return err
			}
			return nil
		},
		WrapUp: func() {},
	}
}

func DisableModPackAction() menu.Action {
	return menu.Action{
		Function: func() error {
			err := installer.DisableModPack()
			if err != nil {
				return err
			}
			return nil
		},
		WrapUp: func() {},
	}
}
