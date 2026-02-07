package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

func GetAllAvailablePackages() (AvailablePackages, error) {
	folders, err := scanPackagesFolder()
	if err != nil {
		return nil, err
	}

	out := make(AvailablePackages)

	for _, folder := range folders {
		pkgs, err := getPackagesFromFolder(folder)
		if err != nil {
			return nil, err
		}

		for _, pkg := range pkgs {
			if out[pkg.Type] == nil {
				out[pkg.Type] = make(map[string]ResolvedPackage)
			}
			out[pkg.Type][pkg.Name] = pkg
		}
	}

	return out, nil
}

func getPackagesFromFolder(folder string) ([]ResolvedPackage, error) {
	var items []GithubContentResponse

	url := fmt.Sprintf("%scontents/packages/%s", netcfg.GithubRepo, folder)
	if err := githubGetJSON(url, &items); err != nil {
		return nil, err
	}

	var result []ResolvedPackage
	for _, it := range items {
		if it.Type != "file" {
			continue
		}

		pkg, err := resolvePackageJSON(it.RawURL)
		if err != nil {
			return nil, fmt.Errorf("folder %s: %w", folder, err)
		}

		result = append(result, pkg)
	}
	return result, nil
}

var githubHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func githubGetJSON(url string, out any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "MyModInstaller/1.0")

	resp, err := githubHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github request failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}
	return nil
}

func scanPackagesFolder() ([]string, error) {
	var items []GithubContentResponse

	url := netcfg.GithubRepo + "contents/packages"
	if err := githubGetJSON(url, &items); err != nil {
		return nil, err
	}

	var folders []string
	for _, it := range items {
		if it.Type == "dir" {
			folders = append(folders, it.Name)
		}
	}

	if len(folders) == 0 {
		return nil, fmt.Errorf("no package folder found")
	}

	return folders, nil
}

func resolvePackageJSON(url string) (ResolvedPackage, error) {
	var pkg ResolvedPackage

	if err := githubGetJSON(url, &pkg); err != nil {
		return pkg, err
	}

	if pkg.Name == "" {
		return pkg, fmt.Errorf("package json missing name")
	}

	if pkg.ListSource == "" {
		pkg.ListSource = url
	}

	return pkg, nil
}
