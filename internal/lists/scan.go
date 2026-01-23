package lists

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

// scanPackagesFolder scans the packages folder in github repo and returns a slice of the paths to each subfolder
func scanPackagesFolder() ([]string, error) {
	req := fmt.Sprintf("%s/contents/packages", netcfg.GithubRepo)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}

	var decodedResp []GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&decodedResp); err != nil {
		return nil, err
	}

	if len(decodedResp) == 0 {
		return nil, fmt.Errorf("scanPackagesFolder: no package subfolders found")
	}

	paths := []string{}
	for _, sub := range decodedResp {
		paths = append(paths, sub.Path)
	}

	return paths, nil
}
