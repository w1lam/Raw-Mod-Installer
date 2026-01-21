package main

import (
	"fmt"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	"github.com/w1lam/Raw-Mod-Installer/internal/lists"
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
	allLists, err := lists.GetAllAvailableLists()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", allLists)
	time.Sleep(time.Hour * 1)

	m := app.Initialize()

	_ = m

	app.Run()
}
