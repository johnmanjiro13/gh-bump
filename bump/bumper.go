package bump

import (
	"bytes"
	"strings"

	"github.com/johnmanjiro13/gh-bump/cmd"
)

type Gh interface {
	ViewRepository() (sout, eout bytes.Buffer, err error)
}

type bumper struct {
	gh         Gh
	repository string
	isCurrent  bool
}

func New(gh Gh) cmd.Bumper {
	return &bumper{gh: gh}
}

func (b *bumper) Bump() error {
	return nil
}

func (b *bumper) WithRepository(repository string) error {
	if repository != "" {
		b.repository = repository
		return nil
	}

	repo, err := b.resolveRepository()
	if err != nil {
		return err
	}
	b.repository = repo
	b.isCurrent = true
	return nil
}

func (b *bumper) resolveRepository() (string, error) {
	sout, _, err := b.gh.ViewRepository()
	if err != nil {
		return "", err
	}
	viewOut := strings.Split(sout.String(), "\n")[0]
	repo := strings.TrimSpace(strings.Split(viewOut, ":")[1])
	return repo, nil
}
