package lists

import (
	"encoding/json"
	"net/http"

	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

// scanPackagesFolder scans the packages folder in github repo and returns a slice of the urls to each subfolder
func scanPackagesFolder() ([]string, error) {
	resp, err := http.Get(netcfg.GithubPackages)
	if err != nil {
		return nil, err
	}

	var decodedResp GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&decodedResp); err != nil {
		return nil, err
	}
}
