package modlist

import "strings"

// ParseModList parses a list of mod entries from strings into ModEntry structs.
func ParseModList(lines []string) ([]ModEntry, error) {
	mods := []ModEntry{}

	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}

		loader := "fabric"
		slug := line

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			loader = parts[0]
			slug = parts[1]
		}

		if strings.Contains(slug, "@") {
			slug = strings.Split(slug, "@")[0]
		}

		mods = append(mods, ModEntry{Loader: loader, Slug: slug})
	}
	return mods, nil
}
