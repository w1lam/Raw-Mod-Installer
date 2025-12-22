package ui

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/modrinthsvc"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

// User Input OLD LOGIC NEW USES EXTERNAL PKG

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

			ExitProgram()

		// Exit
		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			ExitProgram()
		}
	}
}

func InfoInput(modListInfo modrinth.ModInfoList) error {
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
		// Sort by category
		case char == 'c' || char == 'C':
			err := keyboard.Close()
			if err != nil {
				return err
			}

			PrintModInfoList(modListInfo.SortByCategory())

			return nil

		case char == 'n' || char == 'N':
			err := keyboard.Close()
			if err != nil {
				return err
			}

			PrintModInfoList(modListInfo.SortByName())

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

			ExitProgram()

		// Exit
		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			ExitProgram()
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

			fmt.Printf("\nFetching Mod List Info...\n\n")

			modEntryList, errM := modrinthsvc.GetModEntryList(netcfg.ModListURL)
			if errM != nil {
				return errM
			}

			modInfoList, err0 := modrinth.FetchModInfoList(modEntryList, 10)
			if err0 != nil {
				return err0
			}

			PrintModInfoList(modInfoList)

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

			ExitProgram()

		// Exit
		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			ExitProgram()
		}
	}
}
