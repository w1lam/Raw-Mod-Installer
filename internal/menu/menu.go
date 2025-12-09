// Package menu provides functions to handle menu states and user input.
package menu

import (
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

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

func UserInput() (bool, error) {
	if err := keyboard.Open(); err != nil {
		return false, err
	}

	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return false, err
		}

		switch {
		case char == 'y' || char == 'Y':
			err := keyboard.Close()
			if err != nil {
				return false, err
			}
			return true, nil

		case char == 'n' || char == 'N':
			err := keyboard.Close()
			if err != nil {
				return false, err
			}
			return false, nil

		case key == keyboard.KeyEsc:
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...")
			os.Exit(0)

		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...")
			os.Exit(0)
		}
	}
}
