// Package menu provides functions to handle menu states and user input.
package menu

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/w1lam/Packages/pkg/modrinth"
	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/features"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

type ProgramInfo struct {
	ProgramVersion string
	ModListVersion string
}

func MainMenu(programInfo ProgramInfo, width int) {
	fmt.Print(strings.Repeat("━", width))
	fmt.Print("\n\n")

	tui.PrintCentered("┏━━━━━━━━━━━━━━━┓", width)
	fmt.Print("\n")
	tui.PrintCentered("┃ MOD INSTALLER ┃", width)
	fmt.Print("\n")
	tui.PrintCentered("┗━━━━━━━━━━━━━━━┛", width)
	fmt.Print("\n")
	tui.PrintCentered("Program Version: "+programInfo.ProgramVersion, width)
	fmt.Print("\n")
	tui.PrintCentered("Mod List Version: "+programInfo.ModListVersion, width)
	fmt.Print("\n\n")
}

func ExitProgram() {
	fmt.Print("\n\n\n\n\n\n\n\nExiting...\n")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

func PrintModInfoList(modInfoList modrinth.ModInfoList) {
	for _, modInfo := range modInfoList {
		fmt.Printf(" * %s\n  ", modInfo.Title)

		for _, cat := range modInfo.Category {
			fmt.Printf(" 📂 %s ", cat)
		}

		fmt.Printf("\n   ✎  %s\n", modInfo.Description)

		if modInfo.Source != "" {
			fmt.Printf("   🔗 %s\n", modInfo.Source)
		}

		if modInfo.Wiki != "" {
			fmt.Printf("   📖 %s\n\n", modInfo.Wiki)
		} else {
			fmt.Printf("\n")
		}
	}
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

			modEntryList, errM := features.GetModEntryList(paths.ModListURL)
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
