package ui

import (
	"fmt"

	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/install"
)

func BuildMainMenu(ctx *Context) *menu.Menu {
	m := menu.NewMenu("Start Menu", "Main Menu", MainMenuID)

	m.SetRender(func() {
		tui.ClearScreenRaw()
		StartHeader(ctx.Manifest)

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
					err := install.Install(ctx.Manifest, ctx.Paths)
					if err != nil {
						return err
					}
					return nil
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
					return install.Update(ctx.Manifest)
				},
				'u',
				"update",
			)

		default:
			m.AddButton(
				"[R] Reinstall",
				"Reinstall the modpack",
				func() error {
					return install.CleanInstall(ctx.Manifest, ctx.Paths)
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
	})

	return m
}
