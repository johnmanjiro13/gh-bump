package gh

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/cli/safeexec"

	"github.com/johnmanjiro13/gh-bump/bump"
)

type gh struct{}

func New() bump.Gh {
	return &gh{}
}

func (g *gh) ViewRepository() (sout, eout bytes.Buffer, err error) {
	return runGh("repo", "view")
}

func runGh(args ...string) (sout, eout bytes.Buffer, err error) {
	ghBin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("could not find gh. err: %w", err)
		return
	}

	cmd := exec.Command(ghBin, args...)
	cmd.Stdout = &sout
	cmd.Stderr = &eout

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh. err: %w, eout: %s", err, eout.String())
		return
	}
	return
}
