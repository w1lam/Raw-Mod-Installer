package manifest

import (
	"github.com/w1lam/Packages/pkg/utils"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

func Exists() bool {
	path, err := paths.Resolve()
	if err != nil {
		return false
	}
	return utils.CheckFileExists(path.ManifestPath)
}
