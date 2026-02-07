// Package cli handles cli functionality
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcpac",
	Short: "Minecraft package manager",
	Long:  "mcpac manages Minecraft modpacks, resource packs, and shaders.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
