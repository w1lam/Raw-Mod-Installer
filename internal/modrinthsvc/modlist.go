package modrinthsvc

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/w1lam/Packages/pkg/fetch"
	"github.com/w1lam/Raw-Mod-Installer/internal/modlist"
)

func GetModEntryList(modListURL string) ([]modlist.ModEntry, error) {
	rawList, err := fetch.GetList(modListURL)
	if err != nil {
		return nil, err
	}

	modEntryList, err1 := modlist.ParseModList(rawList)
	if err1 != nil {
		return nil, err1
	}

	return modEntryList, nil
}

func GetRemoteVersion(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		cutLine, _ := strings.CutPrefix(line, "# version:")
		return strings.TrimSpace(cutLine), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}
