// Package modpack provides functions to fetch and parse mod lists from remote URLs.
package modpack

import (
	"bufio"
	"net/http"
	"strings"
)

// ResolvedModPack is a resolved mod pack
type ResolvedModPack struct {
	Name        string
	ListSource  string
	ListVersion string
	McVersion   string
	Loader      string
	Description string
	Slugs       []string
}

type availableModPack struct {
	Name string
	URL  string
}

// GetAvailableModPacks gets the url for a modpack from a list
func GetAvailableModPacks(modPacksListURL string) (map[string]ResolvedModPack, error) {
	resp, err := http.Get(modPacksListURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resolvedModPacks := make(map[string]ResolvedModPack)
	var modPack availableModPack

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		splitLines := strings.Split(line, "@")

		modPack.Name = splitLines[0]
		modPack.URL = splitLines[1]

		resolved, err := ResolveModPack(modPack)
		if err != nil {
			return nil, err
		}

		resolvedModPacks[splitLines[0]] = resolved
	}

	return resolvedModPacks, nil
}

// ResolveModPack resolves a modpack list with modpack version, mcversion and loader
func ResolveModPack(modPack availableModPack) (ResolvedModPack, error) {
	resp, err := http.Get(modPack.URL)
	if err != nil {
		return ResolvedModPack{}, err
	}
	defer resp.Body.Close()

	var resolvedModPack ResolvedModPack

	var slugs []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "#") {
			slugs = append(slugs, line)
		}

		if name, ok := strings.CutPrefix(line, "# Name: "); ok {
			resolvedModPack.Name = name
		}

		if version, ok := strings.CutPrefix(line, "# Version: "); ok {
			resolvedModPack.ListVersion = version
		}

		if mcVersion, ok := strings.CutPrefix(line, "# McVersion: "); ok {
			resolvedModPack.McVersion = mcVersion
		}

		if loader, ok := strings.CutPrefix(line, "# Loader: "); ok {
			resolvedModPack.Loader = loader
		}

		if description, ok := strings.CutPrefix(line, "# Description: "); ok {
			resolvedModPack.Description = description
		}
	}

	if err := scanner.Err(); err != nil {
		return ResolvedModPack{}, err
	}

	resolvedModPack.Slugs = slugs

	resolvedModPack.ListSource = modPack.URL

	return resolvedModPack, nil
}
