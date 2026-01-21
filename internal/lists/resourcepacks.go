package lists

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/netcfg"
)

type ResolvedResourceBundle struct {
	Name        string
	ListSource  string
	ListVersion string
	McVersion   string
	Description string
	Hash        string
	Entries     []modrinth.ModrinthListEntry
}

func GetAvailableResourceBundles() (map[string]ResolvedResourceBundle, error) {
	req := fmt.Sprintf("%scontents/resourcebundles", netcfg.GithubRepoAPI)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respJSON []GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return nil, err
	}

	resolvedResourceBundles := make(map[string]ResolvedResourceBundle)
	for _, rb := range respJSON {
		resolved, err := ResolveResourceBundle(rb.RawURL)
		if err != nil {
			return nil, err
		}
		resolvedResourceBundles[rb.Name] = resolved
	}

	return resolvedResourceBundles, nil
}

func ResolveResourceBundle(url string) (ResolvedResourceBundle, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ResolvedResourceBundle{}, err
	}
	defer resp.Body.Close()

	var resolvedResourceBundle ResolvedResourceBundle

	var entries []modrinth.ModrinthListEntry
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if name, ok := strings.CutPrefix(line, "Name: "); ok {
			resolvedResourceBundle.Name = name
			continue
		}

		if version, ok := strings.CutPrefix(line, "Version: "); ok {
			resolvedResourceBundle.ListVersion = version
			continue
		}

		if mcVersion, ok := strings.CutPrefix(line, "McVersion: "); ok {
			resolvedResourceBundle.McVersion = mcVersion
			continue
		}

		if description, ok := strings.CutPrefix(line, "Description: "); ok {
			resolvedResourceBundle.Description = description
			continue
		}

		if hash, ok := strings.CutPrefix(line, "Hash: "); ok {
			resolvedResourceBundle.Hash = hash
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "@", 2)
		entry := modrinth.ModrinthListEntry{
			Slug: parts[0],
		}

		if len(parts) == 2 {
			entry.PinnedVer = parts[1]
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return ResolvedResourceBundle{}, err
	}

	resolvedResourceBundle.Entries = entries

	resolvedResourceBundle.ListSource = url

	return resolvedResourceBundle, nil
}
