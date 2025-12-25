// Package ui provides functions to handle menu states and user input.
package ui

import (
	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/output"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// Menu IDs
const (
	StartMenuID menu.MenuID = iota
	InfoMenuID
)

type Context struct {
	Manifest *manifest.Manifest
	Paths    *paths.Paths
}

func InitializeMenus(ctx Context) (*menu.Menu, *menu.Menu) {
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
		"[P] Print Mod Info to README",
		"Press P to print Mod Info to README.md file.",
		func() error {
			err := output.WriteModInfoListREADME(ctx.Paths.ModsDir, ctx.Manifest.ModsSlice())
			if err != nil {
				return err
			}
			return nil
		},
		'p',
		"printReadme",
	).AddButton(
		"[B] Back",
		"Press B to go back to Main Menu.",
		func() error {
			err := menu.SetCurrent(StartMenuID)
			if err != nil {
				return err
			}
			return nil
		},
		'b',
		"back",
	).SetOnEnter(
		func() {
			PrintModInfoList(ctx.Manifest.ModsSlice())
		})

	StartMenu := menu.NewMenu("Start Menu", "This is the Start Menu.", StartMenuID).AddButton(
		"[S] Start",
		"Press S to start Installation.",
		func() error {
			return nil
		},
		's',
		"start",
	).AddButton(
		"[I] INFO",
		"Press I to show Mod List Info.",
		func() error {
			err := menu.SetCurrent(InfoMenuID)
			if err != nil {
				return err
			}
			return nil
		},
		'i',
		"info",
	).SetOnEnter(
		func() {
			tui.ClearScreenRaw()
			StartHeader(ctx.Manifest)
		})

	return StartMenu, InfoMenu
}
