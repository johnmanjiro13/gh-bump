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

func (g *gh) ListRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = runGh("release", "list")
	} else {
		sout, eout, err = runGh("release", "list", "-R", repo)
	}
	return
}

func (g *gh) ViewRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = runGh("release", "view")
	} else {
		sout, eout, err = runGh("release", "view", "-R", repo)
	}
	return
}

func (g *gh) CreateRelease(version string, repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = runGh("release", "create", version)
	} else {
		sout, eout, err = runGh("release", "create", version, "-R", repo)
	}
	return
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
