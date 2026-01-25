package packages

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/packages"
)

// resolveModPack resolves a package list
func resolvePackage(url string, pkgType packages.PackageType) (packages.ResolvedPackage, error) {
	resp, err := http.Get(url)
	if err != nil {
		return packages.ResolvedPackage{}, err
	}
	defer resp.Body.Close()

	var resolvedPackage packages.ResolvedPackage

	var entries []modrinth.ModrinthListEntry
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if name, ok := strings.CutPrefix(line, "Name: "); ok {
			resolvedPackage.Name = name
			continue
		}

		if env, ok := strings.CutPrefix(line, "Env: "); ok {
			resolvedPackage.Env = env
			continue
		}

		if version, ok := strings.CutPrefix(line, "Version: "); ok {
			resolvedPackage.ListVersion = version
			continue
		}

		if mcVersion, ok := strings.CutPrefix(line, "McVersion: "); ok {
			resolvedPackage.McVersion = mcVersion
			continue
		}

		if loader, ok := strings.CutPrefix(line, "Loader: "); ok {
			resolvedPackage.Loader = loader
			continue
		}

		if description, ok := strings.CutPrefix(line, "Description: "); ok {
			resolvedPackage.Description = description
			continue
		}

		if hash, ok := strings.CutPrefix(line, "Hash: "); ok {
			resolvedPackage.Hash = hash
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
		return packages.ResolvedPackage{}, err
	}

	resolvedPackage.Entries = entries

	resolvedPackage.Type = pkgType
	resolvedPackage.ListSource = url

	return resolvedPackage, nil
}
