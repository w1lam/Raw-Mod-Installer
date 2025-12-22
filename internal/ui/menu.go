// Package ui provides functions to handle menu states and user input.
package ui

import (
	"fmt"
	"os"
	"time"
)

func ExitProgram() {
	fmt.Print("\n\n\n\n\n\n\n\nExiting...\n")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}
