package ui

import (
	"fmt"
	"strings"

	"github.com/w1lam/Packages/download"
	"github.com/w1lam/Raw-Mod-Installer/internal/config"
	"github.com/w1lam/Raw-Mod-Installer/internal/downloader"
)

func RenderProgress(ch <-chan download.Progress, total int) (successFiles []string, failedFiles []string) {
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

func RenderProgressBar(ch <-chan downloader.Progress, total int) (successFiles []string, failedFiles []string) {
	gap := config.Style.Width - 2
	margin := config.Style.Margin * 2
	marg := strings.Repeat(" ", config.Style.Margin)

	success := 0
	failed := 0

	inner := gap - margin
	inner = max(inner, 0)

	fmt.Printf("\r%sDownloading Mods...\n", marg)
	fmt.Printf("%s[%s]%s", marg, strings.Repeat(" ", inner), marg)

	for p := range ch {
		switch p.Status {
		case "success":
			success++
			successFiles = append(successFiles, p.File)
		case "failure":
			failed++
			failedFiles = append(failedFiles, p.File)
		}

		processed := success + failed
		if total == 0 {
			total = 1
		}

		filled := int(float64(processed) / float64(total) * float64(gap))
		if filled > inner {
			filled = inner
		}

		fmt.Printf("\r%s[%s%s]%s", marg, strings.Repeat("▰", filled), strings.Repeat(" ", inner-filled), marg)
	}

	return
}

func SimpleProgress(done, total int, currentFile string) {
	fmt.Printf("\r%d/%d: %s...\n", done, total, currentFile)
}
