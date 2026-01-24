package ui

import (
	"fmt"
)

// Progress represents the download progress of a file
type DownloaderProgress struct {
	File   string
	Status string // "downloading", "success", "failure"
	Err    error
}

// RenderProgress simplified version of progress bar, pure cli
func RenderDownloaderProgress(ch <-chan DownloaderProgress, total int) (successFiles []string, failedFiles []string) {
	success := 0
	failed := 0
	fmt.Println("[Downloading Mods: 0/", total, ", 0%]")

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

		procent := int(float64(processed) / float64(total) * float64(100))
		if procent > 100 {
			procent = 100
		}

		fmt.Println("[Downloading Mods: ", success, "/", total, ", ", procent, "%]")
	}

	return
}

// func RenderProgressBar(ch <-chan DownloaderProgress, total int) (successFiles []string, failedFiles []string) {
// 	gap := config.Style.Width - 2
// 	margin := config.Style.Margin * 2
// 	marg := strings.Repeat(" ", config.Style.Margin)
//
// 	success := 0
// 	failed := 0
//
// 	inner := gap - margin
// 	inner = max(inner, 0)
//
// 	fmt.Printf("\r%sDownloading Mods...\n", marg)
// 	fmt.Printf("%s[%s]%s", marg, strings.Repeat(" ", inner), marg)
//
// 	for p := range ch {
// 		switch p.Status {
// 		case "success":
// 			success++
// 			successFiles = append(successFiles, p.File)
// 		case "failure":
// 			failed++
// 			failedFiles = append(failedFiles, p.File)
// 		}
//
// 		processed := success + failed
// 		if total == 0 {
// 			total = 1
// 		}
//
// 		filled := int(float64(processed) / float64(total) * float64(gap))
// 		if filled > inner {
// 			filled = inner
// 		}
//
// 		fmt.Printf("\r%s[%s%s]%s", marg, strings.Repeat("▰", filled), strings.Repeat(" ", inner-filled), marg)
// 	}
//
// 	return
// }

func SimpleProgress(done, total int, currentFile string) {
	fmt.Printf("\r%d/%d: %s...\n", done, total, currentFile)
}
