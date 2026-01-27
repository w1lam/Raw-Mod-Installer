package app

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/actions"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	pkg "github.com/w1lam/Raw-Mod-Installer/internal/packages/fetch"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
	ui "github.com/w1lam/Raw-Mod-Installer/internal/ui/render"
)

type PackMenuItem struct {
	Name        string
	Type        string
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

type PackageMenuConfig struct {
	Type     packages.PackageType
	Title    string
	Subtitle string
}

var modPackMenuConfig = PackageMenuConfig{
	Type:     packages.PackageModPack,
	Title:    "Mod Packs",
	Subtitle: "Choose a Mod Pack",
}

var resourceBundleMenuConfig = PackageMenuConfig{
	Type:     packages.PackageResourceBundle,
	Title:    "Resource Bundles",
	Subtitle: "Choose a Resource Bundle",
}

// BuildPackageMenu builds a package menu
func BuildPackageMenu(cfg PackageMenuConfig, menuID menu.MenuID) *menu.Menu {
	m := menu.NewMenu(cfg.Title, cfg.Subtitle, menuID)

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
					installed = s.Manifest().InstalledPackages[cfg.Type]
					available = s.AvailablePackages()[cfg.Type]
					enabled = s.Manifest().EnabledPackages[cfg.Type]
				})

				// RESERVED KEYS
				used := map[rune]bool{
					'b': true, // BACK
					'i': true, // INSTALL
					'x': true, // UNINSTALL
					'e': true, // ENABLE
					'd': true, // DISABLE
				}

				// BUILDING AVAILABLE PACKAGE MODEL
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

					rebuildPackageButtons(m, &model)
					menu.RequestRender()
				})
			},
			Async: true,
		})
	})

	renderer := ui.PlainRenderer{Out: os.Stdout}

	m.SetRender(func() {
		view := buildPackageMenuView(cfg, &model, loading, errMsg)
		renderer.RenderPackageMenu(view)
	})

	return m
}

func rebuildPackageButtons(m *menu.Menu, model *PackMenuModel) {
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

					rebuildPackageButtons(m, model)
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
				"Install this Package",
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

						menu.PauseInput()
						defer menu.ResumeInput()
						return installer.PackageInstaller(plan)
					},
					WrapUp: func(err error) {
						if err == nil {
							fmt.Println("Installation Complete!")
							fmt.Println("Returning to Menu...")
							time.Sleep(time.Second * 3)

							rebuildPackageButtons(m, model)
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

					rebuildPackageButtons(m, model)
					menu.RequestRender()
					return nil
				},
			},
			key,
			name,
		)

		if model.Expanded == item.Name {
			if !item.Enabled {
				m.AddButton("Enable",
					"",
					"Enable Package",
					menu.Action{
						Function: func() error {
							err := services.EnablePackage(packages.Pkg{
								Name: item.Name,
								Type: packages.PackageType(item.Type),
							})
							if err == nil {
								item.Enabled = true
							}

							rebuildPackageButtons(m, model)
							menu.RequestRender()
							return err
						},
						Async: true,
					},
					'e',
					"enable"+item.Name,
				)
			} else {
				m.AddButton("Disable",
					"",
					"Disable Package",
					menu.Action{
						Function: func() error {
							err := services.DisablePackage(packages.Pkg{
								Name: item.Name,
								Type: packages.PackageType(item.Type),
							})
							if err == nil {
								item.Enabled = false
							}
							rebuildPackageButtons(m, model)
							menu.RequestRender()
							return err
						},
						Async: true,
					},
					'd',
					"disable"+item.Name,
				)
			}
			m.AddButton("Uninstall",
				"",
				"Uninstall Package",
				menu.Action{
					Function: func() error {
						menu.PauseInput()
						defer menu.ResumeInput()

						return services.UninstallPackage(packages.Pkg{
							Name: item.Name,
							Type: packages.PackageType(item.Type),
						})
					},
					WrapUp: func(err error) {
						fmt.Printf("Package: %s Uninstalled", item.Name)
						time.Sleep(time.Second * 3)

						rebuildPackageButtons(m, model)
						menu.RequestRender()
					},
					Async: true,
				},
				'x',
				"uninstall"+item.Name,
			)
		}
	}
}

func buildPackageMenuView(
	cfg PackageMenuConfig,
	model *PackMenuModel,
	loading bool,
	errMsg string,
) ui.PackageMenuView {
	view := ui.PackageMenuView{
		Title:   cfg.Title,
		Loading: loading,
		Error:   errMsg,
	}

	for _, item := range model.Available {
		view.Available = append(view.Available, ui.PackageMenuItemView{
			Key:         item.Key,
			Name:        item.Name,
			Description: item.Description,
			Version:     item.Version,
			McVersion:   item.McVersion,
			Loader:      item.Loader,
			Expanded:    model.Expanded == item.Name,
		})
	}

	for _, item := range model.Installed {
		view.Installed = append(view.Installed, ui.PackageMenuItemView{
			Key:         item.Key,
			Name:        item.Name,
			Description: item.Description,
			Version:     item.Version,
			McVersion:   item.McVersion,
			Loader:      item.Loader,
			Installed:   true,
			Enabled:     item.Enabled,
			Expanded:    model.Expanded == item.Name,
		})
	}

	return view
}
