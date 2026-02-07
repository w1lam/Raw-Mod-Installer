package main

import (
	"log"
	"os"

	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	"github.com/w1lam/Raw-Mod-Installer/internal/cli"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates

// initiation
func init() {}

func main() {
	if len(os.Args) > 1 {
		err := app.InitCore()
		if err != nil {
			log.Fatal(err)
		}

		cli.Execute()
		return
	}

	if err := app.InitCore(); err != nil {
		log.Fatal(err)
	}

	app.InitTUI()

	app.Run()
}
