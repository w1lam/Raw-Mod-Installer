package main

import (
	"fmt"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	packages "github.com/w1lam/Raw-Mod-Installer/internal/packages/fetch"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates
// Verify installed mods?
// MENU system IS COMIN ALONG MF
// FIX SORT BY CATEGORY

// initiation
func init() {}

func main() {
	all, err := packages.GetAllAvailablePackages()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", all)

	time.Sleep(time.Hour * 1)

	app.Initialize()

	app.Run()
}
