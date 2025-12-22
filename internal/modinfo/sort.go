package modinfo

import (
	"sort"
	"strings"

	"github.com/w1lam/Packages/pkg/modrinth"
)

// ModInfoList is a list of ModInfo objects.
type ModInfoList []modrinth.ModInfo

// SortByCategory sorts the ModInfoList by the first category of each mod.
func (list ModInfoList) SortByCategory() ModInfoList {
	sorted := make(ModInfoList, len(list))
	copy(sorted, list)

	sort.Slice(sorted, func(i, j int) bool {
		ci, cj := "", ""

		if len(sorted[i].Category) > 0 {
			ci = sorted[i].Category[0]
		}
		if len(sorted[j].Category) > 0 {
			cj = sorted[j].Category[0]
		}

		if ci != cj {
			return ci < cj
		}

		// secondary sort (CRITICAL)
		return sorted[i].Title < sorted[j].Title
	})

	return sorted
}

// SortByName sorts the ModInfoList by mod title in a case-insensitive manner.
func (list ModInfoList) SortByName() ModInfoList {
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Title) < strings.ToLower(list[j].Title)
	})

	return list
}
