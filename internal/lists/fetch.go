package lists

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

// GetAllAvailablePackages gets all available packages from github repo
func GetAllAvailablePackages() (AvailablePackages, error) {
	paths, err := scanPackagesFolder()
	if err != nil {
		return nil, err
	}

	availablePackages := make(AvailablePackages)

	for _, path := range paths {
		packageType := strings.TrimPrefix(path, "packages/")

		packages, err := getAvailablePackages(packageType)
		if err != nil {
			return nil, err
		}

		availablePackages[packageType] = packages
	}

	return availablePackages, nil
}

// getAvailablePackages gets all available packages of specified type which is the name of the subfolder inside packages folder in github repo
func getAvailablePackages(packageType string) (map[string]ResolvedPackage, error) {
	req := fmt.Sprintf("%scontents/packages/%s", netcfg.GithubRepo, packageType)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respJSON []GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return nil, err
	}

	resolvedModPacks := make(map[string]ResolvedPackage)
	for _, mp := range respJSON {
		resolved, err := resolvePackage(packageType, mp.RawURL)
		if err != nil {
			return nil, err
		}
		resolvedModPacks[mp.Name] = resolved
	}

	return resolvedModPacks, nil
}
