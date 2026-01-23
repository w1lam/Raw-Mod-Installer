package main

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/app"
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
	// mp, err := lists.GetAvailableModPacks()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("%+v", mp)
	//
	// time.Sleep(time.Hour * 1)

	_ = app.Initialize()

	app.Run()
}
