package main

import (
	"os"

	"github.com/johnmanjiro13/gh-bump/bump"
	"github.com/johnmanjiro13/gh-bump/cmd"
	"github.com/johnmanjiro13/gh-bump/gh"
)

type exitCode int

const (
	// exitStatusOK is status code zero
	exitStatusOK exitCode = iota
	// exitStatusError is status code non-zero
	exitStatusError
)

func main() {
	os.Exit(int(run()))
}

func run() exitCode {
	ghCLI := gh.New()
	bumper := bump.New(ghCLI)
	rootCmd := cmd.New(bumper)
	if err := rootCmd.Execute(); err != nil {
		return exitStatusError
	}
	return exitStatusOK
}
