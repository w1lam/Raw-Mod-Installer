package lists

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

type AvailablePackages struct {
	ModPacks        map[string]ResolvedPackage `json:"modPacks"`
	ResourceBundles map[string]ResolvedPackage `json:"resourcePacks"`
}

// GetAllAvailablePackages gets all available packages
func GetAllAvailablePackages() (AvailablePackages, error) {
	modPacks, err := GetAvailablePackages(PackageModPack)
	if err != nil {
		return AvailablePackages{}, err
	}

	resourceBundles, err := GetAvailableResourceBundles()
	if err != nil {
		return AvailablePackages{}, err
	}

	return AvailablePackages{
		ModPacks:        modPacks,
		ResourceBundles: resourceBundles,
	}, nil
}

// GetAvailablePackages gets all available packages of specified type(PackageModPack, PackageResourceBundle)
func GetAvailablePackages(packageType PackageType) (map[string]ResolvedPackage, error) {
	dest := ""

	switch packageType {
	case PackageModPack:
		dest = "modpacks"
	case PackageResourceBundle:
		dest = "resourcebundles"
	}

	if dest == "" {
		return nil, fmt.Errorf("no package type specified")
	}

	req := fmt.Sprintf("%scontents/%s", netcfg.GithubPackages, dest)

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
		resolved, err := ResolveModPack(mp.RawURL)
		if err != nil {
			return nil, err
		}
		resolvedModPacks[mp.Name] = resolved
	}

	return resolvedModPacks, nil
}
