package app

import (
	"fmt"
	"sort"
	"time"
	"unicode"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/actions"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	pkg "github.com/w1lam/Raw-Mod-Installer/internal/packages/fetch"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

// BuildModPackMenu builds the modPackMenu
func BuildModPackMenu() *menu.Menu {
	m := menu.NewMenu("Mod Packs", "Chose a Mod pack", ModPackMenuID)

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
				gState := state.Get()

				allAvailable, err := pkg.GetAllAvailablePackages()
				if err != nil {
					return err
				}
				state.SetAvailablePackages(allAvailable)

				var installed map[string]manifest.InstalledPackage
				var available map[string]packages.ResolvedPackage
				var enabled string

				gState.Read(func(s *state.State) {
					installed = s.Manifest().InstalledPackages[packages.PackageModPack]
					available = s.AvailablePackages()[packages.PackageModPack]
					enabled = s.Manifest().EnabledPackages[packages.PackageModPack]
				})

				// RESERVED KEYS
				used := map[rune]bool{
					'b': true, // BACK
					'i': true, // INSTALL
				}

				// BUILDING AVAILABLE MODPACKS MODEL
				for _, mp := range available {
					mp := mp
					if _, ok := installed[mp.Name]; ok {
						continue
					}

					model.Available = append(model.Available, PackMenuItem{
						Name:        mp.Name,
						Type:        string(mp.Type),
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
					inst := inst
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
						Type:        string(inst.Type),
						Version:     inst.InstalledVersion,
						McVersion:   inst.McVersion,
						Loader:      inst.Loader,
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

					rebuildModPackButtons(m, &model)
					menu.RequestRender()
				})
			},
			Async: true,
		})
	})

	m.SetRender(func() {
		fmt.Println("  Mod Packs")
		fmt.Println(" ━━━━━━━━━━━\n")

		if loading {
			fmt.Println("  Loading Modpacks...\n")
			fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf(" [B] Back   [Q] Quit\n")
			return
		}

		if errMsg != "" {
			fmt.Printf(" Error: %s\n", errMsg)
		}

		fmt.Println("  Available Mod Packs")
		fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━")
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
					fmt.Printf("    - ModPack Version: %s\n", item.Version)
					fmt.Printf("    - Minecraft Version: %s\n", item.McVersion)
					fmt.Printf("    - Loader: %s\n", item.Loader)
					fmt.Printf("      [I] Install\n")
				}
				fmt.Println()
			}
		}

		fmt.Println("  Installed Mod Packs")
		fmt.Println(" ━━━━━━━━━━━━━━━━━━━━━")
		if len(model.Installed) == 0 {
			fmt.Println("  (none)\n")
		} else {
			for _, item := range model.Installed {
				fmt.Printf("  [%c] %s\n", unicode.ToUpper(item.Key), item.Name)
				if item.Version != "" {
					fmt.Printf("    - ModPack Version: %s\n", item.Version)
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

func rebuildModPackButtons(m *menu.Menu, model *PackMenuModel) {
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

					rebuildModPackButtons(m, model)
					menu.RequestRender()
					return nil
				},
			},
			key,
			name,
		)

		if model.Expanded == item.Name {
			pkgName := item.Name
			pkgType := item.Type

			m.AddButton(
				"Install",
				"",
				"Install this Mod Pack",
				menu.Action{
					Function: func() error {
						var pkg packages.ResolvedPackage

						state.Get().Read(func(s *state.State) {
							ap := s.AvailablePackages()
							if ap == nil {
								return
							}
							pkg = (ap)[packages.PackageType(pkgType)][pkgName]
						})

						plan := installer.InstallPlan{
							RequestedPackage: pkg,
							BackupPolicy:     services.BackupIfExists,
						}

						return installer.PackageInstaller(plan)
					},
					WrapUp: func(err error) {
						if err == nil {
							fmt.Println("Installation Complete!")
							time.Sleep(time.Second * 3)
							menu.RequestRender()
						}
					},
					Async: true,
				},
				'i',
				"install"+item.Name,
			)
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

					rebuildModPackButtons(m, model)
					menu.RequestRender()
					return nil
				},
			},
			key,
			name,
		)

		if model.Expanded == item.Name {
			m.AddButton("Install", "", "Install this Mod Pack", actions.InstallModPackAction(packages.Pkg{Name: item.Name, Type: packages.PackageType(item.Type)}), 'i', "install"+item.Name)
		}
	}
}
