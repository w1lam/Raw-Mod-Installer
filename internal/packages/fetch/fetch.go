package packages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
)

// GetAllAvailablePackages gets all available packages from github repo
func GetAllAvailablePackages() (packages.AvailablePackages, error) {
	subFldrs, err := scanPackagesFolder()
	if err != nil {
		return nil, err
	}

	availablePackages := make(packages.AvailablePackages)

	for _, fldr := range subFldrs {
		resp, err := getAvailablePackages(fldr)
		if err != nil {
			return nil, err
		}

		availablePackages[resp.Type] = resp.Pkgs
	}

	return availablePackages, nil
}

type pkgFetchResponse struct {
	Type packages.PackageType
	Pkgs map[string]packages.ResolvedPackage
}

// getAvailablePackages gets all available packages of specified type which is the name of the subfolder inside packages folder in github repo
func getAvailablePackages(fldrName string) (pkgFetchResponse, error) {
	req := fmt.Sprintf("%scontents/packages/%s", netcfg.GithubRepo, fldrName)

	resp, err := http.Get(req)
	if err != nil {
		return pkgFetchResponse{}, err
	}
	defer resp.Body.Close()

	var respJSON []packages.GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return pkgFetchResponse{}, err
	}

	pkgType := ""
	resolvedPackages := make(map[string]packages.ResolvedPackage)
	for _, mp := range respJSON {
		resolved, err := resolvePackage(mp.RawURL)
		if err != nil {
			return pkgFetchResponse{}, err
		}
		resolvedPackages[mp.Name] = resolved
		if pkgType != string(resolved.Type) {
			pkgType = string(resolved.Type)
		}
	}

	return pkgFetchResponse{
		Type: packages.PackageType(pkgType),
		Pkgs: resolvedPackages,
	}, nil
}
