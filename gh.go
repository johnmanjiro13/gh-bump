package bump

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/cli/safeexec"
)

type ghCLI struct{}

func NewGh() *ghCLI {
	return &ghCLI{}
}

func (g *ghCLI) ViewRepository() (sout, eout *bytes.Buffer, err error) {
	return gh("repo", "view")
}

func (g *ghCLI) ListRelease(repo string, isCurrent bool) (sout, eout *bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = gh("release", "list")
	} else {
		sout, eout, err = gh("release", "list", "-R", repo)
	}
	return
}

func (g *ghCLI) ViewRelease(repo string, isCurrent bool) (sout, eout *bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = gh("release", "view")
	} else {
		sout, eout, err = gh("release", "view", "-R", repo)
	}
	return
}

func (g *ghCLI) CreateRelease(version string, repo string, isCurrent bool, option *ReleaseOption) (sout, eout *bytes.Buffer, err error) {
	args := []string{"release", "create", version}
	if len(option.AssetFiles) > 0 {
		args = append(args, option.AssetFiles...)
	}
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
	sout, eout, err = gh(args...)
	return
}

func gh(args ...string) (sout, eout *bytes.Buffer, err error) {
	sout = new(bytes.Buffer)
	eout = new(bytes.Buffer)
	bin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("could not find gh. err: %w", err)
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Stdout = sout
	cmd.Stderr = eout

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh. err: %w, eout: %s", err, eout.String())
		return
	}
	return
}
