package main

import (
	"fmt"
	"log"
	"os"

	"github.com/w1lam/Packages/pkg/fabric"
)

const mcVersion = "1.21.10"

func main() {
	// ------------------------------------
	//
	// Program Intro
	//
	// ------------------------------------
	fmt.Printf("Raw Mod Installer\n")
	ver, err2 := GetRemoteVersion(ModListURL)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Display Mod List Version
	fmt.Printf("Mod List Version: %s\n", ver)
	fmt.Printf("\nAt any point Press Q or ESC to Exit.\n")

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
		switch input, err := UserInput(); input {
		case true:
			fmt.Printf("Downloading Fabric Installer...\n")
			installerPth, err := fabric.GetLatestFabricInstallerJar()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Installing Fabric...\n")
			err1 := fabric.RunFabricInstaller(installerPth, mcVersion)
			if err1 != nil {
				log.Fatal(err1)
			}

		case false:
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Skipping Fabric Installation.\n")
		}

	case "updateFound":
		fmt.Printf("New Version of Fabric Found: %s\n", ver)
		fmt.Printf("Would you like to install the latest version? Press Y / N to Continue")

		// Get User Input
		switch input, err := UserInput(); input {
		case true:
			fmt.Printf("Downloading Fabric Installer...\n")
			installerPth, err := fabric.GetLatestFabricInstallerJar()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Installing New Fabric Vernsion...\n")
			err1 := fabric.RunFabricInstaller(installerPth, mcVersion)
			if err1 != nil {
				log.Fatal(err1)
			}

		case false:
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
		switch programState, err := GetState(); programState {
		case 0:
			if err != nil {
				log.Fatal(err)
			}

		// Mods Not Installed State Menu
		case StateNotInstalled:

			switch input, err := UserInput(); input {
			case true:
				// Download Mods in Temp Folder
				err := DownloadMods(ModListURL, "mods")
				if err != nil {
					log.Fatal(err)
				}

				// Backup Existing Mods
				err1 := BackupModFolder()
				if err1 != nil {
					log.Fatal(err1)
				}

				// Move Temp Download to Mods Folder
				err2 := os.Rename(TempModDownloadPath, ModFolderPath)
				if err2 != nil {
					log.Fatal(err2)
				}

				fmt.Printf("\n\nPress ESC or Q to Exit")
				_, err3 := UserInput()
				if err3 != nil {
					log.Fatal(err3)
				}

			case false:
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\nBruh u Deadass\n\nPress ESC or Q to Exit")
				_, err := UserInput()
				if err != nil {
					log.Fatal(err)
				}
			}

		// Mods Update Found State Menu
		case StateUpdateFound:

			fmt.Printf("\n\nMods update found. Would you like to update them? Press Y / N to Continue")

			switch input, err := UserInput(); input {
			case true:
				// Download Mods in Temp Folder
				err := DownloadMods(ModListURL, "mods")
				if err != nil {
					log.Fatal(err)
				}

				// Uninstall Existing Mods
				err1 := UninstallMods()
				if err1 != nil {
					log.Fatal(err1)
				}

				// Move Temp Download to Mods Folder
				err2 := os.Rename(TempModDownloadPath, ModFolderPath)
				if err2 != nil {
					log.Fatal(err2)
				}

				fmt.Printf("\n\nPress ESC or Q to Exit")
				_, err3 := UserInput()
				if err3 != nil {
					log.Fatal(err3)
				}

			case false:
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\nBruh u Deadass\n\nPress ESC or Q to Exit")
				_, err := UserInput()
				if err != nil {
					log.Fatal(err)
				}

			}
		// Mods Up to Date State Menu
		case StateUpToDate:
			fmt.Printf("\n\n Mods are up to date. Would you like to uninstall them? Press Y / N to Continue")

			switch input, err := UserInput(); input {
			case true:
				fmt.Printf("\n Uninstalling Mods...")

				err := UninstallMods()
				if err != nil {
					log.Fatal(err)
				}

				err1 := RestoreModBackup()
				if err1 != nil {
					log.Fatal(err1)
				}

				fmt.Printf("\n Mods Uninstalled Successfully!")
				fmt.Printf("\n\n Press ESC or Q to Exit")
				_, err2 := UserInput()
				if err2 != nil {
					log.Fatal(err2)
				}

			case false:
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\n\n Press ESC or Q to Exit")
				_, err1 := UserInput()
				if err1 != nil {
					log.Fatal(err1)
				}
			}
		}
	}
}
