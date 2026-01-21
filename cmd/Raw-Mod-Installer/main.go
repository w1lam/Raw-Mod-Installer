package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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
	indented, err := json.MarshalIndent(allLists, "", "  ")
	if err != nil {
		panic(err)
	}

	m := app.Initialize()

	fmt.Printf("%+v\n\n\n", string(indented))

	clientSide, err := lists.ComputeDirHash(filepath.Join(m.Paths.ModPacksDir, "SwagPackClientSide"))
	if err != nil {
		panic(err)
	}
	serverSide, err := lists.ComputeDirHash(filepath.Join(m.Paths.ModPacksDir, "SwagPackServerSide"))
	if err != nil {
		panic(err)
	}
	full, err := lists.ComputeDirHash(filepath.Join(m.Paths.ModPacksDir, "SwagPack"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("clientSide: %s\n\nserverSide: %s\n\nfull: %s", clientSide, serverSide, full)

	time.Sleep(time.Hour * 1)

	_ = m

	app.Run()
}
