package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBumper struct {
	repository string
	title      string
}

func (b *mockBumper) Bump() error {
	return nil
}

func (b *mockBumper) WithRepository(repository string) error {
	b.repository = repository
	return nil
}

func (b *mockBumper) WithTitle(title string) {
	b.title = title
}

func TestNewCmd(t *testing.T) {
	tests := map[string]struct {
		command   string
		wantRepo  string
		wantTitle string
	}{
		"repository given": {
			command:  "bump -R johnmanjiro13/gh-bump",
			wantRepo: "johnmanjiro13/gh-bump",
		},
		"current repository": {
			command:  "bump",
			wantRepo: "",
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
			assert.Equal(t, tt.wantTitle, bumper.title)
		})
	}
}
