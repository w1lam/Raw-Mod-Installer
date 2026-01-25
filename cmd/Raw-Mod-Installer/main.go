package main

import (
	"fmt"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	fetch "github.com/w1lam/Raw-Mod-Installer/internal/packages/fetch"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
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
	app.Initialize()

	all, err := fetch.GetAllAvailablePackages()
	if err != nil {
		panic(err)
	}

	fmt.Println(all["modpack"]["SwagPack"])
	time.Sleep(time.Hour * 1)

	// fmt.Printf("%+v\n\n\n", all[packages.PackageResourceBundle])
	//
	// filter := modrinth.EntryFilter{
	// 	McVersion:   "1.21.10",
	// 	ProjectType: string(packages.PackageModPack),
	// 	Loader:      "fabric",
	// }
	//
	// ver := modrinth.FetchBestVersions(all[packages.PackageModPack]["SwagPack"].Entries, filter)
	//
	// dlItem := downloader.DownloadItem{
	// 	ID:       "sodium",
	// 	FileName: filepath.Base(ver["sodium"].Files[0].URL),
	// 	URL:      ver["sodium"].Files[0].URL,
	// 	Sha1:     ver["sodium"].Files[0].Hashes.Sha1,
	// 	Sha512:   ver["sodium"].Files[0].Hashes.Sha512,
	// 	Version:  ver["sodium"].VersionNumber,
	// }
	// fmt.Printf("%+v", dlItem)
	//
	// p, err := paths.Resolve()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// dlMap := map[string]downloader.DownloadItem{
	// 	"sodium": dlItem,
	// }
	// dlr, err := downloader.DownloadEntries(dlMap, p)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("%+v", dlr)

	plan := installer.InstallPlan{
		RequestedPackage: all[packages.PackageModPack]["SwagPack"],
		BackupPolicy:     services.BackupOnce,
	}

	if err := installer.PackageInstaller(plan); err != nil {
		panic(err)
	}

	time.Sleep(time.Hour * 1)

	app.Run()
}
