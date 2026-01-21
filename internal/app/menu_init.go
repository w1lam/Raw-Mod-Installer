package app

import (
	"fmt"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/actions"
	"github.com/w1lam/Raw-Mod-Installer/internal/env"
	"github.com/w1lam/Raw-Mod-Installer/internal/lists"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

// Menu IDs
const (
	MainMenuID menu.MenuID = iota
	ModPackMenuID
	UpdateMenuID
	HelpMenuID
)

// InitializeMenus initializes the empty menus for the program
func InitializeMenus(m *manifest.Manifest) {
	if m == nil {
		panic("InitializeMenus: Manifest is nil")
	}

	// MAIN MENU
	mainMenu := menu.NewMenu("Main Menu", "This is the Main Menu.", MainMenuID)
	mainMenu.AddButton("Mod Packs", "", "Press M to view available Mod Packs", menu.ChangeMenu(ModPackMenuID), 'm', "modpacks")
	mainMenu.AddButton("Updates", "", "Press U for Update menu", menu.ChangeMenu(UpdateMenuID), 'u', "updateMenu")
	mainMenu.AddButton("Help", "", "Press H for help menu", menu.ChangeMenu(HelpMenuID), 'h', "help")

	// MODPACK MENU
	modPackMenu := menu.NewMenu("Mod Packs", "This is the mod packs menu", ModPackMenuID)
	modPackMenu.AddButton("Back", "", "Press B to go Back", menu.ChangeMenu(MainMenuID), 'b', "back")
	modPackMenu.SetOnEnter(
		func() {
			m := env.GlobalManifest

			modPackMenu.ClearButtons()

			var err error
			env.AvailableModPacks, err = lists.GetAvailableModPacks()
			if err != nil {
				fmt.Printf("failed %s, ID: %d. onEnter func: %s", modPackMenu.Header, modPackMenu.ID, err)
				return
			}

			used := map[rune]bool{}
			// AvailableModPacks
			for _, mp := range env.AvailableModPacks {
				if _, ok := m.InstalledModPacks[mp.Name]; !ok {
					key := menu.AssignKey(mp.Name, used)
					modPackMenu.AddButton(
						mp.Name,
						"*",
						mp.Description,
						actions.InstallModPackAction(mp.Name),
						key,
						mp.Name,
					)
				} else {
					continue
				}
			}
			// InstalledModPacks
			for _, installed := range m.InstalledModPacks {
				key := menu.AssignKey(installed.Name, used)
				title := installed.Name
				action := actions.EnableModPackAction(installed.Name)

				if installed.Name == m.EnabledModPack {
					title = fmt.Sprintf("%s (Enabled)", installed.Name)
					action = actions.DisableModPackAction()
				}

				modPackMenu.AddButton(
					title,
					"",
					env.AvailableModPacks[installed.Name].Description,
					action,
					key,
					installed.Name,
				)
			}
			modPackMenu.AddButton("Back", "<", "Press B to go Back", menu.ChangeMenu(MainMenuID), 'b', "back")
		})

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
