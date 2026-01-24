// Package actions holds actions
package actions

import (
	"fmt"
	"time"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Packages/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

func InstallModPackAction(pkg packages.Pkg) menu.Action {
	gState := state.Get()

	plan := installer.InstallPlan{}
	gState.Read(func(s *state.State) {
		plan = installer.InstallPlan{
			RequestedPackage: s.AvailablePackages()[pkg.Type][pkg.Name],
			BackupPolicy:     services.BackupOnce,
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
			if err == nil {
				fmt.Printf("\n* %s Installation Complete!\n", pkg.Name)
				time.Sleep(time.Second * 3)
				tui.ClearScreenRaw()
				menu.RenderCurrentMenu()
			} else {
				fmt.Printf("\n* %s Installation Failed!\n", pkg.Name)
				time.Sleep(time.Second * 3)
				tui.ClearScreenRaw()
				menu.RenderCurrentMenu()
			}
		},
		Async: true,
	}
}

func EnablePackageAction(pkg packages.Pkg) menu.Action {
	return menu.Action{
		Function: func() error {
			err := services.EnablePackage(pkg)
			if err != nil {
				return err
			}
			return nil
		},
		WrapUp: func(err error) {},
		Async:  true,
	}
}

func DisablePackageAction(pkgType packages.PackageType) menu.Action {
	gState := state.Get()

	var pkg packages.Pkg
	gState.Read(func(s *state.State) {
		pkg = packages.Pkg{
			Name: s.Manifest().EnabledPackages[pkgType],
			Type: pkgType,
		}
	})

	return menu.Action{
		Function: func() error {
			err := services.DisablePackage(pkg)
			if err != nil {
				return err
			}
			return nil
		},
		WrapUp: func(err error) {},
		Async:  true,
	}
}
