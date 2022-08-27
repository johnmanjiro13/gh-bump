package gh

import (
	"bytes"

	"github.com/cli/go-gh"

	"github.com/johnmanjiro13/gh-bump/bump"
)

type ghCLI struct{}

func New() *ghCLI {
	return &ghCLI{}
}

func (g *ghCLI) ViewRepository() (sout, eout bytes.Buffer, err error) {
	return gh.Exec("repo", "view")
}

func (g *ghCLI) ListRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = gh.Exec("release", "list")
	} else {
		sout, eout, err = gh.Exec("release", "list", "-R", repo)
	}
	return
}

func (g *ghCLI) ViewRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	if isCurrent {
		sout, eout, err = gh.Exec("release", "view")
	} else {
		sout, eout, err = gh.Exec("release", "view", "-R", repo)
	}
	return
}

func (g *ghCLI) CreateRelease(version string, repo string, isCurrent bool, option *bump.ReleaseOption) (sout, eout bytes.Buffer, err error) {
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
	sout, eout, err = gh.Exec(args...)
	return
}
