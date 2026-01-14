package app

import (
	"fmt"
	"os"
	"time"

	"github.com/w1lam/Packages/tui"
)

func Exit() {
	tui.ClearScreenRaw()
	fmt.Printf("Exiting...")
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
}
