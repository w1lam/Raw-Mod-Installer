// Package lists provides functions to fetch and parse mod lists from remote URLs.
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

// ResolvedModPack is a resolved mod pack
type ResolvedModPack struct {
	Name        string
	ListSource  string
	Env         string
	ListVersion string
	McVersion   string
	Loader      string
	Description string
	Entries     []modrinth.ModrinthListEntry
}

func GetAvailableModPacks() (map[string]ResolvedModPack, error) {
	req := fmt.Sprintf("%s/contents/modpacks", netcfg.GithubRepoAPI)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respJSON []GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return nil, err
	}

	resolvedModPacks := make(map[string]ResolvedModPack)
	for _, mp := range respJSON {
		resolved, err := ResolveModPack(mp.RawURL)
		if err != nil {
			return nil, err
		}
		resolvedModPacks[mp.Name] = resolved
	}

	return resolvedModPacks, nil
}

// ResolveModPack resolves a modpack list with modpack version, mcversion and loader
func ResolveModPack(modPackURL string) (ResolvedModPack, error) {
	resp, err := http.Get(modPackURL)
	if err != nil {
		return ResolvedModPack{}, err
	}
	defer resp.Body.Close()

	var resolvedModPack ResolvedModPack

	var entries []modrinth.ModrinthListEntry
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if name, ok := strings.CutPrefix(line, "# Name: "); ok {
			resolvedModPack.Name = name
			continue
		}

		if env, ok := strings.CutPrefix(line, "# Env: "); ok {
			resolvedModPack.Env = env
			continue
		}

		if version, ok := strings.CutPrefix(line, "# Version: "); ok {
			resolvedModPack.ListVersion = version
			continue
		}

		if mcVersion, ok := strings.CutPrefix(line, "# McVersion: "); ok {
			resolvedModPack.McVersion = mcVersion
			continue
		}

		if loader, ok := strings.CutPrefix(line, "# Loader: "); ok {
			resolvedModPack.Loader = loader
			continue
		}

		if description, ok := strings.CutPrefix(line, "# Description: "); ok {
			resolvedModPack.Description = description
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
		return ResolvedModPack{}, err
	}

	resolvedModPack.Entries = entries

	resolvedModPack.ListSource = modPackURL

	return resolvedModPack, nil
}
