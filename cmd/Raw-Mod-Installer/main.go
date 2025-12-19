package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/downloadmods"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates
// Verify installed mods?
// MENU system IS COMIN ALONG MF
// FIX SORT BY CATEGORY

var ModListInfo modrinth.ModInfoList

var ver string

// initiation
func init() {
	tui.ClearScreenRaw()

	fmt.Printf("Starting Up...\n")

	tui.SetProgramExitFunc(func() {
		tui.ClearScreenRaw()
		fmt.Printf("Exiting...")
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	})

	fmt.Printf("Fetching Mod List Info...\n")

	var err error
	ModListInfo, err = modrinth.FetchModInfoList(paths.ModListURL, 10)
	if err != nil {
		fmt.Printf("Failed to fetch Mod Info List: %v\n", err)
	}

	var err1 error
	ver, err1 = features.GetRemoteVersion(paths.ModListURL)
	if err1 != nil {
		log.Fatal(err1)
	}
}

const (
	MainMenu tui.CurrentMenu = iota
	InfoMenu
)

var CurrentMenu tui.CurrentMenu

func main() {
	// TEMP TESTING CODE

	infoMenu := tui.NewRawMenu("Mod List Info", "Menu for Mod List Info").AddButton(
		"C] Category",
		"Press C to sort by Category.",
		func(any) error {
			tui.ClearScreenRaw()
			menu.PrintModInfoList(ModListInfo.SortByCategory())
			return nil
		},
		'c',
	).AddButton(
		"N] Name",
		"Press N to sort by Name.",
		func(any) error {
			tui.ClearScreenRaw()
			menu.PrintModInfoList(ModListInfo.SortByName())
			return nil
		},
		'n',
	).AddButton(
		"B] Back",
		"Press B to go back to Main Menu.",
		func(any) error {
			CurrentMenu = MainMenu
			return nil
		},
		'b',
	)

	mainMenu := tui.NewRawMenu("Main Menu", "This is the Main Menu.").AddButton(
		"I] INFO",
		"Press I to show Mod List Info.",
		func(any) error {
			CurrentMenu = InfoMenu
			return nil
		},
		'i',
	).AddButton(
		"S] Start",
		"Press S to start Installation.",
		func(any) error {
			return nil
		},
		'b',
	)

	// MAIN SYSTEM LOOP
	for {
		switch CurrentMenu {
		case MainMenu:
			tui.ClearScreenRaw()

			fmt.Printf("MOD INSTALLER\n")
			fmt.Printf("Mod List Version: %s\n\n", ver)

			mainMenu.RenderMenu(60)

			err := mainMenu.GetInput()
			if err != nil {
				log.Fatal(err)
			}
		case InfoMenu:
			menu.PrintModInfoList(ModListInfo)

			infoMenu.RenderMenu(60)

			err := infoMenu.GetInput()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//modInfoList, errM := modrinth.FetchModInfoList(paths.ModListURL, 10)
	//if errM != nil {
	//	log.Fatal(errM)
	//}
	//
	//	err0 := menu.PrintModInfoList(modInfoList)
	//	if err0 != nil {
	//		log.Fatal(err0)
	//	}
	time.Sleep(5 * time.Hour)

	// ------------------------------------
	//
	// Program Intro
	//
	// ------------------------------------
	fmt.Printf("-------------------\n")
	fmt.Printf(" Raw Mod Installer \n")
	fmt.Printf("-------------------\n")
	ver, err2 := features.GetRemoteVersion(paths.ModListURL)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Display Mod List Version
	fmt.Printf("Mod List Version: %s\n", ver)

	fmt.Printf("\nAt any point Press Q or ESC to Exit.\n\n")
	fmt.Printf("Press Enter to Continue or Press I for Mod List Info.\n")
	err := menu.InitInput()
	if err != nil {
		log.Fatal(err)
	}

	// Check Fabric version and install if req
	fmt.Print("Checking Fabric version...\n")
	fabricState, err := fabric.CheckFabricVersions(paths.McVersion)
	if err != nil {
		log.Fatal(err)
	}

	// ------------------------------------
	//
	// Fabric Install / Update Menu
	//
	// ------------------------------------
	switch fabricState {
	case "notInstalled":
		fmt.Printf("Fabric not Installed...\n")
		fmt.Printf("Would you like to install Fabric now? Press Y / N to Continue\n")

		// Get User Input
		switch input, err := menu.UserInput(); input {
		case "yes":
			fmt.Printf("\nDownloading Fabric Installer...\n")
			installerPth, err := fabric.GetLatestFabricInstallerJar()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Installing Fabric...\n")
			err1 := fabric.RunFabricInstaller(installerPth, paths.McVersion)
			if err1 != nil {
				log.Fatal(err1)
			}
			fmt.Printf("Fabric Installed Successfully!\n")

		case "no":
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Skipping Fabric Installation.\n")
		}

	case "updateFound":
		fmt.Printf("New Version of Fabric Found\n")
		fmt.Printf("Would you like to install the latest version? Press Y / N to Continue\n")

		// Get User Input
		switch input, err := menu.UserInput(); input {
		case "yes":
			fmt.Printf("\nDownloading Fabric Installer...\n")
			installerPth, err := fabric.GetLatestFabricInstallerJar()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Installing New Fabric Vernsion...\n")
			err1 := fabric.RunFabricInstaller(installerPth, paths.McVersion)
			if err1 != nil {
				log.Fatal(err1)
			}
			fmt.Printf("Fabric Updated Successfully!\n")

		case "no":
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Skipping Fabric Update.\n")
		}

	case "upToDate":
		fmt.Printf("Fabric is up to date.\n")

	case "localNewer":
		log.Fatal("Local Fabric version is newer than the latest available version. Please check your installation.")
	}

	// ------------------------------------
	//
	// Main Program Loop
	//
	// ------------------------------------
	for {
		// Get Program State
		switch programState, err := menu.GetState(); programState {
		case 0:
			if err != nil {
				log.Fatal(err)
			}

		// Mods Not Installed State Menu
		case menu.StateNotInstalled:

			fmt.Printf("\n\nMods are not installed. Would you like to install them? Press Y / N to Continue")

			switch input, err := menu.UserInput(); input {
			case "yes":
				// Download Mods in Temp Folder
				err := downloadmods.DownloadMods(paths.ModListURL, "mods")
				if err != nil {
					log.Fatal(err)
				}

				// Backup Existing Mods
				err1 := features.BackupModFolder()
				if err1 != nil {
					log.Fatal(err1)
				}

				// Move Temp Download to Mods Folder
				err2 := os.Rename(paths.TempModDownloadPath, paths.ModFolderPath)
				if err2 != nil {
					log.Fatal(err2)
				}

				fmt.Printf("\n\nPress ESC or Q to Exit")
				_, err3 := menu.UserInput()
				if err3 != nil {
					log.Fatal(err3)
				}

			case "no":
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\nBruh u Deadass\n\nPress ESC or Q to Exit")
				_, err := menu.UserInput()
				if err != nil {
					log.Fatal(err)
				}
			}

		// Mods Update Found State Menu
		case menu.StateUpdateFound:

			fmt.Printf("\n\nMods update found. Would you like to update them? Press Y / N to Continue")

			switch input, err := menu.UserInput(); input {
			case "yes":
				// Download Mods in Temp Folder
				err := downloadmods.DownloadMods(paths.ModListURL, "mods")
				if err != nil {
					log.Fatal(err)
				}

				// Uninstall Existing Mods
				err1 := features.UninstallMods()
				if err1 != nil {
					log.Fatal(err1)
				}

				// Move Temp Download to Mods Folder
				err2 := os.Rename(paths.TempModDownloadPath, paths.ModFolderPath)
				if err2 != nil {
					log.Fatal(err2)
				}

				fmt.Printf("\n\nPress ESC or Q to Exit")
				_, err3 := menu.UserInput()
				if err3 != nil {
					log.Fatal(err3)
				}

			case "no":
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\nBruh u Deadass\n\nPress ESC or Q to Exit")
				_, err := menu.UserInput()
				if err != nil {
					log.Fatal(err)
				}

			}
		// Mods Up to Date State Menu
		case menu.StateUpToDate:
			fmt.Printf("\nMods are up to date. Would you like to uninstall them? Press Y / N to Continue")

			switch input, err := menu.UserInput(); input {
			case "yes":
				fmt.Printf("\n Uninstalling Mods...")

				err := features.UninstallMods()
				if err != nil {
					log.Fatal(err)
				}

				err1 := features.RestoreModBackup()
				if err1 != nil {
					log.Fatal(err1)
				}

				fmt.Printf("\nMods Uninstalled Successfully!")
				fmt.Printf("\n\nPress ESC or Q to Exit")
				_, err2 := menu.UserInput()
				if err2 != nil {
					log.Fatal(err2)
				}

			case "no":
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\n\nPress ESC or Q to Exit")
				_, err1 := menu.UserInput()
				if err1 != nil {
					log.Fatal(err1)
				}
			}
		}
	}
}
