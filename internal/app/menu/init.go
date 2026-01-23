package app

import (
	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

type PackMenuItem struct {
	Name        string
	Version     string
	McVersion   string
	Loader      string
	Description string
	Installed   bool
	Enabled     bool
	Key         rune
	Action      menu.Action
}

type PackMenuModel struct {
	Available []PackMenuItem
	Installed []PackMenuItem

	Expanded string
}

// Menu IDs
const (
	MainMenuID menu.MenuID = iota
	ModPackMenuID
	UpdateMenuID
	HelpMenuID
	ResourceMenuID
)

// InitializeMenus initializes the empty menus for the program
func InitializeMenus(m *manifest.Manifest) {
	if m == nil {
		panic("InitializeMenus: Manifest is nil")
	}

	// MAIN MENU
	mainMenu := menu.NewMenu("Main Menu", "This is the Main Menu.", MainMenuID)
	mainMenu.AddButton("Mod Packs", "", "Press M to view available Mod Packs", menu.ChangeMenu(ModPackMenuID), 'm', "modpacks")
	mainMenu.AddButton("Resource Bundles", "", "Press R to view available Resource Bundles", menu.ChangeMenu(ResourceMenuID), 'r', "resourceBundles")
	mainMenu.AddButton("Updates", "", "Press U for Update menu", menu.ChangeMenu(UpdateMenuID), 'u', "updateMenu")
	mainMenu.AddButton("Help", "", "Press H for help menu", menu.ChangeMenu(HelpMenuID), 'h', "help")

	// RESOURCE BUNDLE MENU
	resourceMenu := menu.NewMenu("Resource Bundles", "This is the resource bundles menu", ResourceMenuID)
	resourceMenu.AddButton("Back", "<", "Go Back", menu.ChangeMenu(MainMenuID), 'b', "back")

	// MODPACK MENU
	modPackMenu := BuildModPackMenu()
	_ = modPackMenu

	// UPDATE MENU
	updateMenu := menu.NewMenu("Update Menu", "This is the Update Menu(CURRENTLY NOT IMPLEMENTED)", UpdateMenuID)
	updateMenu.AddButton("Check for Updates", "", "Press C to check for Updates", menu.Action{}, 'c', "updateCheck")
	updateMenu.AddButton("Update All", "", "Press U to Install Available Updates", menu.Action{}, 'u', "updateAll")
	updateMenu.AddButton("Back", "", "Press B to go Back", menu.ChangeMenu(MainMenuID), 'b', "back")

	// HELP MENU
	helpMenu := menu.NewMenu("Help Menu", "This is the Help Menu", HelpMenuID)
	helpMenu.AddButton("Back", "", "Press B to go Back", menu.ChangeMenu(MainMenuID), 'b', "back")

	menu.MustSetCurrent(MainMenuID)
}
