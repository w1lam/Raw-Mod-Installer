package app

import (
	"fmt"
	"sort"
	"unicode"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/actions"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	pkg "github.com/w1lam/Raw-Mod-Installer/internal/packages/fetch"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// BuildResourceBundleMenu builds resourceBundleMenu
func BuildResourceBundleMenu() *menu.Menu {
	m := menu.NewMenu("Resource Bundles", "Chose a Resource Bundle", ResourceMenuID)

	var (
		model   PackMenuModel
		loading bool
		errMsg  string
	)

	_ = loading
	_ = errMsg

	m.SetOnEnter(func() {
		loading = true
		errMsg = ""
		model = PackMenuModel{}

		m.ClearButtons()
		m.AddButton("Back", "<", "Go Back", menu.ChangeMenu(MainMenuID), 'b', "back")

		menu.Queue(menu.Action{
			Function: func() error {
				allAvailable, err := pkg.GetAllAvailablePackages()
				if err != nil {
					return err
				}
				state.SetAvailablePackages(allAvailable)

				available := allAvailable[packages.PackageResourceBundle]

				var installed map[string]manifest.InstalledPackage
				var enabled string

				state.Get().Read(func(s *state.State) {
					installed = s.Manifest().InstalledPackages[packages.PackageResourceBundle]
					enabled = s.Manifest().EnabledPackages[packages.PackageResourceBundle]
				})

				// RESERVED KEYS
				used := map[rune]bool{
					'b': true, // BACK
					'i': true, // INSTALL
				}

				// BUILDING AVAILABLE MODPACKS MODEL
				for _, mp := range available {
					if _, ok := installed[mp.Name]; ok {
						continue
					}

					model.Available = append(model.Available, PackMenuItem{
						Name:        mp.Name,
						Version:     mp.ListVersion,
						McVersion:   mp.McVersion,
						Loader:      mp.Loader,
						Description: mp.Description,
						Installed:   false,
						Enabled:     false,
						Action:      menu.Action{},
					})
				}
				sort.Slice(model.Available, func(i, j int) bool { return model.Available[i].Name < model.Available[j].Name })

				// ASSIGN KEYS
				for i := range model.Available {
					model.Available[i].Key = menu.AssignKey(model.Available[i].Name, used)
				}

				// BUILDING INSTALLED MODPACKS MODEL
				for _, inst := range installed {
					enabledNow := inst.Name == enabled
					title := inst.Name
					action := actions.EnablePackageAction(packages.Pkg{Name: inst.Name, Type: inst.Type})

					if enabledNow {
						title = fmt.Sprintf("%s (Enabled)", inst.Name)
						action = actions.DisablePackageAction(inst.Type)
					}

					desc := ""
					if ap, ok := available[inst.Name]; ok {
						desc = ap.Description
					}

					model.Installed = append(model.Installed, PackMenuItem{
						Name:        title,
						Version:     inst.InstalledVersion,
						McVersion:   inst.McVersion,
						Loader:      "",
						Description: desc,
						Installed:   true,
						Enabled:     enabledNow,
						Action:      action,
					})
				}
				sort.Slice(model.Installed, func(i, j int) bool { return model.Installed[i].Name < model.Installed[j].Name })

				// ASSIGN KEYS
				for i := range model.Installed {
					model.Installed[i].Key = menu.AssignKey(model.Installed[i].Name, used)
				}

				return nil
			},
			WrapUp: func(err error) {
				menu.DispatchUI(func() {
					loading = false

					if err != nil {
						errMsg = err.Error()
					}

					rebuildResourceBundleButtons(m, &model)
					menu.RequestRender()
				})
			},
			Async: true,
		})
	})

	m.SetRender(func() {
		fmt.Println("  Resource Bundles")
		fmt.Println(" ━━━━━━━━━━━━━━━━━━\n")

		if loading {
			fmt.Println("  Loading Resource Bundles...\n")
			fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf(" [B] Back   [Q] Quit\n")
			return
		}

		if errMsg != "" {
			fmt.Printf(" Error: %s\n", errMsg)
		}

		fmt.Println("  Available Resource Bundles")
		fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		if len(model.Available) == 0 {
			fmt.Println("  (none)\n")
		} else {
			for _, item := range model.Available {
				item := item

				fmt.Printf("  [%c] %s\n", unicode.ToUpper(item.Key), item.Name)
				if item.Description != "" {
					fmt.Printf("    - %s\n", item.Description)
				}

				if model.Expanded == item.Name {
					fmt.Printf("    - ResourceBundle Version: %s\n", item.Version)
					fmt.Printf("    - Minecraft Version: %s\n", item.McVersion)
					fmt.Printf("    - Loader: %s\n", item.Loader)
					fmt.Printf("      [I] Install\n")
				}
				fmt.Println()
			}
		}

		fmt.Println("  Installed Resource Bundles")
		fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		if len(model.Installed) == 0 {
			fmt.Println("  (none)\n")
		} else {
			for _, item := range model.Installed {
				fmt.Printf("  [%c] %s\n", unicode.ToUpper(item.Key), item.Name)
				if item.Version != "" {
					fmt.Printf("    - ResourceBundle Version: %s\n", item.Version)
				}
				if item.McVersion != "" {
					fmt.Printf("    - Minecraft Version: %s\n", item.McVersion)
				}
				if item.Description != "" {
					fmt.Printf("    - %s\n", item.Description)
				}
				fmt.Println()
			}
		}

		fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  [B] Back   [Q] Quit\n")
	})

	return m
}

func rebuildResourceBundleButtons(m *menu.Menu, model *PackMenuModel) {
	m.ClearButtons()

	m.AddButton("Back", "", "Go Back", menu.ChangeMenu(MainMenuID), 'b', "back")

	// BUILD AVAILABLE BUTTONS
	for i := range model.Available {
		item := &model.Available[i]
		name := item.Name
		key := item.Key

		m.AddButton(
			item.Name,
			"",
			item.Description,
			menu.Action{
				Function: func() error {
					if model.Expanded == name {
						model.Expanded = ""
					} else {
						model.Expanded = name
					}

					rebuildResourceBundleButtons(m, model)
					menu.RequestRender()
					return nil
				},
			},
			key,
			name,
		)

		if model.Expanded == item.Name {
			m.AddButton("Install", "", "Install this Resource Bundle", menu.Action{}, 'i', "install"+item.Name)
		}
	}

	// BUILD INSTALLED BUTTONS
	for i := range model.Installed {
		item := &model.Installed[i]
		name := item.Name
		key := item.Key

		m.AddButton(
			item.Name,
			"",
			item.Description,
			menu.Action{
				Function: func() error {
					if model.Expanded == name {
						model.Expanded = ""
					} else {
						model.Expanded = name
					}

					rebuildResourceBundleButtons(m, model)
					menu.RequestRender()
					return nil
				},
			},
			key,
			name,
		)

		if model.Expanded == item.Name {
			m.AddButton("Install", "", "Install this Resource Bundle", menu.Action{}, 'i', "install"+item.Name)
		}
	}
}
