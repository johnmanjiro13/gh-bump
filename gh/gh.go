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

func (g *gh) CreateRelease(version string, repo string, isCurrent bool, option *bump.ReleaseOption) (sout, eout bytes.Buffer, err error) {
	args := []string{"release", "create", version}
	if !isCurrent {
		args = append(args, []string{"-R", repo}...)
	}
	if option.IsDraft {
		args = append(args, "--draft")
	}
	if option.IsPrerelease {
		args = append(args, "--prerelease")
	}
	if option.DiscussionCategory != "" {
		args = append(args, []string{"--discussion-category", option.DiscussionCategory}...)
	}
	if option.GenerateNotes {
		args = append(args, "--generate-notes")
	}
	if option.Notes != "" {
		args = append(args, []string{"--notes", option.Notes}...)
	}
	if option.NotesFilename != "" {
		args = append(args, []string{"--notes-file", option.NotesFilename}...)
	}
	if option.Target != "" {
		args = append(args, []string{"--target", option.Target}...)
	}
	if option.Title != "" {
		args = append(args, []string{"-t", option.Title}...)
	}
	sout, eout, err = runGh(args...)
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
