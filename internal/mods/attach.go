package mods

import (
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

func AttachLocalVersions(
	resolved resolve.ResolvedModList,
	localMods []LocalMod,
) {
	for i := range resolved {
		for _, local := range localMods {
			if resolve.NormalizeID(resolved[i].FabricID) == resolve.NormalizeID(local.ID) {
				resolved[i].LocalVer = local.Version
				break
			}

			if resolve.NormalizeID(resolved[i].Slug) == resolve.NormalizeID(local.ID) {
				resolved[i].LocalVer = local.Version
				break
			}

			if strings.Contains(resolve.NormalizeID(local.File), resolve.NormalizeID(resolved[i].Slug)) {
				resolved[i].LocalVer = local.Version
				break
			}

		}
	}
}
