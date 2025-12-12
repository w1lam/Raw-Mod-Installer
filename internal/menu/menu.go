// Package menu provides functions to handle menu states and user input.
package menu

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func PrintModInfoList(modListURL string) error {
	fmt.Printf("\nFetching Mod List Info...\n\n")

	modInfoList, err := modrinth.FetchModInfoList(modListURL, 10)
	if err != nil {
		return err
	}

	fmt.Printf("MOD LIST INFO:\n\n")

	for _, modInfo := range modInfoList {
		fmt.Printf(" 🛠  %s\n ✎  %s\n", modInfo.Title, modInfo.Description)

		if modInfo.Source != "" {
			fmt.Printf(" 🔗 %s\n", modInfo.Source)
		}

		if modInfo.Wiki != "" {
			fmt.Printf(" 📖 %s\n\n", modInfo.Wiki)
		} else {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n\nPress Enter to Continue...\n")

	err1 := InitInput()
	if err1 != nil {
		return err1
	}

	return nil
}

type State int

const (
	_ State = iota
	StateNotInstalled
	StateUpdateFound
	StateUpToDate
)

func GetState() (State, error) {
	if _, err := os.Stat(paths.VerFilePath); err != nil {
		return StateNotInstalled, nil
	} else {

		// Update check
		upToDate, err := features.CheckForModlistUpdate()
		switch {
		case err != nil:
			return 0, err

		case upToDate:
			return StateUpdateFound, nil

		default:
			return StateUpToDate, nil
		}
	}
}

// User Input

func UserInput() (string, error) {
	if err := keyboard.Open(); err != nil {
		return "", err
	}

	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return "", err
		}

		switch {
		// Yes
		case char == 'y' || char == 'Y':
			err := keyboard.Close()
			if err != nil {
				return "", err
			}
			return "yes", nil

			// No
		case char == 'n' || char == 'N':
			err := keyboard.Close()
			if err != nil {
				return "", err
			}
			return "no", nil

			// Exit
		case key == keyboard.KeyEsc:
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...\n")
			time.Sleep(1 * time.Second)
			os.Exit(0)

		// Exit
		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...\n")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}
}

func InitInput() error {
	if err := keyboard.Open(); err != nil {
		return err
	}

	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch {
		// Mod Info
		case char == 'i' || char == 'I':
			err := keyboard.Close()
			if err != nil {
				return err
			}

			err1 := PrintModInfoList(paths.ModListURL)
			if err1 != nil {
				return err1
			}

			return nil

		// Continue
		case key == keyboard.KeyEnter:
			err := keyboard.Close()
			if err != nil {
				return err
			}
			return nil

			// Exit
		case key == keyboard.KeyEsc:
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...\n")
			time.Sleep(1 * time.Second)
			os.Exit(0)

		// Exit
		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...\n")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}
}
