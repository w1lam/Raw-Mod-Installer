package actions

import (
	"github.com/w1lam/Raw-Mod-Installer/internal/installer"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/services"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
)

func InstallPackage(pkg packages.Pkg) error {
	var resolvedPkg packages.ResolvedPackage

	state.Get().Read(func(s *state.State) {
		ap := s.AvailablePackages()
		if ap == nil {
			return
		}
		resolvedPkg = (ap)[pkg.Type][pkg.Name]
	})

	plan := installer.InstallPlan{
		RequestedPackage: resolvedPkg,
		BackupPolicy:     services.BackupIfExists,
	}

	err := installer.PackageInstaller(plan)
	if err != nil {
		return err
	}

	return services.EnablePackage(packages.Pkg{Name: plan.RequestedPackage.Name, Type: plan.RequestedPackage.Type})
}
