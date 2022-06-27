package main

import (
	"os"

	"github.com/johnmanjiro13/gh-bump/bump"
	"github.com/johnmanjiro13/gh-bump/cmd"
	"github.com/johnmanjiro13/gh-bump/gh"
)

const (
	// exitStatusOK is status code zero
	exitStatusOK int = iota
	// exitStatusError is status code non-zero
	exitStatusError
)

func main() {
	os.Exit(run())
}

func run() int {
	ghCLI := gh.New()
	bumper := bump.New(ghCLI)
	rootCmd := cmd.New(bumper)
	if err := rootCmd.Execute(); err != nil {
		return exitStatusError
	}
	return exitStatusOK
}
