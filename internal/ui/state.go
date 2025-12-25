package ui

import (
	"os"
)

// OLD SHIT GARBAGE
type State int

// OLD SHIT GARBAGE
const (
	_ State = iota
	StateNotInstalled
	StateUpdateFound
	StateUpToDate
)

// OLD SHIT GARBAGE
func GetState() (State, error) {
	if _, err := os.Stat("BLANK"); err != nil {
		return StateNotInstalled, nil
	} else {

		// Update check
		upToDate, err := false, err
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
