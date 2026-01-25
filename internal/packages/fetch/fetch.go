package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

var folderToPkgType = map[string]packages.PackageType{
	"modpacks":        packages.PackageModPack,
	"resourcebundles": packages.PackageResourceBundle,
	"shaderbundles":   packages.PackageShaderBundle,
}

// GetAllAvailablePackages gets all available packages from github repo
func GetAllAvailablePackages() (packages.AvailablePackages, error) {
	subFldrs, err := scanPackagesFolder()
	if err != nil {
		return nil, err
	}

	availablePackages := make(packages.AvailablePackages)

	for _, fldr := range subFldrs {
		pkgs, err := getAvailablePackages(fldr)
		if err != nil {
			return nil, err
		}

		if len(pkgs) == 0 {
			continue
		}

		var pkgType packages.PackageType
		for _, p := range pkgs {
			pkgType = p.Type
			break
		}

		availablePackages[pkgType] = pkgs
	}

	return availablePackages, nil
}

// getAvailablePackages gets all available packages from subfolder in repo/packages
func getAvailablePackages(fldrName string) (map[string]packages.ResolvedPackage, error) {
	path, err := paths.Resolve()
	req := fmt.Sprintf("%scontents/packages/%s", netcfg.GithubRepo, fldrName)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respJSON []packages.GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return nil, err
	}

	// Map folder name to pkgType
	pkgType, ok := folderToPkgType[fldrName]
	if !ok {
		return nil, fmt.Errorf("package type of %s not found in index", fldrName)
	}

	resolvedPackages := make(map[string]packages.ResolvedPackage)
	for _, p := range respJSON {
		if p.Type != "file" {
			continue
		}

		resolved, err := resolvePackage(p.RawURL, pkgType)
		if err != nil {
			return nil, err
		}

		if resolved.Name == "" {
			return nil, fmt.Errorf("package %s jas no Name", resolved.Name)
		}

		if marshaled, err := json.MarshalIndent(resolved, "", " "); err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		} else {
			outFile := filepath.Join(path.ProgramFilesDir, fmt.Sprintf("%s.json", resolved.Name))
			err := os.WriteFile(outFile, marshaled, 0o755)
			if err != nil {
				return nil, fmt.Errorf("failed to write json file: %w", err)
			}
		}

		resolvedPackages[resolved.Name] = resolved
	}

	return resolvedPackages, nil
}
