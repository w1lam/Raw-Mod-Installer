package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/w1lam/Packages/pkg/dl"
)

var (
	UserProfile, _      = os.UserHomeDir()
	ModFolderPath       = filepath.Join(UserProfile, "AppData", "Roaming", ".minecraft", "mods")
	ModBackupPath       = filepath.Join(UserProfile, "AppData", "Roaming", ".minecraft", "mods_old")
	VerFilePath         = filepath.Join(UserProfile, "AppData", "Roaming", ".minecraft", "mods", "ver.txt")
	ModListURL          = "https://raw.githubusercontent.com/w1lam/mods/refs/heads/main/mod-list.txt"
	TempModDownloadPath = filepath.Join(os.TempDir(), "temp-mod-downloads")
)

type State int

const (
	_ State = iota
	StateNotInstalled
	StateUpdateFound
	StateUpToDate
)

func GetState() (State, error) {
	if _, err := os.Stat(VerFilePath); err != nil {
		return StateNotInstalled, nil
	} else {

		// Update check
		upToDate, err := CheckForModlistUpdate()
		switch {
		case err != nil:
			return 0, err

		case upToDate:
			return StateUpdateFound, nil

		default:
			return StateUpToDate, nil
		}
	}
}

// User Input

func UserInput() (bool, error) {
	if err := keyboard.Open(); err != nil {
		return false, err
	}

	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return false, err
		}

		switch {
		case char == 'y' || char == 'Y':
			err := keyboard.Close()
			if err != nil {
				return false, err
			}
			return true, nil

		case char == 'n' || char == 'N':
			err := keyboard.Close()
			if err != nil {
				return false, err
			}
			return false, nil

		case key == keyboard.KeyEsc:
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...")
			os.Exit(0)

		case char == 'q':
			err := keyboard.Close()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n\nExiting...")
			os.Exit(0)
		}
	}
}

// Fabric Installation

func InstallFabric(URL string) error {
	fileName := filepath.Base(URL)
	if err := os.Chdir(os.TempDir()); err != nil {
		return err
	}
	err := dl.DownloadFile(fileName, URL)
	if err != nil {
		return err
	} else {
		javaPath, err := exec.LookPath("java.exe")
		if err != nil {
			return err
		} else {
			cmd := exec.Command(javaPath, "-jar", fileName, "client", "-mcversion", "1.21.10")
			err := cmd.Start()
			time.Sleep(time.Duration(1) * time.Second)
			if err != nil {
				return err
			} else {
				os.Remove(filepath.Join(UserProfile, "AppData", "Local", "Temp", fileName))
			}
		}
	}
	return nil
}

// Mod list Version Check Functions

func CheckForModlistUpdate() (bool, error) {
	if _, err := os.Stat(VerFilePath); err == nil {

		remoteVer, err := GetRemoteVersion(ModListURL)
		localVer := GetLocalVersion()

		switch {
		case err != nil:
			return false, err

		case remoteVer != localVer:
			return true, nil

		case remoteVer == localVer:
			return false, nil

		default:
			return false, err
		}
	}
	return false, nil
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

func GetLocalVersion() string {
	data, err := os.ReadFile(VerFilePath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// Progress Display and Download Functions

type Progress struct {
	File   string
	Status string // "downloading", "success", "failure"
	Err    error
}

func DisplayProgress(progressCh <-chan Progress, fopts string, listURL string) {
	success, failures, active := 0, 0, 0

	for p := range progressCh {
		switch p.Status {
		case "downloading":
			active++
			fmt.Print("\n ◌ ", p.File, "...")
		case "success":
			active--
			success++
			fmt.Print("\n ● ", p.File, " ✓")
		case "failure":
			active--
			failures++
			fmt.Print("\n ✗ ", p.File, ": ", p.Err)
		}
		fmt.Print(" [Active: ", active, " | Success: ", success, " | Failures: ", failures, "]")
	}

	if failures == 0 {
		version, _ := GetRemoteVersion(listURL)
		os.WriteFile("ver.txt", []byte(version), 0o755)
		fmt.Printf("\n\nAll %d %s installed successfully! ✓", success, fopts)
	} else {
		fmt.Printf("\n\n%d %s failed to install! ✗ \n", failures, fopts)
	}
}

func DownloadFromListConcurrent(list []string, progressCh chan<- Progress) {
	var wg sync.WaitGroup

	for _, uri := range list {
		wg.Add(1)
		go func(uri string) {
			defer wg.Done()
			fileName := filepath.Base(uri)
			progressCh <- Progress{File: fileName, Status: "downloading"}

			err := dl.DownloadFile(fileName, uri)
			if err != nil {
				progressCh <- Progress{File: fileName, Status: "failure", Err: err}
				return
			}
			progressCh <- Progress{File: fileName, Status: "success"}
		}(uri)
	}

	wg.Wait()
	close(progressCh)
}

// Modrinth API fetching and functions

type MRFile struct {
	URL     string `json:"url"`
	Primary bool   `json:"primary"`
}

type MRVersion struct {
	Files []MRFile `json:"files"`
}

func FetchLatestModrinthDownload(slug, mcVersion, loader string) (string, error) {
	fetch := func(params url.Values) ([]MRVersion, error) {
		base := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", slug)
		finalURL := base + "?" + params.Encode()

		req, _ := http.NewRequest("GET", finalURL, nil)
		req.Header.Set("User-Agent", "MyModInstaller/1.0")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var versions []MRVersion
		if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
			return nil, err
		}
		return versions, nil
	}

	params := url.Values{}
	params.Set("game_versions", fmt.Sprintf("[\"%s\"]", mcVersion))
	params.Set("loaders", fmt.Sprintf("[\"%s\"]", loader))

	versions, err := fetch(params)
	if err == nil && len(versions) > 0 && len(versions[0].Files) > 0 {
		return versions[0].Files[0].URL, nil
	}

	paramsFallback := url.Values{}
	paramsFallback.Set("loaders", fmt.Sprintf("[\"%s\"]", loader))

	versions, err = fetch(paramsFallback)
	if err == nil && len(versions) > 0 && len(versions[0].Files) > 0 {
		return versions[0].Files[0].URL, nil
	}

	return "", fmt.Errorf("no compatible or fallback versions found for %s", slug)
}

type ModEntry struct {
	Loader string
	Slug   string
}

func ParseModList(lines []string) ([]ModEntry, error) {
	mods := []ModEntry{}

	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}

		loader := "fabric"
		slug := line

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			loader = parts[0]
			slug = parts[1]
		}

		if strings.Contains(slug, "@") {
			slug = strings.Split(slug, "@")[0]
		}

		mods = append(mods, ModEntry{Loader: loader, Slug: slug})
	}
	return mods, nil
}

func FetchAllDownloadURLs(mods []ModEntry, mcVersion string) ([]string, error) {
	var urls []string

	for _, mod := range mods {
		url, err := FetchLatestModrinthDownload(mod.Slug, mcVersion, mod.Loader)
		if err != nil {
			return nil, fmt.Errorf("mod %s: %v", mod.Slug, err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func FetchAllConcurrent(
	mods []ModEntry,
	mcVersion string,
	progressFunc func(done, total int, currentMod string),
) ([]string, error) {
	total := len(mods)
	results := make([]string, total)
	errChan := make(chan error, total)

	var wg sync.WaitGroup
	var done int32 = 0

	for i, mod := range mods {
		wg.Add(1)

		go func(i int, mod ModEntry) {
			defer wg.Done()

			url, err := FetchLatestModrinthDownload(mod.Slug, mcVersion, mod.Loader)
			if err != nil {
				errChan <- fmt.Errorf("%s: %w", mod.Slug, err)
				return
			}

			results[i] = url

			atomic.AddInt32(&done, 1)
			progressFunc(int(done), total, mod.Slug)
		}(i, mod)
	}

	wg.Wait()
	close(errChan)

	var combined strings.Builder
	for e := range errChan {
		combined.WriteString(e.Error() + "\n")
	}

	if combined.Len() > 0 {
		return nil, errors.New(combined.String())
	}
	return results, nil
}

func SimpleProgress(done, total int, mod string) {
	fmt.Printf("Fetched %d/%d -> %s\n", done, total, mod)
}

// Main Mod Download Function

func DownloadMods(listURL string, fopts string) error {
	fmt.Printf("Fetching Mod List...\n")
	slugList, err := dl.GetList(listURL)
	if err != nil {
		return err
	}

	fmt.Printf("Parsing Mod List...\n")
	parsedList, err := ParseModList(slugList)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching Download URLs...\n")
	fetchedURLs, err := FetchAllConcurrent(parsedList, "1.21.10", SimpleProgress)
	if err != nil {
		return err
	}

	// Creates temp mod download dir
	err2 := os.MkdirAll(TempModDownloadPath, os.ModePerm)
	if err2 != nil {
		return fmt.Errorf("ERROR: failed to change to temp dir: %v", err2)
	}
	err3 := os.Chdir(TempModDownloadPath)
	if err3 != nil {
		return fmt.Errorf("ERROR: failed to change to temp mod download dir: %v", err3)
	}

	progressCh := make(chan Progress)

	var wg sync.WaitGroup

	// Start downloads concurrently
	wg.Go(func() {
		DownloadFromListConcurrent(fetchedURLs, progressCh)
	})

	// Handle UI printing in main goroutine
	DisplayProgress(progressCh, fopts, listURL)
	wg.Wait()

	os.Chdir(UserProfile)

	return nil
}

// Mod Backup and Restore Functions

func BackupModFolder() error {
	if _, err := os.Stat(ModFolderPath); err == nil {
		if _, err2 := os.Stat(ModBackupPath); err2 == nil {
			timestamp := time.Now().Format("20060102150405")

			err3 := os.Rename(ModBackupPath, ModBackupPath+"_"+timestamp)
			if err3 != nil {
				return fmt.Errorf("failed to backup existing mod backup folder: %v", err3)
			}
		}

		err := os.Rename(ModFolderPath, ModBackupPath)
		if err != nil {
			timestamp := time.Now().Format("20060102150405")

			err2 := os.Rename(ModFolderPath, ModBackupPath+"_"+timestamp+"bruh")
			if err2 != nil {
				return fmt.Errorf("failed to backup existing mod backup folder: %v", err2)
			}
		}
	}

	return nil
}

func RestoreModBackup() error {
	err := os.Rename(ModBackupPath, ModFolderPath)
	if err != nil {
		return fmt.Errorf("failed to restore backup mods: %v", err)
	}
	return nil
}

func UninstallMods() error {
	err := os.RemoveAll(ModFolderPath)
	if err != nil {
		return fmt.Errorf("failed to uninstall mods: %v", err)
	}
	return nil
}
