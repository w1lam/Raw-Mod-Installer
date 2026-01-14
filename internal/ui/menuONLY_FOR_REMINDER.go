package ui

// InfoMenu.AddButton(
// 	// SORT BY CATEGORY CURRENTLY NOT WORKING
// 	"[C] Category",
// 	"Press C to sort by Category.",
// 	func() error {
// 		PrintModInfoList(manifest.ModList(ctx.Manifest.ModsSlice()).SortedByCategory())
// 		return nil
// 	},
// 	'c',
// 	"sortCategory",
// ).AddButton(
// 	"[N] Name",
// 	"Press N to sort by Name.",
// 	func() error {
// 		PrintModInfoList(manifest.ModList(ctx.Manifest.ModsSlice()).SortedByName())
// 		return nil
// 	},
// 	'n',
// 	"sortName",
// ).AddButton(
// 	"[B] Back",
// 	"Press B to go back to Main Menu.",
// 	func() error {
// 		err := menu.SetCurrent(MainMenuID)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	},
// 	'b',
// 	"back",
// ).SetRender(
// 	func() {
// 		PrintModInfoList(ctx.Manifest.ModsSlice())
// 	})
//
// MainMenu.AddButton(
// 	"[I] Install",
// 	"Press I to install Modpack.",
// 	func() error {
// 		m, err := install.ExecutePlan(ctx.Manifest, ctx.Paths, install.InstallPlan{
// 			Intent:       install.IntentInstall,
// 			EnsureFabric: true,
// 			BackupPolicy: install.BackupOnce,
// 			EnableAfter:  true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		if err := m.Save(ctx.Paths.ManifestPath); err != nil {
// 			return err
// 		}
// 		return nil
// 	},
// 	'i',
// 	"install",
// ).AddButton(
// 	"[U] Update",
// 	"Press U to Update Modpack",
// 	func() error {
// 		m, err := install.ExecutePlan(ctx.Manifest, ctx.Paths, install.InstallPlan{
// 			Intent:       install.IntentUpdate,
// 			EnsureFabric: true,
// 			BackupPolicy: install.BackupIfExists,
// 			EnableAfter:  true,
// 		})
// 		if err != nil {
// 			return err
// 		}
//
// 		return m.Save(ctx.Paths.ManifestPath)
// 	},
// 	'u',
// 	"update",
// ).AddButton(
// 	"[E] Enable",
// 	"Press E to Enable Mods",
// 	func() error {
// 		err := install.EnableMods(ctx.Paths)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	},
// 	'e',
// 	"enable",
// ).AddButton(
// 	"[D] Disable",
// 	"Press D to Disable Mods",
// 	func() error {
// 		err := install.DisableMods(ctx.Paths)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	},
// 	'd',
// 	"disable",
// ).AddButton(
// 	"[H] Help/Info",
// 	"Press H to show Help/Info.",
// 	func() error {
// 		err := menu.SetCurrent(InfoMenuID)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	},
// 	'h',
// 	"help",
// ).SetRender(
// 	func() {
// 		tui.ClearScreenRaw()
// 		StartHeader(ctx.Manifest)
// 	})
//
// menu.MustSetCurrent(MainMenuID)
//
// return MainMenu, InfoMenu
