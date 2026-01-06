// Package ui provides functions to handle menu states and user input.
package ui

import (
	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/install"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Menu IDs
const (
	MainMenuID menu.MenuID = iota
	InfoMenuID
)

type Context struct {
	Manifest *manifest.Manifest
	Paths    *paths.Paths
}

// InitializeMenus initializes the default menus for the program
func InitializeMenus(ctx Context) (*menu.Menu, *menu.Menu) {
	if ctx.Manifest == nil {
		panic("InitializeMenus: ctx.Manifest is nil")
	}

	InfoMenu := menu.NewMenu("Mod List Info", "Menu for Mod List Info", InfoMenuID).AddButton(
		// SORT BY CATEGORY CURRENTLY NOT WORKING
		"[C] Category",
		"Press C to sort by Category.",
		func() error {
			PrintModInfoList(manifest.ModList(ctx.Manifest.ModsSlice()).SortedByCategory())
			return nil
		},
		'c',
		"sortCategory",
	).AddButton(
		"[N] Name",
		"Press N to sort by Name.",
		func() error {
			PrintModInfoList(manifest.ModList(ctx.Manifest.ModsSlice()).SortedByName())
			return nil
		},
		'n',
		"sortName",
	).AddButton(
		"[B] Back",
		"Press B to go back to Main Menu.",
		func() error {
			err := menu.SetCurrent(MainMenuID)
			if err != nil {
				return err
			}
			return nil
		},
		'b',
		"back",
	).SetRender(
		func() {
			PrintModInfoList(ctx.Manifest.ModsSlice())
		})

	MainMenu := menu.NewMenu("Main Menu", "This is the Main Menu.", MainMenuID).AddButton(
		"[I] Install",
		"Press I to install Modpack.",
		func() error {
			err := install.CleanInstall(ctx.Manifest, ctx.Paths)
			if err != nil {
				return err
			}
			return nil
		},
		'i',
		"install",
	).AddButton(
		"[H] Help/Info",
		"Press H to show Help/Info.",
		func() error {
			err := menu.SetCurrent(InfoMenuID)
			if err != nil {
				return err
			}
			return nil
		},
		'h',
		"help",
	).SetRender(
		func() {
			tui.ClearScreenRaw()
			StartHeader(ctx.Manifest)
		})

	menu.MustSetCurrent(MainMenuID)

	return MainMenu, InfoMenu
}
