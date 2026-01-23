package actions

import (
	"fmt"
	"time"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Packages/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

func InstallModPackAction(modPackName string) menu.Action {
	gState := state.Get()

	plan := installer.InstallPlan{}
	gState.Read(func(s *state.State) {
		plan = installer.InstallPlan{
			Intent:           installer.Install,
			RequestedPackage: s.Packages()["modpacks"][modPackName],
			EnsureFabric:     true,
			BackupPolicy:     filesystem.BackupOnce,
			EnableAfter:      true,
		}
	})

	return menu.Action{
		Function: func() error {
			err := installer.PackageInstaller(plan)
			if err != nil {
				return err
			}
			return nil
		},
		WrapUp: func(err error) {
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
		WrapUp: func(err error) {},
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
		WrapUp: func(err error) {},
	}
}
