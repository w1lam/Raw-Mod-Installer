package main

import (
	"github.com/w1lam/Packages/pkg/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/app"
	"github.com/w1lam/Raw-Mod-Installer/internal/ui"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates
// Verify installed mods?
// MENU system IS COMIN ALONG MF
// FIX SORT BY CATEGORY

// initiation
func init() {
}

func main() {
	ctx := app.Initialize()

	ui.InitializeMenus(ctx)

	menu.Run()
}
