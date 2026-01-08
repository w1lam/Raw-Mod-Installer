// Package ui provides functions to handle menu states and user input.
package ui

import (
	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

// Menu IDs
const (
	MainMenuID menu.MenuID = iota
	ModPacksMenuID
	HelpMenuID
)

// InitializeMenus initializes the empty menus for the program
func InitializeMenus(m *manifest.Manifest) []*menu.Menu {
	if m == nil {
		panic("InitializeMenus: ctx.Manifest is nil")
	}

	var menus []*menu.Menu

	menus = append(menus, menu.NewMenu("Main Menu", "This is the Main Menu.", MainMenuID))
	menus = append(menus, menu.NewMenu("Mod Packs", "This is the Mod Packs Menu", ModPacksMenuID))
	menus = append(menus, menu.NewMenu("Help Menu", "This is the Help Menu", HelpMenuID))

	menu.MustSetCurrent(MainMenuID)

	return menus
}
