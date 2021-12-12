package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBumper struct {
	repository   string
	isPrerelease bool
	target       string
	title        string
}

func (b *mockBumper) Bump() error {
	return nil
}

func (b *mockBumper) WithRepository(repository string) error {
	b.repository = repository
	return nil
}

func (b *mockBumper) WithPrerelease() {
	b.isPrerelease = true
}

func (b *mockBumper) WithTitle(title string) {
	b.title = title
}

func (b *mockBumper) WithTarget(target string) {
	b.target = target
}

func TestNewCmd(t *testing.T) {
	tests := map[string]struct {
		command        string
		wantRepo       string
		wantPrerelease bool
		wantTarget     string
		wantTitle      string
	}{
		"repository given": {
			command:  "bump -R johnmanjiro13/gh-bump",
			wantRepo: "johnmanjiro13/gh-bump",
		},
		"current repository": {
			command:  "bump",
			wantRepo: "",
		},
		"with prerelease": {
			command:        "bump --prerelease",
			wantPrerelease: true,
		},
		"with target": {
			command:    "bump --target feature",
			wantTarget: "feature",
		},
		"with title": {
			command:   "bump -t test_title",
			wantTitle: "test_title",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bumper := &mockBumper{}
			cmd := NewCmd(bumper)
			cmd.SetArgs(strings.Split(tt.command, " ")[1:])

			assert.NoError(t, cmd.Execute())
			assert.Equal(t, tt.wantRepo, bumper.repository)
			assert.Equal(t, tt.wantPrerelease, bumper.isPrerelease)
			assert.Equal(t, tt.wantTarget, bumper.target)
			assert.Equal(t, tt.wantTitle, bumper.title)
		})
	}
}
