package verify

import (
	"fmt"

	"github.com/w1lam/Raw-Mod-Installer/internal/errors"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
)

func ReconcileActiveWithManifest(m *manifest.Manifest, found []manifest.InstalledPackage) {
	for _, pkg := range found {
		_, ok := m.InstalledPackages[pkg.Type][pkg.Name]

		if !ok {
			errors.ReportCtx("startup.reconile",
				fmt.Errorf("active package not in manifest"),
				map[string]string{
					"name": pkg.Name,
					"type": string(pkg.Type),
				},
			)
			m.InstalledPackages[pkg.Type][pkg.Name] = pkg
		}

		// LAST SCANNED WINS MIGHT NEED CHANGING IN THE FUTURE
		m.EnabledPackages[pkg.Type] = pkg.Name
	}
}
