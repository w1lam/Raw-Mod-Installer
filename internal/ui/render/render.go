// Package ui holds ui renders and headers ONLY visuals
package ui

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type PackageMenuView struct {
	Title   string
	Loading bool
	Error   string

	Available []PackageMenuItemView
	Installed []PackageMenuItemView
}

type PackageMenuItemView struct {
	Key         rune
	Name        string
	Description string

	Version   string
	McVersion string
	Loader    string

	Installed bool
	Enabled   bool
	Expanded  bool
}

type PackageMenuRenderer interface {
	RenderPackageMenu(view PackageMenuView)
}

type PlainRenderer struct {
	Out io.Writer
}

func (r *PlainRenderer) RenderPackageMenu(v PackageMenuView) {
	titleLen := len(v.Title)

	fmt.Fprintf(r.Out, "  %s\n", v.Title)
	fmt.Fprintf(r.Out, " ━%s━\n\n", strings.Repeat("━", titleLen))

	if v.Loading {
		fmt.Fprintf(r.Out, "  Loading %s...\n\n", v.Title)
		fmt.Fprintf(r.Out, " ━━━━━━━━━%s━━━━\n", strings.Repeat("━", titleLen))
		fmt.Fprintf(r.Out, " [B] Back   [Q] Quit\n")
		return
	}

	if v.Error != "" {
		fmt.Fprintf(r.Out, " Error: %s\n", v.Error)
	}

	renderSection := func(title string, items []PackageMenuItemView) {
		fmt.Fprintf(r.Out, "  %s\n", title)
		fmt.Fprintf(r.Out, " ━━━━━━━━━━━%s━\n", strings.Repeat("━", titleLen))

		if len(items) == 0 {
			fmt.Fprintf(r.Out, "  (none)\n")
			return
		}

		for _, item := range items {
			fmt.Fprintf(r.Out, "  [%c] %s", unicode.ToUpper(item.Key), item.Name)
			if item.Enabled {
				fmt.Fprintf(r.Out, " (enabled)")
			}
			fmt.Fprintln(r.Out)

			if item.Description != "" {
				fmt.Fprintf(r.Out, "    - %s\n", item.Description)
			}

			if item.Expanded {
				fmt.Fprintf(r.Out, "    - Package Version: %s\n", item.Version)
				fmt.Fprintf(r.Out, "    - Minecraft Version: %s\n", item.McVersion)
				if item.Loader != "" {
					fmt.Fprintf(r.Out, "    - Loader: %s\n", item.Loader)
				}

				if item.Installed {
					if item.Enabled {
						fmt.Fprintf(r.Out, "      [D] Disable   [X] Uninstall\n")
					} else {
						fmt.Fprintf(r.Out, "      [E] Enabled   [X] Uninstall\n")
					}
				} else {
					fmt.Fprintf(r.Out, "      [I] Install\n")
				}
			}
			fmt.Fprintln(r.Out)
		}
	}

	renderSection("Available "+v.Title, v.Available)
	renderSection("Installed "+v.Title, v.Installed)

	fmt.Fprintf(r.Out, " ━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(r.Out, "  [B] Back   [Q] Quit\n")
}
