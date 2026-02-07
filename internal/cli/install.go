package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w1lam/Raw-Mod-Installer/internal/actions"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
)

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Install a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		pkgType, _ := cmd.Flags().GetString("type")

		pkg := packages.Pkg{
			Name: name,
			Type: packages.PackageType(pkgType),
		}

		if err := actions.InstallPackage(pkg); err != nil {
			return err
		}

		fmt.Println("Installed:", name)
		return nil
	},
}

func init() {
	installCmd.Flags().StringP("type", "t", "modpack", "package type")
	rootCmd.AddCommand(installCmd)
}
