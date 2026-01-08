package ui

import (
	"fmt"

	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/install"
)

func BuildMainMenu(ctx *Context) *menu.Menu {
	m := menu.NewMenu("Main Menu", "Main Menu", MainMenuID)

	m.SetOnEnter(func() {
		state, err := ComputeMainMenuState(ctx.Manifest)
		if err != nil {
			fmt.Println("⚠ Failed to check updates")
			return
		}

		m.ClearButtons()

		switch {
		case !state.HasModsInstalled:
			m.AddButton(
				"[I] Install",
				"Install the modpack",
				func() error {
					m, err := install.ExecutePlan(ctx.Manifest, ctx.Paths, install.InstallPlan{
						Intent:       install.IntentInstall,
						EnsureFabric: true,
						BackupPolicy: install.BackupOnce,
						EnableAfter:  true,
					})
					if err != nil {
						return err
					}
					return m.Save(ctx.Paths.ManifestPath)
				},
				'i',
				"install",
			)

		case state.ModlistUpdateAvailable || state.FabricUpdateAvailable || len(state.ModUpdates) > 0:
			m.AddButton(
				"[U] Update",
				fmt.Sprintf(
					"Updates available (%d mods)",
					len(state.ModUpdates),
				),
				func() error {
					m, err := install.ExecutePlan(ctx.Manifest, ctx.Paths, install.InstallPlan{
						Intent:       install.IntentUpdate,
						EnsureFabric: true,
						BackupPolicy: install.BackupIfExists,
						EnableAfter:  true,
					})
					if err != nil {
						return err
					}

					return m.Save(ctx.Paths.ManifestPath)
				},
				'u',
				"update",
			)

		default:
			m.AddButton(
				"[R] Reinstall",
				"Reinstall the modpack",
				func() error {
					m, err := install.ExecutePlan(ctx.Manifest, ctx.Paths, install.InstallPlan{
						Intent:       install.IntentReinstall,
						EnsureFabric: false,
						BackupPolicy: install.BackupNever,
						EnableAfter:  true,
					})
					if err != nil {
						return err
					}
					return m.Save(ctx.Paths.ManifestPath)
				},
				'r',
				"reinstall",
			)
		}

		m.AddButton(
			"[H] Help / Info",
			"Show mod list info",
			func() error {
				return menu.SetCurrent(InfoMenuID)
			},
			'h',
			"help",
		)

		m.SetRender(func() {
			tui.ClearScreenRaw()
			StartHeader(ctx.Manifest)
		})
	})

	return m
}
