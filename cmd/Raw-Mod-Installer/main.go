package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/w1lam/Packages/pkg/fabric"
	"github.com/w1lam/Raw-Mod-Installer/internal/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

const mcVersion = "1.21.10"

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add ModInstallerProgram version to mod-list.txt and check for program updates
// Print mod list info to README.txt or a separate MODLIST-INFO.txt file
// Merge README and ver.txt and mod-list into one file that contains all information about versions and mod list info
// Keep catgegory names from mod list and print them in mod info list
// Verify installed mods?
// !!!!!!Fix Fabric installer detection to check for exact version match instead of just latest version!!!!!! IMPORTANT

func main() {
	// TEMP TESTING CODE

	fabVer, err0 := fabric.GetLatestLocalFabricVersion(mcVersion)
	if err0 != nil {
		log.Fatal(err0)
	}
	fmt.Printf("Latest Fabric Loader Version for MC %s: %s\n", mcVersion, fabVer)
	//	err := menu.PrintModInfoList(paths.ModListURL)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	time.Sleep(5 * time.Hour)

	// ------------------------------------
	//
	// Program Intro
	//
	// ------------------------------------
	fmt.Printf("Raw Mod Installer\n")
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
	fabricState, err := fabric.CheckFabricVersions(mcVersion)
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
			err1 := fabric.RunFabricInstaller(installerPth, mcVersion)
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
			err1 := fabric.RunFabricInstaller(installerPth, mcVersion)
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
				err := download.DownloadMods(paths.ModListURL, "mods")
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
				err := download.DownloadMods(paths.ModListURL, "mods")
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
