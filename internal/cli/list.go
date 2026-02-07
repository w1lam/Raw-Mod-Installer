package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching packages...")

		availablePkgs, err := packages.GetAllAvailablePackages()
		if err != nil {
			return err
		}

		for pType, pkgs := range availablePkgs {
			fmt.Println()
			fmt.Println(strings.ToUpper(string(pType)) + "S:")
			fmt.Println(strings.Repeat("-", len(pType)) + "---")
			for _, pkg := range pkgs {
				fmt.Println(" " + pkg.Name)
				fmt.Println(" - List Version:", pkg.ListVersion)
				fmt.Println(" - Minecraft Version:", pkg.McVersion)
				if pkg.Description != "" {
					fmt.Println(" - Desctription:", pkg.Description)
				}
				fmt.Println(" - Env:", pkg.Env)
				if pkg.Loader != "" {
					fmt.Println(" - Loader:", pkg.Loader)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
