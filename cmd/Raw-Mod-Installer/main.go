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

	// modrinth.EnableDevMode()

	// pkgs, err := packages.GetAllAvailablePackages()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", pkgs)
	//
	// vers := modrinth.FetchBestVersions(pkgs["resourcebundles"]["Visual Enhancements Bundle"].Entries, modrinth.Filter{McVersion: "1.21.10", Loader: ""})
	// fmt.Printf("%+v", vers)
	//
	// time.Sleep(time.Hour * 1)

	app.Run()
}
