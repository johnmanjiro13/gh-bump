package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBumper struct {
	repository         string
	isDraft            bool
	isPrerelease       bool
	discussionCategory string
	generateNotes      bool
	notes              string
	notesFilename      string
	target             string
	title              string
}

func (b *mockBumper) Bump() error {
	return nil
}

func (b *mockBumper) WithRepository(repository string) error {
	b.repository = repository
	return nil
}

func (b *mockBumper) WithDraft() {
	b.isDraft = true
}

func (b *mockBumper) WithPrerelease() {
	b.isPrerelease = true
}

func (b *mockBumper) WithDiscussionCategory(category string) {
	b.discussionCategory = category
}

func (b *mockBumper) WithGenerateNotes() {
	b.generateNotes = true
}

func (b *mockBumper) WithNotes(notes string) {
	b.notes = notes
}

func (b *mockBumper) WithNotesFile(filename string) {
	b.notesFilename = filename
}

func (b *mockBumper) WithTitle(title string) {
	b.title = title
}

func (b *mockBumper) WithTarget(target string) {
	b.target = target
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		command                string
		wantRepo               string
		wantDraft              bool
		wantPrerelease         bool
		wantDiscussionCategory string
		wantGenerateNotes      bool
		wantNotes              string
		wantNotesFilename      string
		wantTarget             string
		wantTitle              string
	}{
		"repository given": {
			command:  "bump -R johnmanjiro13/gh-bump",
			wantRepo: "johnmanjiro13/gh-bump",
		},
		"current repository": {
			command:  "bump",
			wantRepo: "",
		},
		"with draft": {
			command:   "bump --draft",
			wantDraft: true,
		},
		"with prerelease": {
			command:        "bump --prerelease",
			wantPrerelease: true,
		},
		"with discussion category": {
			command:                "bump --discussion-category category!",
			wantDiscussionCategory: "category!",
		},
		"with generate-notes": {
			command:           "bump --generate-notes",
			wantGenerateNotes: true,
		},
		"with notes": {
			command:   "bump --notes release",
			wantNotes: "release",
		},
		"with notes file": {
			command:           "bump --notes-file filename",
			wantNotesFilename: "filename",
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
			cmd := New(bumper)
			cmd.SetArgs(strings.Split(tt.command, " ")[1:])

			assert.NoError(t, cmd.Execute())
			assert.Equal(t, tt.wantRepo, bumper.repository)
			assert.Equal(t, tt.wantDraft, bumper.isDraft)
			assert.Equal(t, tt.wantPrerelease, bumper.isPrerelease)
			assert.Equal(t, tt.wantDiscussionCategory, bumper.discussionCategory)
			assert.Equal(t, tt.wantGenerateNotes, bumper.generateNotes)
			assert.Equal(t, tt.wantNotes, bumper.notes)
			assert.Equal(t, tt.wantNotesFilename, bumper.notesFilename)
			assert.Equal(t, tt.wantTarget, bumper.target)
			assert.Equal(t, tt.wantTitle, bumper.title)
		})
	}
}
