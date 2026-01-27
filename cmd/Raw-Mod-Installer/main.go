package main

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/app"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates

// initiation
func init() {}

func main() {
	app.Initialize()

	app.Run()
}
