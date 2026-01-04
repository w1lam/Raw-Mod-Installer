package install

import (
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/resolve"
)

type installContext struct {
	FileToMod map[string]resolve.ResolvedMod
}

func buildInstallContext(mods []resolve.ResolvedMod) installContext {
	m := make(map[string]resolve.ResolvedMod)

	for _, mod := range mods {
		filename := filepath.Base(mod.DownloadURL)
		m[filename] = mod
	}

	return installContext{FileToMod: m}
}

func UpdateManifestInstalledVersions(
	manifest *manifest.Manifest,
	ctx installContext,
	successFiles []string,
) {
	for _, file := range successFiles {
		mod, ok := ctx.FileToMod[file]
		if !ok {
			continue
		}

		entry := manifest.Mods[mod.Slug]
		entry.InstalledVersion = mod.LatestVer
		manifest.Mods[mod.Slug] = entry
	}
}
