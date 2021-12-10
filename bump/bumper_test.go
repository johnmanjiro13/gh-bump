package bump

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	repoDocs = `name:   johnmanjiro13/gh-bump
description:    gh extension for bumping version of a repository`
	tagList = `v0.2.1  Latest  v0.2.1  2021-12-08T04:19:16Z`
)

type mockGh struct{}

func (g *mockGh) ViewRepository() (sout, eout bytes.Buffer, err error) {
	sout.WriteString(repoDocs)
	return
}

func (g *mockGh) ListRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	sout.WriteString(tagList)
	return
}

func (g *mockGh) ViewRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	return
}

func (g *mockGh) CreateRelease(version string, repo string, isCurrent bool) (sout, eout bytes.Buffer, err error) {
	sout.WriteString(version)
	return
}

func TestBumper_WithRepository(t *testing.T) {
	tests := map[string]struct {
		repository     string
		wantRepository string
		wantIsCurrent  bool
	}{
		"repository was given": {
			repository:     "johnmanjiro13/gh-bump",
			wantRepository: "johnmanjiro13/gh-bump",
			wantIsCurrent:  false,
		},
		"repository was not given": {
			repository:     "",
			wantRepository: "johnmanjiro13/gh-bump",
			wantIsCurrent:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &bumper{gh: &mockGh{}}
			assert.NoError(t, b.WithRepository(tt.repository))

			assert.Equal(t, tt.wantRepository, b.repository)
			assert.Equal(t, tt.wantIsCurrent, b.isCurrent)
		})
	}
}

func TestBumper_ResolveRepository(t *testing.T) {
	b := &bumper{gh: &mockGh{}}
	got, err := b.resolveRepository()
	assert.NoError(t, err)
	assert.Equal(t, "johnmanjiro13/gh-bump", got)
}

func TestBumper_listReleases(t *testing.T) {
	b := &bumper{gh: &mockGh{}}
	got, err := b.listReleases()
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Tags:\n%s", tagList), got)
}

func TestBumper_createRelease(t *testing.T) {
	b := &bumper{gh: &mockGh{}}
	got, err := b.createRelease("v1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, "v1.0.0", got)
}
