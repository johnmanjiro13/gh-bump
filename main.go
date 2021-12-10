package main

import (
	"os"

	"github.com/johnmanjiro13/gh-bump/bump"
	"github.com/johnmanjiro13/gh-bump/cmd"
	"github.com/johnmanjiro13/gh-bump/gh"
)

func main() {
	os.Exit(run())
}

func run() int {
	ghCLI := gh.New()
	bumper := bump.New(ghCLI)
	rootCmd := cmd.NewCmd(bumper)
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
