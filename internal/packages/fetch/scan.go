package packages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
)

// scanPackagesFolder scans the packages folder in github repo and returns a slice of the names of each subfolder
func scanPackagesFolder() ([]string, error) {
	req := fmt.Sprintf("%scontents/packages", netcfg.GithubRepo)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}

	var decodedResp []packages.GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&decodedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal packages folder resp: %w", err)
	}

	if len(decodedResp) == 0 {
		return nil, fmt.Errorf("scanPackagesFolder: no package subfolders found")
	}

	names := []string{}
	for _, sub := range decodedResp {
		if sub.Type != "dir" {
			continue
		}
		names = append(names, sub.Name)
	}

	return names, nil
}
