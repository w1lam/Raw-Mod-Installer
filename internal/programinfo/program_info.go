// Package programinfo has information about the program
package programinfo

type ProgramInformation struct {
	Version   string
	McVersion string
}

var ProgramInfo = ProgramInformation{
	Version:   "ALPHA 0.0.1",
	McVersion: "1.21.10",
}
