package app

import (
	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/actions"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/ui"
)

// Menu IDs
const (
	MainMenuID menu.MenuID = iota
	InstallerMenuID
	UpdateMenuID
	HelpMenuID
)

// InitializeMenus initializes the empty menus for the program
func InitializeMenus(m *manifest.Manifest) {
	if m == nil {
		panic("InitializeMenus: Manifest is nil")
	}

	mainMenu := menu.NewMenu("Main Menu", "This is the Main Menu.", MainMenuID)
	installerMenu := menu.NewMenu("Installer Menu", "This is the Installer Menu", InstallerMenuID)
	updateMenu := menu.NewMenu("Update Menu", "This is the Update Menu", UpdateMenuID)
	helpMenu := menu.NewMenu("Help Menu", "This is the Help Menu", HelpMenuID)

	mainMenu.AddButton("BTN", "⌂", "btn btn", func() error { return actions.InstallModPackAction("SwagPack", m) }, 'b', "btn")
	mainMenu.SetRender(ui.CustomMenuRender)

	_ = mainMenu
	_ = installerMenu
	_ = updateMenu
	_ = helpMenu

	menu.MustSetCurrent(MainMenuID)
}
