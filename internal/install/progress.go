package install

import (
	"fmt"
	"os"

	dl "github.com/w1lam/Packages/pkg/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/modrinthsvc"
)

func DisplayDownloadProgress(progressCh <-chan dl.Progress, fopts string, listURL string) {
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
		version, _ := modrinthsvc.GetRemoteVersion(listURL)

		err := os.WriteFile("ver.txt", []byte(version), 0o755)
		if err != nil {
			return
		}

		fmt.Printf("\n\nAll %d %s installed successfully! ✓", success, fopts)
	} else {
		fmt.Printf("\n\n%d %s failed to install! ✗ \n", failures, fopts)
	}
}
