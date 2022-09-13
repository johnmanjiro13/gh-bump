package main

import (
	"os"

	bump "github.com/johnmanjiro13/gh-bump"
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
	gh := bump.NewGh()
	bumper := bump.NewBumper(gh)
	rootCmd := bump.NewRootCmd(bumper)
	if err := rootCmd.Execute(); err != nil {
		return exitStatusError
	}
	return exitStatusOK
}
