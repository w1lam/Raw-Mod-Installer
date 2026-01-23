// Package ui holds ui renders and headers ONLY visuals
package ui

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/w1lam/Packages/menu"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
)

func PrintInfoField(msg string) {
	margin := strings.Repeat(" ", config.Style.Margin)
	marginBottom := strings.Repeat("\n", config.Style.Margin/3)
	trailSpace := strings.Repeat(" ", config.Style.Width-(config.Style.Margin+utf8.RuneCountInString(msg)))
	fmt.Printf("\r%s%s%s%s", margin, msg, trailSpace, marginBottom)
}

// RenderModPackModList renders a list of the mods inside a modpack with Title, Categories, Source and Wiki if available
func RenderModPackModList(modPack manifest.InstalledModPack, meta *meta.MetaData) {
	for _, mod := range modPack.Mods {
		fmt.Printf(" * %s\n  ", meta.Mods[mod.Slug].Title)

		for _, cat := range meta.Mods[mod.Slug].Categories {
			fmt.Printf(" 📂 %s ", cat)
		}

		fmt.Printf("\n   ✎  %s\n", meta.Mods[mod.Slug].Description)

		if meta.Mods[mod.Slug].Source != "" {
			fmt.Printf("   🔗 %s\n", meta.Mods[mod.Slug].Source)
		}

		if meta.Mods[mod.Slug].Wiki != "" {
			fmt.Printf("   📖 %s\n\n", meta.Mods[mod.Slug].Wiki)
		} else {
			fmt.Printf("\n")
		}
	}
}

// customBtnRender draws the button to the console.
func customBtnRender(b *menu.Button) {
	gapSize := config.Style.Width - (utf8.RuneCountInString(b.Title) + utf8.RuneCountInString(b.Description) + utf8.RuneCountInString(b.Icon) + (config.Style.Margin * 2))
	if gapSize < 1 {
		gapSize = 1
	}

	margin := strings.Repeat(" ", config.Style.Margin)
	gap := strings.Repeat(" ", gapSize-(config.Style.Padding*2))
	if b.Icon == "" {
		fmt.Printf("%s┃%s[%c] %s%s%s%s┃%s", margin, config.Style.PaddingStr(), unicode.ToUpper(b.Key), b.Title, gap, b.Description, config.Style.PaddingStr(), margin)
	} else {
		gap = strings.Repeat(" ", gapSize-(config.Style.Padding*5))
		fmt.Printf("%s┃%s%s%s[%c] %s%s%s%s┃%s", margin, config.Style.PaddingStr(), b.Icon, config.Style.PaddingStr(), unicode.ToUpper(b.Key), b.Title, gap, b.Description, config.Style.PaddingStr(), margin)
	}
}

// CustomMenuRender renders a menu to the console
func CustomMenuRender() {
	m := menu.CurrentActiveMenu
	margin := strings.Repeat(" ", config.Style.Margin)

	if config.Style.RenderHeaders {
		fmt.Print(margin, "┏", strings.Repeat("━", config.Style.Width-2-(config.Style.Margin*2)), "┓", margin)
		fmt.Print(margin, "┃", strings.Repeat(" ", config.Style.Width-2-(config.Style.Margin*2)), "┃", margin)

		gap := config.Style.Width - 2 - utf8.RuneCountInString(m.Header) - (config.Style.Margin * 2)

		left := gap / 2
		right := gap - left
		fmt.Printf("%s┃%s%s%s┃%s", margin, strings.Repeat(" ", left), m.Header, strings.Repeat(" ", right), margin)
		fmt.Print(margin, "┃", strings.Repeat(" ", config.Style.Width-2-(config.Style.Margin*2)), "┃", margin)
		fmt.Print(margin, "┠", strings.Repeat("─", config.Style.Width-2-(config.Style.Margin*2)), "┨", margin)
	} else {
		fmt.Print(margin, "┏", strings.Repeat("━", config.Style.Width-2-(config.Style.Margin*2)), "┓", margin)
	}

	fmt.Print(margin, "┃", strings.Repeat(" ", config.Style.Width-2-(config.Style.Margin*2)), "┃", margin)

	for btn := range m.Buttons {
		customBtnRender(&m.Buttons[btn])
	}
	gap := strings.Repeat(" ", config.Style.Width-(config.Style.Padding*2)-11-(config.Style.Margin*2))
	fmt.Printf("%s┃%s⛝%s[Q] Exit%s┃%s", margin, config.Style.PaddingStr(), config.Style.PaddingStr(), gap, margin)

	fmt.Print(margin, "┃", strings.Repeat(" ", config.Style.Width-2-(config.Style.Margin*2)), "┃", margin)
	fmt.Print(margin, "┗", strings.Repeat("━", config.Style.Width-2-(config.Style.Margin*2)), "┛", margin)

	for range config.Style.Margin / 3 {
		fmt.Print("\n")
	}
}
