// Package verify holds verification functions
package verify

import (
	"fmt"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/errors"
	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

func VerifyAndReconcile(path *paths.Paths) {
	gs := state.Get()

	var (
		m     *manifest.Manifest
		found []manifest.InstalledPackage
	)

	gs.Read(func(s *state.State) {
		m = s.Manifest()
	})

	var err error
	found, err = ScanActitvePackages(path)
	if err != nil {
		errors.Report("startup.scan", err)
	}

	ReconcileActiveWithManifest(m, found)

	verifyManifestPackages()

	_ = gs.Write(func(s *state.State) error {
		return s.Manifest().Save()
	})
}

func verifyManifestPackages() {
	gs := state.Get()

	gs.Read(func(s *state.State) {
		for pkgType, pkgs := range s.Manifest().InstalledPackages {
			for name := range pkgs {
				ok, err := VerifyPackageIntegrity(packages.Pkg{
					Type: pkgType,
					Name: name,
				})
				if err != nil {
					errors.ReportCtx(
						"startup.verify",
						err,
						map[string]string{
							"name": name,
							"type": string(pkgType),
						},
					)
					continue
				}

				if !ok {
					errors.ReportCtx(
						"startup.verify",
						fmt.Errorf("package integrity mismath"),
						map[string]string{
							"name": name,
							"type": string(pkgType),
						},
					)
				}
			}
		}
	})
}

func VerifyPackageIntegrity(pkg packages.Pkg) (bool, error) {
	gs := state.Get()

	var (
		iPkg    manifest.InstalledPackage
		enabled bool
		found   bool
	)

	gs.Read(func(s *state.State) {
		if p, ok := s.Manifest().InstalledPackages[pkg.Type][pkg.Name]; ok {
			iPkg = p
			found = true
		}
		enabled = s.Manifest().EnabledPackages[pkg.Type] == pkg.Name
	})

	if !found {
		return false, fmt.Errorf("package not installed: %s", pkg.Name)
	}

	var path string
	if enabled {
		path = iPkg.ActivePath
	} else {
		path = iPkg.StoragePath
	}

	hash, err := filesystem.ComputeDirHash(path)
	if err != nil {
		return false, err
	}

	pkgIDPath := filepath.Join(path, pkg.Name+".id.json")
	pkgID, err := manifest.ReadPackageIDFile(pkgIDPath)
	if err != nil {
		return false, err
	}

	return hash == pkgID.Hash, nil
}
