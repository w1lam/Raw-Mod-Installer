package install

import (
	"fmt"

	"github.com/w1lam/Packages/pkg/download"
)

func RenderProgress(ch <-chan download.Progress) (successFiles []string, failedFiles []string) {
	for p := range ch {
		switch p.Status {
		case "downloading":
			fmt.Print("\n ◌ ", p.File, "...")

		case "success":
			fmt.Print("\n ● ", p.File, " ✓")
			successFiles = append(successFiles, p.File)

		case "failure":
			fmt.Print("\n ✗ ", p.File, ": ", p.Err)
			failedFiles = append(failedFiles, p.File)
		}
	}
	return
}

func SimpleProgress(done, total int, currentFile string) {
	fmt.Printf("%d/%d: %s...\n", done, total, currentFile)
}
