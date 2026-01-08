package install

import (
	"fmt"
	"os"
	"time"

	"github.com/w1lam/Packages/pkg/tui"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

func FullUninstall(m *manifest.Manifest) error {
	if err := os.RemoveAll(m.Paths.ModsDir); err != nil {
		return err
	}

	if err := filesystem.RestoreBackup(manifest.InstalledModPack{Name: "DEFAULT"}, m); err != nil {
		return err
	}

	if err := os.RemoveAll(m.Paths.ProgramFilesDir); err != nil {
		return err
	}

	tui.ClearScreenRaw()

	fmt.Println("Installer Files and Mods Uninstalled. Exiting Program...")
	time.Sleep(time.Second * 3)

	os.Exit(0)
	return nil
}
