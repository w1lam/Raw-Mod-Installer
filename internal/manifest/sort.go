package manifest

import (
	"sort"
	"strings"
)

// SortedByCategory sorts the ModList by the first category of each mod.
func (list ModList) SortedByCategory() ModList {
	sorted := make(ModList, len(list))
	copy(sorted, list)

	sort.Slice(sorted, func(i, j int) bool {
		ci, cj := "", ""

		if len(sorted[i].Categories) > 0 {
			ci = sorted[i].Categories[0]
		}
		if len(sorted[j].Categories) > 0 {
			cj = sorted[j].Categories[0]
		}

		if ci != cj {
			return ci < cj
		}

		return strings.ToLower(sorted[i].Title) <
			strings.ToLower(sorted[j].Title)
	})

	return sorted
}

// SortedByName sorts the ModList by mod title in a case-insensitive manner.
func (list ModList) SortedByName() ModList {
	sorted := make(ModList, len(list))
	copy(sorted, list)

	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].Title) <
			strings.ToLower(sorted[j].Title)
	})

	return sorted
}
